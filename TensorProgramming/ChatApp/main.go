package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/Guilospanck/gRPC/TensorProgramming/ChatApp/proto"
	"google.golang.org/grpc"
	glog "google.golang.org/grpc/grpclog"
)

var (
	grpcLog glog.LoggerV2
	port    = flag.Int("port", 8080, "gRPC server port")
)

func init() {
	// when the applications starts, init is called and then the logger is set up
	grpcLog = glog.NewLoggerV2(os.Stdout, os.Stdout, os.Stdout)
}

type Connection struct {
}

type server struct{}

func newServer() *server {
	return &server{}
}

func main() {
	flag.Parse()

	listener, err := net.Listen("tcp", fmt.Sprint("localhost:%d", *port))
	if err != nil {
		log.Fatalf("Error trying to create tcp listener: %v", err)
	}

	var opts []grpc.ServerOption

	// create gRPC server
	grpcServer := grpc.NewServer(opts...)

	// register our server to gRPC
	proto.RegisterBroadcastServer(grpcServer, newServer())

	// listens
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Error trying to serve gRPC: %v", err)
	}
}
