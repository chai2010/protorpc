// Copyright 2013 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package service

import (
	"fmt"
	"log"
	"net/rpc"
	"sync"
	"testing"
	"unicode/utf8"

	"code.google.com/p/goprotobuf/proto"
)

var (
	echoHost = "127.0.0.1"
	echoPort = 2015

	echoRequest         = "Hello, new gopher!"
	echoResponse        = echoRequest + echoRequest
	echoMassiveRequest  = makeMassive("Hello, 世界.")
	echoMassiveResponse = echoMassiveRequest + echoMassiveRequest

	onceEcho sync.Once
)

func makeMassive(args string) string {
	runeLen := utf8.RuneCountInString(args)
	runeBuf := make([]rune, runeLen*1024*1024*2)
	for i := 0; i < 1024*1024*2; i++ {
		offset := i * runeLen
		j := 0
		for _, r := range args {
			runeBuf[offset+j] = r
			j++
		}
	}
	return string(runeBuf)
}

func setupEchoServer() {
	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Wait()

	go func() {
		wg.Done()

		addr := fmt.Sprintf(":%d", echoPort)
		err := ListenAndServeEchoService("tcp", addr, new(Echo))
		if err != nil {
			log.Fatalf("ListenAndServeEchoService: %v", err)
		}
	}()
}

func TestEchoService(t *testing.T) {
	onceEcho.Do(setupEchoServer)

	addr := fmt.Sprintf("%s:%d", echoHost, echoPort)
	_, client, err := DialEchoService("tcp", addr)
	if err != nil {
		t.Fatalf(
			`net.Dial("tcp", "%s:%d"): %v`,
			echoHost, echoPort,
			err,
		)
	}
	defer client.Close()

	testEchoService(t, client)
}

func testEchoService(t *testing.T, client *rpc.Client) {
	var args EchoRequest
	var reply EchoResponse
	var err error

	// EchoService.EchoTwice
	args.Msg = proto.String(echoRequest)
	err = client.Call("EchoService.EchoTwice", &args, &reply)
	if err != nil {
		t.Fatalf(`EchoService.EchoTwice: %v`, err)
	}
	if reply.GetMsg() != echoResponse {
		t.Fatalf(
			`EchoService.EchoTwice: expected = "%s", got = "%s"`,
			echoResponse, reply.GetMsg(),
		)
	}

	// EchoService.EchoTwice (Massive)
	args.Msg = proto.String(echoMassiveRequest)
	err = client.Call("EchoService.EchoTwice", &args, &reply)
	if err != nil {
		t.Fatalf(`EchoService.EchoTwice: %v`, err)
	}
	if reply.GetMsg() != echoMassiveResponse {
		got := reply.GetMsg()
		if len(got) > 8 {
			got = got[:8] + "..."
		}
		t.Fatalf(`EchoService.EchoTwice: len = %d, got = %v`,
			len(reply.GetMsg()), got,
		)
	}
}

func TestClientSyncEcho(t *testing.T) {
	onceEcho.Do(setupEchoServer)

	addr := fmt.Sprintf("%s:%d", echoHost, echoPort)
	echoClient, client, err := DialEchoService("tcp", addr)
	if err != nil {
		t.Fatalf(
			`net.Dial("tcp", "%s:%d"): %v`,
			echoHost, echoPort,
			err,
		)
	}
	defer client.Close()

	var args EchoRequest
	var reply EchoResponse

	// EchoService.EchoTwice
	args.Msg = proto.String("abc")
	err = echoClient.EchoTwice(&args, &reply)
	if err != nil {
		t.Fatalf(`EchoService.EchoTwice: %v`, err)
	}
	if reply.GetMsg() != args.GetMsg()+args.GetMsg() {
		t.Fatalf(
			`EchoService.EchoTwice: expected = "%s", got = "%s"`,
			args.GetMsg()+args.GetMsg(), reply.GetMsg(),
		)
	}

	// EchoService.EchoTwice
	args.Msg = proto.String("你好, 世界")
	err = echoClient.EchoTwice(&args, &reply)
	if err != nil {
		t.Fatalf(`EchoService.EchoTwice: %v`, err)
	}
	if reply.GetMsg() != args.GetMsg()+args.GetMsg() {
		t.Fatalf(
			`EchoService.EchoTwice: expected = "%s", got = "%s"`,
			args.GetMsg()+args.GetMsg(), reply.GetMsg(),
		)
	}
}

