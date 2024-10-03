package app

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"golang.org/x/sync/errgroup"
	"studentgit.kata.academy/quant/torque/internal/domain"
	"studentgit.kata.academy/quant/torque/internal/domain/observation"

	"studentgit.kata.academy/quant/torque/internal/domain/f"
	"studentgit.kata.academy/quant/torque/pkg/logster"
)

//nolint:lll // no linebreaks here is ok
type okxStorage interface {
	Now() time.Time
	GenerateULID() domain.ClientOrderID

	GetOKXBookUpdateMessages() <-chan domain.OKXOrderBook
	GetOKXOrderUpdates() <-chan domain.OKXUserStreamOrderUpdate

	GetOKXOrderBookSnapshot(context.Context, domain.OKXSpotInstrument, int) (domain.OKXOrderBook, error)
	GetOKXBalance(context.Context) ([]domain.Balance, error)

	OKXCreateOrder(context.Context, domain.OrderRequest) (domain.OKXOrderResponseAck, error)
	OKXCancelOrder(
		context.Context, domain.OKXSpotInstrument, domain.OrderID, domain.ClientOrderID,
	) (domain.OKXOrderResponseAck, error)
	OKXMoveOrder(
		context.Context, domain.OKXSpotInstrument, domain.OrderID, domain.ClientOrderID, decimal.Decimal, decimal.Decimal,
	) error
	OKXGetSpotOpenOrders(context.Context) ([]domain.OpenOrder, error)
	OKXGetSpotFillsHistory(context.Context, time.Time, time.Time) ([]domain.Order, error)
	OKXSubscribeInstruments(context.Context, []domain.OKXInstrument) error
	OKXUnsubscribeInstruments(context.Context, []domain.OKXInstrument) error

	GetLocalInstrument(context.Context, domain.ExchangeInstrument) (domain.LocalInstrument, bool)
	GetAllInstrumentsDetailsByExchange(context.Context, domain.Exchange) []domain.InstrumentDetails
	GetExchangeInstrument(context.Context, domain.LocalInstrument) (domain.ExchangeInstrument, bool)
	GetInstrumentsDetails(context.Context, domain.LocalInstrument) (domain.InstrumentDetails, bool)
	GetOrderBook(context.Context, domain.LocalInstrument) (domain.OrderBook, bool)
	SaveOrderBook(context.Context, domain.OrderBook)
	DeleteOrderBook(context.Context, domain.LocalInstrument)

	GetPistonCreateOrderMessages() <-chan domain.AddOrder
	GetPistonCancelOrderMessages() <-chan domain.CancelOrder
	GetPistonMoveOrderMessages() <-chan domain.MoveOrder
	GetPistonSubscribeMessageRequests() <-chan domain.PistonSubscribeMarketDataMessage
	GetPistonBalancesMessageRequests() <-chan domain.PistonBalancesRequestMessage
	GetPistonOpenOrdersMessageRequests() <-chan domain.PistonOpenOrderRequestMessage
	GetPistonOrderHistoryMessageRequests() <-chan domain.PistonOrderHistoryRequestMessage
	GetPistonInstrumentDetailsMessageRequests() <-chan domain.PistonInstrumentsDetailsRequestMessage

	PistonSendOrderBook(context.Context, domain.OrderBook) error
	PistonSendBalances(context.Context, domain.PistonBalancesMessage) error
	PistonSendOpenOrders(context.Context, domain.PistonOpenOrdersMessage) error
	PistonSendInstrumentDetails(context.Context, []domain.InstrumentDetails) error
	PistonRateLimitExceedError(ctx context.Context) error

	PistonSendOrderAdded(context.Context, domain.OrderAdded) error
	PistonSendOrderAddRejected(context.Context, domain.OrderAddReject) error
	PistonSendOrderMoved(context.Context, domain.OrderMoved) error
	PistonSendOrderMoveReject(context.Context, domain.OrderMoveReject) error
	PistonSendOrderCancelled(context.Context, domain.OrderCancelled) error
	PistonSendOrderCancelReject(context.Context, domain.OrderCancelReject) error
	PistonSendOrderFilled(context.Context, domain.OrderFilled) error
	PistonSendOrderExecuted(context.Context, domain.OrderExecuted) error

	PistonGetIDKey(context.Context, domain.ClientOrderID) (domain.PistonIDKey, bool)
	PistonAddID(context.Context, domain.PistonIDKey, domain.ClientOrderID)
	GetExchangeClientOrderID(context.Context, domain.PistonIDKey) (domain.ClientOrderID, bool)
	PistonDeleteID(context.Context, domain.PistonIDKey, domain.ClientOrderID)
	ContainsPistonIDKey(context.Context, domain.PistonIDKey) bool

	GetOKXSpotCreateOrderAckMessages() <-chan domain.OKXOrderResponseAck
	GetOKXSpotCreateOrderFailedMessages() <-chan domain.OKXOrderResponseAck
	GetOKXSpotModifyOrderAckMessages() <-chan domain.OKXOrderResponseAck
	GetOKXSpotModifyOrderFailedMessages() <-chan domain.OKXOrderResponseAck
	GetOKXSpotCancelOrderAckMessages() <-chan domain.OKXOrderResponseAck
	GetOKXSpotCancelOrderFailedMessages() <-chan domain.OKXOrderResponseAck
	GetOKXSpotTradingWebsocketStatus() bool
	CancelOKXSpotTradingWebsocketOrder(context.Context, domain.OKXSpotInstrument, domain.OrderID, domain.ClientOrderID) error
	ModifyOKXSpotTradingWebsocketOrder(context.Context, domain.OKXSpotInstrument, domain.OrderID, domain.ClientOrderID, decimal.Decimal, decimal.Decimal) error
	CreateOKXSpotTradingWebsocketOrder(context.Context, domain.OrderRequest) error

	SetLabelValue(observation.Name, string, float64)
}

