package storage

import (
	"context"
	"time"

	"studentgit.kata.academy/quant/torque/config"

	"github.com/shopspring/decimal"
	"studentgit.kata.academy/quant/torque/internal/domain"
	"studentgit.kata.academy/quant/torque/pkg/checksum"
	"studentgit.kata.academy/quant/torque/pkg/logster"
	"studentgit.kata.academy/quant/torque/pkg/metrics"
)

const (
	cronIntervalChecksum = time.Second * 20
	cronInterval         = time.Second * 25
	deadMensSwitch       = time.Second * 30
	cancelOrderWaitTime  = time.Millisecond * 100
)

func (b *Builder) WithOKXAPIClient(config config.HTTPClient) *Builder {
	b.buildCap(CapabilityOKXAPIClient)
	b.build(func(s *Storage) {
		s.clients.RegisterOKXAPIClient(config)
	})
	return b
}

func (s *Storage) GetOKXOrderBookSnapshot(
	ctx context.Context, symbol domain.OKXSpotInstrument, limit int,
) (domain.OKXOrderBook, error) {
	s.mustHave(CapabilityOKXAPIClient)
	response, err := s.clients.OKXApi.GetOrderBook(ctx, symbol, limit)
	return response, err
}

func (s *Storage) GetOKXBalance(ctx context.Context) ([]domain.Balance, error) {
	s.mustHave(CapabilityOKXAPIClient)
	response, err := s.clients.OKXApi.GetBalances(ctx)
	return response, err
}

func (s *Storage) OKXCreateOrder(ctx context.Context, order domain.OrderRequest) (domain.OKXOrderResponseAck, error) {
	s.mustHave(CapabilityOKXAPIClient)
	return s.clients.OKXApi.CreateSpotOrder(ctx, order)
}

func (s *Storage) OKXCreateFuturesOrder(
	ctx context.Context, order domain.OrderRequest,
) (domain.OKXOrderResponseAck, error) {
	s.mustHave(CapabilityOKXAPIClient)
	return s.clients.OKXApi.CreateFuturesOrder(ctx, order)
}

func (s *Storage) OKXCancelOrder(
	ctx context.Context, symbol domain.OKXSpotInstrument, orderID domain.OrderID, clientOrderID domain.ClientOrderID,
) (domain.OKXOrderResponseAck, error) {
	s.mustHave(CapabilityOKXAPIClient)
	return s.clients.OKXApi.CancelOrder(ctx, symbol, orderID, clientOrderID)
}

func (s *Storage) OKXGetSpotOpenOrders(ctx context.Context) ([]domain.OpenOrder, error) {
	s.mustHave(CapabilityOKXAPIClient)
	return s.clients.OKXApi.GetSpotOpenOrders(ctx)
}

func (s *Storage) OKXGetFuturesOpenOrders(ctx context.Context) ([]domain.OpenOrder, error) {
	s.mustHave(CapabilityOKXAPIClient)
	return s.clients.OKXApi.GetFuturesOpenOrders(ctx)
}

func (s *Storage) OKXGetOpenOrder(
	ctx context.Context, symbol domain.OKXSpotInstrument, orderID domain.OrderID, clientOrderID domain.ClientOrderID,
) (domain.OpenOrder, error) {
	s.mustHave(CapabilityOKXAPIClient)
	return s.clients.OKXApi.GetOrder(ctx, symbol, orderID, clientOrderID)
}

func (s *Storage) OKXMoveOrder(
	ctx context.Context, symbol domain.OKXSpotInstrument, orderID domain.OrderID, clientOrderID domain.ClientOrderID,
	newSize, newPrice decimal.Decimal,
) error {
	s.mustHave(CapabilityOKXAPIClient)
	return s.clients.OKXApi.AmendOrder(ctx, symbol, orderID, clientOrderID, newSize, newPrice)
}

func (s *Storage) OKXGetSpotOrderHistory(ctx context.Context, startTime, endTime time.Time) ([]domain.Order, error) {
	s.mustHave(CapabilityOKXAPIClient)
	return s.clients.OKXApi.GetSpotOrderHistory(ctx, startTime, endTime)
}

func (s *Storage) OKXGetSpotFillsHistory(ctx context.Context, startTime, endTime time.Time) ([]domain.Order, error) {
	s.mustHave(CapabilityOKXAPIClient)
	return s.clients.OKXApi.GetSpotFillsHistory(ctx, startTime, endTime)
}

