package streamer

import (
	"context"

	"google.golang.org/grpc"
	"studentgit.kata.academy/quant/torque/internal/adapters/grpc_clients/types/piston"
)

type StreamFactory[SendType any, RecvType any] interface {
	createStream(ctx context.Context, conn *grpc.ClientConn,
		opts ...grpc.CallOption) (Streamer[SendType, RecvType], error)
}

type MDStream struct{}

//nolint:unused // она используется в connectAndCreateStream, не знаю почему линтер не видит
func (MDStream) createStream(
	ctx context.Context,
	conn *grpc.ClientConn,
	opts ...grpc.CallOption,
) (Streamer[*piston.MarketdataMessage, *piston.MarketdataMessage], error) {
	stream, err := piston.NewPistonClient(conn).ConnectMDGateway(ctx, opts...)
	if err != nil {
		return nil, err
	}
	return stream, nil
}

type TDStream struct{}

//nolint:unused // а этот метод на будущее
func (TDStream) createStream(
	ctx context.Context,
	conn *grpc.ClientConn,
	opts ...grpc.CallOption,
) (Streamer[*piston.TradingMessage, *piston.TradingMessage], error) {
	stream, err := piston.NewPistonClient(conn).ConnectTradingGateway(ctx, opts...)
	if err != nil {
		return nil, err
	}
	return stream, nil
}
