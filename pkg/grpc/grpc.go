package grpc

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"studentgit.kata.academy/quant/torque/pkg/logster"
)

func NewGRPCServer(logger logster.Logger) *grpc.Server {
	return grpc.NewServer(
		grpc.ChainUnaryInterceptor(UnaryLogger(logger), UnaryValidator, UnaryInternalErrorSetter),
		grpc.ChainStreamInterceptor(StreamLogger(logger), StreamInternalErrorSetter),
	)
}

//nolint:nonamedreturns // for clarity
func UnaryValidator(
	ctx context.Context,
	req interface{},
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	if v, ok := req.(interface{ Validate() error }); ok {
		if errV := v.Validate(); errV != nil {
			return nil, status.Error(codes.InvalidArgument, errV.Error())
		}
	}
	resp, err = handler(ctx, req)
	return resp, err
}

//nolint:nonamedreturns // for clarity
func UnaryInternalErrorSetter(
	ctx context.Context,
	req interface{},
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	resp, err = handler(ctx, req)
	if _, ok := status.FromError(err); !ok {
		err = status.Error(codes.Internal, err.Error())
	}
	return resp, err
}

func UnaryLogger(logger logster.Logger) grpc.UnaryServerInterceptor {
	//nolint:nonamedreturns // for clarity
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		l := logger.WithField("method", info.FullMethod)
		l.Debugf("start grpc unary request handling")

		start := time.Now()
		resp, err = handler(ctx, req)
		end := time.Since(start).Milliseconds()

		l = l.WithField("elapsed_ms", end)

		if err != nil {
			l.WithError(err).WithField("request", req).Errorf("grpc unary request error")
		} else {
			l.Infof("grpc unary request executed")
		}

		return resp, err
	}
}

func StreamInternalErrorSetter(
	srv interface{},
	ss grpc.ServerStream,
	_ *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	err := handler(srv, ss)
	if _, ok := status.FromError(err); !ok {
		err = status.Error(codes.Internal, err.Error())
	}
	return err
}

func StreamLogger(logger logster.Logger) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		l := logger.WithField("method", info.FullMethod)
		l.Debugf("start grpc stream handling")

		start := time.Now()
		err := handler(srv, ss)
		end := time.Since(start).Milliseconds()

		l = l.WithField("elapsed_ms", end)

		if err != nil {
			l.WithError(err).Errorf("grpc stream error")
		} else {
			l.Infof("grpc stream ended")
		}
		return err
	}
}
