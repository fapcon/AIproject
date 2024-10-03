package marketdata

import (
	"context"
	"time"

	"studentgit.kata.academy/quant/torque/config"

	"github.com/mailru/easyjson"
	"studentgit.kata.academy/quant/torque/internal/adapters/websocket_clients/types/okx"
	"studentgit.kata.academy/quant/torque/internal/domain"
	"studentgit.kata.academy/quant/torque/internal/domain/f"
	"studentgit.kata.academy/quant/torque/pkg/fp/slices"
	"studentgit.kata.academy/quant/torque/pkg/logster"
	websocketclient "studentgit.kata.academy/quant/torque/pkg/websocket_client"
)

type OKXMarketDataWebsocket struct {
	logger                 logster.Logger
	client                 *websocketclient.Client
	instrumentsToSubscribe []domain.OKXInstrument
	orderBookUpdate        chan domain.OKXOrderBook
}

func NewOKXMarketDataWebsocket(logger logster.Logger, config config.MarketDataConfig) *OKXMarketDataWebsocket {
	return &OKXMarketDataWebsocket{
		logger: logger.WithField(f.Module, "okx_market_data_websocket"),
		client: websocketclient.NewClient(config.URL, logger,
			//nolint:gomnd // We want to reconnect if we don't have market data for 10 seconds
			websocketclient.WithReadDeadline(10*time.Second),
		),
		instrumentsToSubscribe: slices.Map(config.InstrumentsToSubscribe, func(i string) domain.OKXInstrument {
			return domain.OKXSpotInstrument(i)
		}),
		orderBookUpdate: make(chan domain.OKXOrderBook),
	}
}

func (c *OKXMarketDataWebsocket) Run(ctx context.Context) error {
	err := c.SubscribeInstruments(c.instrumentsToSubscribe)
	if err != nil {
		return err
	}
	go c.dispatchRequests(ctx)
	c.logger.Infof("market data websocket is running")
	return nil
}

func (c *OKXMarketDataWebsocket) Init(ctx context.Context) error {
	err := c.client.Connect(ctx)
	if err != nil {
		return err
	}
	c.setReconnectInstruments()
	return nil
}

func (c *OKXMarketDataWebsocket) Close() {
	c.logger.Infof("closing market data websocket...")
	c.client.Shutdown()
	close(c.orderBookUpdate)
	c.logger.Infof("market data websocket closed")
}

func (c *OKXMarketDataWebsocket) setReconnectInstruments() {
	c.client.SendMessageAfterReconnect([][]byte{c.makeMDRequest(c.instrumentsToSubscribe, true)})
}

func (c *OKXMarketDataWebsocket) SubscribeInstruments(instruments []domain.OKXInstrument) error {
	return c.client.Send(c.makeMDRequest(instruments, true))
}

func (c *OKXMarketDataWebsocket) UnsubscribeInstruments(instruments []domain.OKXInstrument) error {
	return c.client.Send(c.makeMDRequest(instruments, false))
}

func (c *OKXMarketDataWebsocket) GetOrderBookUpdate() <-chan domain.OKXOrderBook {
	return c.orderBookUpdate
}

func (c *OKXMarketDataWebsocket) dispatchRequests(ctx context.Context) {
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

func (c *OKXMarketDataWebsocket) makeMDRequest(instruments []domain.OKXInstrument, subscribe bool) []byte {
	var args []okx.Arg
	for _, instrument := range instruments {
		args = append(args, okx.Arg{
			Channel: okx.BooksChannel,
			InstID:  string(instrument.Symbol()),
		})
	}
	eventType := okx.EventTypeUnsubscribe
	if subscribe {
		eventType = okx.EventTypeSubscribe
	}
	msg := okx.MDRequest{
		Event: eventType,
		Arg:   args,
	}
	msgBytes, err := easyjson.Marshal(msg)
	if err != nil {
		c.logger.Errorf("error while marshaling message: %s", err)
		return nil
	}
	return msgBytes
}

func (c *OKXMarketDataWebsocket) decodeEvent(msg []byte) {
	var orderBookUpdate okx.OrderBookUpdate
	err := easyjson.Unmarshal(msg, &orderBookUpdate)
	if err != nil {
		c.logger.WithError(err).Errorf("error while unmarshalling message")
	}
	if orderBookUpdate.Action == okx.EventTypeUpdate || orderBookUpdate.Action == okx.EventTypeSnapshot {
		domainBook := okx.OrderBookToDomain(orderBookUpdate)
		c.orderBookUpdate <- domainBook
		return
	}
	var mdResponse okx.MDResponse
	err = easyjson.Unmarshal(msg, &mdResponse)
	if err != nil {
		c.logger.WithError(err).Errorf("error while unmarshalling message")
	}
	if mdResponse.Event == okx.EventTypeError {
		c.logger.WithField(f.MDResponse, mdResponse).Errorf("md request failed")
	} else {
		c.logger.WithField(f.MDResponse, mdResponse).Infof("md request succeeded")
	}
}
