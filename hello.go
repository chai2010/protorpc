// Copyright 2014 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

package main

import (
	"fmt"
	"log"

	"github.com/chai2010/protorpc"
	service "github.com/chai2010/protorpc/internal/service.pb"
	"github.com/golang/protobuf/proto"
)

type Echo int

func (t *Echo) Echo(args *service.EchoRequest, reply *service.EchoResponse) error {
	reply.Msg = proto.String(args.GetMsg())
	return nil
}

func (t *Echo) EchoTwice(args *service.EchoRequest, reply *service.EchoResponse) error {
	reply.Msg = proto.String(args.GetMsg() + args.GetMsg())
	return nil
}

func init() {
	go service.ListenAndServeEchoService("tcp", `127.0.0.1:9527`, new(Echo))
}

func main() {
	echoClient, err := service.DialEchoService("tcp", `127.0.0.1:9527`)
	if err != nil {
		log.Fatalf("service.DialEchoService: %v", err)
	}
	defer echoClient.Close()

	args := &service.EchoRequest{Msg: proto.String("你好, 世界!")}
	reply := &service.EchoResponse{}
	err = echoClient.EchoTwice(args, reply)
	if err != nil {
		log.Fatalf("echoClient.EchoTwice: %v", err)
	}
	fmt.Println(reply.GetMsg())

	// or use normal client
	client, err := protorpc.Dial("tcp", `127.0.0.1:9527`)
	if err != nil {
		log.Fatalf("protorpc.Dial: %v", err)
	}
	defer client.Close()

	echoClient1 := &service.EchoServiceClient{client}
	echoClient2 := &service.EchoServiceClient{client}
	echoClient1.EchoTwice(args, reply)
	echoClient2.EchoTwice(args, reply)

	// Output:
	// 你好, 世界!你好, 世界!
}
