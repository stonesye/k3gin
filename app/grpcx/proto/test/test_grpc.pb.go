// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.11.0
// source: test/test.proto

package test

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// TestInfoClient is the client API for TestInfo service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TestInfoClient interface {
	AddSomething(ctx context.Context, in *Test, opts ...grpc.CallOption) (*TestID, error)
}

type testInfoClient struct {
	cc grpc.ClientConnInterface
}

func NewTestInfoClient(cc grpc.ClientConnInterface) TestInfoClient {
	return &testInfoClient{cc}
}

func (c *testInfoClient) AddSomething(ctx context.Context, in *Test, opts ...grpc.CallOption) (*TestID, error) {
	out := new(TestID)
	err := c.cc.Invoke(ctx, "/proto.TestInfo/AddSomething", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TestInfoServer is the server API for TestInfo service.
// All implementations must embed UnimplementedTestInfoServer
// for forward compatibility
type TestInfoServer interface {
	AddSomething(context.Context, *Test) (*TestID, error)
	mustEmbedUnimplementedTestInfoServer()
}

// UnimplementedTestInfoServer must be embedded to have forward compatible implementations.
type UnimplementedTestInfoServer struct {
}

func (UnimplementedTestInfoServer) AddSomething(context.Context, *Test) (*TestID, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddSomething not implemented")
}
func (UnimplementedTestInfoServer) mustEmbedUnimplementedTestInfoServer() {}

// UnsafeTestInfoServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TestInfoServer will
// result in compilation errors.
type UnsafeTestInfoServer interface {
	mustEmbedUnimplementedTestInfoServer()
}

func RegisterTestInfoServer(s grpc.ServiceRegistrar, srv TestInfoServer) {
	s.RegisterService(&TestInfo_ServiceDesc, srv)
}

func _TestInfo_AddSomething_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Test)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TestInfoServer).AddSomething(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.TestInfo/AddSomething",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TestInfoServer).AddSomething(ctx, req.(*Test))
	}
	return interceptor(ctx, in, info, handler)
}

// TestInfo_ServiceDesc is the grpc.ServiceDesc for TestInfo service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TestInfo_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.TestInfo",
	HandlerType: (*TestInfoServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddSomething",
			Handler:    _TestInfo_AddSomething_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "test/test.proto",
}