package main

import (
	"context"

	"github.com/Guilospanck/gRPC/TensorProgramming/BasicAPI/proto"
)

type server struct {
	proto.uni
}

func (s *server) Add(ctx context.Context, request *proto.Request) (*proto.Response, error) {

}

func (s *server) Multiply(ctx context.Context, request *proto.Request) (*proto.Response, error) {

}
