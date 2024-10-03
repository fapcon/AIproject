package storage

import (
	"context"

	"studentgit.kata.academy/quant/torque/internal/domain"
)

func (b *Builder) WithPistonClient(addr string) *Builder {
	b.buildCap(CapabilityPistonClient)
	b.build(func(s *Storage) {
		s.grpcClients.RegisterPistonClient(addr)
	})
	return b
}

func (b *Builder) InitPistonClient() *Builder {
	b.init(func(s *Storage, ctx context.Context) error {
		s.mustHave(CapabilityPistonClient)
		err := s.grpcClients.PistonClient.Init(ctx)
		if err != nil {
			return err
		}
		return nil
	})
	return b
}

func (s *Storage) GetPistonCancelOrderMessages() <-chan domain.CancelOrder {
	s.mustHave(CapabilityPistonClient)
	return s.grpcClients.PistonClient.GetCancelOrderMessagesRequests()
}

func (s *Storage) GetPistonCreateOrderMessages() <-chan domain.AddOrder {
	s.mustHave(CapabilityPistonClient)
	return s.grpcClients.PistonClient.GetCreateOrderMessagesRequests()
}

func (s *Storage) GetPistonMoveOrderMessages() <-chan domain.MoveOrder {
	s.mustHave(CapabilityPistonClient)
	return s.grpcClients.PistonClient.GetMoveOrderMessagesRequests()
}

func (s *Storage) GetPistonSubscribeMessageRequests() <-chan domain.PistonSubscribeMarketDataMessage {
	s.mustHave(CapabilityPistonClient)
	return s.grpcClients.PistonClient.GetSubscribeMessageRequests()
}

func (s *Storage) GetPistonOpenOrdersMessageRequests() <-chan domain.PistonOpenOrderRequestMessage {
	s.mustHave(CapabilityPistonClient)
	return s.grpcClients.PistonClient.GetOpenOrdersMessageRequests()
}

func (s *Storage) GetPistonBalancesMessageRequests() <-chan domain.PistonBalancesRequestMessage {
	s.mustHave(CapabilityPistonClient)
	return s.grpcClients.PistonClient.GetBalancesMessageRequests()
}

func (s *Storage) GetPistonOrderHistoryMessageRequests() <-chan domain.PistonOrderHistoryRequestMessage {
	s.mustHave(CapabilityPistonClient)
	return s.grpcClients.PistonClient.GetOrderHistoryMessageRequests()
}

func (s *Storage) GetPistonInstrumentDetailsMessageRequests() <-chan domain.PistonInstrumentsDetailsRequestMessage {
	s.mustHave(CapabilityPistonClient)
	return s.grpcClients.PistonClient.GetInstrumentDetailsMessageRequests()
}

func (s *Storage) PistonSendOrderBook(ctx context.Context, ob domain.OrderBook) error {
	s.mustHave(CapabilityPistonClient)
	return s.grpcClients.PistonClient.SendOrderBook(ctx, ob)
}

func (s *Storage) PistonSendBalances(ctx context.Context, b domain.PistonBalancesMessage) error {
	s.mustHave(CapabilityPistonClient)
	return s.grpcClients.PistonClient.SendBalances(ctx, b)
}

func (s *Storage) PistonSendOpenOrders(ctx context.Context, o domain.PistonOpenOrdersMessage) error {
	s.mustHave(CapabilityPistonClient)
	return s.grpcClients.PistonClient.SendOpenOrders(ctx, o)
}

func (s *Storage) PistonSendOrderAdded(ctx context.Context, o domain.OrderAdded) error {
	s.mustHave(CapabilityPistonClient)
	return s.grpcClients.PistonClient.SendOrderAdded(ctx, o)
}

func (s *Storage) PistonSendOrderAddRejected(ctx context.Context, o domain.OrderAddReject) error {
	s.mustHave(CapabilityPistonClient)
	return s.grpcClients.PistonClient.SendOrderAddRejected(ctx, o)
}

func (s *Storage) PistonSendOrderMoved(ctx context.Context, o domain.OrderMoved) error {
	s.mustHave(CapabilityPistonClient)
	return s.grpcClients.PistonClient.SendOrderMoved(ctx, o)
}

func (s *Storage) PistonSendOrderMoveReject(ctx context.Context, o domain.OrderMoveReject) error {
	s.mustHave(CapabilityPistonClient)
	return s.grpcClients.PistonClient.SendOrderMoveReject(ctx, o)
}

