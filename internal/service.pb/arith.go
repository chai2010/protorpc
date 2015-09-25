// Copyright 2013 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package service

import (
	"errors"

	"github.com/golang/protobuf/proto"
)

type Arith int

func (t *Arith) Add(args *ArithRequest, reply *ArithResponse) error {
	reply.C = proto.Int32(args.GetA() + args.GetB())
	return nil
}

func (t *Arith) Mul(args *ArithRequest, reply *ArithResponse) error {
	reply.C = proto.Int32(args.GetA() * args.GetB())
	return nil
}

func (t *Arith) Div(args *ArithRequest, reply *ArithResponse) error {
	if args.GetB() == 0 {
		return errors.New("divide by zero")
	}
	reply.C = proto.Int32(args.GetA() / args.GetB())
	return nil
}

func (t *Arith) Error(args *ArithRequest, reply *ArithResponse) error {
	return errors.New("ArithError")
}