type OKX struct {
	logger                     logster.Logger
	storage                    okxStorage
	okxAccount                 domain.ExchangeAccount
	depthUpdate                <-chan domain.OKXOrderBook
	orderUpdate                <-chan domain.OKXUserStreamOrderUpdate
	createOrderAck             <-chan domain.OKXOrderResponseAck
	modifyOrderAck             <-chan domain.OKXOrderResponseAck
	cancelOrderAck             <-chan domain.OKXOrderResponseAck
	createOrderFailed          <-chan domain.OKXOrderResponseAck
	modifyOrderFailed          <-chan domain.OKXOrderResponseAck
	cancelOrderFailed          <-chan domain.OKXOrderResponseAck
	pistonSubscribeRequests    <-chan domain.PistonSubscribeMarketDataMessage
	pistonOpenOrdersRequests   <-chan domain.PistonOpenOrderRequestMessage
	pistonBalancesRequests     <-chan domain.PistonBalancesRequestMessage
	pistonOrderHistoryRequests <-chan domain.PistonOrderHistoryRequestMessage
	pistonInstrumentDetails    <-chan domain.PistonInstrumentsDetailsRequestMessage
	pistonCreateOrderRequests  <-chan domain.AddOrder
	pistonMoveOrderRequests    <-chan domain.MoveOrder
	pistonCancelOrderRequests  <-chan domain.CancelOrder
}

func NewOKX(logger logster.Logger, storage okxStorage, okxAccount domain.ExchangeAccount) *OKX {
	return &OKX{
		logger:                     logger.WithField(f.Module, "okx_gateway"),
		storage:                    storage,
		okxAccount:                 okxAccount,
		depthUpdate:                storage.GetOKXBookUpdateMessages(),
		orderUpdate:                storage.GetOKXOrderUpdates(),
		pistonSubscribeRequests:    storage.GetPistonSubscribeMessageRequests(),
		pistonCreateOrderRequests:  storage.GetPistonCreateOrderMessages(),
		pistonCancelOrderRequests:  storage.GetPistonCancelOrderMessages(),
		pistonOpenOrdersRequests:   storage.GetPistonOpenOrdersMessageRequests(),
		pistonBalancesRequests:     storage.GetPistonBalancesMessageRequests(),
		pistonOrderHistoryRequests: storage.GetPistonOrderHistoryMessageRequests(),
		pistonMoveOrderRequests:    storage.GetPistonMoveOrderMessages(),
		pistonInstrumentDetails:    storage.GetPistonInstrumentDetailsMessageRequests(),
		createOrderAck:             storage.GetOKXSpotCreateOrderAckMessages(),
		modifyOrderAck:             storage.GetOKXSpotModifyOrderAckMessages(),
		cancelOrderAck:             storage.GetOKXSpotCancelOrderAckMessages(),
		createOrderFailed:          storage.GetOKXSpotCreateOrderFailedMessages(),
		modifyOrderFailed:          storage.GetOKXSpotModifyOrderFailedMessages(),
		cancelOrderFailed:          storage.GetOKXSpotCancelOrderFailedMessages(),
	}
}

