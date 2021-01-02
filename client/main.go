package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	pb "hello/protos"
	sv "hello/server/services"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	req := pb.HelloRequest{Name: "gavin"}
	cert, err := tls.LoadX509KeyPair("../certs/client.pem", "../certs/client.key")
	//certPool := x509.NewCertPool()
	// ca, _ := ioutil.ReadFile("ca.pem")
	// certPool.AppendCertsFromPEM(ca)

	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ServerName:   "localhost",
		//RootCAs:      certPool,
	})
	// GRPC
	conn, err := grpc.Dial("localhost:9090", grpc.WithTransportCredentials(creds))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	c := sv.NewServer()
	r, err := c.SayHello(context.Background(), &req)
	if err != nil {
		fmt.Println(err)
	}
	// 打印返回值
	fmt.Println(r)
	fmt.Println("http Start......................")
	//http
	/*
		t := &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{cert},
				// RootCAs:      caCertPool,
			},
		}
	*/
	requestByte, _ := json.Marshal(req)
	client := http.Client{Timeout: 15 * time.Second}
	resp, err := client.Post("http://localhost:8080/hello_world", "application/json", strings.NewReader(string(requestByte)))

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	fmt.Println(string(bodyBytes))
}
