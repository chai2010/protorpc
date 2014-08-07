// Copyright 2013 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package protorpc implements a Protobuf-RPC ClientCodec and ServerCodec
for the rpc package.

To install it, you must first have Go (version 1) installed
(see http://golang.org/doc/install). Next, install the standard
protocol buffer implementation from http://code.google.com/p/protobuf/;
you must be running version 2.3 or higher.

Finally run

	go get code.google.com/p/protorpc
	go get code.google.com/p/protorpc/protoc-gen-go

to install the support library and protocol compiler.

Here is a simple proto file("arith.pb/arith.proto"):

	package arith;

	// go use cc_generic_services option
	option cc_generic_services = true;

	message ArithRequest {
		optional int32 a = 1;
		optional int32 b = 2;
	}

	message ArithResponse {
		optional int32 val = 1;
		optional int32 quo = 2;
		optional int32 rem = 3;
	}

	service ArithService {
		rpc multiply (ArithRequest) returns (ArithResponse);
		rpc divide (ArithRequest) returns (ArithResponse);
	}

Then use "protoc-gen-go" to generate "arith.pb.go" file(include rpc stub):

	cd arith.pb && protoc --go_out=. arith.proto

The server calls (for TCP service):

	package server

	import (
		"errors"

		"code.google.com/p/goprotobuf/proto"

		"./arith.pb"
	)

	type Arith int

	func (t *Arith) Multiply(args *arith.ArithRequest, reply *arith.ArithResponse) error {
		reply.Val = proto.Int32(args.GetA() * args.GetB())
		return nil
	}

	func (t *Arith) Divide(args *arith.ArithRequest, reply *arith.ArithResponse) error {
		if args.GetB() == 0 {
			return errors.New("divide by zero")
		}
		reply.Quo = proto.Int32(args.GetA() / args.GetB())
		reply.Rem = proto.Int32(args.GetA() % args.GetB())
		return nil
	}

	func main() {
		arith.ListenAndServeArithService("tcp", ":1984", new(Arith))
	}

At this point, clients can see a service "Arith" with methods "ArithService.Multiply" and
"ArithService.Divide". To invoke one, a client first dials the server:

	stub, client, err := arith.DialArithService("tcp", "127.0.0.1:1984")
	if err != nil {
		log.Fatal(`arith.DialArithService("tcp", "127.0.0.1:1984"):`, err)
	}
	defer client.Close()

Then it can make a remote call with stub:

	var args ArithRequest
	var reply ArithResponse

	args.A = proto.Int32(7)
	args.B = proto.Int32(8)
	if err = stub.Multiply(&args, &reply); err != nil {
		log.Fatal("arith error:", err)
	}
	fmt.Printf("Arith: %d*%d=%d", args.GetA(), args.GetB(), reply.GetVal())

More example:

	go test code.google.com/p/protorpc/service.pb

It's very simple to use "Protobuf-RPC" with "protoc-gen-go" tool.
Try it out.
*/
package protorpc
