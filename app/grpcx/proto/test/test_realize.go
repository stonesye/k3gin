package test

import (
	"context"
	"errors"
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
