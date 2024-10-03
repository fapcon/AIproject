package storage

import (
	"context"
	"time"

	"studentgit.kata.academy/quant/torque/config"

	"github.com/shopspring/decimal"
	"studentgit.kata.academy/quant/torque/internal/domain"
)

func (b *Builder) WithBybitAPIClient(config config.HTTPClient) *Builder {
	b.buildCap(CapabilityBybitAPIClient)
	b.build(func(s *Storage) {
		s.clients.RegisterBybitAPIClient(config)
	})
	return b
}

func (s *Storage) GetBybitOrderBookSnapshot(
	ctx context.Context, symbol domain.Symbol, limit int,
) (domain.OrderBook, error) {
	s.mustHave(CapabilityBybitAPIClient)
	response, err := s.clients.BybitAPI.GetOrderBook(ctx, symbol, limit)
	return response, err
}

func (s *Storage) GetBybitBalance(ctx context.Context) ([]domain.Balance, error) {
	s.mustHave(CapabilityBybitAPIClient)
	return s.clients.BybitAPI.GetBalances(ctx)
}

func (s *Storage) BybitCreateSpotOrder(
	ctx context.Context, order domain.OrderRequest,
) (domain.BybitOrderResponseAck, error) {
	s.mustHave(CapabilityBybitAPIClient)
	return s.clients.BybitAPI.CreateSpotOrder(ctx, order)
}

func (s *Storage) BybitCreateFuturesOrder(
	ctx context.Context, order domain.OrderRequest,
) (domain.BybitOrderResponseAck, error) {
	s.mustHave(CapabilityBybitAPIClient)
	return s.clients.BybitAPI.CreateFuturesOrder(ctx, order)
}

func (s *Storage) BybitCancelOrder(
	ctx context.Context, symbol domain.Symbol, orderID domain.OrderID, clientOrderID domain.ClientOrderID,
) (domain.BybitOrderResponseAck, error) {
	s.mustHave(CapabilityBybitAPIClient)
	return s.clients.BybitAPI.CancelOrder(ctx, symbol, orderID, clientOrderID)
}

func (s *Storage) BybitCancelSpotOrder(
	ctx context.Context, symbol domain.Symbol,
) error {
	s.mustHave(CapabilityBybitAPIClient)
	return s.clients.BybitAPI.CancelSpotOrders(ctx, symbol)
}

func (s *Storage) BybitCancelFuturesOrder(
	ctx context.Context, symbol domain.Symbol,
) error {
	s.mustHave(CapabilityBybitAPIClient)
	return s.clients.BybitAPI.CancelFuturesOrders(ctx, symbol)
}

func (s *Storage) BybitGetSpotOpenOrders(ctx context.Context, symbol domain.Symbol) ([]domain.OpenOrder, error) {
	s.mustHave(CapabilityBybitAPIClient)
	return s.clients.BybitAPI.GetSpotOpenOrders(ctx, symbol)
}

func (s *Storage) BybitGetFuturesOpenOrders(ctx context.Context, symbol domain.Symbol) ([]domain.OpenOrder, error) {
	s.mustHave(CapabilityBybitAPIClient)
	return s.clients.BybitAPI.GetFuturesOpenOrders(ctx, symbol)
}

func (s *Storage) BybitAmendSpotOrder(
	ctx context.Context, symbol domain.Symbol, orderID domain.OrderID, clientOrderID domain.ClientOrderID,
	newSize, newPrice decimal.Decimal,
) error {
	s.mustHave(CapabilityBybitAPIClient)
	return s.clients.BybitAPI.AmendSpotOrder(ctx, symbol, orderID, clientOrderID, newSize, newPrice)
}

func (s *Storage) BybitAmendFuturesOrder(
	ctx context.Context, symbol domain.Symbol, orderID domain.OrderID, clientOrderID domain.ClientOrderID,
	newSize, newPrice decimal.Decimal,
) error {
	s.mustHave(CapabilityBybitAPIClient)
	return s.clients.BybitAPI.AmendFuturesOrder(ctx, symbol, orderID, clientOrderID, newSize, newPrice)
}

func (s *Storage) BybitGetSpotOrderHistory(
	ctx context.Context, startTime, endTime time.Time,
) ([]domain.Order, error) {
	s.mustHave(CapabilityBybitAPIClient)
	return s.clients.BybitAPI.GetSpotOrderHistory(ctx, startTime, endTime)
}

func (s *Storage) BybitGetFuturesOrderHistory(
	ctx context.Context, startTime, endTime time.Time,
) ([]domain.Order, error) {
	s.mustHave(CapabilityBybitAPIClient)
	return s.clients.BybitAPI.GetFuturesOrderHistory(ctx, startTime, endTime)
}

func (b *Builder) WithBybitUpdateUserStreamSubscribeMessage() *Builder {
	b.runCron(cronInterval, func(s *Storage, _ context.Context) {
		s.mustHave(CapabilityBybitUserStreamWebsocketClient)
		s.websocketClients.BybitUserStreamWebsocket.SetReconnectMessage()
	})
	return b
}

func (b *Builder) WithBybitUserStreamWebsocketClient(config config.UserStreamConfig) *Builder {
	b.buildCap(CapabilityBybitUserStreamWebsocketClient)
	b.build(func(s *Storage) {
		s.websocketClients.RegisterBybitUserStreamWebsocketClient(config)
	})
	return b
}

func (s *Storage) GetBybitOrderUpdates() <-chan domain.BybitUserStreamOrderUpdate {
	s.mustHave(CapabilityBybitUserStreamWebsocketClient)
	return s.websocketClients.BybitUserStreamWebsocket.GetOrderUpdates()
}

func (b *Builder) WithBybitMarketDataWebsocketClient(config config.MarketDataConfig) *Builder {
	b.buildCap(CapabilityBybitMarketDataWebsocketClient)
	b.build(func(s *Storage) {
		s.websocketClients.RegisterBybitMarketDataWebsocketClient(config)
	})
	return b
}

func (s *Storage) BybitSubscribeInstruments(_ context.Context, instruments []domain.BybitInstrument) error {
	s.mustHave(CapabilityBybitMarketDataWebsocketClient)
	return s.websocketClients.BybitMarketDataWebsocket.SubscribeInstruments(instruments)
}

func (s *Storage) BybitUnsubscribeInstruments(_ context.Context, instruments []domain.BybitInstrument) error {
	s.mustHave(CapabilityBybitMarketDataWebsocketClient)
	return s.websocketClients.BybitMarketDataWebsocket.UnsubscribeInstruments(instruments)
}

func (s *Storage) GetBybitBookUpdateMessages() <-chan domain.BybitOrderBook {
	s.mustHave(CapabilityBybitMarketDataWebsocketClient)
	return s.websocketClients.BybitMarketDataWebsocket.GetOrderBookUpdate()
}
