protorpc
========

```
██████╗ ██████╗  ██████╗ ████████╗ ██████╗       ██████╗ ██████╗  ██████╗
██╔══██╗██╔══██╗██╔═══██╗╚══██╔══╝██╔═══██╗      ██╔══██╗██╔══██╗██╔════╝
██████╔╝██████╔╝██║   ██║   ██║   ██║   ██║█████╗██████╔╝██████╔╝██║     
██╔═══╝ ██╔══██╗██║   ██║   ██║   ██║   ██║╚════╝██╔══██╗██╔═══╝ ██║     
██║     ██║  ██║╚██████╔╝   ██║   ╚██████╔╝      ██║  ██║██║     ╚██████╗
╚═╝     ╚═╝  ╚═╝ ╚═════╝    ╚═╝    ╚═════╝       ╚═╝  ╚═╝╚═╝      ╚═════╝
```

[![Build Status](https://travis-ci.org/chai2010/protorpc.svg)](https://travis-ci.org/chai2010/protorpc)
[![GoDoc](https://godoc.org/github.com/chai2010/protorpc?status.svg)](https://godoc.org/github.com/chai2010/protorpc)

C++ Version: [https://github.com/chai2010/protorpc.cxx](https://github.com/chai2010/protorpc.cxx)

Talks: [Go/C++语言Protobuf-RPC简介](http://go-talks.appspot.com/github.com/chai2010/talks/chai2010-protorpc-intro.slide)

Install
=======

1. `go get github.com/chai2010/protorpc`
2. `go get github.com/chai2010/protorpc/protoc-gen-go`
3. `go run hello.go`

Example
=======

First, create [echo.proto](https://github.com/chai2010/protorpc/blob/master/internal/service.pb/echo.proto):

```Proto
package service;

message EchoRequest {
	optional string msg = 1;
}

message EchoResponse {
	optional string msg = 1;
}

service EchoService {
	rpc Echo (EchoRequest) returns (EchoResponse);
	rpc EchoTwice (EchoRequest) returns (EchoResponse);
}
```

Second, generate [echo.pb.go](https://github.com/chai2010/protorpc/blob/master/internal/service.pb/echo.pb.go)
from [echo.proto](https://github.com/chai2010/protorpc/blob/master/internal/service.pb/echo.proto) (we can use `go generate` to invoke this command, see [proto.go](https://github.com/chai2010/protorpc/blob/master/internal/service.pb/proto.go)).

	protoc --go_out=plugins=protorpc:. echo.proto


Now, we can use the stub code like this: 

```Go
package main

import (
	"fmt"
	"log"

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
}
```

BUGS
====

Report bugs to <chaishushan@gmail.com>.

Thanks!