//nolint:lll // linebreak looks ugly
func (b *OKX) Run() error {
	g := errgroup.Group{}
	g.Go(process(b.logger, "Process depth update", b.depthUpdate, b.processDepthUpdates))
	g.Go(process(b.logger, "Process order update", b.orderUpdate, b.processOrderUpdate))
	g.Go(process(b.logger, "Process piston subscribe requests", b.pistonSubscribeRequests, b.processPistonSubscribeRequests))
	g.Go(process(b.logger, "Process piston create order messages", b.pistonCreateOrderRequests, b.processPistonCreateOrderMessages))
	g.Go(process(b.logger, "Process piston move order messages", b.pistonMoveOrderRequests, b.processPistonMoveOrderMessages))
	g.Go(process(b.logger, "Process piston cancel order messages", b.pistonCancelOrderRequests, b.processPistonCancelOrderMessages))
	g.Go(process(b.logger, "Process piston instrument details requests", b.pistonInstrumentDetails, b.processPistonInstrumentDetailsRequests))
	g.Go(process(b.logger, "Process piston open orders requests", b.pistonOpenOrdersRequests, b.processPistonOpenOrdersRequests))
	g.Go(process(b.logger, "Process piston balances requests", b.pistonBalancesRequests, b.processPistonBalancesRequests))
	g.Go(process(b.logger, "Process okx create order failed", b.createOrderFailed, b.processOKXCreateOrderFailedMessages))
	g.Go(process(b.logger, "Process okx cancel order failed", b.cancelOrderFailed, b.processOKXCancelOrderFailedMessages))
	g.Go(process(b.logger, "Process okx modify order failed", b.modifyOrderFailed, b.processOKXModifyOrderFailedMessages))
	g.Go(process(b.logger, "Process okx cancel order ack", b.cancelOrderAck, b.processOKXCancelOrderMessages))
	g.Go(process(b.logger, "Process okx create order ack", b.createOrderAck, b.processOKXCreateOrderAckMessages))
	g.Go(process(b.logger, "Process okx modify order ack", b.modifyOrderAck, b.processOKXModifyOrderAckMessages))
	return g.Wait()
}

func (b *OKX) processOKXModifyOrderFailedMessages(ctx context.Context, order domain.OKXOrderResponseAck) error {
	b.logger.WithField(f.ModifyOrderFailed, order).Infof("ModifyOrderFailed received")
	if order.ClientOrderID == domain.FakeOrder {
		return nil // do not process orders with this ClientOrderID, because we only use them
		// to support websocket user stream and avoid reconnects
	}
	pistonID, ok := b.storage.PistonGetIDKey(ctx, order.ClientOrderID)
	if !ok {
		b.logger.WithField(f.ExchangeClientOrderID, order.ClientOrderID).Warnf("PistonID not found")
	}
	err := b.storage.PistonSendOrderMoveReject(ctx, domain.OrderMoveReject{
		OrderRequestID:        pistonID.OrderRequestID,
		PistonID:              pistonID.PistonID,
		ExchangeID:            order.OrderID,
		ExchangeClientOrderID: order.ClientOrderID,
	})
	if err != nil {
		return err
	}
	return nil
}

func (b *OKX) processOKXCancelOrderFailedMessages(ctx context.Context, order domain.OKXOrderResponseAck) error {
	b.logger.WithField(f.CancelOrderFailed, order).Infof("CancelOrderFailed received")
	if order.ClientOrderID == domain.FakeOrder {
		b.logger.Warnf("failed to cancel fake order")
		return nil // do not process orders with this ClientOrderID, because we only use them
		// to support websocket user stream and avoid reconnects
	}
	pistonID, ok := b.storage.PistonGetIDKey(ctx, order.ClientOrderID)
	if !ok {
		b.logger.WithField(f.ExchangeClientOrderID, order.ClientOrderID).Warnf("PistonID not found")
	}
	err := b.storage.PistonSendOrderCancelReject(ctx, domain.OrderCancelReject{
		OrderRequestID:        pistonID.OrderRequestID,
		PistonID:              pistonID.PistonID,
		ExchangeID:            order.OrderID,
		ExchangeClientOrderID: order.ClientOrderID,
	})
	if err != nil {
		return err
	}
	return nil
}

func (b *OKX) processOKXCreateOrderFailedMessages(ctx context.Context, order domain.OKXOrderResponseAck) error {
	b.logger.WithField(f.CreateOrderFailed, order).Infof("CreateOrderFailed received")
	if order.ClientOrderID == domain.FakeOrder {
		b.logger.Warnf("failed to place fake order")
		return nil // do not process orders with this ClientOrderID, because we only use them
		// to support websocket user stream and avoid reconnects
	}
	pistonIDKey, ok := b.storage.PistonGetIDKey(ctx, order.ClientOrderID)
	if !ok {
		b.logger.WithField(f.ExchangeClientOrderID, order.ClientOrderID).Warnf("PistonID not found")
	}
	err := b.storage.PistonSendOrderAddRejected(ctx, domain.OrderAddReject{
		OrderRequestID:        pistonIDKey.OrderRequestID,
		PistonID:              pistonIDKey.PistonID,
		ExchangeClientOrderID: order.ClientOrderID,
	})
	if err != nil {
		return err
	}
	return nil
}

func (b *OKX) processOKXCancelOrderMessages(_ context.Context, order domain.OKXOrderResponseAck) error {
	b.logger.WithField(f.OrderCanceled, order).Infof("CancelOrder received")
	// We process cancel orders wia UserStreamUpdates
	return nil
}

func (b *OKX) processOKXCreateOrderAckMessages(_ context.Context, order domain.OKXOrderResponseAck) error {
	b.logger.WithField(f.OrderCreatedAck, order).Infof("CreateOrderAck received")
	// We process add orders wia UserStreamUpdates
	return nil
}

