package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"grpcpool/greeter-grpc/proto"
	"log"
	"net"
)

var (
	port = flag.Int("port", 50051, "")
)

type server struct {
	proto.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *proto.HelloRequest) (*proto.HelloReply, error) {
	fmt.Println("Server Recv:", in.Msg)
	return &proto.HelloReply{Msg: "hello Client"}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer()
	proto.RegisterGreeterServer(s, &server{})
	s.Serve(lis)
}
