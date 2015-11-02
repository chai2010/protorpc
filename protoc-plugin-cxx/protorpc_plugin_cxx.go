// Copyright 2015 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// grpc_cpp_plugin is a plugin for the Google protocol buffer compiler to generate
// protorpc.cxx code.
package main

import (
	"github.com/golang/protobuf/protoc-gen-go/generator"
)

type CxxGenerator struct {
	*generator.Generator
}

func NewCxxGenerator() *CxxGenerator {
	return &CxxGenerator{generator.New()}
}

func (p *CxxGenerator) WrapTypes() {
	//
}

func (p *CxxGenerator) SetPackageNames() {
	//
}
func (p *CxxGenerator) BuildTypeNameMap() {
	//
}

func (p *CxxGenerator) GenerateAllFiles() {
	//
}