func (b *OKX) processOKXModifyOrderAckMessages(_ context.Context, order domain.OKXOrderResponseAck) error {
	b.logger.WithField(f.OrderMoveAck, order).Infof("ModifyOrderAck received")
	// We process add orders wia UserStreamUpdates
	return nil
}

func (b *OKX) processPistonInstrumentDetailsRequests(
	ctx context.Context, _ domain.PistonInstrumentsDetailsRequestMessage,
) error {
	b.logger.Infof("Instrument details request from Piston")
	details := b.storage.GetAllInstrumentsDetailsByExchange(ctx, domain.OKXExchange)
	err := b.storage.PistonSendInstrumentDetails(ctx, details)
	if err != nil {
		return err
	}
	return nil
}

func (b *OKX) processPistonOpenOrdersRequests(
	ctx context.Context, request domain.PistonOpenOrderRequestMessage,
) error {
	b.logger.Infof("Open orders request from Piston")
	openOrders, err := b.storage.OKXGetSpotOpenOrders(ctx)
	// Limited to 100 orders, account for this in future !!!
	if err != nil {
		return err
	}
	pistonOrders := make([]domain.PistonOpenOrder, 0, len(openOrders))
	for _, order := range openOrders {
		if order.ClientOrderID == domain.FakeOrder {
			continue
		}
		localInstrument, ok := b.storage.GetLocalInstrument(ctx, domain.ExchangeInstrument{
			Exchange: domain.OKXExchange,
			Symbol:   order.Instrument,
		})
		if !ok {
			b.logger.WithField(f.Instrument, order.Instrument).Warnf("Instrument not found")
			localInstrument = domain.LocalInstrument{
				Symbol:   order.Instrument,
				Exchange: domain.OKXExchange,
			}
		}
		pistonID, ok := b.storage.PistonGetIDKey(ctx, order.ClientOrderID)
		if !ok {
			b.logger.WithField(f.ExchangeClientOrderID, order.ClientOrderID).Warnf("PistonID not found")
		}
		pistonOrder := domain.PistonOpenOrder{
			PistonID:              pistonID.PistonID,
			Exchange:              localInstrument.Exchange,
			ExchangeID:            order.OrderID,
			ExchangeClientOrderID: order.ClientOrderID,
			Instrument:            localInstrument.Symbol,
			Side:                  order.Side,
			OrderType:             order.Type,
			Size:                  order.OrigQty,
			Price:                 order.Price,
			RemainingSize:         order.OrigQty.Sub(order.ExecutedQty),
			Created:               order.Created,
			LastUpdated:           order.UpdatedAt,
		}
		pistonOrders = append(pistonOrders, pistonOrder)
	}
	err = b.storage.PistonSendOpenOrders(ctx, domain.PistonOpenOrdersMessage{
		RequestID: request.RequestID,
		Orders:    pistonOrders,
	})
	if err != nil {
		return err
	}
	return nil
}

func (b *OKX) processPistonBalancesRequests(ctx context.Context, _ domain.PistonBalancesRequestMessage) error {
	b.logger.Infof("Balances request from Piston")
	balances, err := b.storage.GetOKXBalance(ctx)
	if err != nil {
		return err
	}
	err = b.storage.PistonSendBalances(ctx, domain.PistonBalancesMessage{
		Exchange:        domain.OKXExchange,
		ExchangeAccount: b.okxAccount,
		Balances:        balances,
		Timestamp:       b.storage.Now(),
	})
	if err != nil {
		return err
	}
	return nil
}

func (b *OKX) processPistonSubscribeRequests(ctx context.Context, list domain.PistonSubscribeMarketDataMessage) error {
	b.logger.WithField(f.SubscribeInstruments, list).Infof("Subscribe to instruments")
	subscribeList := make([]domain.OKXInstrument, 0, len(list.Instruments))
	for _, symbol := range list.Instruments {
		exchangeInstrument, ok := b.storage.GetExchangeInstrument(ctx, domain.LocalInstrument{
			Exchange: domain.OKXExchange,
			Symbol:   symbol,
		})
		if !ok {
			b.logger.WithField(f.Instrument, symbol).Errorf("Instrument not found")
			continue
		}
		subscribeList = append(subscribeList, domain.OKXSpotInstrument(exchangeInstrument.Symbol))
	}
	if len(subscribeList) == 0 {
		b.logger.Errorf("empty instrument list")
		return nil
	}
	return b.storage.OKXSubscribeInstruments(ctx, subscribeList)
}

