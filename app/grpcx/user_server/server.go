package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"k3gin/app/gin_grpc/user_pb"
	"log"
	"net"
)

type server struct {
	// 其实是 user.proto 的 UserInfo, 可以理解就是user.proto, 如果有多个proto，可以写多个
	user_pb.UnimplementedUserInfoServer
}

// 实现函数
func (s *server) AddUser(ctx context.Context, in *user_pb.User) (*user_pb.UserID, error) {
	log.Printf("Received %v", in)
	return &user_pb.UserID{UserID: "数据:" + in.UserID}, nil
}

func (s *server) GetUser(ctx context.Context, in *user_pb.UserID) (*user_pb.User, error) {
	log.Printf("Received %v", in)
	return &user_pb.User{
		UserID:      in.UserID,
		Name:        "hhhhhhh",
		Description: "测试一下而已",
	}, nil
}

var (
	port     = flag.Int("port", 50051, "server port")
	certFile = flag.String("cert_file", "", "The TLS cert file")
	keyFile  = flag.String("key_file", "", "The TLS key file")
	tls      = flag.Bool("tls", true, "This is TLS option")
)

func main() {
	flag.Parse()

	// 创建监听
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption

	if *tls { // 如果是TLS 需要给option新增证书
		*certFile = "/Users/yelei/data/code/go-projects/k3gin/app/gin_grpc/x509/server_cert.pem"
		*keyFile = "/Users/yelei/data/code/go-projects/k3gin/app/gin_grpc/x509/server_key.pem"
		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		if err != nil {
			log.Fatalf("Fialed to generate credentials:%v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}

	// 创建grpc服务
	s := grpc.NewServer(opts...)

	// 将前面的proto注册到grpc服务，这里的&server{}里面包含了user.proto
	user_pb.RegisterUserInfoServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())

	// 启动grpc服务，并监听
	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to server: %v", err)
	}
}
