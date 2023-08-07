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
	user_pb.UnimplementedUserInfoServer
}

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
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption

	*certFile = "/Users/yelei/data/code/go-projects/k3gin/app/gin_grpc/x509/server_cert.pem"
	*keyFile = "/Users/yelei/data/code/go-projects/k3gin/app/gin_grpc/x509/server_key.pem"
	creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)

	if err != nil {
		log.Fatalf("Fialed to generate credentials:%v", err)
	}

	opts = []grpc.ServerOption{grpc.Creds(creds)}
	s := grpc.NewServer(opts...)
	user_pb.RegisterUserInfoServer(s, &server{})

	log.Printf("server listening at %v", lis.Addr())
	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to server: %v", err)
	}
}
