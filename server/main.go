package main
 
import (
	"fmt"
	pb "hello/protos"
	"hello/server/services"
	"log"
	"net"
	"hello/gateway"
	"google.golang.org/grpc"
)
 
func main() {
	grpcPort := ":9090"
	httpPort := ":8080"
	lis, err := net.Listen("tcp", grpcPort)
 
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	go gateway.HttpRun(grpcPort, httpPort)
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, services.NewServer())
	fmt.Println("rpc services started, listen on localhost" + grpcPort)
	s.Serve(lis)
}