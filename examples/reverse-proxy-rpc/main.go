// Copyright 2018 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// https://github.com/grpc/grpc-go/issues/484

package main

import (
	"fmt"
	"log"
	"net"

	pb "github.com/chai2010/protorpc/examples/stdrpc.pb"
)

var _ pb.EchoService

type EchoService struct{}

func (p *EchoService) Echo(in *pb.EchoRequest, out *pb.EchoResponse) error {
	out = &pb.EchoResponse{Msg: in.GetMsg()}
	return nil
}
func (p *EchoService) EchoTwice(in *pb.EchoRequest, out *pb.EchoResponse) error {
	out = &pb.EchoResponse{Msg: in.GetMsg() + in.GetMsg()}
	return nil
}

func main() {
	fmt.Println("Listening on :9527")

	listener, err := net.Listen("tcp", ":9527")
	if err != nil {
		log.Fatal(err)
	}

	//defer fmt.Println("close tcp server")
	//defer listener.Close()

	go func() {
		fmt.Println("dial tcp localhost:9527")

		conn, err := net.Dial("tcp", "localhost:9527")
		println(444)

		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		pb.ServeEchoService(conn, new(EchoService))
		fmt.Println("close rpc server")
	}()

	for {
		println(11)
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		println(22)

		client := pb.NewEchoServiceClient(conn)
		reply, err := client.Echo(&pb.EchoRequest{Msg: "hello rpc"})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("echo:", reply.GetMsg())
	}
}