func TestClientSyncMassive(t *testing.T) {
	onceEcho.Do(setupEchoServer)

	addr := fmt.Sprintf("%s:%d", echoHost, echoPort)
	echoClient, client, err := DialEchoService("tcp", addr)
	if err != nil {
		t.Fatalf(
			`net.Dial("tcp", "%s:%d"): %v`,
			echoHost, echoPort,
			err,
		)
	}
	defer client.Close()

	var args EchoRequest
	var reply EchoResponse

	// EchoService.EchoTwice
	args.Msg = proto.String(echoMassiveRequest + "abc")
	err = echoClient.EchoTwice(&args, &reply)
	if err != nil {
		t.Fatalf(`EchoService.EchoTwice: %v`, err)
	}
	if reply.GetMsg() != args.GetMsg()+args.GetMsg() {
		got := reply.GetMsg()
		if len(got) > 8 {
			got = got[:8] + "..."
		}
		t.Fatalf(`EchoService.EchoTwice: len = %d, got = %v`,
			len(reply.GetMsg()), got,
		)
	}

	// EchoService.EchoTwice
	args.Msg = proto.String(echoMassiveRequest + "你好, 世界")
	err = echoClient.EchoTwice(&args, &reply)
	if err != nil {
		t.Fatalf(`EchoService.EchoTwice: %v`, err)
	}
	if reply.GetMsg() != args.GetMsg()+args.GetMsg() {
		got := reply.GetMsg()
		if len(got) > 8 {
			got = got[:8] + "..."
		}
		t.Fatalf(`EchoService.EchoTwice: len = %d, got = %v`,
			len(reply.GetMsg()), got,
		)
	}
}

func TestClientAsyncEcho(t *testing.T) {
	onceEcho.Do(setupEchoServer)

	addr := fmt.Sprintf("%s:%d", echoHost, echoPort)
	_, client, err := DialEchoService("tcp", addr)
	if err != nil {
		t.Fatalf(
			`net.Dial("tcp", "%s:%d"): %v`,
			echoHost, echoPort,
			err,
		)
	}
	defer client.Close()

	var args EchoRequest
	var reply EchoResponse

	// EchoService.EchoTwice
	args.Msg = proto.String(echoRequest)
	call := client.Go("EchoService.EchoTwice", &args, &reply, nil)

	call = <-call.Done
	if call.Error != nil {
		t.Fatalf(`EchoService.EchoTwice: %v`, call.Error)
	}
	if call.Reply.(*EchoResponse).GetMsg() != echoResponse {
		t.Fatalf(
			`EchoService.EchoTwice: expected = "%s", got = "%s"`,
			echoResponse, call.Reply.(*EchoResponse).GetMsg(),
		)
	}
}

