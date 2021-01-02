package services

import (
	"context"
	"fmt"
	pb "hello/protos"
)

type server struct{}

func NewServer() *server {
	return &server{}
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	fmt.Println("request: ", in.Name)
	return &pb.HelloReply{Message: "hello, " + in.Name}, nil
}
