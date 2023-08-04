package main

import (
	"context"
	"flag"
	"google.golang.org/grpc"
	"k3gin/app/gin_grpc/user_pb"
	"log"
	"time"
)

var (
	caFile             = flag.String("ca_file", "", "The file containing the CA root cert file")
	serverAddr         = flag.String("addr", "localhost:50051", "The server address in the format of host:port")
	serverHostOverride = flag.String("server_host_override", "x.test.example.com", "The server name used to verify the hostname returned by the TLS handshake")
)

func main() {
	flag.Parse()

	var opts []grpc.DialOption

	*caFile =

	c := user_pb.NewUserInfoClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := c.AddUser(ctx, &user_pb.User{
		UserID:      "1",
		Name:        "stones",
		Description: "this is stones",
	})
	if err != nil {
		log.Fatalf("could not add user: %v", err)
	}

	log.Printf("User Add Res : %v", res)

	r, err := c.GetUser(ctx, &user_pb.UserID{UserID: "1"})
	if err != nil {
		log.Fatalf("could not get user : %v", err)
	}

	log.Printf("User Get Res : %v", r)

}
