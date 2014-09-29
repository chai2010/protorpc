// Copyright 2013 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package service

import (
	"github.com/chai2010/protorpc/proto"
)

type Echo int

func (t *Echo) Echo(args *EchoRequest, reply *EchoResponse) error {
	reply.Msg = proto.String(args.GetMsg())
	return nil
}

func (t *Echo) EchoTwice(args *EchoRequest, reply *EchoResponse) error {
	reply.Msg = proto.String(args.GetMsg() + args.GetMsg())
	return nil
}
