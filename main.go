package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	pb "hello/protos"
	sv "hello/server/services"
	"io/ioutil"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

func GPRCServer() {
	// 监听本地端口
	listener, err := net.Listen("tcp", "localhost:8181")
	if err != nil {
		return
	}

	//证书
	cert, _ := tls.LoadX509KeyPair("server.pem", "server.key")
	certPool := x509.NewCertPool()
	ca, _ := ioutil.ReadFile("ca.pem")
	certPool.AppendCertsFromPEM(ca)

	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	})
	s := grpc.NewServer(grpc.Creds(creds))      // 创建GRPC
	pb.RegisterGreeterServer(s, sv.NewServer()) // 在GRPC服务端注册服务

	reflection.Register(s) // 在GRPC服务器注册服务器反射服务
	// Serve方法接收监听的端口,每到一个连接创建一个ServerTransport和server的grroutine
	// 这个goroutine读取GRPC请求,调用已注册的处理程序进行响应
	err = s.Serve(listener)
	if err != nil {
		return
	}
}

func GPRCClient() {
	cert, _ := tls.LoadX509KeyPair("cert/client.pem", "cert/client.key")
	certPool := x509.NewCertPool()
	ca, _ := ioutil.ReadFile("cert/ca.pem")
	certPool.AppendCertsFromPEM(ca)

	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ServerName:   "localhost",
		RootCAs:      certPool,
	})
	// 连接服务器
	conn, err := grpc.Dial("localhost:8181", grpc.WithTransportCredentials(creds))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	defer conn.Close()
	c := sv.NewServer()
	req := pb.HelloRequest{Name: "gavin"}
	r, err := c.SayHello(context.Background(), &req)
	if err != nil {
		fmt.Println(err)
	}
	// 打印返回值
	fmt.Println(r)
}

func main() {
	go GPRCServer()
	time.Sleep(1000)
	go GPRCClient()
	var s string
	fmt.Scan(&s)

}
