// Copyright 2017 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a Apache
// license that can be found in the LICENSE file.

// protoc-gen-protorpc is a plugin for the Google protocol buffer compiler to generate
// Go code.  Run it by building this program and putting it in your path with
// the name
// 	protoc-gen-go
// That word 'go' at the end becomes part of the option string set for the
// protocol compiler, so once the protocol compiler (protoc) is installed
// you can run
// 	protoc --protorpc_out=output_directory input_directory/file.proto
// to generate Go bindings for the protocol defined by file.proto.
// With that input, the output will be written to
// 	output_directory/file.pb.go
//
// The generated code is documented in the package comment for
// the library.
//
// See the README and documentation for protocol buffers to learn more:
// 	https://developers.google.com/protocol-buffers/
package main

import (
	"github.com/chai2010/protorpc/protoc-gen-protorpc/plugin"
	_ "github.com/chai2010/protorpc/protoc-gen-protorpc/protorpc-go"
)

func main() {
	plugin.Main()
}
