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
	ServerGetTestID(ctx context.Context, in *Test, opts ...grpc.CallOption) (*TestID, error)
	ServerStreamEcho(ctx context.Context, opts ...grpc.CallOption) (TestInfo_ServerStreamEchoClient, error)
}

type testInfoClient struct {
	cc grpc.ClientConnInterface
}

func NewTestInfoClient(cc grpc.ClientConnInterface) TestInfoClient {
	return &testInfoClient{cc}
}

func (c *testInfoClient) ServerGetTestID(ctx context.Context, in *Test, opts ...grpc.CallOption) (*TestID, error) {
	out := new(TestID)
	err := c.cc.Invoke(ctx, "/proto.TestInfo/ServerGetTestID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *testInfoClient) ServerStreamEcho(ctx context.Context, opts ...grpc.CallOption) (TestInfo_ServerStreamEchoClient, error) {
	stream, err := c.cc.NewStream(ctx, &TestInfo_ServiceDesc.Streams[0], "/proto.TestInfo/ServerStreamEcho", opts...)
	if err != nil {
		return nil, err
	}
	x := &testInfoServerStreamEchoClient{stream}
	return x, nil
}

type TestInfo_ServerStreamEchoClient interface {
	Send(*TestRequest) error
	Recv() (*TestResponse, error)
	grpc.ClientStream
}

type testInfoServerStreamEchoClient struct {
	grpc.ClientStream
}

func (x *testInfoServerStreamEchoClient) Send(m *TestRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *testInfoServerStreamEchoClient) Recv() (*TestResponse, error) {
	m := new(TestResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// TestInfoServer is the server API for TestInfo service.
// All implementations must embed UnimplementedTestInfoServer
// for forward compatibility
type TestInfoServer interface {
	ServerGetTestID(context.Context, *Test) (*TestID, error)
	ServerStreamEcho(TestInfo_ServerStreamEchoServer) error
	mustEmbedUnimplementedTestInfoServer()
}

// UnimplementedTestInfoServer must be embedded to have forward compatible implementations.
type UnimplementedTestInfoServer struct {
}

func (UnimplementedTestInfoServer) ServerGetTestID(context.Context, *Test) (*TestID, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ServerGetTestID not implemented")
}
func (UnimplementedTestInfoServer) ServerStreamEcho(TestInfo_ServerStreamEchoServer) error {
	return status.Errorf(codes.Unimplemented, "method ServerStreamEcho not implemented")
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

func _TestInfo_ServerGetTestID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Test)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TestInfoServer).ServerGetTestID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.TestInfo/ServerGetTestID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TestInfoServer).ServerGetTestID(ctx, req.(*Test))
	}
	return interceptor(ctx, in, info, handler)
}

func _TestInfo_ServerStreamEcho_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(TestInfoServer).ServerStreamEcho(&testInfoServerStreamEchoServer{stream})
}

type TestInfo_ServerStreamEchoServer interface {
	Send(*TestResponse) error
	Recv() (*TestRequest, error)
	grpc.ServerStream
}

type testInfoServerStreamEchoServer struct {
	grpc.ServerStream
}

func (x *testInfoServerStreamEchoServer) Send(m *TestResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *testInfoServerStreamEchoServer) Recv() (*TestRequest, error) {
	m := new(TestRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// TestInfo_ServiceDesc is the grpc.ServiceDesc for TestInfo service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TestInfo_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.TestInfo",
	HandlerType: (*TestInfoServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ServerGetTestID",
			Handler:    _TestInfo_ServerGetTestID_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ServerStreamEcho",
			Handler:       _TestInfo_ServerStreamEcho_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "test/test.proto",
}