func TestClientAsyncEchoBatches(t *testing.T) {
	onceEcho.Do(setupEchoServer)

	addr := fmt.Sprintf("%s:%d", echoHost, echoPort)
	_, client, err := DialEchoService("tcp", addr)
	if err != nil {
		t.Fatalf(
			`net.Dial("tcp", "%s:%d"): %v`,
			echoHost, echoPort,
			err,
		)
	}
	defer client.Close()

	var args1 EchoRequest
	var reply1 EchoResponse
	var args2 EchoRequest
	var reply2 EchoResponse
	var args3 EchoRequest
	var reply3 EchoResponse

	// EchoService.EchoTwice
	args1.Msg = proto.String("abc")
	call1 := client.Go("EchoService.EchoTwice", &args1, &reply1, nil)
	args2.Msg = proto.String("你好, 世界")
	call2 := client.Go("EchoService.EchoTwice", &args2, &reply2, nil)
	args3.Msg = proto.String("Hello, 世界")
	call3 := client.Go("EchoService.EchoTwice", &args3, &reply3, nil)

	call1 = <-call1.Done
	call2 = <-call2.Done
	call3 = <-call3.Done

	// call1
	if call1.Error != nil {
		t.Fatalf(`EchoService.EchoTwice: %v`, call1.Error)
	}
	if call1.Reply.(*EchoResponse).GetMsg() != args1.GetMsg()+args1.GetMsg() {
		t.Fatalf(
			`EchoService.EchoTwice: expected = "%s", got = "%s"`,
			args1.GetMsg()+args1.GetMsg(),
			call1.Reply.(*EchoResponse).GetMsg(),
		)
	}

	// call2
	if call2.Error != nil {
		t.Fatalf(`EchoService.EchoTwice: %v`, call2.Error)
	}
	if call2.Reply.(*EchoResponse).GetMsg() != args2.GetMsg()+args2.GetMsg() {
		t.Fatalf(
			`EchoService.EchoTwice: expected = "%s", got = "%s"`,
			args2.GetMsg()+args2.GetMsg(),
			call2.Reply.(*EchoResponse).GetMsg(),
		)
	}

	// call3
	if call3.Error != nil {
		t.Fatalf(`EchoService.EchoTwice: %v`, call3.Error)
	}
	if call3.Reply.(*EchoResponse).GetMsg() != args3.GetMsg()+args3.GetMsg() {
		t.Fatalf(
			`EchoService.EchoTwice: expected = "%s", got = "%s"`,
			args3.GetMsg()+args3.GetMsg(),
			call3.Reply.(*EchoResponse).GetMsg(),
		)
	}
}

func TestClientAsyncMassive(t *testing.T) {
	onceEcho.Do(setupEchoServer)

	addr := fmt.Sprintf("%s:%d", echoHost, echoPort)
	_, client, err := DialEchoService("tcp", addr)
	if err != nil {
		t.Fatalf(
			`net.Dial("tcp", "%s:%d"): %v`,
			echoHost, echoPort,
			err,
		)
	}
	defer client.Close()

	var args EchoRequest
	var reply EchoResponse

	// EchoService.EchoTwice
	args.Msg = proto.String(echoMassiveRequest)
	call := client.Go("EchoService.EchoTwice", &args, &reply, nil)

	call = <-call.Done
	if call.Error != nil {
		t.Fatalf(`EchoService.EchoTwice: %v`, call.Error)
	}
	if call.Reply.(*EchoResponse).GetMsg() != echoMassiveResponse {
		got := call.Reply.(*EchoResponse).GetMsg()
		if len(got) > 8 {
			got = got[:8] + "..."
		}
		t.Fatalf(`EchoService.EchoTwice: len = %d, got = %v`,
			len(call.Reply.(*EchoResponse).GetMsg()), got,
		)
	}
}

