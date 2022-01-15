package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync"

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
	stream proto.Broadcast_CreateStreamServer
	id     string
	active bool
	error  chan error
}

type server struct {
	proto.UnimplementedBroadcastServer
	connections []*Connection
}

func (s *server) CreateStream(connect *proto.Connect, stream proto.Broadcast_CreateStreamServer) error {
	conn := &Connection{
		stream: stream,
		id:     connect.User.Id,
		active: true,
		error:  make(chan error),
	}

	s.connections = append(s.connections, conn)

	return <-conn.error
}

func (s *server) BroadcastMessage(ctx context.Context, message *proto.Message) (*proto.Close, error) {
	wait := sync.WaitGroup{}
	done := make(chan int)

	for _, conn := range s.connections {
		wait.Add(1)

		go func(message *proto.Message, conn *Connection) {
			defer wait.Done()

			if conn.active {
				err := conn.stream.Send(message)
				grpcLog.Info("Sending message to: ", conn.stream)

				if err != nil {
					grpcLog.Errorf("Error trying to send message to stream %s - Error %v", conn.stream, err)
					conn.active = false
					conn.error <- err
				}
			}

		}(message, conn)

	}

	go func() {
		wait.Wait()
		close(done)
	}()

	<-done
	return &proto.Close{}, nil
}

func newServer() *server {
	var connections []*Connection
	return &server{connections: connections}
}

func main() {
	flag.Parse()

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("Error trying to create tcp listener: %v", err)
	}

	var opts []grpc.ServerOption

	// create gRPC server
	grpcServer := grpc.NewServer(opts...)

	grpcLog.Info("Starting server at port :", *port)

	// register our server to gRPC
	proto.RegisterBroadcastServer(grpcServer, newServer())

	// listens
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Error trying to serve gRPC: %v", err)
	}
}
