//	Package streamer представляет собой абстракцию над стримами gRPC, реализованную c помощью дженериков.
//	Для корректного пересоздания стримов и реконнектов пришлось реализовать абстрактную фабрику вида:
//
//	type StreamFactory[SendType any, RecvType any] interface {
//	CreateStream(ctx context.Context, conn *grpc.ClientConn) (Streamer[SendType, RecvType], error)
//	}
//
//	Теперь при реализации нового стрима нужно будет написать простую имплементацию:

//	type MDStream struct{}
//
//	func (MDStream) createStream(
//		ctx context.Context,
//		conn *grpc.ClientConn,
//		opts ...grpc.CallOption,
//	) (Streamer[*piston.MarketdataMessage, *piston.MarketdataMessage], error) {
//		stream, err := piston.NewPistonClient(conn).ConnectMDGateway(ctx, opts...)
//		if err != nil {
//			return nil, err
//		}
//		return stream, nil
//	}
//
// Пример использования с piston:
//
//	logger := logster.New(os.Stdout, logster.Config{
//		Project:           "x",
//		Format:            "text",
//		Level:             "debug",
//		Env:               "local",
//		DisableStackTrace: true,
//		System:            "",
//		Inst:              "",
//	})
//
//	logger.Infof("Starting server")
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//
//	piston := piston_test_client.NewPistonClientMarketTest(logger, "localhost:5800")
//	if err := piston.Init(ctx); err != nil {
//		logger.WithError(err).Errorf("failed to init")
//		return
//	}
//
//	// Receive messages.
//	go func() {
//		for msg := range piston.MarketDataSubscribeRequests {
//			logger.Infof("Received message: %+v", msg)
//		}
//	}()
//
//	// Send messages.
//	go func() {
//		for i := 0; i < 5; i++ {
//			<-time.After(5 * time.Second)
//			logger.Infof("Sending orderbook")
//			err := piston.SendOrderBook(ctx, domain.OrderBook{
//				Exchange:   domain.OKXExchange,
//				Instrument: "BTC/USDT",
//				Asks:       map[string]float64{"123": 1, "456": 2, "789": 3, "321": 4, "654": 5},
//				Bids:       map[string]float64{"123": 1, "456": 2, "789": 3, "321": 4, "654": 5},
//			})
//			if err != nil {
//				logger.WithError(err).Errorf("failed to send orderbook")
//			}
//		}
//	}()
//
//	// Run piston.
//	if err := piston.Run(ctx); err != nil {
//		logger.WithError(err).Errorf("failed to run")
//		return
//	}

package streamer

import (
	"context"
	"errors"
	"sync/atomic"
	"time"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"studentgit.kata.academy/quant/torque/pkg/logster"
)

const (
	op = "GenericStreamHandler"
)

type Streamer[SendType any, RecvType any] interface {
	Send(m SendType) error
	Recv() (RecvType, error)
	CloseSend() error
}

func (g *genericStreamHandler[SendType, RecvType]) Send(m SendType) error {
	return g.stream.Send(m)
}

func (g *genericStreamHandler[SendType, RecvType]) Recv() (RecvType, error) {
	return g.stream.Recv()
}

func (g *genericStreamHandler[SendType, RecvType]) CloseSend() error {
	g.streamClosed = true
	return g.stream.CloseSend()
}

type GenericHandler[SendType any, RecvType any] interface {
	RunWithReconnect(ctx context.Context) error
	GetInputChan() chan SendType
	GetOutputChan() chan RecvType
	GetErrChanChannel() chan error
	WithGrpcDialOpts(opts ...grpc.DialOption)
	WithGrpcCallOpts(opts ...grpc.CallOption)
	CloseSend() error
	IsClosed() bool
}

type genericStreamHandler[SendType any, RecvType any] struct {
	ctx      context.Context
	stream   Streamer[SendType, RecvType]
	factory  StreamFactory[SendType, RecvType]
	conn     *grpc.ClientConn
	config   *Config
	channels *channels[SendType, RecvType]
	log      logster.Logger

	streamClosed  bool
	errChanClosed bool
}