func TestClientAsyncMassiveBatches(t *testing.T) {
	onceEcho.Do(setupEchoServer)

	addr := fmt.Sprintf("%s:%d", echoHost, echoPort)
	_, client, err := DialEchoService("tcp", addr)
	if err != nil {
		t.Fatalf(
			`net.Dial("tcp", "%s:%d"): %v`,
			echoHost, echoPort,
			err,
		)
	}
	defer client.Close()

	var args1 EchoRequest
	var reply1 EchoResponse
	var args2 EchoRequest
	var reply2 EchoResponse
	var args3 EchoRequest
	var reply3 EchoResponse

	// EchoService.EchoTwice
	args1.Msg = proto.String(echoMassiveRequest + "abc")
	call1 := client.Go("EchoService.EchoTwice", &args1, &reply1, nil)
	args2.Msg = proto.String(echoMassiveRequest + "你好, 世界")
	call2 := client.Go("EchoService.EchoTwice", &args2, &reply2, nil)
	args3.Msg = proto.String(echoMassiveRequest + "Hello, 世界")
	call3 := client.Go("EchoService.EchoTwice", &args3, &reply3, nil)

	call1 = <-call1.Done
	call2 = <-call2.Done
	call3 = <-call3.Done

	// call1
	if call1.Error != nil {
		t.Fatalf(`EchoService.EchoTwice: %v`, call1.Error)
	}
	if call1.Reply.(*EchoResponse).GetMsg() != args1.GetMsg()+args1.GetMsg() {
		got := call1.Reply.(*EchoResponse).GetMsg()
		if len(got) > 8 {
			got = got[:8] + "..."
		}
		t.Fatalf(`EchoService.EchoTwice: len = %d, got = %v`,
			len(call1.Reply.(*EchoResponse).GetMsg()), got,
		)
	}

	// call2
	if call2.Error != nil {
		t.Fatalf(`EchoService.EchoTwice: %v`, call2.Error)
	}
	if call2.Reply.(*EchoResponse).GetMsg() != args2.GetMsg()+args2.GetMsg() {
		got := call2.Reply.(*EchoResponse).GetMsg()
		if len(got) > 8 {
			got = got[:8] + "..."
		}
		t.Fatalf(`EchoService.EchoTwice: len = %d, got = %v`,
			len(call2.Reply.(*EchoResponse).GetMsg()), got,
		)
	}

	// call3
	if call3.Error != nil {
		t.Fatalf(`EchoService.EchoTwice: %v`, call3.Error)
	}
	if call3.Reply.(*EchoResponse).GetMsg() != args3.GetMsg()+args3.GetMsg() {
		got := call3.Reply.(*EchoResponse).GetMsg()
		if len(got) > 8 {
			got = got[:8] + "..."
		}
		t.Fatalf(`EchoService.EchoTwice: len = %d, got = %v`,
			len(call3.Reply.(*EchoResponse).GetMsg()), got,
		)
	}
}

func BenchmarkSyncEcho(b *testing.B) {
	onceEcho.Do(setupEchoServer)

	addr := fmt.Sprintf("%s:%d", echoHost, echoPort)
	echoClient, client, err := DialEchoService("tcp", addr)
	if err != nil {
		b.Fatalf(
			`net.Dial("tcp", "%s:%d"): %v`,
			echoHost, echoPort,
			err,
		)
	}
	defer client.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var args EchoRequest
		var reply EchoResponse

		// EchoService.EchoTwice
		args.Msg = proto.String("abc")
		err = echoClient.EchoTwice(&args, &reply)
		if err != nil {
			b.Fatalf(`EchoService.EchoTwice: %v`, err)
		}
		if reply.GetMsg() != args.GetMsg()+args.GetMsg() {
			b.Fatalf(
				`EchoService.EchoTwice: expected = "%s", got = "%s"`,
				args.GetMsg()+args.GetMsg(), reply.GetMsg(),
			)
		}

		// EchoService.EchoTwice
		args.Msg = proto.String("你好, 世界")
		err = echoClient.EchoTwice(&args, &reply)
		if err != nil {
			b.Fatalf(`EchoService.EchoTwice: %v`, err)
		}
		if reply.GetMsg() != args.GetMsg()+args.GetMsg() {
			b.Fatalf(
				`EchoService.EchoTwice: expected = "%s", got = "%s"`,
				args.GetMsg()+args.GetMsg(), reply.GetMsg(),
			)
		}

		// EchoService.EchoTwice
		args.Msg = proto.String("Hello, 世界")
		err = echoClient.EchoTwice(&args, &reply)
		if err != nil {
			b.Fatalf(`EchoService.EchoTwice: %v`, err)
		}
		if reply.GetMsg() != args.GetMsg()+args.GetMsg() {
			b.Fatalf(
				`EchoService.EchoTwice: expected = "%s", got = "%s"`,
				args.GetMsg()+args.GetMsg(), reply.GetMsg(),
			)
		}
	}
}

