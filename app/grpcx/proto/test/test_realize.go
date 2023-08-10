package test

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"io"
	"k3gin/app/logger"
)

type TestRealize struct {
	UnimplementedTestInfoServer // 匿名结构体，此匿名结构体已经实现了TestInfoServer ,意味着 TestRealize也实现了TestInfoServer
}

// 下面重写TestInfoServer 的接口函数

func (x *TestRealize) ServerGetTestID(ctx context.Context, test *Test) (*TestID, error) {
	logger.WithContext(ctx).Infof("Recv client testinfo :%v", test)

	if test.Flag == true {
		return &TestID{
			ID: test.ID,
		}, nil
	}
	return nil, errors.New("Test info flag is false ")
}

func (x *TestRealize) ServerStreamEcho(server TestInfo_ServerStreamEchoServer) error {
	for { // 循环获取信息
		in, err := server.Recv()
		if err != nil {
			if err == io.EOF { // 数据获取完了
				return nil
			}
			logger.WithContext(context.TODO()).Infof("Receiving message from steam : %v", err)
			return err
		}

		logger.WithContext(context.TODO()).Infof("Receiving message : %v", in.Message)
		err = server.Send(&TestResponse{
			Message: in.Message,
		})

		if err != nil {
			return err
		}
	}
}

func CallServerGetTestID(ctx context.Context, client TestInfoClient, test *Test, opts ...grpc.CallOption) error {
	id, err := client.ServerGetTestID(ctx, test, opts...)
	if err != nil {
		return status.Errorf(status.Code(err), "ServerGetTestID RPC failed : %v", err)
	}

	logger.WithContext(ctx).Infof("ServerGetTestID: %v", id)

	return nil
}

func CallServerStreamEcho(ctx context.Context, client TestInfoClient, request *TestRequest, opts ...grpc.CallOption) error {

	c, err := client.ServerStreamEcho(ctx, opts...)

	if err != nil {
		return status.Errorf(status.Code(err), "ServerStreamEcho RPC failed : %v", err)
	}

	for i := 0; i < 5; i++ { // 为了测试 发5次信息给服务端
		err = c.Send(request)
		if err == io.EOF {
			break
		}

		if err != nil {
			return status.Errorf(status.Code(err), "sending  message: %v", err)
		}
	}
	c.CloseSend()

	for {
		resp, err := c.Recv() // 一直收数据
		if err == io.EOF {
			break
		}
		if err != nil {
			return status.Errorf(status.Code(err), "receiving  message: %v", err)
		}
		fmt.Println("Receive: ", resp.Message)
	}
	return nil
}
