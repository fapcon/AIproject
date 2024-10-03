package marketdata

import (
	"context"
	"fmt"
	"time"

	"studentgit.kata.academy/quant/torque/config"

	"github.com/mailru/easyjson"
	"studentgit.kata.academy/quant/torque/internal/adapters/websocket_clients/types/bybit"
	"studentgit.kata.academy/quant/torque/internal/domain"
	"studentgit.kata.academy/quant/torque/internal/domain/f"
	"studentgit.kata.academy/quant/torque/pkg/fp/slices"
	"studentgit.kata.academy/quant/torque/pkg/logster"
	websocketclient "studentgit.kata.academy/quant/torque/pkg/websocket_client"
)

type BybitMarketDataWebsocket struct {
	logger                 logster.Logger
	client                 *websocketclient.Client
	instrumentsToSubscribe []domain.BybitInstrument
	orderBookUpdate        chan domain.BybitOrderBook
}

func NewBybitMarketDataWebsocket(logger logster.Logger, config config.MarketDataConfig) *BybitMarketDataWebsocket {
	return &BybitMarketDataWebsocket{
		logger: logger.WithField(f.Module, "bybit_market_data_websocket"),
		client: websocketclient.NewClient(config.URL, logger,
			//nolint:gomnd // We want to reconnect if we don't have market data for 10 seconds
			websocketclient.WithReadDeadline(10*time.Second),
		),
		instrumentsToSubscribe: slices.Map(config.InstrumentsToSubscribe, func(i string) domain.BybitInstrument {
			return domain.BybitSpotInstrument(i)
		}),
		orderBookUpdate: make(chan domain.BybitOrderBook),
	}
}

func (c *BybitMarketDataWebsocket) Run(ctx context.Context) error {
	err := c.SubscribeInstruments(c.instrumentsToSubscribe)
	if err != nil {
		return err
	}
	go c.dispatchRequests(ctx)
	c.logger.Infof("market data websocket is running")
	return nil
}

func (c *BybitMarketDataWebsocket) Init(ctx context.Context) error {
	err := c.client.Connect(ctx)
	if err != nil {
		return err
	}
	c.setReconnectInstruments()
	return nil
}

func (c *BybitMarketDataWebsocket) Close() {
	c.logger.Infof("closing market data websocket...")
	c.client.Shutdown()
	close(c.orderBookUpdate)
	c.logger.Infof("market data websocket closed")
}

func (c *BybitMarketDataWebsocket) setReconnectInstruments() {
	c.client.SendMessageAfterReconnect([][]byte{c.makeMDRequest(c.instrumentsToSubscribe, true)})
}

func (c *BybitMarketDataWebsocket) SubscribeInstruments(instruments []domain.BybitInstrument) error {
	return c.client.Send(c.makeMDRequest(instruments, true))
}

func (c *BybitMarketDataWebsocket) UnsubscribeInstruments(instruments []domain.BybitInstrument) error {
	return c.client.Send(c.makeMDRequest(instruments, false))
}

func (c *BybitMarketDataWebsocket) GetOrderBookUpdate() <-chan domain.BybitOrderBook {
	return c.orderBookUpdate
}

func (c *BybitMarketDataWebsocket) dispatchRequests(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			c.logger.Infof("stopping dispatching requests")
			c.Close()
			return
		case msg, ok := <-c.client.GetMessages():
			if !ok {
				c.logger.Infof("message channel closed, stopping dispatching requests")
				return
			}
			c.decodeEvent(msg)
		}
	}
}

func (c *BybitMarketDataWebsocket) makeMDRequest(instruments []domain.BybitInstrument, subscribe bool) []byte {
	args := make([]string, 0, len(instruments))
	for _, instrument := range instruments {
		args = append(args, fmt.Sprintf("orderbook.%s.%s", bybit.Level50, instrument.Symbol()))
	}

	eventType := bybit.OpTypeUnsubscribe
	if subscribe {
		eventType = bybit.OpTypeSubscribe
	}
	msg := bybit.MDRequest{
		Op:   eventType,
		Args: args,
	}
	msgBytes, err := easyjson.Marshal(msg)
	if err != nil {
		c.logger.Errorf("error while marshaling message: %s", err)
		return nil
	}
	return msgBytes
}

func (c *BybitMarketDataWebsocket) decodeEvent(msg []byte) {
	var orderBookUpdate bybit.OrderBookUpdate
	err := easyjson.Unmarshal(msg, &orderBookUpdate)
	if err != nil {
		c.logger.WithError(err).Errorf("error while unmarshalling message")
	}
	if orderBookUpdate.Type == bybit.OpTypeSnapshot || orderBookUpdate.Type == bybit.OpTypeDelta {
		domainBook := bybit.OrderBookToDomain(orderBookUpdate)
		c.orderBookUpdate <- domainBook
		return
	}
	var mdResponse bybit.MDResponse
	err = easyjson.Unmarshal(msg, &mdResponse)
	if err != nil {
		c.logger.WithError(err).Errorf("error while unmarshalling message")
	}
	if !mdResponse.Success {
		c.logger.WithField(f.MDResponse, mdResponse).Errorf("md request failed")
	} else {
		c.logger.WithField(f.MDResponse, mdResponse).Infof("md request succeeded")
	}
}
