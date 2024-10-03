package main

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"studentgit.kata.academy/quant/torque/internal/adapters/grpc_clients/types/piston"
	"studentgit.kata.academy/quant/torque/internal/domain"
	"studentgit.kata.academy/quant/torque/pkg/grpc/streamer"
	"studentgit.kata.academy/quant/torque/pkg/logster"
)

const timeout = time.Second * 5
const reconnectInterval = time.Second * 5

type PistonClientMarketTest struct {
	logger                      logster.Logger
	address                     string
	MarketDataSubscribeRequests chan domain.PistonSubscribeMarketDataMessage
	mdStreamHandler             streamer.GenericHandler[*piston.MarketdataMessage, *piston.MarketdataMessage]
	closed                      bool
	mdStreamIncomingMessage     chan *piston.MarketdataMessage
}

func NewPistonClientMarketTest(logger logster.Logger, address string) PistonClientMarketTest {
	return PistonClientMarketTest{
		logger:  logger.WithField("module", "piston_client"),
		address: address,
		// TODO: make unexported when ready.
		MarketDataSubscribeRequests: make(chan domain.PistonSubscribeMarketDataMessage),
		mdStreamIncomingMessage:     make(chan *piston.MarketdataMessage),
		closed:                      false,
		mdStreamHandler:             nil,
	}
}

func (c *PistonClientMarketTest) close() {
	c.logger.Infof("closing Piston client...")
	c.closed = true
	// TODO: handle error?
	_ = c.mdStreamHandler.CloseSend()
	close(c.MarketDataSubscribeRequests)
	close(c.mdStreamIncomingMessage)
	c.logger.Infof("Piston client closed")
}

func (c *PistonClientMarketTest) Init(ctx context.Context) error {
	cfg := streamer.NewConfig(c.address, 5, timeout, reconnectInterval)
	handler := streamer.NewGenericStreamHandler[*piston.MarketdataMessage, *piston.MarketdataMessage](
		ctx,
		new(streamer.MDStream),
		cfg,
		c.logger,
	)
	handler.WithGrpcDialOpts(grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	handler.WithGrpcCallOpts(grpc.WaitForReady(true))

	c.mdStreamHandler = handler

	c.logger.Infof("Piston client is initialized")
	return nil
}

func (c *PistonClientMarketTest) Run(ctx context.Context) error {
	c.logger.Infof("Piston client is running")
	var errGroup errgroup.Group
	errGroup.Go(func() error {
		if err := c.mdStreamHandler.RunWithReconnect(ctx); err != nil {
			c.logger.WithError(err).Errorf("failed to runWithReconnect")
			return err
		}
		return nil
	})
	go c.dispatchRequests(ctx)
	return errGroup.Wait()
}

func (c *PistonClientMarketTest) dispatchRequests(ctx context.Context) {
	fromStream := c.mdStreamHandler.GetOutputChan()
	for {
		select {
		case <-ctx.Done(): // external context
			c.logger.Infof("stopping dispatching requests")
			c.close()
			return
		case mdRecvMsg := <-fromStream:
			c.MarketDataSubscribeRequests <- piston.MarketDataMessageToDomain(mdRecvMsg)
		}
	}
}

func (c *PistonClientMarketTest) sendMarketData(_ context.Context, mdMsg *piston.MarketdataMessage) error {
	if !c.mdStreamHandler.IsClosed() {
		inStream := c.mdStreamHandler.GetInputChan()
		inStream <- mdMsg
		return nil
	}
	return fmt.Errorf("stream is closed")
}

func (c *PistonClientMarketTest) SendOrderBook(ctx context.Context, ob domain.OrderBook) error {
	mdMsg := piston.OrderBookToMarketDataMessage(ob)
	return c.sendMarketData(ctx, mdMsg)
}

func (c *PistonClientMarketTest) GetSubscribeMessageRequests() <-chan domain.PistonSubscribeMarketDataMessage {
	return c.MarketDataSubscribeRequests
}