//nolint:gocognit //we have what we have, this func is complex indeed
func (b *OKX) processPistonMoveOrderMessages(ctx context.Context, moveOrder domain.MoveOrder) error {
	b.logger.WithField(f.OrderMoveRequest, moveOrder).Infof("Piston move order message")
	localInstrument := domain.LocalInstrument{
		Exchange: domain.OKXExchange,
		Symbol:   moveOrder.Instrument,
	}
	exchangeInstrument, ok := b.storage.GetExchangeInstrument(ctx, localInstrument)
	if !ok {
		return fmt.Errorf("failed to get exchange instrument for symbol %v", moveOrder.Instrument)
	}
	pistonIDKey := domain.PistonIDKey{
		OrderRequestID: moveOrder.OrderRequestID,
		PistonID:       moveOrder.PistonID,
	}
	clientOrderID, clientOrderIDExists := b.storage.GetExchangeClientOrderID(ctx, pistonIDKey)
	if !clientOrderIDExists {
		b.logger.WithField(f.PistonID, pistonIDKey).Warnf("failed to find origClientOrderID")
		errX := b.storage.PistonSendOrderMoveReject(ctx, domain.OrderMoveReject{
			OrderRequestID:        pistonIDKey.OrderRequestID,
			PistonID:              pistonIDKey.PistonID,
			ExchangeID:            moveOrder.ExchangeID,
			ExchangeClientOrderID: clientOrderID,
		})
		if errX != nil {
			return errX
		}
		return nil
	}
	instrumentDetails, ok := b.storage.GetInstrumentsDetails(ctx, localInstrument)
	if !ok {
		return fmt.Errorf("failed to get instrumentDetails for symbol %v", moveOrder.Instrument)
	}
	if moveOrder.Size.LessThan(instrumentDetails.MinLot) {
		b.logger.Errorf("orderSize violates MinLot constraint createOrderSize: %s, MinLot: %s",
			moveOrder.Size.String(), instrumentDetails.MinLot.String())
		errX := b.storage.PistonSendOrderMoveReject(ctx, domain.OrderMoveReject{
			OrderRequestID:        pistonIDKey.OrderRequestID,
			PistonID:              pistonIDKey.PistonID,
			ExchangeID:            moveOrder.ExchangeID,
			ExchangeClientOrderID: clientOrderID,
		})
		if errX != nil {
			return errX
		}
		return nil
	}
	//nolint:nestif // Probably need to refactor at some point
	if b.storage.GetOKXSpotTradingWebsocketStatus() {
		err := b.storage.ModifyOKXSpotTradingWebsocketOrder(
			ctx, domain.OKXSpotInstrument(exchangeInstrument.Symbol),
			moveOrder.ExchangeID, clientOrderID, moveOrder.Size, moveOrder.Price,
		)
		if err != nil {
			b.logger.WithError(err).Errorf("error creating websocket message")
			errX := b.storage.PistonSendOrderMoveReject(ctx, domain.OrderMoveReject{
				OrderRequestID:        moveOrder.OrderRequestID,
				PistonID:              moveOrder.PistonID,
				ExchangeClientOrderID: moveOrder.ExchangeClientOrderID,
				ExchangeID:            moveOrder.ExchangeID,
			})
			if errX != nil {
				b.logger.WithError(err).Errorf("failed to send OrderMoveReject")
			}
			if errors.Is(err, domain.ErrOrderRateLimitExceed) {
				return b.storage.PistonRateLimitExceedError(ctx)
			}
			return err
		}
	} else {
		b.logger.Warnf("websocket is unavailable, sending modify via rest")
		err := b.storage.OKXMoveOrder(ctx, domain.OKXSpotInstrument(exchangeInstrument.Symbol),
			moveOrder.ExchangeID,
			clientOrderID,
			moveOrder.Size.Round(int32(instrumentDetails.SizePrecision)),
			moveOrder.Price.Round(int32(instrumentDetails.PricePrecision)),
		)
		if err != nil {
			b.logger.WithError(err).Errorf("error moving order")
			errX := b.storage.PistonSendOrderMoveReject(ctx, domain.OrderMoveReject{
				OrderRequestID:        pistonIDKey.OrderRequestID,
				PistonID:              pistonIDKey.PistonID,
				ExchangeID:            moveOrder.ExchangeID,
				ExchangeClientOrderID: clientOrderID,
			})
			if errX != nil {
				return errX
			}
			return err
		}
	}
	return nil
}