func NewGenericStreamHandler[SendType any, RecvType any](
	ctx context.Context,
	factory StreamFactory[SendType, RecvType],
	cfg *Config,
	log logster.Logger,
) GenericHandler[SendType, RecvType] {
	handler := &genericStreamHandler[SendType, RecvType]{
		ctx:     ctx,
		stream:  nil,
		factory: factory,
		conn:    nil,
		config:  cfg,
		channels: &channels[SendType, RecvType]{
			input:     make(chan SendType),
			output:    make(chan RecvType),
			errorChan: make(chan error),
		},
		log:           log.WithField("op", op),
		streamClosed:  false,
		errChanClosed: false,
	}
	return handler
}

type Config struct {
	address              string
	reconnectMaxAttempts int32
	timeout              time.Duration
	reconnectInterval    time.Duration
	dialOpts             []grpc.DialOption
	callOpts             []grpc.CallOption
}

type channels[SendType any, RecvType any] struct {
	input     chan SendType
	output    chan RecvType
	errorChan chan error
}

func NewConfig(
	address string,
	reconnectMaxAttempts int32,
	timeout time.Duration,
	reconnectInterval time.Duration,
) *Config {
	return &Config{
		address:              address,
		reconnectMaxAttempts: reconnectMaxAttempts,
		timeout:              timeout,
		reconnectInterval:    reconnectInterval,
		dialOpts:             nil,
		callOpts:             nil,
	}
}

// RunWithReconnect runs the generic stream handler with the ability to reconnect in case stopping of grpc.Server.
//
// It takes a context.Context as a parameter and returns an error. Other options is in Config
//
//nolint:gocognit // complex, yes
func (g *genericStreamHandler[SendType, RecvType]) RunWithReconnect(ctx context.Context) error {
	const localOp = "RunWithReconnect"
	var errGroup errgroup.Group
	defer close(g.channels.errorChan)
	defer close(g.channels.output)
	if err := g.connectAndCreateStream(); err != nil {
		return err
	}
	ctxWithCancel, cancel := context.WithCancel(ctx)
	defer cancel()
	go g.readMessages(ctxWithCancel)
	go g.writeMessages(ctxWithCancel)

	g.log.WithField("op", localOp).Infof("started successfully")

	errGroup.Go(func() error {
		for err := range g.channels.errorChan {
			if errors.Is(err, ErrChannelClosed) {
				cancel()
				_ = g.CloseSend()
				return err
			}
			g.log.WithField("op", localOp).WithError(err).Errorf("error in stream: %v", err)

			g.log.WithField("op", localOp).Infof("waiting for reconnect...")

			<-time.After(g.config.reconnectInterval)

			g.log.WithField("op", localOp).Infof("reconnecting...")

			var counter = new(int32)
			for *counter < g.config.reconnectMaxAttempts {
				if connErr := g.connectAndCreateStream(); connErr == nil {
					g.log.WithField(
						"op", localOp).Infof("reconnected successfully")
					go g.readMessages(ctxWithCancel)
					go g.writeMessages(ctxWithCancel)
					break
				}
				atomic.AddInt32(counter, 1)
				if *counter == g.config.reconnectMaxAttempts {
					g.log.WithField(
						"op", localOp).WithError(ErrMaxAttempts).Errorf("can't reconnect, exiting")
					if err = g.CloseSend(); err != nil {
						g.log.WithField("op", localOp).
							WithError(err).
							Errorf("error in stream.CloseSend()")
						return err
					}
					cancel()
					return ErrMaxAttempts
				}
				g.log.WithField(
					"attempts_counter", counter).Errorf("can't reconnect, retrying")
			}
		}
		return nil
	})
	return errGroup.Wait()
}

// GetInputChan returns the input channel. It is used to send messages of type SendType to the stream.
//
// No parameters.
// Returns a channel of type SendType.
func (g *genericStreamHandler[SendType, RecvType]) GetInputChan() chan SendType {
	return g.channels.input
}

// GetOutputChan returns the output channel. It contains messages of RecvType received from the stream.
// ***It is necessary to read messages from this channel, or goroutine will be blocked.***
//
// No parameters.
// Returns a channel of type RecvType.
func (g *genericStreamHandler[SendType, RecvType]) GetOutputChan() chan RecvType {
	return g.channels.output
}

