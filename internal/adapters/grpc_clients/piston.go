package grpcclients

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"studentgit.kata.academy/quant/torque/internal/adapters/grpc_clients/types/piston"
	"studentgit.kata.academy/quant/torque/internal/domain"
	"studentgit.kata.academy/quant/torque/internal/domain/f"
	"studentgit.kata.academy/quant/torque/pkg/logster"
)

const timeout = time.Second * 5
const reconnectInterval = time.Second * 5

type PistonClient struct {
	logger  logster.Logger
	address string

	marketDataSubscribeRequests chan domain.PistonSubscribeMarketDataMessage
	balancesRequests            chan domain.PistonBalancesRequestMessage
	openOrdersRequests          chan domain.PistonOpenOrderRequestMessage
	orderHistoryRequests        chan domain.PistonOrderHistoryRequestMessage
	instrumentDetails           chan domain.PistonInstrumentsDetailsRequestMessage
	createOrder                 chan domain.AddOrder
	cancelOrder                 chan domain.CancelOrder
	moveOrder                   chan domain.MoveOrder

	conn          *grpc.ClientConn
	mdStream      piston.Piston_ConnectMDGatewayClient
	tradingStream piston.Piston_ConnectTradingGatewayClient
	errChan       chan error
	closed        bool

	mdStreamIncomingMessage      chan *piston.MarketdataMessage
	tradingStreamIncomingMessage chan *piston.TradingMessage
}

func NewPistonClient(logger logster.Logger, address string) *PistonClient {
	return &PistonClient{
		logger:  logger.WithField("module", "piston_client"),
		address: address,

		marketDataSubscribeRequests: make(chan domain.PistonSubscribeMarketDataMessage),
		createOrder:                 make(chan domain.AddOrder),
		cancelOrder:                 make(chan domain.CancelOrder),
		balancesRequests:            make(chan domain.PistonBalancesRequestMessage),
		instrumentDetails:           make(chan domain.PistonInstrumentsDetailsRequestMessage),
		openOrdersRequests:          make(chan domain.PistonOpenOrderRequestMessage),
		orderHistoryRequests:        make(chan domain.PistonOrderHistoryRequestMessage),
		moveOrder:                   make(chan domain.MoveOrder),

		mdStreamIncomingMessage:      make(chan *piston.MarketdataMessage),
		tradingStreamIncomingMessage: make(chan *piston.TradingMessage),
		errChan:                      make(chan error),
		closed:                       false,

		conn:          nil,
		mdStream:      nil,
		tradingStream: nil,
	}
}

func (c *PistonClient) close() {
	c.logger.Infof("closing Piston client...")
	c.closed = true
	close(c.marketDataSubscribeRequests)
	close(c.createOrder)
	close(c.cancelOrder)
	close(c.moveOrder)
	close(c.mdStreamIncomingMessage)
	close(c.tradingStreamIncomingMessage)
	close(c.balancesRequests)
	close(c.openOrdersRequests)
	close(c.orderHistoryRequests)
	close(c.instrumentDetails)
	c.logger.Infof("Piston client closed")
}

func (c *PistonClient) Run(ctx context.Context) error {
	c.logger.Infof("Piston client is running")
	go c.readTradingMessages(ctx)
	go c.readMarketDataMessages(ctx)
	go c.dispatchRequests(ctx)
	go c.reconnectAndRestart(ctx)
	return nil
}

func (c *PistonClient) Init(ctx context.Context) error {
	dialCtx, cancelDial := context.WithTimeout(ctx, timeout)
	defer cancelDial()
	var err error
	c.conn, err = grpc.DialContext(
		dialCtx, c.address, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock(),
	)
	if err != nil {
		c.logger.WithError(err).Errorf("failed to dial Piston")
		return err
	}
	client := piston.NewPistonClient(c.conn)
	c.mdStream, err = client.ConnectMDGateway(ctx, grpc.WaitForReady(true))
	if err != nil {
		c.logger.Errorf("error here")
		return err
	}
	c.tradingStream, err = client.ConnectTradingGateway(ctx, grpc.WaitForReady(true))
	if err != nil {
		c.logger.Errorf("error here")
		return err
	}
	c.logger.Infof("Piston client is initialized")
	return nil
}