func (s *Storage) OKXGetFuturesOrderHistory(ctx context.Context, startTime, endTime time.Time) ([]domain.Order, error) {
	s.mustHave(CapabilityOKXAPIClient)
	return s.clients.OKXApi.GetFuturesOrderHistory(ctx, startTime, endTime)
}

func (b *Builder) WithOKXMarketDataWebsocketClient(config config.MarketDataConfig) *Builder {
	b.buildCap(CapabilityOKXMarketDataWebsocketClient)
	b.build(func(s *Storage) {
		s.websocketClients.RegisterOKXMarketDataWebsocketClient(config)
	})
	return b
}

func (b *Builder) WithOKXUserStreamWebsocketClient(config config.UserStreamConfig) *Builder {
	b.buildCap(CapabilityOKXUserStreamWebsocketClient)
	b.build(func(s *Storage) {
		s.websocketClients.RegisterOKXUserStreamWebsocketClient(config)
	})
	return b
}

func (b *Builder) WithOKXUpdateUserStreamSubscribeMessage() *Builder {
	b.runCron(cronInterval, func(s *Storage, _ context.Context) {
		s.mustHave(CapabilityOKXUserStreamWebsocketClient)
		s.websocketClients.OKXUserStreamWebsocket.SetReconnectMessage()
	})
	return b
}

//nolint:gomnd // we want to send order to keep our websocket alive
func (b *Builder) WithOKXFakeTradeMessage() *Builder {
	b.runCron(cronInterval, func(s *Storage, ctx context.Context) {
		s.mustHave(CapabilityOKXAPIClient)
		s.mustHave(CapabilityOKXSpotTradingWebsocketClient)
		err := s.websocketClients.OKXSpotTradingClient.CreateOrder(ctx, domain.OrderRequest{
			Symbol:        "BTC-USDT",
			Side:          domain.SideBuy,
			Type:          domain.OrderTypeLimit,
			Quantity:      decimal.NewFromFloat(0.00001), //
			Price:         decimal.NewFromInt(10000),
			ClientOrderID: domain.FakeOrder,
			TimeInForce:   "",
		})
		if err != nil {
			s.logger.WithError(err).Infof("Failed to create fake order")
			return
		}
		time.Sleep(cancelOrderWaitTime)
		err = s.websocketClients.OKXSpotTradingClient.CancelOrder(
			ctx, domain.OKXSpotInstrument("BTC-USDT"), "", domain.FakeOrder,
		)
		if err != nil {
			s.logger.WithError(err).Infof("Failed to cancel fake order")
		}
	})
	return b
}

func (b *Builder) WithOKXDeadMansSwitch() *Builder {
	b.runCron(cronInterval, func(s *Storage, ctx context.Context) {
		s.mustHave(CapabilityOKXAPIClient)
		cancelTime, err := s.clients.OKXApi.DeadMensSwitch(ctx, deadMensSwitch)
		if err != nil {
			s.logger.WithError(err).Errorf("failed to send dead mens switch request")
			return
		}
		s.logger.Debugf("dead mens switch will be triggered at %s", cancelTime)
	})
	return b
}

func (b *Builder) WithOKXChecksumCheck() *Builder {
	b.runCron(cronIntervalChecksum, func(s *Storage, ctx context.Context) {
		books := s.GetOrderBookList(ctx)
		for _, book := range books {
			if book.Exchange != domain.OKXExchange {
				continue
			}

			localChecksum := checksum.CalculateLocalChecksum(book.Asks, book.Bids)
			if localChecksum != uint32(book.Checksum) {
				s.logger.WithFields(logster.Fields{
					"exchange":   book.Exchange,
					"instrument": book.Instrument,
				}).Errorf("checksum mismatch")
				s.DeleteOrderBook(ctx, domain.LocalInstrument{
					Exchange: book.Exchange,
					Symbol:   book.Instrument,
				})
			}
		}
	})
	return b
}

func (s *Storage) OKXSubscribeInstruments(_ context.Context, instruments []domain.OKXInstrument) error {
	s.mustHave(CapabilityOKXMarketDataWebsocketClient)
	return s.websocketClients.OKXMarketDataWebsocket.SubscribeInstruments(instruments)
}

func (s *Storage) OKXUnsubscribeInstruments(_ context.Context, instruments []domain.OKXInstrument) error {
	s.mustHave(CapabilityOKXMarketDataWebsocketClient)
	return s.websocketClients.OKXMarketDataWebsocket.UnsubscribeInstruments(instruments)
}