func (b *OKX) processPistonCreateOrderMessages(ctx context.Context, createOrder domain.AddOrder) error {
	b.logger.WithField(f.OrderCreateRequest, createOrder).Infof("Piston create order message")
	localInstrument := domain.LocalInstrument{
		Exchange: domain.OKXExchange,
		Symbol:   createOrder.Instrument,
	}
	exchangeInstrument, ok := b.storage.GetExchangeInstrument(ctx, localInstrument)
	if !ok {
		return fmt.Errorf("failed to get exchange instrument for symbol %v", createOrder.Instrument)
	}
	instrumentDetails, ok := b.storage.GetInstrumentsDetails(ctx, localInstrument)
	if !ok {
		return fmt.Errorf("failed to get instrumentDetails for symbol %v", createOrder.Instrument)
	}
	if createOrder.Size.LessThan(instrumentDetails.MinLot) {
		err := b.storage.PistonSendOrderAddRejected(ctx, domain.OrderAddReject{
			OrderRequestID:        createOrder.OrderRequestID,
			PistonID:              createOrder.PistonID,
			ExchangeClientOrderID: "",
		})
		if err != nil {
			return err
		}
		b.logger.Errorf("orderSize violates MinLot constraint createOrderSize: %s, MinLot: %s",
			createOrder.Size.String(), instrumentDetails.MinLot.String())
		return nil
	}
	exchangeID := b.storage.GenerateULID()
	orderRequest := domain.OrderRequest{
		Symbol:        exchangeInstrument.Symbol,
		Side:          createOrder.Side,
		Type:          createOrder.OrderType,
		TimeInForce:   domain.TimeInForceGTC,
		Quantity:      createOrder.Size.Round(int32(instrumentDetails.SizePrecision)),
		Price:         createOrder.Price.Round(int32(instrumentDetails.PricePrecision)),
		ClientOrderID: exchangeID,
	}
	pistonIDKey := domain.PistonIDKey{
		OrderRequestID: createOrder.OrderRequestID,
		PistonID:       createOrder.PistonID,
	}
	b.storage.PistonAddID(ctx, pistonIDKey, exchangeID)
	//nolint:nestif // Probably need to refactor at some point
	if b.storage.GetOKXSpotTradingWebsocketStatus() {
		err := b.storage.CreateOKXSpotTradingWebsocketOrder(ctx, orderRequest)
		if err != nil {
			defer b.storage.PistonDeleteID(ctx, pistonIDKey, exchangeID)
			b.logger.WithError(err).Errorf("error creating websocket message")
			errX := b.storage.PistonSendOrderAddRejected(ctx, domain.OrderAddReject{
				OrderRequestID:        createOrder.OrderRequestID,
				PistonID:              createOrder.PistonID,
				ExchangeClientOrderID: exchangeID,
			})
			if errX != nil {
				b.logger.WithError(err).Errorf("failed to send OrderAddReject")
			}
			if errors.Is(err, domain.ErrOrderRateLimitExceed) {
				return b.storage.PistonRateLimitExceedError(ctx)
			}
			return err
		}
	} else {
		b.logger.Warnf("websocket is unavailable, sending create order via rest")
		ack, err := b.storage.OKXCreateOrder(ctx, orderRequest)
		if err != nil {
			defer b.storage.PistonDeleteID(ctx, pistonIDKey, exchangeID)
			errX := b.storage.PistonSendOrderAddRejected(ctx, domain.OrderAddReject{
				OrderRequestID:        createOrder.OrderRequestID,
				PistonID:              createOrder.PistonID,
				ExchangeClientOrderID: exchangeID,
			})
			if errX != nil {
				b.logger.WithError(err).Errorf("failed to send OrderAddReject")
			}
			if errors.Is(err, domain.ErrPlaceOrderRateLimitExceed) {
				return b.storage.PistonRateLimitExceedError(ctx)
			}
			return err
		}
		b.logger.WithField(f.OrderCreatedAck, ack).Infof("CreateOrderAck received")
	}

	return nil
}