func BenchmarkSyncMassive(b *testing.B) {
	onceEcho.Do(setupEchoServer)

	addr := fmt.Sprintf("%s:%d", echoHost, echoPort)
	echoClient, client, err := DialEchoService("tcp", addr)
	if err != nil {
		b.Fatalf(
			`net.Dial("tcp", "%s:%d"): %v`,
			echoHost, echoPort,
			err,
		)
	}
	defer client.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var args EchoRequest
		var reply EchoResponse

		// EchoService.EchoTwice
		args.Msg = proto.String(echoMassiveRequest + "abc")
		err = echoClient.EchoTwice(&args, &reply)
		if err != nil {
			b.Fatalf(`EchoService.EchoTwice: %v`, err)
		}
		if reply.GetMsg() != args.GetMsg()+args.GetMsg() {
			got := reply.GetMsg()
			if len(got) > 8 {
				got = got[:8] + "..."
			}
			b.Fatalf(`EchoService.EchoTwice: len = %d, got = %v`,
				len(reply.GetMsg()), got,
			)
		}

		// EchoService.EchoTwice
		args.Msg = proto.String(echoMassiveRequest + "你好, 世界")
		err = echoClient.EchoTwice(&args, &reply)
		if err != nil {
			b.Fatalf(`EchoService.EchoTwice: %v`, err)
		}
		if reply.GetMsg() != args.GetMsg()+args.GetMsg() {
			got := reply.GetMsg()
			if len(got) > 8 {
				got = got[:8] + "..."
			}
			b.Fatalf(`EchoService.EchoTwice: len = %d, got = %v`,
				len(reply.GetMsg()), got,
			)
		}

		// EchoService.EchoTwice
		args.Msg = proto.String(echoMassiveRequest + "Hello, 世界")
		err = echoClient.EchoTwice(&args, &reply)
		if err != nil {
			b.Fatalf(`EchoService.EchoTwice: %v`, err)
		}
		if reply.GetMsg() != args.GetMsg()+args.GetMsg() {
			got := reply.GetMsg()
			if len(got) > 8 {
				got = got[:8] + "..."
			}
			b.Fatalf(`EchoService.EchoTwice: len = %d, got = %v`,
				len(reply.GetMsg()), got,
			)
		}
	}
}

func BenchmarkAsyncEcho(b *testing.B) {
	onceEcho.Do(setupEchoServer)

	addr := fmt.Sprintf("%s:%d", echoHost, echoPort)
	_, client, err := DialEchoService("tcp", addr)
	if err != nil {
		b.Fatalf(
			`net.Dial("tcp", "%s:%d"): %v`,
			echoHost, echoPort,
			err,
		)
	}
	defer client.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var args1 EchoRequest
		var reply1 EchoResponse
		var args2 EchoRequest
		var reply2 EchoResponse
		var args3 EchoRequest
		var reply3 EchoResponse

		// EchoService.EchoTwice
		args1.Msg = proto.String("abc")
		call1 := client.Go("EchoService.EchoTwice", &args1, &reply1, nil)
		args2.Msg = proto.String("你好, 世界")
		call2 := client.Go("EchoService.EchoTwice", &args2, &reply2, nil)
		args3.Msg = proto.String("Hello, 世界")
		call3 := client.Go("EchoService.EchoTwice", &args3, &reply3, nil)

		call1 = <-call1.Done
		call2 = <-call2.Done
		call3 = <-call3.Done

		// call1
		if call1.Error != nil {
			b.Fatalf(`EchoService.EchoTwice: %v`, call1.Error)
		}
		if call1.Reply.(*EchoResponse).GetMsg() != args1.GetMsg()+args1.GetMsg() {
			b.Fatalf(
				`EchoService.EchoTwice: expected = "%s", got = "%s"`,
				args1.GetMsg()+args1.GetMsg(),
				call1.Reply.(*EchoResponse).GetMsg(),
			)
		}

		// call2
		if call2.Error != nil {
			b.Fatalf(`EchoService.EchoTwice: %v`, call2.Error)
		}
		if call2.Reply.(*EchoResponse).GetMsg() != args2.GetMsg()+args2.GetMsg() {
			b.Fatalf(
				`EchoService.EchoTwice: expected = "%s", got = "%s"`,
				args2.GetMsg()+args2.GetMsg(),
				call2.Reply.(*EchoResponse).GetMsg(),
			)
		}

		// call3
		if call3.Error != nil {
			b.Fatalf(`EchoService.EchoTwice: %v`, call3.Error)
		}
		if call3.Reply.(*EchoResponse).GetMsg() != args3.GetMsg()+args3.GetMsg() {
			b.Fatalf(
				`EchoService.EchoTwice: expected = "%s", got = "%s"`,
				args3.GetMsg()+args3.GetMsg(),
				call3.Reply.(*EchoResponse).GetMsg(),
			)
		}
	}
}