func (s *Storage) GetOKXBookUpdateMessages() <-chan domain.OKXOrderBook {
	s.mustHave(CapabilityOKXMarketDataWebsocketClient)
	return s.websocketClients.OKXMarketDataWebsocket.GetOrderBookUpdate()
}

func (s *Storage) GetOKXOrderUpdates() <-chan domain.OKXUserStreamOrderUpdate {
	s.mustHave(CapabilityOKXUserStreamWebsocketClient)
	return s.websocketClients.OKXUserStreamWebsocket.GetOrderUpdates()
}

func (b *Builder) WithOKXTradingWebsocketClient(config config.UserStreamConfig) *Builder {
	b.buildCap(CapabilityOKXSpotTradingWebsocketClient)
	b.build(func(s *Storage) {
		s.websocketClients.RegisterOKXTradingWebsocketClient(config)
	})
	return b
}

func (b *Builder) WithOKXTradingWebsocketClientsMetrics(collector *metrics.Collector) *Builder {
	b.build(func(s *Storage) {
		s.mustHave(CapabilityOKXSpotTradingWebsocketClient)
		s.websocketClients.OKXSpotTradingClient.WithMetrics(collector)
	})
	return b
}

func (s *Storage) GetOKXSpotCreateOrderAckMessages() <-chan domain.OKXOrderResponseAck {
	s.mustHave(CapabilityOKXSpotTradingWebsocketClient)
	return s.websocketClients.OKXSpotTradingClient.GetCreateOrderAck()
}

func (s *Storage) GetOKXSpotCreateOrderFailedMessages() <-chan domain.OKXOrderResponseAck {
	s.mustHave(CapabilityOKXSpotTradingWebsocketClient)
	return s.websocketClients.OKXSpotTradingClient.GetCreateOrderFailed()
}

func (s *Storage) GetOKXSpotModifyOrderAckMessages() <-chan domain.OKXOrderResponseAck {
	s.mustHave(CapabilityOKXSpotTradingWebsocketClient)
	return s.websocketClients.OKXSpotTradingClient.GetModifyOrderAck()
}

func (s *Storage) GetOKXSpotModifyOrderFailedMessages() <-chan domain.OKXOrderResponseAck {
	s.mustHave(CapabilityOKXSpotTradingWebsocketClient)
	return s.websocketClients.OKXSpotTradingClient.GetModifyOrderFailed()
}

func (s *Storage) GetOKXSpotCancelOrderAckMessages() <-chan domain.OKXOrderResponseAck {
	s.mustHave(CapabilityOKXSpotTradingWebsocketClient)
	return s.websocketClients.OKXSpotTradingClient.GetCancelOrderAck()
}

func (s *Storage) GetOKXSpotCancelOrderFailedMessages() <-chan domain.OKXOrderResponseAck {
	s.mustHave(CapabilityOKXSpotTradingWebsocketClient)
	return s.websocketClients.OKXSpotTradingClient.GetCancelOrderFailed()
}

func (s *Storage) GetOKXSpotTradingWebsocketStatus() bool {
	s.mustHave(CapabilityOKXSpotTradingWebsocketClient)
	return s.websocketClients.OKXSpotTradingClient.IsConnected()
}

func (s *Storage) CreateOKXSpotTradingWebsocketOrder(
	ctx context.Context, order domain.OrderRequest,
) error {
	s.mustHave(CapabilityOKXSpotTradingWebsocketClient)
	return s.websocketClients.OKXSpotTradingClient.CreateOrder(ctx, order)
}

func (s *Storage) CancelOKXSpotTradingWebsocketOrder(
	ctx context.Context, symbol domain.OKXSpotInstrument, orderID domain.OrderID, clientOrderID domain.ClientOrderID,
) error {
	s.mustHave(CapabilityOKXSpotTradingWebsocketClient)
	return s.websocketClients.OKXSpotTradingClient.CancelOrder(ctx, symbol, orderID, clientOrderID)
}

func (s *Storage) ModifyOKXSpotTradingWebsocketOrder(
	ctx context.Context, symbol domain.OKXSpotInstrument, orderID domain.OrderID, clientOrderID domain.ClientOrderID,
	newSize, newPrice decimal.Decimal,
) error {
	s.mustHave(CapabilityOKXSpotTradingWebsocketClient)
	return s.websocketClients.OKXSpotTradingClient.ModifyOrder(ctx, symbol, orderID, clientOrderID, newSize, newPrice)
}
