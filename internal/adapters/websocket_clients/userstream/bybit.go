package userstream

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"studentgit.kata.academy/quant/torque/config"

	"github.com/mailru/easyjson"
	"studentgit.kata.academy/quant/torque/internal/adapters/websocket_clients/types/bybit"
	"studentgit.kata.academy/quant/torque/internal/domain"
	"studentgit.kata.academy/quant/torque/internal/domain/f"
	"studentgit.kata.academy/quant/torque/pkg/logster"
	websocketclient "studentgit.kata.academy/quant/torque/pkg/websocket_client"
)

type BybitUserStreamWebsocket struct {
	logger          logster.Logger
	client          *websocketclient.Client
	heartbeatTicker *time.Ticker
	apiKey          string
	apiSecret       string
	orderUpdates    chan domain.BybitUserStreamOrderUpdate
}

const (
	freqHeartbeatInSec = 20
	expiration         = 1000000000
)

func NewBybitUserStreamWebsocket(
	logger logster.Logger, config config.UserStreamConfig,
) *BybitUserStreamWebsocket {
	return &BybitUserStreamWebsocket{
		logger:          logger.WithField(f.Module, "bybit_user_stream_websocket"),
		client:          websocketclient.NewClient(config.URL, logger),
		heartbeatTicker: time.NewTicker(freqHeartbeatInSec * time.Second),
		apiKey:          config.APIKey,
		apiSecret:       config.APISecret,
		orderUpdates:    make(chan domain.BybitUserStreamOrderUpdate),
	}
}

func (b *BybitUserStreamWebsocket) Run(_ context.Context) error {
	defer b.logger.Infof("user stream websocket is running")
	return b.SubscribeChannels()
}

func (b *BybitUserStreamWebsocket) SetReconnectMessage() {
	b.client.SendMessageAfterReconnect([][]byte{b.getAuth(), b.getSubscribeChannels()})
}

func (b *BybitUserStreamWebsocket) Init(ctx context.Context) error {
	err := b.client.Connect(ctx)
	if err != nil {
		return err
	}
	err = b.SendAuth()
	if err != nil {
		return err
	}
	time.Sleep(1 * time.Second) // need to wait for auth response
	go func() {
		for range b.heartbeatTicker.C {
			err = b.SendHeartbeat()
			if err != nil {
				b.logger.WithError(err).Errorf("Failed to send ping")
			}
		}
	}()
	go b.dispatchRequests(ctx)
	b.SetReconnectMessage()
	return nil
}

func (b *BybitUserStreamWebsocket) Close() {
	b.logger.Infof("closing user stream websocket...")
	b.heartbeatTicker.Stop()
	close(b.orderUpdates)
	b.client.Shutdown()
	b.logger.Infof("user stream websocket closed")
}

func (b *BybitUserStreamWebsocket) GetOrderUpdates() <-chan domain.BybitUserStreamOrderUpdate {
	return b.orderUpdates
}

func (b *BybitUserStreamWebsocket) dispatchRequests(ctx context.Context) {
	for {
		select {
		case msg, ok := <-b.client.GetMessages():
			if !ok {
				b.logger.Infof("message channel closed, stopping dispatching requests")
				return
			}

			b.decodeMessage(msg)

		case <-ctx.Done():
			b.logger.Infof("stopping dispatching requests")
			b.Close()
			return
		}
	}
}

func (b *BybitUserStreamWebsocket) getAuth() []byte {
	expires := time.Now().UnixNano() + expiration
	signature := b.generateSignature(fmt.Sprintf("GET/realtime%d", expires))
	data := bybit.AuthMessage{
		Op:   bybit.OpTypeAuth,
		Args: []any{b.apiKey, expires, signature},
	}
	jsonData, _ := easyjson.Marshal(data)
	return jsonData
}

func (b *BybitUserStreamWebsocket) getSubscribeChannels() []byte {
	data := bybit.SubscriptionMessage{
		Op:   bybit.OpTypeSubscribe,
		Args: []bybit.Topic{bybit.OrderTopic},
	}
	jsonData, _ := easyjson.Marshal(data)
	return jsonData
}

func (b *BybitUserStreamWebsocket) getHeartbeat() []byte {
	data := bybit.HeartbeatMessage{
		Op: string(bybit.OpTypeRequestHeartbeat),
	}
	jsonData, _ := easyjson.Marshal(data)
	return jsonData
}

func (b *BybitUserStreamWebsocket) SendAuth() error {
	return b.client.Send(b.getAuth())
}

func (b *BybitUserStreamWebsocket) SubscribeChannels() error {
	return b.client.Send(b.getSubscribeChannels())
}

func (b *BybitUserStreamWebsocket) SendHeartbeat() error {
	return b.client.Send(b.getHeartbeat())
}

func (b *BybitUserStreamWebsocket) generateSignature(data string) string {
	mac := hmac.New(sha256.New, []byte(b.apiSecret))
	mac.Write([]byte(data))
	return hex.EncodeToString(mac.Sum(nil))
}

func (b *BybitUserStreamWebsocket) decodeMessage(msg []byte) {
	var orderUpdate bybit.OrderUpdate
	err := easyjson.Unmarshal(msg, &orderUpdate)
	if err != nil {
		b.logger.WithError(err).Errorf("error while unmarshalling message")
	}
	switch {
	case orderUpdate.Op == bybit.OpTypeResponseHeartbeat:
		b.logger.WithField(f.BybitHeartbeatResponse, orderUpdate).Infof("heartbeat successful")
	case orderUpdate.Topic == bybit.OrderTopic && len(orderUpdate.Data) > 0:
		b.logger.WithField(f.OrderUpdate, orderUpdate).Infof("received order update")
		domainOrders := bybit.OrderToDomain(orderUpdate)
		for _, order := range domainOrders {
			b.orderUpdates <- order
		}
		return
	case !orderUpdate.Success:
		b.logger.WithField(f.BybitError, orderUpdate).Errorf("request failed")
	case orderUpdate.Op == bybit.OpTypeAuth:
		b.logger.WithField(f.BybitAuthResponse, orderUpdate).Infof("auth successful")
	case orderUpdate.Op == bybit.OpTypeSubscribe:
		b.logger.WithField(f.BybitUserStreamSubscribeResponse, orderUpdate).Infof("subscribed successfully")
	}
}