// GetErrChanChannel returns the error channel. For inside use and testing purposes only.
//
// No parameters.
// Returns a channel of error type.
func (g *genericStreamHandler[SendType, RecvType]) GetErrChanChannel() chan error {
	return g.channels.errorChan
}

// WithGrpcDialOpts sets the gRPC dial options for the generic stream handler.
//
// It takes variadic grpc.DialOption as input and does not return anything.
func (g *genericStreamHandler[SendType, RecvType]) WithGrpcDialOpts(opts ...grpc.DialOption) {
	if len(opts) == 0 {
		return
	}
	g.config.dialOpts = append(g.config.dialOpts, opts...)
}

// WithGrpcCallOpts updates the gRPC call options for the generic stream handler.
//
// It takes variadic grpc.CallOption as input and does not return anything.
func (g *genericStreamHandler[SendType, RecvType]) WithGrpcCallOpts(opts ...grpc.CallOption) {
	if len(opts) == 0 {
		return
	}
	g.config.callOpts = append(g.config.callOpts, opts...)
}

// IsClosed checks if the genericStreamHandler is closed.
//
// None.
// boolean - true if the stream is closed, false otherwise.
func (g *genericStreamHandler[SendType, RecvType]) IsClosed() bool {
	return g.streamClosed
}

// connectAndCreateStream connects and creates a Streamer which is an abstraction of grpc.Stream.
//
// No parameters.
// Returns an error.
func (g *genericStreamHandler[SendType, RecvType]) connectAndCreateStream() error {
	var err error
	ctxTimeout, timeoutCancel := context.WithTimeout(g.ctx, g.config.timeout)
	defer timeoutCancel()
	g.conn, err = grpc.DialContext(ctxTimeout, g.config.address, g.config.dialOpts...)
	if err != nil {
		return ErrContextTimout
	}
	g.stream, err = g.factory.createStream(g.ctx, g.conn, g.config.callOpts...)
	if err != nil {
		return ErrConnectionToServer
	}
	g.errChanClosed = false
	g.streamClosed = false
	return nil
}

// readMessages reads and handles messages from grpc.Stream.
//
// ctx context.Context
// No return.
func (g *genericStreamHandler[SendType, RecvType]) readMessages(ctx context.Context) {
	const opLocal = "readMessages"
	for {
		select {
		case <-ctx.Done():
			g.log.WithField("op", opLocal).Infof("context is done, exiting")
			if err := g.CloseSend(); err != nil {
				g.log.WithField(
					"op", opLocal).WithError(err).Errorf("error in CloseSend()")
				return
			}
			return
		default:
			if g.streamClosed {
				return
			}
			g.log.WithField("op", opLocal).Infof("waiting for message...")
			msg, err := g.stream.Recv()
			if err != nil {
				if g.errChanClosed {
					return
				}
				g.streamClosed = true
				g.channels.errorChan <- err
				g.errChanClosed = true
				return
			}
			g.channels.output <- msg
		}
	}
}

// writeMessages handles writing messages into grpc.Stream.
//
// ctx: context.Context
//
//nolint:gocognit // complex, yes
func (g *genericStreamHandler[SendType, RecvType]) writeMessages(ctx context.Context) {
	const opLocal = "writeMessages"
	for {
		select {
		case <-ctx.Done():
			g.log.WithField("op", opLocal).Infof("context is done, exiting")
			if err := g.CloseSend(); err != nil {
				g.log.WithField(
					"op", opLocal).WithError(err).Errorf("error in stream.CloseSend()")
				return
			}
			return
		default:
			if g.streamClosed {
				return
			}
			//nolint:nestif // complex, yes
			if msg, ok := <-g.channels.input; !ok {
				g.log.WithField(
					"op", opLocal).WithError(ErrChannelClosed).Errorf(
					"input channel is closed, exiting",
				)
				if g.errChanClosed {
					return
				}
				g.streamClosed = true
				g.channels.errorChan <- ErrChannelClosed
				g.errChanClosed = true
				return
			} else {
				if err := g.stream.Send(msg); err != nil {
					if g.errChanClosed {
						return
					}
					g.streamClosed = true
					g.channels.errorChan <- ErrChannelClosed
					g.errChanClosed = true
					return
				}
			}
		}
	}
}
