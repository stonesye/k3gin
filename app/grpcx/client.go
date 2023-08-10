package grpcx

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"k3gin/app/config"
	"k3gin/app/logger"
	"strconv"
)

func InitClientRPC() (*grpc.ClientConn, func(), error) {
	var C = config.C.GRPC

	var opts []grpc.DialOption

	if C.CACert != "" {

		creds, err := credentials.NewClientTLSFromFile(C.CACert, C.HostOverride)

		if err != nil {
			logger.WithContext(context.TODO()).Errorf("failed to load credentials: %v", err)
			return nil, nil, err
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", C.Host, strconv.Itoa(C.Port)), opts...)

	if err != nil {
		logger.WithContext(context.TODO()).Errorf("failed to dial : %v", err)
		return nil, nil, err
	}

	return conn, func() {
		conn.Close()
	}, nil
}