func (s *Storage) PistonSendOrderCancelled(ctx context.Context, o domain.OrderCancelled) error {
	s.mustHave(CapabilityPistonClient)
	return s.grpcClients.PistonClient.SendOrderCancelled(ctx, o)
}

func (s *Storage) PistonSendOrderCancelReject(ctx context.Context, o domain.OrderCancelReject) error {
	s.mustHave(CapabilityPistonClient)
	return s.grpcClients.PistonClient.SendOrderCancelReject(ctx, o)
}

func (s *Storage) PistonSendOrderFilled(ctx context.Context, o domain.OrderFilled) error {
	s.mustHave(CapabilityPistonClient)
	return s.grpcClients.PistonClient.SendOrderFilled(ctx, o)
}

func (s *Storage) PistonSendOrderExecuted(ctx context.Context, o domain.OrderExecuted) error {
	s.mustHave(CapabilityPistonClient)
	return s.grpcClients.PistonClient.SendOrderExecuted(ctx, o)
}

func (s *Storage) PistonSendInstrumentDetails(ctx context.Context, details []domain.InstrumentDetails) error {
	s.mustHave(CapabilityPistonClient)
	return s.grpcClients.PistonClient.SendInstrumentDetails(ctx, details)
}

func (s *Storage) PistonRateLimitExceedError(ctx context.Context) error {
	s.mustHave(CapabilityPistonClient)
	return s.grpcClients.PistonClient.SendRateLimitExceed(ctx)
}

func (s *Storage) PistonGetIDKey(_ context.Context, i domain.ClientOrderID) (domain.PistonIDKey, bool) {
	s.mustHave(CapabilityPistonIDCache)
	s.logger.Debugf(s.cache.PistonIDCache.PrettyPrint())
	v, ok := s.cache.PistonIDCache.GetInverse(i)
	if ok {
		return v, true
	}
	return domain.PistonIDKey{}, false //nolint:exhaustruct // OK
}

func (s *Storage) PistonAddID(_ context.Context, pistonID domain.PistonIDKey, exchangeID domain.ClientOrderID) {
	s.mustHave(CapabilityPistonIDCache)
	s.cache.PistonIDCache.Add(pistonID, exchangeID)
}

func (s *Storage) PistonDeleteID(_ context.Context, pistonID domain.PistonIDKey, exchangeID domain.ClientOrderID) {
	s.mustHave(CapabilityPistonIDCache)
	s.cache.PistonIDCache.DeleteInverse(exchangeID)
	s.cache.PistonIDCache.DeleteDirect(pistonID)
}

func (s *Storage) PistonReplaceID(
	_ context.Context, pistonID domain.PistonIDKey, oldID domain.ClientOrderID, newID domain.ClientOrderID,
) {
	s.mustHave(CapabilityPistonIDCache)
	s.cache.PistonIDCache.DeleteInverse(oldID)
	s.cache.PistonIDCache.Add(pistonID, newID)
}

func (b *Builder) WithPistonIDCache() *Builder {
	b.buildCap(CapabilityPistonIDCache)
	return b
}

func (s *Storage) GetExchangeClientOrderID(_ context.Context, i domain.PistonIDKey) (domain.ClientOrderID, bool) {
	s.mustHave(CapabilityPistonIDCache)
	v, ok := s.cache.PistonIDCache.GetDirect(i)
	if ok {
		return v, true
	}
	return "", false
}

func (s *Storage) ContainsPistonIDKey(_ context.Context, i domain.PistonIDKey) bool {
	s.mustHave(CapabilityPistonIDCache)
	return s.cache.PistonIDCache.DirectContains(i)
}

func (b *Builder) WithOrderMoveCache() *Builder {
	b.buildCap(CapabilityOrderMoveCache)
	return b
}

func (s *Storage) AddCommand(_ context.Context, cmd domain.ActiveCommand) {
	s.mustHave(CapabilityOrderMoveCache)
	s.cache.ActiveCommandsCache.Add(cmd)
}

func (s *Storage) GetCommand(_ context.Context, id domain.ClientOrderID) (domain.ActiveCommand, bool) {
	s.mustHave(CapabilityOrderMoveCache)
	s.logger.Debugf(s.cache.ActiveCommandsCache.PrettyPrint())
	v, ok := s.cache.ActiveCommandsCache.Get(id)
	if ok {
		return v, true
	}
	return nil, false
}

func (s *Storage) DeleteCommand(_ context.Context, id domain.ClientOrderID) {
	s.mustHave(CapabilityOrderMoveCache)
	s.cache.ActiveCommandsCache.Delete(id)
}
