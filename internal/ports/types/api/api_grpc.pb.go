// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v5.26.0
// source: api.proto

package api

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	APIService_SubscribeInstruments_FullMethodName    = "/api.APIService/SubscribeInstruments"
	APIService_UnsubscribeInstruments_FullMethodName  = "/api.APIService/UnsubscribeInstruments"
	APIService_ModifyNameResolvers_FullMethodName     = "/api.APIService/ModifyNameResolvers"
	APIService_ModifyInstrumentDetails_FullMethodName = "/api.APIService/ModifyInstrumentDetails"
)

// APIServiceClient is the client API for APIService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type APIServiceClient interface {
	SubscribeInstruments(ctx context.Context, in *Instruments, opts ...grpc.CallOption) (*emptypb.Empty, error)
	UnsubscribeInstruments(ctx context.Context, in *Instruments, opts ...grpc.CallOption) (*emptypb.Empty, error)
	ModifyNameResolvers(ctx context.Context, in *Instruments, opts ...grpc.CallOption) (*emptypb.Empty, error)
	ModifyInstrumentDetails(ctx context.Context, in *InstrumentDetails, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type aPIServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAPIServiceClient(cc grpc.ClientConnInterface) APIServiceClient {
	return &aPIServiceClient{cc}
}

func (c *aPIServiceClient) SubscribeInstruments(ctx context.Context, in *Instruments, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, APIService_SubscribeInstruments_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aPIServiceClient) UnsubscribeInstruments(ctx context.Context, in *Instruments, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, APIService_UnsubscribeInstruments_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aPIServiceClient) ModifyNameResolvers(ctx context.Context, in *Instruments, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, APIService_ModifyNameResolvers_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aPIServiceClient) ModifyInstrumentDetails(ctx context.Context, in *InstrumentDetails, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, APIService_ModifyInstrumentDetails_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// APIServiceServer is the server API for APIService service.
// All implementations must embed UnimplementedAPIServiceServer
// for forward compatibility
type APIServiceServer interface {
	SubscribeInstruments(context.Context, *Instruments) (*emptypb.Empty, error)
	UnsubscribeInstruments(context.Context, *Instruments) (*emptypb.Empty, error)
	ModifyNameResolvers(context.Context, *Instruments) (*emptypb.Empty, error)
	ModifyInstrumentDetails(context.Context, *InstrumentDetails) (*emptypb.Empty, error)
	mustEmbedUnimplementedAPIServiceServer()
}

// UnimplementedAPIServiceServer must be embedded to have forward compatible implementations.
type UnimplementedAPIServiceServer struct {
}

func (UnimplementedAPIServiceServer) SubscribeInstruments(context.Context, *Instruments) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SubscribeInstruments not implemented")
}
func (UnimplementedAPIServiceServer) UnsubscribeInstruments(context.Context, *Instruments) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnsubscribeInstruments not implemented")
}
func (UnimplementedAPIServiceServer) ModifyNameResolvers(context.Context, *Instruments) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ModifyNameResolvers not implemented")
}
func (UnimplementedAPIServiceServer) ModifyInstrumentDetails(context.Context, *InstrumentDetails) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ModifyInstrumentDetails not implemented")
}
func (UnimplementedAPIServiceServer) mustEmbedUnimplementedAPIServiceServer() {}

// UnsafeAPIServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to APIServiceServer will
// result in compilation errors.
type UnsafeAPIServiceServer interface {
	mustEmbedUnimplementedAPIServiceServer()
}

func RegisterAPIServiceServer(s grpc.ServiceRegistrar, srv APIServiceServer) {
	s.RegisterService(&APIService_ServiceDesc, srv)
}

func _APIService_SubscribeInstruments_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Instruments)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(APIServiceServer).SubscribeInstruments(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: APIService_SubscribeInstruments_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(APIServiceServer).SubscribeInstruments(ctx, req.(*Instruments))
	}
	return interceptor(ctx, in, info, handler)
}

func _APIService_UnsubscribeInstruments_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Instruments)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(APIServiceServer).UnsubscribeInstruments(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: APIService_UnsubscribeInstruments_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(APIServiceServer).UnsubscribeInstruments(ctx, req.(*Instruments))
	}
	return interceptor(ctx, in, info, handler)
}

func _APIService_ModifyNameResolvers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Instruments)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(APIServiceServer).ModifyNameResolvers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: APIService_ModifyNameResolvers_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(APIServiceServer).ModifyNameResolvers(ctx, req.(*Instruments))
	}
	return interceptor(ctx, in, info, handler)
}

func _APIService_ModifyInstrumentDetails_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InstrumentDetails)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(APIServiceServer).ModifyInstrumentDetails(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: APIService_ModifyInstrumentDetails_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(APIServiceServer).ModifyInstrumentDetails(ctx, req.(*InstrumentDetails))
	}
	return interceptor(ctx, in, info, handler)
}

// APIService_ServiceDesc is the grpc.ServiceDesc for APIService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var APIService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.APIService",
	HandlerType: (*APIServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SubscribeInstruments",
			Handler:    _APIService_SubscribeInstruments_Handler,
		},
		{
			MethodName: "UnsubscribeInstruments",
			Handler:    _APIService_UnsubscribeInstruments_Handler,
		},
		{
			MethodName: "ModifyNameResolvers",
			Handler:    _APIService_ModifyNameResolvers_Handler,
		},
		{
			MethodName: "ModifyInstrumentDetails",
			Handler:    _APIService_ModifyInstrumentDetails_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api.proto",
}