//nolint:gocognit // Need to breakdown into separate functions
func (b *OKX) processPistonCancelOrderMessages(ctx context.Context, cancelOrder domain.CancelOrder) error {
	b.logger.WithField(f.CancelOrderRequest, cancelOrder).Infof("Piston cancel order message")
	// TODO Convert Created to Time
	inst := domain.LocalInstrument{
		Exchange: domain.OKXExchange,
		Symbol:   cancelOrder.Instrument,
	}
	instrument, ok := b.storage.GetExchangeInstrument(ctx, inst)
	if !ok {
		return fmt.Errorf("failed to get exchange instrument for symbol %v", cancelOrder.Instrument)
	}
	pistonID := domain.PistonIDKey{
		OrderRequestID: cancelOrder.OrderRequestID,
		PistonID:       cancelOrder.PistonID,
	}
	if pistonID.IsEmpty() {
		containsIDKey := b.storage.ContainsPistonIDKey(ctx, pistonID)
		if !containsIDKey {
			err := b.storage.PistonSendOrderCancelReject(ctx, domain.OrderCancelReject{
				OrderRequestID:        pistonID.OrderRequestID,
				PistonID:              pistonID.PistonID,
				ExchangeID:            "",
				ExchangeClientOrderID: "",
			})
			if err != nil {
				return err
			}
			return nil
		}
	}
	clientOrderID, ok := b.storage.GetExchangeClientOrderID(ctx, pistonID)
	if !ok {
		b.logger.WithField(f.PistonID, pistonID).Warnf("failed to find origClientOrderID")
		if cancelOrder.ExchangeID == "" {
			b.logger.Warnf("tried to cancel order without exchangeID and exchangeClientOrderID")
			err := b.storage.PistonSendOrderCancelReject(ctx, domain.OrderCancelReject{
				OrderRequestID:        pistonID.OrderRequestID,
				PistonID:              pistonID.PistonID,
				ExchangeID:            "",
				ExchangeClientOrderID: "",
			})
			if err != nil {
				return err
			}
			return nil
		}
	}
	//nolint:nestif // Probably need to refactor at some point
	if b.storage.GetOKXSpotTradingWebsocketStatus() {
		err := b.storage.CancelOKXSpotTradingWebsocketOrder(
			ctx, domain.OKXSpotInstrument(instrument.Symbol), cancelOrder.ExchangeID, clientOrderID,
		)
		if err != nil {
			b.logger.WithError(err).Errorf("error creating websocket message")
			errX := b.storage.PistonSendOrderCancelReject(ctx, domain.OrderCancelReject{
				OrderRequestID:        cancelOrder.OrderRequestID,
				PistonID:              cancelOrder.PistonID,
				ExchangeID:            cancelOrder.ExchangeID,
				ExchangeClientOrderID: clientOrderID,
			})
			if errX != nil {
				b.logger.WithError(err).Errorf("failed to send OrderCancelReject")
			}
			if errors.Is(err, domain.ErrOrderRateLimitExceed) {
				return b.storage.PistonRateLimitExceedError(ctx)
			}
			return err
		}
	} else {
		b.logger.Warnf("websocket is unavailable, sending cancel order via rest")
		cancel, err := b.storage.OKXCancelOrder(
			ctx, domain.OKXSpotInstrument(instrument.Symbol), cancelOrder.ExchangeID, clientOrderID,
		)
		if err != nil {
			errX := b.storage.PistonSendOrderCancelReject(ctx, domain.OrderCancelReject{
				OrderRequestID:        cancelOrder.OrderRequestID,
				PistonID:              cancelOrder.PistonID,
				ExchangeID:            cancelOrder.ExchangeID,
				ExchangeClientOrderID: clientOrderID,
			})
			if errX != nil {
				b.logger.WithError(err).Errorf("failed to send OrderAddReject")
			}
			if errors.Is(err, domain.ErrCancelOrderRateLimitExceed) {
				return b.storage.PistonRateLimitExceedError(ctx)
			}
			return err
		}
		err = b.storage.PistonSendOrderCancelled(ctx, domain.OrderCancelled{
			OrderRequestID:        cancelOrder.OrderRequestID,
			PistonID:              cancelOrder.PistonID,
			ExchangeID:            cancel.OrderID,
			ExchangeClientOrderID: cancel.ClientOrderID,
		})
		if err != nil {
			b.logger.WithError(err).Errorf("failed to send OrderCancelled")
		}
	}
	return nil
}