//nolint:exhaustive // We dont want to do switch for all message types, we only choose what to process
func (c *PistonClient) dispatchRequests(ctx context.Context) {
	for {
		select {
		case <-ctx.Done(): // external context
			c.logger.Infof("stopping dispatching requests")
			c.close()
			return
		case mdRecvMsg := <-c.mdStreamIncomingMessage:
			c.marketDataSubscribeRequests <- piston.MarketDataMessageToDomain(mdRecvMsg)
		case tradingRecvMsg := <-c.tradingStreamIncomingMessage:
			c.logger.Debugf("tradingRecvMsg: %v", tradingRecvMsg)
			switch tradingRecvMsg.Type {
			case piston.TradingMessageType_RequestOpenOrders:
				c.openOrdersRequests <- domain.PistonOpenOrderRequestMessage{
					RequestID: tradingRecvMsg.GetOpenOrdersRequest().GetRequestID(),
				}
			case piston.TradingMessageType_RequestBalance:
				c.balancesRequests <- domain.PistonBalancesRequestMessage{}
			case piston.TradingMessageType_AddOrderRequest:
				c.createOrder <- piston.CreateOrderTradingMessageToDomain(tradingRecvMsg.GetAddOrder())
			case piston.TradingMessageType_CancelOrderRequest:
				c.cancelOrder <- piston.CancelOrderTradingMessageToDomain(tradingRecvMsg.GetCancelOrder())
			case piston.TradingMessageType_MoveOrderRequest:
				c.moveOrder <- piston.MoveOrderTradingMessageToDomain(tradingRecvMsg.GetMoveOrder())
			case piston.TradingMessageType_RequestInstrumentsDetails:
				c.instrumentDetails <- domain.PistonInstrumentsDetailsRequestMessage{}
			default:
				c.logger.Infof("unknown trading message type: %v", tradingRecvMsg.GetType())
			}
		}
	}
}

