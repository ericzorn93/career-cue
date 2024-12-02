// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: webhooks/inboundwebhooksapi/v1/api.proto

package inboundwebhooksapiv1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	v1 "packages/proto-gen/go/common/v1"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	InboundWebhooksAPI_UserCreated_FullMethodName = "/webhooks.inboundwebhooksapi.v1.InboundWebhooksAPI/UserCreated"
)

// InboundWebhooksAPIClient is the client API for InboundWebhooksAPI service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type InboundWebhooksAPIClient interface {
	// User APIs
	UserCreated(ctx context.Context, in *v1.Empty, opts ...grpc.CallOption) (*v1.Empty, error)
}

type inboundWebhooksAPIClient struct {
	cc grpc.ClientConnInterface
}

func NewInboundWebhooksAPIClient(cc grpc.ClientConnInterface) InboundWebhooksAPIClient {
	return &inboundWebhooksAPIClient{cc}
}

func (c *inboundWebhooksAPIClient) UserCreated(ctx context.Context, in *v1.Empty, opts ...grpc.CallOption) (*v1.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(v1.Empty)
	err := c.cc.Invoke(ctx, InboundWebhooksAPI_UserCreated_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// InboundWebhooksAPIServer is the server API for InboundWebhooksAPI service.
// All implementations must embed UnimplementedInboundWebhooksAPIServer
// for forward compatibility.
type InboundWebhooksAPIServer interface {
	// User APIs
	UserCreated(context.Context, *v1.Empty) (*v1.Empty, error)
	mustEmbedUnimplementedInboundWebhooksAPIServer()
}

// UnimplementedInboundWebhooksAPIServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedInboundWebhooksAPIServer struct{}

func (UnimplementedInboundWebhooksAPIServer) UserCreated(context.Context, *v1.Empty) (*v1.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserCreated not implemented")
}
func (UnimplementedInboundWebhooksAPIServer) mustEmbedUnimplementedInboundWebhooksAPIServer() {}
func (UnimplementedInboundWebhooksAPIServer) testEmbeddedByValue()                            {}

// UnsafeInboundWebhooksAPIServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to InboundWebhooksAPIServer will
// result in compilation errors.
type UnsafeInboundWebhooksAPIServer interface {
	mustEmbedUnimplementedInboundWebhooksAPIServer()
}

func RegisterInboundWebhooksAPIServer(s grpc.ServiceRegistrar, srv InboundWebhooksAPIServer) {
	// If the following call pancis, it indicates UnimplementedInboundWebhooksAPIServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&InboundWebhooksAPI_ServiceDesc, srv)
}

func _InboundWebhooksAPI_UserCreated_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(v1.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InboundWebhooksAPIServer).UserCreated(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: InboundWebhooksAPI_UserCreated_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InboundWebhooksAPIServer).UserCreated(ctx, req.(*v1.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// InboundWebhooksAPI_ServiceDesc is the grpc.ServiceDesc for InboundWebhooksAPI service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var InboundWebhooksAPI_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "webhooks.inboundwebhooksapi.v1.InboundWebhooksAPI",
	HandlerType: (*InboundWebhooksAPIServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UserCreated",
			Handler:    _InboundWebhooksAPI_UserCreated_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "webhooks/inboundwebhooksapi/v1/api.proto",
}
