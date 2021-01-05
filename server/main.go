package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"hello/gateway"
	pb "hello/protos"
	"hello/server/services"
	"io/ioutil"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

func main() {
	//grpc
	grpctslPort := ":9090"
	grpcPort := ":8081"
	httpPort := ":8080"
	///grpc tsl 用于双向认证
	go GrpcTslServer(grpctslPort)

	///普通的主要是便于gateway使用
	lis, err := net.Listen("tcp", grpcPort)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	//gateway
	go gateway.HttpRun(grpcPort, httpPort)
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, services.NewServer())
	fmt.Println("rpc services started, listen on localhost" + grpcPort)
	s.Serve(lis)
}
func GrpcTslServer(grpctslPort string) error {
	//证书
	cert, a := tls.LoadX509KeyPair("../certs/server.pem", "../certs/server.key")
	if a != nil {
		fmt.Println(a)
	}
	certPool := x509.NewCertPool()
	ca, _ := ioutil.ReadFile("../certs/ca.pem")
	certPool.AppendCertsFromPEM(ca)

	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	})
	sl := grpc.NewServer(grpc.Creds(creds)) // 创建GRPC

	pb.RegisterGreeterServer(sl, services.NewServer())

	reflection.Register(sl) // 在GRPC服务器注册服务器反射服务
	listener, err := net.Listen("tcp", grpctslPort)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("rpc tsl services started, listen on localhost" + grpctslPort)
	return sl.Serve(listener)

}
