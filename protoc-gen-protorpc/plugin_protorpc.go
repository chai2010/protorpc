// Copyright 2017 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a Apache
// license that can be found in the LICENSE file.

package main

import (
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/golang/protobuf/protoc-gen-go/generator"
)

func init() {
	RegisterCodeGenerator(nil)
}

type protorpcCodeGenerator struct {
	//
}

func (p *protorpcCodeGenerator) Name() string {
	return "protorpc-go"
}
func (p *protorpcCodeGenerator) FileNameExt() string {
	return ".pb.protorpc.go"
}

func (p *protorpcCodeGenerator) HeaderCode(g *generator.Generator, file *generator.FileDescriptor) string {
	return ""
}
func (p *protorpcCodeGenerator) ServiceCode(g *generator.Generator, file *generator.FileDescriptor, svc *descriptor.ServiceDescriptorProto) string {
	return ""
}

func (p *protorpcCodeGenerator) MessageCode(g *generator.Generator, file *generator.FileDescriptor, msg *descriptor.DescriptorProto) string {
	return ""
}