func BenchmarkAsyncMassive(b *testing.B) {
	onceEcho.Do(setupEchoServer)

	addr := fmt.Sprintf("%s:%d", echoHost, echoPort)
	_, client, err := DialEchoService("tcp", addr)
	if err != nil {
		b.Fatalf(
			`net.Dial("tcp", "%s:%d"): %v`,
			echoHost, echoPort,
			err,
		)
	}
	defer client.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var args1 EchoRequest
		var reply1 EchoResponse
		var args2 EchoRequest
		var reply2 EchoResponse
		var args3 EchoRequest
		var reply3 EchoResponse

		// EchoService.EchoTwice
		args1.Msg = proto.String(echoMassiveRequest + "abc")
		call1 := client.Go("EchoService.EchoTwice", &args1, &reply1, nil)
		args2.Msg = proto.String(echoMassiveRequest + "你好, 世界")
		call2 := client.Go("EchoService.EchoTwice", &args2, &reply2, nil)
		args3.Msg = proto.String(echoMassiveRequest + "Hello, 世界")
		call3 := client.Go("EchoService.EchoTwice", &args3, &reply3, nil)

		call1 = <-call1.Done
		call2 = <-call2.Done
		call3 = <-call3.Done

		// call1
		if call1.Error != nil {
			b.Fatalf(`EchoService.EchoTwice: %v`, call1.Error)
		}
		if call1.Reply.(*EchoResponse).GetMsg() != args1.GetMsg()+args1.GetMsg() {
			got := call1.Reply.(*EchoResponse).GetMsg()
			if len(got) > 8 {
				got = got[:8] + "..."
			}
			b.Fatalf(`EchoService.EchoTwice: len = %d, got = %v`,
				len(call1.Reply.(*EchoResponse).GetMsg()), got,
			)
		}

		// call2
		if call2.Error != nil {
			b.Fatalf(`EchoService.EchoTwice: %v`, call2.Error)
		}
		if call2.Reply.(*EchoResponse).GetMsg() != args2.GetMsg()+args2.GetMsg() {
			got := call2.Reply.(*EchoResponse).GetMsg()
			if len(got) > 8 {
				got = got[:8] + "..."
			}
			b.Fatalf(`EchoService.EchoTwice: len = %d, got = %v`,
				len(call2.Reply.(*EchoResponse).GetMsg()), got,
			)
		}

		// call3
		if call3.Error != nil {
			b.Fatalf(`EchoService.EchoTwice: %v`, call3.Error)
		}
		if call3.Reply.(*EchoResponse).GetMsg() != args3.GetMsg()+args3.GetMsg() {
			got := call3.Reply.(*EchoResponse).GetMsg()
			if len(got) > 8 {
				got = got[:8] + "..."
			}
			b.Fatalf(`EchoService.EchoTwice: len = %d, got = %v`,
				len(call3.Reply.(*EchoResponse).GetMsg()), got,
			)
		}
	}
}
