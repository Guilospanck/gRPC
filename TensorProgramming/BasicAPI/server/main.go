package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/Guilospanck/gRPC/TensorProgramming/BasicAPI/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	port = flag.Int("port", 8080, "The server port")
)

type server struct {
	proto.UnimplementedAddServiceServer
}

func (s *server) Add(ctx context.Context, request *proto.Request) (*proto.Response, error) {
	firstNum := request.A
	secondNum := request.B

	result := firstNum + secondNum

	return &proto.Response{Result: result}, nil
}

func (s *server) Multiply(ctx context.Context, request *proto.Request) (*proto.Response, error) {
	firstNum, secondNum := request.A, request.B

	result := firstNum * secondNum

	return &proto.Response{Result: result}, nil
}

func newServer() *server {
	return &server{}
}

func main() {
	flag.Parse()

	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("Failed to start TCP server: %v", err)
	}

	var opts []grpc.ServerOption

	// creates an instance of gRPC server
	grpcServer := grpc.NewServer(opts...)

	// register service implementation with the gRPC server
	proto.RegisterAddServiceServer(grpcServer, newServer())

	// reflection to serialize data
	reflection.Register(grpcServer)

	// listens
	if err := grpcServer.Serve(listener); err != nil {
		panic(err)
	}

}
