package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"grpcpool/greeter-grpc/proto"
	"log"
)

var (
	addr = flag.String("addr", "localhost:50051", "")
)

func main() {
	flag.Parse()
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	sayHello(conn)
}

func sayHello(conn *grpc.ClientConn) {
	c := proto.NewGreeterClient(conn)
	ctx := context.Background()
	in := &proto.HelloRequest{
		Msg: "hello Server",
	}
	r, err := c.SayHello(ctx, in)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("client Recv:", r.Msg)
}