//nolint:funlen,gocognit // OK
func (b *OKX) processOrderUpdate(ctx context.Context, orderUpdate domain.OKXUserStreamOrderUpdate) error {
	b.logger.WithField(f.UserStreamOrderResponse, orderUpdate).Infof("Order update received")
	if orderUpdate.ClientOrderID == domain.FakeOrder {
		return nil // do not process orders with this ClientOrderID, because we only use them
		// to support websocket user stream and avoid reconnects
	}
	exchangeInstrument := domain.ExchangeInstrument{
		Exchange: domain.OKXExchange,
		Symbol:   orderUpdate.Instrument.Symbol(),
	}
	localInstrument, ok := b.storage.GetLocalInstrument(ctx, exchangeInstrument)
	if !ok {
		b.logger.WithField(f.Instrument, exchangeInstrument).Errorf("Instrument not found")
	}
	pistonIDKey, ok := b.storage.PistonGetIDKey(ctx, orderUpdate.ClientOrderID)
	if !ok {
		b.logger.WithField(f.ExchangeClientOrderID, orderUpdate.ClientOrderID).Warnf("PistonID not found")
	}

	if !orderUpdate.AmendResult {
		b.logger.Warnf("failed to move order")
		err := b.storage.PistonSendOrderMoveReject(ctx, domain.OrderMoveReject{
			OrderRequestID:        pistonIDKey.OrderRequestID,
			PistonID:              pistonIDKey.PistonID,
			ExchangeClientOrderID: orderUpdate.ClientOrderID,
			ExchangeID:            orderUpdate.OrderID,
		})
		if err != nil {
			return err
		}
		return nil
	}

	switch orderUpdate.Status { //nolint:exhaustive // we process all statuses we need events later
	case domain.OrderStatusCreated:
		if orderUpdate.RequestID == orderUpdate.ClientOrderID {
			err := b.storage.PistonSendOrderMoved(ctx, domain.OrderMoved{
				OrderRequestID:        pistonIDKey.OrderRequestID,
				PistonID:              pistonIDKey.PistonID,
				Exchange:              domain.OKXExchange,
				ExchangeID:            orderUpdate.OrderID,
				ExchangeClientOrderID: orderUpdate.ClientOrderID,
				Size:                  orderUpdate.OrigQty,
				Price:                 orderUpdate.Price,
				Created:               orderUpdate.TransactionTime,
			})
			if err != nil {
				return err
			}
			return nil
		}
		err := b.storage.PistonSendOrderAdded(ctx, domain.OrderAdded{
			OrderRequestID:        pistonIDKey.OrderRequestID,
			PistonID:              pistonIDKey.PistonID,
			Exchange:              domain.OKXExchange,
			ExchangeID:            orderUpdate.OrderID,
			ExchangeClientOrderID: orderUpdate.ClientOrderID,
			Instrument:            localInstrument.Symbol,
			Side:                  orderUpdate.Side,
			OrderType:             orderUpdate.Type,
			Size:                  orderUpdate.OrigQty,
			Price:                 orderUpdate.Price,
			LastUpdated:           orderUpdate.TransactionTime,
		})
		if err != nil {
			return err
		}
	case domain.OrderStatusPartiallyFilled, domain.OrderStatusFilled:
		err := b.storage.PistonSendOrderFilled(ctx, domain.OrderFilled{
			OrderRequestID:        pistonIDKey.OrderRequestID,
			PistonID:              pistonIDKey.PistonID,
			Exchange:              domain.OKXExchange,
			ExchangeID:            orderUpdate.OrderID,
			ExchangeClientOrderID: orderUpdate.ClientOrderID,
			TradeID:               orderUpdate.TradeID,
			Instrument:            localInstrument.Symbol,
			Side:                  orderUpdate.Side,
			OrderType:             orderUpdate.Type,
			Size:                  orderUpdate.FillSize,
			Price:                 orderUpdate.FillPrice,
			Fee:                   orderUpdate.FeeAmount,
			FeeCurrency:           orderUpdate.FeeCurrency,
			Timestamp:             orderUpdate.FillTime,
		})
		if err != nil {
			return err
		}
		if orderUpdate.ExecutedQty.Sub(orderUpdate.OrigQty).Abs().LessThan(decimal.NewFromFloat(domain.Precision)) {
			defer b.storage.PistonDeleteID(ctx, pistonIDKey, orderUpdate.ClientOrderID)
			errX := b.storage.PistonSendOrderExecuted(ctx, domain.OrderExecuted{
				OrderRequestID:        pistonIDKey.OrderRequestID,
				PistonID:              pistonIDKey.PistonID,
				ExchangeID:            orderUpdate.OrderID,
				ExchangeClientOrderID: orderUpdate.ClientOrderID,
			})
			if errX != nil {
				return err
			}
		}
	case domain.OrderStatusCanceled:
		defer b.storage.PistonDeleteID(ctx, pistonIDKey, orderUpdate.ClientOrderID)
		err := b.storage.PistonSendOrderCancelled(ctx, domain.OrderCancelled{
			OrderRequestID:        pistonIDKey.OrderRequestID,
			PistonID:              pistonIDKey.PistonID,
			ExchangeID:            orderUpdate.OrderID,
			ExchangeClientOrderID: orderUpdate.ClientOrderID,
		})
		if err != nil {
			return err
		}
	default:
		b.logger.Warnf("unhandled order update")
	}
	return nil
}

func (b *OKX) processDepthUpdates(ctx context.Context, update domain.OKXOrderBook) error {
	// TODO Add Check Sum

	exchangeInstrument := domain.ExchangeInstrument{
		Exchange: domain.OKXExchange,
		Symbol:   update.Symbol.Symbol(),
	}

	instrument, instrumentExists := b.storage.GetLocalInstrument(ctx, exchangeInstrument)
	if !instrumentExists {
		b.logger.WithField(f.Instrument, exchangeInstrument).Errorf("Instrument not found")
	}

	switch update.Type {
	case domain.Snapshot:
		book := domain.OrderBook{
			Exchange:   instrument.Exchange,
			Instrument: instrument.Symbol,
			Asks:       make(map[string]float64),
			Bids:       make(map[string]float64),
			Checksum:   update.Checksum,
		}
		book.UpdatePriceLevelsString(true, update.Bids)  // true for bids
		book.UpdatePriceLevelsString(false, update.Asks) // false for asks
		b.storage.SaveOrderBook(ctx, book)
		err := b.storage.PistonSendOrderBook(ctx, book)
		if err != nil {
			return err
		}
	case domain.Update:
		book, bookExists := b.storage.GetOrderBook(ctx, instrument)
		if !bookExists {
			// Order book not found, resubscribing
			b.logger.WithField(f.Instrument, instrument).Warnf("Order book not found")
			err := b.storage.OKXSubscribeInstruments(ctx, []domain.OKXInstrument{update.Symbol})
			if err != nil {
				return err
			}
			return nil
		} else {
			book.UpdatePriceLevelsString(true, update.Bids)  // true for bids
			book.UpdatePriceLevelsString(false, update.Asks) // false for asks
			book.Checksum = update.Checksum
			b.storage.SaveOrderBook(ctx, book)
			b.storage.SetLabelValue(observation.OrderBook,
				string(instrument.Symbol), float64(len(book.Asks)+len(book.Bids)))
			err := b.storage.PistonSendOrderBook(ctx, book)
			if err != nil {
				return err
			}
		}
	case domain.Unknown:
		b.logger.Warnf("got unknown update %+v", update)
	}
	return nil
}
