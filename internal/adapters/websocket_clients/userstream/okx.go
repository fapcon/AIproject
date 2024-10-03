package userstream

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"

	"studentgit.kata.academy/quant/torque/config"

	"github.com/mailru/easyjson"
	"studentgit.kata.academy/quant/torque/internal/adapters/websocket_clients/types/okx"
	"studentgit.kata.academy/quant/torque/internal/domain"
	"studentgit.kata.academy/quant/torque/internal/domain/f"
	"studentgit.kata.academy/quant/torque/pkg/logster"
	websocketclient "studentgit.kata.academy/quant/torque/pkg/websocket_client"
)

type OKXUserStreamWebsocket struct {
	logger        logster.Logger
	client        *websocketclient.Client
	apiKey        string
	apiSecret     string
	apiPassphrase string
	orderUpdates  chan domain.OKXUserStreamOrderUpdate
}

func NewOkxUserStreamWebsocket(
	logger logster.Logger, config config.UserStreamConfig,
) *OKXUserStreamWebsocket {
	return &OKXUserStreamWebsocket{
		logger:        logger.WithField(f.Module, "okx_user_stream_websocket"),
		client:        websocketclient.NewClient(config.URL, logger),
		apiKey:        config.APIKey,
		apiSecret:     config.APISecret,
		apiPassphrase: config.APIPassphrase,
		orderUpdates:  make(chan domain.OKXUserStreamOrderUpdate),
	}
}

func (c *OKXUserStreamWebsocket) Run(_ context.Context) error {
	defer c.logger.Infof("user stream websocket is running")
	return c.SubscribeChannels(c.getSubscribeChannels())
}

func (c *OKXUserStreamWebsocket) SetReconnectMessage() {
	c.client.SendMessageAfterReconnect([][]byte{c.getAuth(), c.getSubscribeChannels()})
}

func (c *OKXUserStreamWebsocket) Init(ctx context.Context) error {
	err := c.client.Connect(ctx)
	if err != nil {
		return err
	}
	err = c.SendAuth(c.getAuth())
	if err != nil {
		return err
	}
	time.Sleep(1 * time.Second) // need to wait for auth response
	go c.dispatchRequests(ctx)
	c.SetReconnectMessage()
	return nil
}

func (c *OKXUserStreamWebsocket) Close() {
	c.logger.Infof("closing user stream websocket...")
	close(c.orderUpdates)
	c.client.Shutdown()
	c.logger.Infof("user stream websocket closed")
}

func (c *OKXUserStreamWebsocket) GetOrderUpdates() <-chan domain.OKXUserStreamOrderUpdate {
	return c.orderUpdates
}

func (c *OKXUserStreamWebsocket) dispatchRequests(ctx context.Context) {
	for {
		select {
		case msg, ok := <-c.client.GetMessages():
			if !ok {
				c.logger.Infof("message channel closed, stopping dispatching requests")
				return
			}

			c.decodeMessage(msg)

		case <-ctx.Done():
			c.logger.Infof("stopping dispatching requests")
			c.Close()
			return
		}
	}
}

func (c *OKXUserStreamWebsocket) getAuth() []byte {
	timestamp := time.Now().Unix()
	data := okx.LoginMessage{
		Op: okx.EventTypeLogin,
		Args: []okx.LoginArgs{
			{
				APIKey:     c.apiKey,
				Passphrase: c.apiPassphrase,
				Timestamp:  timestamp,
				Sign:       c.generateSignature(fmt.Sprintf("%dGET/users/self/verify", timestamp)),
			},
		},
	}
	jsonData, _ := easyjson.Marshal(data)
	return jsonData
}

func (c *OKXUserStreamWebsocket) getSubscribeChannels() []byte {
	data := okx.SubscriptionMessage{
		Op: okx.EventTypeSubscribe,
		Args: []okx.SubscriptionArgs{
			{
				Channel:  okx.OrdersChannel,
				InstType: okx.SPOT,
			},
		},
	}
	jsonData, _ := easyjson.Marshal(data)
	return jsonData
}

func (c *OKXUserStreamWebsocket) SendAuth(msg []byte) error {
	return c.client.Send(msg)
}

func (c *OKXUserStreamWebsocket) SubscribeChannels(msg []byte) error {
	return c.client.Send(msg)
}

func (c *OKXUserStreamWebsocket) generateSignature(data string) string {
	mac := hmac.New(sha256.New, []byte(c.apiSecret))
	mac.Write([]byte(data))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func (c *OKXUserStreamWebsocket) decodeMessage(msg []byte) {
	var orderUpdate okx.OrderUpdate
	err := easyjson.Unmarshal(msg, &orderUpdate)
	if err != nil {
		c.logger.WithError(err).Errorf("error while unmarshalling message")
	}
	switch {
	case orderUpdate.Event == okx.EventTypeError:
		c.logger.WithField(f.OKXError, orderUpdate).Errorf("request failed")
	case orderUpdate.Event == okx.EventTypeSubscribe:
		c.logger.WithField(f.OKXUserStreamSubscribeResponse, orderUpdate).Infof("request succeeded")
	case orderUpdate.Arg.Channel == okx.OrdersChannel && len(orderUpdate.Data) > 0:
		c.logger.WithField(f.OrderUpdate, orderUpdate).Infof("received order update")
		domainOrder := okx.OrderToDomain(orderUpdate)
		c.orderUpdates <- domainOrder
		return
	}
}
