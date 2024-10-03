package marketdata

import (
	"context"
	"time"

	"studentgit.kata.academy/quant/torque/config"

	"github.com/mailru/easyjson"
	"studentgit.kata.academy/quant/torque/internal/adapters/websocket_clients/types/gateio"
	"studentgit.kata.academy/quant/torque/internal/domain"
	"studentgit.kata.academy/quant/torque/internal/domain/f"
	"studentgit.kata.academy/quant/torque/pkg/fp/slices"
	"studentgit.kata.academy/quant/torque/pkg/logster"
	websocketclient "studentgit.kata.academy/quant/torque/pkg/websocket_client"
)

type GateIOMarketDataWebsocket struct {
	logger                 logster.Logger
	client                 *websocketclient.Client
	instrumentsToSubscribe []domain.GateIOInstrument
	level                  gateio.OrderBookLevel
	interval               gateio.OrderBookInterval
	orderBookUpdate        chan domain.GateIOOrderBook
}

func NewGateIOMarketDataWebsocket(logger logster.Logger, config config.MarketDataConfig) *GateIOMarketDataWebsocket {
	return &GateIOMarketDataWebsocket{
		logger: logger.WithField(f.Module, "gateio_market_data_websocket"),
		client: websocketclient.NewClient(config.URL, logger,
			//nolint:gomnd // We want to reconnect if we don't have market data for 10 seconds
			websocketclient.WithReadDeadline(10*time.Second),
		),
		instrumentsToSubscribe: slices.Map(config.InstrumentsToSubscribe, func(i string) domain.GateIOInstrument {
			return domain.GateIOSpotInstrument(i)
		}),
		level:           gateio.Level50,       //  config.Level,
		interval:        gateio.Interval100ms, //  config.Interval,
		orderBookUpdate: make(chan domain.GateIOOrderBook),
	}
}

func (c *GateIOMarketDataWebsocket) Run(ctx context.Context) error {
	err := c.SubscribeInstruments(c.instrumentsToSubscribe)
	if err != nil {
		return err
	}
	go c.dispatchRequests(ctx)
	c.logger.Infof("gateio market data websocket is running")
	return nil
}

func (c *GateIOMarketDataWebsocket) Init(ctx context.Context) error {
	err := c.client.Connect(ctx)
	if err != nil {
		return err
	}
	c.setReconnectInstruments()
	return nil
}

func (c *GateIOMarketDataWebsocket) Close() {
	c.logger.Infof("closing gateio market data websocket...")
	c.client.Shutdown()
	close(c.orderBookUpdate)
	c.logger.Infof("gateio market data websocket closed")
}

func (c *GateIOMarketDataWebsocket) setReconnectInstruments() {
	c.client.SendMessageAfterReconnect([][]byte{c.makeMDRequest(c.instrumentsToSubscribe, true)})
}

func (c *GateIOMarketDataWebsocket) SubscribeInstruments(instruments []domain.GateIOInstrument) error {
	return c.client.Send(c.makeMDRequest(instruments, true))
}

func (c *GateIOMarketDataWebsocket) UnsubscribeInstruments(instruments []domain.GateIOInstrument) error {
	return c.client.Send(c.makeMDRequest(instruments, false))
}

func (c *GateIOMarketDataWebsocket) GetOrderBookUpdate() <-chan domain.GateIOOrderBook {
	return c.orderBookUpdate
}

func (c *GateIOMarketDataWebsocket) dispatchRequests(ctx context.Context) {
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

func (c *GateIOMarketDataWebsocket) makeMDRequest(instruments []domain.GateIOInstrument, subscribe bool) []byte {
	payload := make([][3]string, 0, len(instruments))
	for _, instrument := range instruments {
		batch := [3]string{
			string(instrument.Symbol()),
			string(c.level),
			string(c.interval),
		}
		payload = append(payload, batch)
	}

	eventType := gateio.EventTypeUnsubscribe
	if subscribe {
		eventType = gateio.EventTypeSubscribe
	}
	msg := gateio.MDRequest{
		Time:    time.Now().Unix(),
		Channel: gateio.SpotOrderBook,
		Event:   eventType,
		Payload: payload,
	}
	msgBytes, err := easyjson.Marshal(msg)
	if err != nil {
		c.logger.Errorf("error while marshaling message: %s", err)
		return nil
	}
	return msgBytes
}

func (c *GateIOMarketDataWebsocket) decodeEvent(msg []byte) {
	var orderBookUpdate gateio.OrderBookUpdate
	err := easyjson.Unmarshal(msg, &orderBookUpdate)
	if err != nil {
		c.logger.WithError(err).Errorf("error while unmarshalling message")
	}
	if orderBookUpdate.Event == gateio.EventTypeUpdate {
		domainBook := gateio.OrderBookToDomain(orderBookUpdate)
		c.orderBookUpdate <- domainBook
		return
	}
	var mdResponse gateio.CommonResponse
	err = easyjson.Unmarshal(msg, &mdResponse)
	if err != nil {
		c.logger.WithError(err).Errorf("error while unmarshalling message")
	}
	if mdResponse.Error.Message != "" {
		c.logger.WithField(f.MDResponse, mdResponse).Errorf("md request failed")
	} else {
		c.logger.WithField(f.MDResponse, mdResponse).Infof("md request succeeded")
	}
}
