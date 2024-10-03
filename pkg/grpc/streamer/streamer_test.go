package streamer_test

import (
	"context"
	"io"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"studentgit.kata.academy/quant/torque/internal/domain"
	"studentgit.kata.academy/quant/torque/pkg/grpc/streamer"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	pb "studentgit.kata.academy/quant/torque/internal/adapters/grpc_clients/types/piston"
	"studentgit.kata.academy/quant/torque/pkg/logster"
)

const bufSize = 1024 * 1024

//nolint:gochecknoglobals //for testing purposes
var lis *bufconn.Listener

type mockPistonServer struct {
	pb.UnimplementedPistonServer
}

type StreamerTestSuite struct {
	suite.Suite
	srv *grpc.Server
	wg  sync.WaitGroup
}

func (s *StreamerTestSuite) SetupTest() {
	lis = bufconn.Listen(bufSize)
	s.srv = grpc.NewServer()
	pb.RegisterPistonServer(s.srv, &mockPistonServer{
		UnimplementedPistonServer: pb.UnimplementedPistonServer{},
	})
	go func() {
		_ = s.srv.Serve(lis)
	}()
}

func (s *StreamerTestSuite) TearDownSuite() {
	s.srv.GracefulStop()
	if lis != nil {
		_ = lis.Close()
	}
}

func TestStreamerTestSuite(t *testing.T) {
	suite.Run(t, new(StreamerTestSuite))
}

func (s *mockPistonServer) ConnectMDGateway(stream pb.Piston_ConnectMDGatewayServer) error {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		for i := 0; i < 5; i++ {
			err := stream.Send(&pb.MarketdataMessage{
				MsgType: pb.MarketdataMsgType_SUBSCRIBE,
				Message: &pb.MarketdataMessage_SubscribeMessage{
					SubscribeMessage: &pb.SubscribeMessage{
						Instruments: []string{"BTC/USDT"},
					},
				},
			})
			if err != nil {
				return
			}
		}
		wg.Done()
	}()
	go func() {
		for i := 0; i < 5; i++ {
			in, err := stream.Recv()
			if err != nil {
				return
			}
			_ = stream.Send(in)
		}
		wg.Done()
	}()
	wg.Wait()
	_ = lis.Close()
	return io.EOF
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func (s *StreamerTestSuite) TestMarketDataStreamerReconnect() {
	logger := logster.Discard
	ctx := context.Background()
	cfg := streamer.NewConfig("bufnet", 5, 5*time.Millisecond, 5*time.Millisecond)
	handler := streamer.
		NewGenericStreamHandler[*pb.MarketdataMessage, *pb.MarketdataMessage](
		ctx, new(streamer.MDStream), cfg, logger,
	)
	handler.
		WithGrpcDialOpts(
			grpc.WithContextDialer(bufDialer),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithBlock())
	handler.WithGrpcCallOpts(grpc.WaitForReady(true))

	input := handler.GetInputChan()
	output := handler.GetOutputChan()
	orderbook := domain.OrderBook{
		Exchange:   domain.OKXExchange,
		Instrument: "BTC/USDT",
		Asks:       map[string]float64{"123": 1, "456": 2, "789": 3, "321": 4, "654": 5},
		Bids:       map[string]float64{"123": 1, "456": 2, "789": 3, "321": 4, "654": 5},
		Checksum:   0,
	}
	orderbookToSent := pb.OrderBookToMarketDataMessage(orderbook)
	s.wg.Add(3)
	go func() {
		for i := 0; i < 5; i++ {
			input <- orderbookToSent
		}
		s.wg.Done()
	}()
	go func() {
		for msg := range output {
			assert.IsType(s.T(), &pb.MarketdataMessage{
				MsgType: 0,
				Message: nil,
			}, msg)
		}
		s.wg.Done()
	}()

	go func() {
		if err := handler.RunWithReconnect(ctx); err != nil {
			require.Error(s.T(), err, streamer.ErrMaxAttempts, streamer.ErrContextTimout)
			s.wg.Done()
		}
	}()
	s.wg.Wait()
}

func (s *StreamerTestSuite) TestMarketDataStreamerCtxShutdown() {
	logger := logster.Discard
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*5)
	cfg := streamer.NewConfig("bufnet", 5, 5*time.Millisecond, 5*time.Millisecond)
	handler := streamer.
		NewGenericStreamHandler[*pb.MarketdataMessage, *pb.MarketdataMessage](
		ctx, new(streamer.MDStream), cfg, logger,
	)
	handler.
		WithGrpcDialOpts(
			grpc.WithContextDialer(bufDialer),
			grpc.WithTransportCredentials(insecure.NewCredentials()))

	go func() {
		err := handler.RunWithReconnect(ctx)
		assert.ErrorIs(s.T(), err, streamer.ErrContextTimout)
	}()
	cancel()
}
