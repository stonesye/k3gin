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

/**
放在wire初始化里面
*/

// InitClientRPC
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

type ClientRPC struct {
	TLS          bool
	CaFile       string
	Address      string
	HostOverride string
}

func WithClientRPCTLS(tls bool) func(*ClientRPC) {
	return func(rpc *ClientRPC) {
		rpc.TLS = tls
	}
}

func WithClientCaFile(caFile string) func(*ClientRPC) {
	return func(rpc *ClientRPC) {
		rpc.CaFile = caFile
	}
}

func WithClientAddress(address string) func(*ClientRPC) {
	return func(rpc *ClientRPC) {
		rpc.Address = address
	}
}

func WithHostOverride(hostOverride string) func(*ClientRPC) {
	return func(rpc *ClientRPC) {
		rpc.HostOverride = hostOverride
	}
}

func NewClientRPC(client *ClientRPC) (*grpc.ClientConn, func(), error) {
	var opts []grpc.DialOption
	if client.TLS {
		creds, err := credentials.NewClientTLSFromFile(client.CaFile, client.HostOverride)
		if err != nil {
			return nil, nil, err
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	conn, err := grpc.Dial(client.Address, opts...)
	if err != nil {

		return nil, nil, err
	}
	return conn, func() {
		conn.Close()
	}, nil
}
