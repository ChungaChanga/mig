// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.28.2
// source: phone/service.proto

package phone

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	VerificationService_VerifyCode_FullMethodName = "/phone.VerificationService/VerifyCode"
	VerificationService_SendCode_FullMethodName   = "/phone.VerificationService/SendCode"
)

// VerificationServiceClient is the client API for VerificationService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type VerificationServiceClient interface {
	VerifyCode(ctx context.Context, in *VerifyRequest, opts ...grpc.CallOption) (*Nil, error)
	SendCode(ctx context.Context, in *SendRequest, opts ...grpc.CallOption) (*Nil, error)
}

type verificationServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewVerificationServiceClient(cc grpc.ClientConnInterface) VerificationServiceClient {
	return &verificationServiceClient{cc}
}

func (c *verificationServiceClient) VerifyCode(ctx context.Context, in *VerifyRequest, opts ...grpc.CallOption) (*Nil, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Nil)
	err := c.cc.Invoke(ctx, VerificationService_VerifyCode_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *verificationServiceClient) SendCode(ctx context.Context, in *SendRequest, opts ...grpc.CallOption) (*Nil, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Nil)
	err := c.cc.Invoke(ctx, VerificationService_SendCode_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// VerificationServiceServer is the server API for VerificationService service.
// All implementations must embed UnimplementedVerificationServiceServer
// for forward compatibility.
type VerificationServiceServer interface {
	VerifyCode(context.Context, *VerifyRequest) (*Nil, error)
	SendCode(context.Context, *SendRequest) (*Nil, error)
	mustEmbedUnimplementedVerificationServiceServer()
}

// UnimplementedVerificationServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedVerificationServiceServer struct{}

func (UnimplementedVerificationServiceServer) VerifyCode(context.Context, *VerifyRequest) (*Nil, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VerifyCode not implemented")
}
func (UnimplementedVerificationServiceServer) SendCode(context.Context, *SendRequest) (*Nil, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendCode not implemented")
}
func (UnimplementedVerificationServiceServer) mustEmbedUnimplementedVerificationServiceServer() {}
func (UnimplementedVerificationServiceServer) testEmbeddedByValue()                             {}

// UnsafeVerificationServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to VerificationServiceServer will
// result in compilation errors.
type UnsafeVerificationServiceServer interface {
	mustEmbedUnimplementedVerificationServiceServer()
}

func RegisterVerificationServiceServer(s grpc.ServiceRegistrar, srv VerificationServiceServer) {
	// If the following call pancis, it indicates UnimplementedVerificationServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&VerificationService_ServiceDesc, srv)
}

func _VerificationService_VerifyCode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VerifyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VerificationServiceServer).VerifyCode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: VerificationService_VerifyCode_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VerificationServiceServer).VerifyCode(ctx, req.(*VerifyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VerificationService_SendCode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VerificationServiceServer).SendCode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: VerificationService_SendCode_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VerificationServiceServer).SendCode(ctx, req.(*SendRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// VerificationService_ServiceDesc is the grpc.ServiceDesc for VerificationService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var VerificationService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "phone.VerificationService",
	HandlerType: (*VerificationServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "VerifyCode",
			Handler:    _VerificationService_VerifyCode_Handler,
		},
		{
			MethodName: "SendCode",
			Handler:    _VerificationService_SendCode_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "phone/service.proto",
}