//nolint:gocognit,nestif // complex, yes, should be simplified in futures via abstracting streams via pkg
func (c *PistonClient) readTradingMessages(ctx context.Context) {
	resultChan := make(chan *piston.TradingMessage)

	go func() {
		for {
			if c.tradingStream != nil {
				in, err := c.tradingStream.Recv()
				if err != nil {
					c.logger.WithError(err).Infof("trading stream is closed")
					close(resultChan)
					errX := c.tradingStream.CloseSend()
					if errX != nil {
						c.logger.WithError(errX).Errorf("failed to send CloseSend to trading stream")
					}
					return
				}
				if !c.closed {
					resultChan <- in
				} else {
					c.logger.Warnf("sending to close channel")
				}
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			c.logger.Infof("context is done, exiting readTradingMessages")
			return
		case in, ok := <-resultChan:
			if !ok {
				c.logger.Debugf("trading stream result chan is closed")
				return
			}
			c.logger.WithField(f.PistonTradingMessage, in).Infof("got trading message")
			c.tradingStreamIncomingMessage <- in
		}
	}
}

//nolint:gocognit,nestif // complex, yes, should be simplified in future via abstracting streams via pkg
func (c *PistonClient) readMarketDataMessages(ctx context.Context) {
	resultChan := make(chan *piston.MarketdataMessage)

	// This is done in a separate goroutine because Recv() is blocking.

	// Lines from docs.
	// RecvMsg blocks until it receives a message into m or the stream is
	// done. It returns io.EOF when the stream completes successfully. On
	// any other error, the stream is aborted and the error contains the RPC
	// status. https://github.com/grpc/grpc-go/blob/v1.57.0/stream.go
	go func() {
		for {
			if c.mdStream != nil {
				in, err := c.mdStream.Recv()
				if err != nil {
					c.logger.WithError(err).Infof("md stream is closed")
					c.errChan <- err
					close(resultChan)
					errX := c.mdStream.CloseSend()
					if errX != nil {
						c.logger.WithError(errX).Errorf("failed to send CloseSend to md stream")
					}
					return
				}
				if !c.closed {
					resultChan <- in
				} else {
					c.logger.Warnf("sending to close channel")
				}
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			c.logger.Infof("context is done, exiting readMarketDataMessages")
			return
			// TODO CHECK if resultChan is closed and exit if it is to avoid goroutine leak
		case in, ok := <-resultChan:
			if !ok {
				c.logger.Debugf("market data result chan is closed")
				return
			}
			c.logger.Debugf("received market data response: %v", in)
			c.mdStreamIncomingMessage <- in
		}
	}
}

//nolint:gocognit // complex, yes, should be simplified in futures via abstracting streams via pkg
func (c *PistonClient) reconnectAndRestart(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			c.logger.Infof("Context done, stopping reconnection attempts")
			return
		case err, ok := <-c.errChan:
			if !ok {
				c.logger.Debugf("err chan is closed")
				return
			}
			c.logger.WithError(err).Infof("reconnectAndRestart process error")

			err = c.conn.Close()
			if err != nil {
				c.logger.WithError(err).Errorf("error closing grpc connection")
			}

			dialCtx, cancelDial := context.WithTimeout(ctx, timeout)
			conn, err := grpc.DialContext(dialCtx, c.address,
				grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
			cancelDial() // Cancel the dialing context once the DialContext is done

			if err != nil {
				c.logger.WithError(err).Errorf("Failed to reconnect, retrying in %v seconds...", reconnectInterval)
				time.Sleep(reconnectInterval)
				go func() {
					if !c.closed {
						c.errChan <- err
					} else {
						c.logger.Warnf("sending to closed channel when reconnecting")
					}
				}()
				continue
			}
			c.conn = conn

			client := piston.NewPistonClient(c.conn)
			mdStream, mdError := client.ConnectMDGateway(ctx, grpc.WaitForReady(true))
			tradingStream, tradingError := client.ConnectTradingGateway(ctx, grpc.WaitForReady(true))

			if mdError != nil || tradingError != nil {
				c.logger.Errorf("Failed to restart streams, retrying in %v seconds...", reconnectInterval)
				time.Sleep(reconnectInterval)
				go func() {
					if mdError != nil {
						c.logger.WithError(err).Errorf("Market data stream error")
						c.errChan <- mdError
					} else {
						c.logger.WithError(err).Errorf("Trading stream error")
						c.errChan <- tradingError
					}
				}()
				continue
			}
			c.mdStream = mdStream
			c.tradingStream = tradingStream

			c.logger.Infof("Reconnected successfully!")
			go c.readTradingMessages(ctx)
			go c.readMarketDataMessages(ctx)
		}
	}
}

func (c *PistonClient) sendMarketData(_ context.Context, mdMsg *piston.MarketdataMessage) error {
	if c.conn.GetState() == connectivity.Ready {
		err := c.mdStream.Send(mdMsg)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *PistonClient) SendOrderBook(ctx context.Context, ob domain.OrderBook) error {
	mdMsg := piston.OrderBookToMarketDataMessage(ob)
	return c.sendMarketData(ctx, mdMsg)
}

func (c *PistonClient) SendBalances(_ context.Context, b domain.PistonBalancesMessage) error {
	tradingMsg := piston.BalanceToTradingMessage(b)
	return c.sendTradingData(tradingMsg)
}

func (c *PistonClient) SendOpenOrders(_ context.Context, m domain.PistonOpenOrdersMessage) error {
	tradingMsg := piston.OpenOrdersToTradingMessage(m)
	return c.sendTradingData(tradingMsg)
}

func (c *PistonClient) SendOrderAdded(_ context.Context, order domain.OrderAdded) error {
	tradingMsg := piston.OrderAddedToTradingMessage(order)
	return c.sendTradingData(tradingMsg)
}

func (c *PistonClient) SendOrderAddRejected(_ context.Context, order domain.OrderAddReject) error {
	tradingMsg := piston.OrderAddRejectToTradingMessage(order)
	return c.sendTradingData(tradingMsg)
}

func (c *PistonClient) SendOrderMoved(_ context.Context, order domain.OrderMoved) error {
	tradingMsg := piston.OrderMovedToTradingMessage(order)
	return c.sendTradingData(tradingMsg)
}

func (c *PistonClient) SendOrderMoveReject(_ context.Context, order domain.OrderMoveReject) error {
	tradingMsg := piston.OrderMoveRejectToTradingMessage(order)
	return c.sendTradingData(tradingMsg)
}

func (c *PistonClient) SendOrderCancelled(_ context.Context, order domain.OrderCancelled) error {
	tradingMsg := piston.OrderCancelledToTradingMessage(order)
	return c.sendTradingData(tradingMsg)
}

func (c *PistonClient) SendOrderCancelReject(_ context.Context, order domain.OrderCancelReject) error {
	tradingMsg := piston.OrderCancelRejectToTradingMessage(order)
	return c.sendTradingData(tradingMsg)
}

func (c *PistonClient) SendOrderFilled(_ context.Context, order domain.OrderFilled) error {
	tradingMsg := piston.OrderFilledToTradingMessage(order)
	return c.sendTradingData(tradingMsg)
}

func (c *PistonClient) SendOrderExecuted(_ context.Context, order domain.OrderExecuted) error {
	tradingMsg := piston.OrderExecutedToTradingMessage(order)
	return c.sendTradingData(tradingMsg)
}

func (c *PistonClient) SendInstrumentDetails(_ context.Context, details []domain.InstrumentDetails) error {
	tradingMsg := piston.InstrumentDetailsToTradingMessage(details)
	return c.sendTradingData(tradingMsg)
}

func (c *PistonClient) SendRateLimitExceed(_ context.Context) error {
	return c.sendTradingData(&piston.TradingMessage{
		Type: piston.TradingMessageType_RateLimitExceedError,
		Data: &piston.TradingMessage_RateLimitExceed{}, //nolint:exhaustruct // OK
	})
}

func (c *PistonClient) sendTradingData(mdMsg *piston.TradingMessage) error {
	// TODO data race, conn is being closed in reconnectAndRestart and state is checked here
	if c.conn.GetState() == connectivity.Ready {
		c.logger.WithField(f.PistonTradingMessageType, mdMsg.Type).
			WithField(f.PistonTradingMessage, mdMsg.Data).Infof("sending Piston trading message")
		err := c.tradingStream.Send(mdMsg)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *PistonClient) GetCreateOrderMessagesRequests() <-chan domain.AddOrder {
	return c.createOrder
}

func (c *PistonClient) GetMoveOrderMessagesRequests() <-chan domain.MoveOrder {
	return c.moveOrder
}

func (c *PistonClient) GetCancelOrderMessagesRequests() <-chan domain.CancelOrder {
	return c.cancelOrder
}

func (c *PistonClient) GetSubscribeMessageRequests() <-chan domain.PistonSubscribeMarketDataMessage {
	return c.marketDataSubscribeRequests
}

func (c *PistonClient) GetOpenOrdersMessageRequests() <-chan domain.PistonOpenOrderRequestMessage {
	return c.openOrdersRequests
}

func (c *PistonClient) GetBalancesMessageRequests() <-chan domain.PistonBalancesRequestMessage {
	return c.balancesRequests
}

func (c *PistonClient) GetOrderHistoryMessageRequests() <-chan domain.PistonOrderHistoryRequestMessage {
	return c.orderHistoryRequests
}

func (c *PistonClient) GetInstrumentDetailsMessageRequests() <-chan domain.PistonInstrumentsDetailsRequestMessage {
	return c.instrumentDetails
}
