// Copyright 2018 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"log"
	"os"
	"text/template"

	plugin "github.com/chai2010/protorpc/protoc-gen-plugin"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/golang/protobuf/protoc-gen-go/generator"
)

var flagPrefix = os.Getenv(ENV_PROTOC_GEN_PROTORPC_FLAG_PREFIX)

func main() {
	plugin.Main()
}

func init() {
	plugin.RegisterCodeGenerator(new(protorpcPlugin))
}

type protorpcPlugin struct{}

func (p *protorpcPlugin) Name() string        { return "protorpc-go" }
func (p *protorpcPlugin) FileNameExt() string { return ".pb.protorpc.go" }

func (p *protorpcPlugin) HeaderCode(g *generator.Generator, file *generator.FileDescriptor) string {
	const tmpl = `
{{- $G := .G -}}
{{- $File := .File -}}

// Code generated by protoc-gen-protorpc. DO NOT EDIT.
//
// plugin: https://github.com/chai2010/protorpc/tree/master/protoc-gen-plugin
// plugin: https://github.com/chai2010/protorpc/tree/master/protoc-gen-protorpc
//
// source: {{$File.GetName}}

package {{$File.PackageName}}

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/rpc"
	"time"

	"github.com/chai2010/protorpc"
	"github.com/golang/protobuf/proto"
)

var (
	_ = fmt.Sprint
	_ = io.Reader(nil)
	_ = log.Print
	_ = net.Addr(nil)
	_ = rpc.Call{}
	_ = time.Second

	_ = proto.String
	_ = protorpc.Dial
)
`
	var buf bytes.Buffer
	t := template.Must(template.New("").Parse(tmpl))
	err := t.Execute(&buf,
		struct {
			G      *generator.Generator
			File   *generator.FileDescriptor
			Prefix string
		}{
			G:      g,
			File:   file,
			Prefix: flagPrefix,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	return buf.String()
}

func (p *protorpcPlugin) ServiceCode(g *generator.Generator, file *generator.FileDescriptor, svc *descriptor.ServiceDescriptorProto) string {
	var code string
	code += p.genServiceInterface(g, file, svc)
	code += p.genServiceServer(g, file, svc)
	code += p.genServiceClient(g, file, svc)
	return code
}

func (p *protorpcPlugin) MessageCode(g *generator.Generator, file *generator.FileDescriptor, msg *descriptor.DescriptorProto) string {
	return ""
}

func (p *protorpcPlugin) genServiceInterface(
	g *generator.Generator,
	file *generator.FileDescriptor,
	svc *descriptor.ServiceDescriptorProto,
) string {
	const serviceInterfaceTmpl = `
type {{.Prefix}}{{.ServiceName}} interface {
	{{.CallMethodList}}
}
`
	const callMethodTmpl = `
{{.MethodName}}(in *{{.ArgsType}}, out *{{.ReplyType}}) error`

	// gen call method list
	var callMethodList string
	for _, m := range svc.Method {
		out := bytes.NewBuffer([]byte{})
		t := template.Must(template.New("").Parse(callMethodTmpl))
		t.Execute(out, &struct {
			Prefix      string
			ServiceName string
			MethodName  string
			ArgsType    string
			ReplyType   string
		}{
			Prefix:      flagPrefix,
			ServiceName: generator.CamelCase(svc.GetName()),
			MethodName:  generator.CamelCase(m.GetName()),
			ArgsType:    g.TypeName(g.ObjectNamed(m.GetInputType())),
			ReplyType:   g.TypeName(g.ObjectNamed(m.GetOutputType())),
		})
		callMethodList += out.String()

		g.RecordTypeUse(m.GetInputType())
		g.RecordTypeUse(m.GetOutputType())
	}

	// gen all interface code
	{
		out := bytes.NewBuffer([]byte{})
		t := template.Must(template.New("").Parse(serviceInterfaceTmpl))
		t.Execute(out, &struct {
			Prefix         string
			ServiceName    string
			CallMethodList string
		}{
			Prefix:         flagPrefix,
			ServiceName:    generator.CamelCase(svc.GetName()),
			CallMethodList: callMethodList,
		})

		return out.String()
	}
}

func (p *protorpcPlugin) genServiceServer(
	g *generator.Generator,
	file *generator.FileDescriptor,
	svc *descriptor.ServiceDescriptorProto,
) string {
	const serviceHelperFunTmpl = `
// {{.Prefix}}Accept{{.ServiceName}}Client accepts connections on the listener and serves requests
// for each incoming connection.  Accept blocks; the caller typically
// invokes it in a go statement.
func {{.Prefix}}Accept{{.ServiceName}}Client(lis net.Listener, x {{.Prefix}}{{.ServiceName}}) {
	srv := rpc.NewServer()
	if err := srv.RegisterName("{{.ServiceRegisterName}}", x); err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Fatalf("lis.Accept(): %v\n", err)
		}
		go srv.ServeCodec(protorpc.NewServerCodec(conn))
	}
}

// {{.Prefix}}Register{{.ServiceName}} publish the given {{.Prefix}}{{.ServiceName}} implementation on the server.
func {{.Prefix}}Register{{.ServiceName}}(srv *rpc.Server, x {{.Prefix}}{{.ServiceName}}) error {
	if err := srv.RegisterName("{{.ServiceRegisterName}}", x); err != nil {
		return err
	}
	return nil
}

// {{.Prefix}}New{{.ServiceName}}Server returns a new {{.Prefix}}{{.ServiceName}} Server.
func {{.Prefix}}New{{.ServiceName}}Server(x {{.Prefix}}{{.ServiceName}}) *rpc.Server {
	srv := rpc.NewServer()
	if err := srv.RegisterName("{{.ServiceRegisterName}}", x); err != nil {
		log.Fatal(err)
	}
	return srv
}

// {{.Prefix}}ListenAndServe{{.ServiceName}} listen announces on the local network address laddr
// and serves the given {{.ServiceName}} implementation.
func {{.Prefix}}ListenAndServe{{.ServiceName}}(network, addr string, x {{.Prefix}}{{.ServiceName}}) error {
	lis, err := net.Listen(network, addr)
	if err != nil {
		return err
	}
	defer lis.Close()

	srv := rpc.NewServer()
	if err := srv.RegisterName("{{.ServiceRegisterName}}", x); err != nil {
		return err
	}

	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Fatalf("lis.Accept(): %v\n", err)
		}
		go srv.ServeCodec(protorpc.NewServerCodec(conn))
	}
}

// {{.Prefix}}Serve{{.ServiceName}} serves the given {{.Prefix}}{{.ServiceName}} implementation.
func {{.Prefix}}Serve{{.ServiceName}}(conn io.ReadWriteCloser, x {{.Prefix}}{{.ServiceName}}) {
	srv := rpc.NewServer()
	if err := srv.RegisterName("{{.ServiceRegisterName}}", x); err != nil {
		log.Fatal(err)
	}
	srv.ServeCodec(protorpc.NewServerCodec(conn))
}
`
	{
		out := bytes.NewBuffer([]byte{})
		t := template.Must(template.New("").Parse(serviceHelperFunTmpl))
		t.Execute(out, &struct {
			Prefix              string
			PackageName         string
			ServiceName         string
			ServiceRegisterName string
		}{
			Prefix:      flagPrefix,
			PackageName: file.GetPackage(),
			ServiceName: generator.CamelCase(svc.GetName()),
			ServiceRegisterName: p.makeServiceRegisterName(
				file, file.GetPackage(), generator.CamelCase(svc.GetName()),
			),
		})

		return out.String()
	}
}

func (p *protorpcPlugin) genServiceClient(
	g *generator.Generator,
	file *generator.FileDescriptor,
	svc *descriptor.ServiceDescriptorProto,
) string {
	const clientHelperFuncTmpl = `
type {{.Prefix}}{{.ServiceName}}Client struct {
	*rpc.Client
}

// {{.Prefix}}New{{.ServiceName}}Client returns a {{.Prefix}}{{.ServiceName}} stub to handle
// requests to the set of {{.Prefix}}{{.ServiceName}} at the other end of the connection.
func {{.Prefix}}New{{.ServiceName}}Client(conn io.ReadWriteCloser) (*{{.Prefix}}{{.ServiceName}}Client) {
	c := rpc.NewClientWithCodec(protorpc.NewClientCodec(conn))
	return &{{.Prefix}}{{.ServiceName}}Client{c}
}

{{.MethodList}}

// {{.Prefix}}Dial{{.ServiceName}} connects to an {{.Prefix}}{{.ServiceName}} at the specified network address.
func {{.Prefix}}Dial{{.ServiceName}}(network, addr string) (*{{.Prefix}}{{.ServiceName}}Client, error) {
	c, err := protorpc.Dial(network, addr)
	if err != nil {
		return nil, err
	}
	return &{{.Prefix}}{{.ServiceName}}Client{c}, nil
}

// {{.Prefix}}Dial{{.ServiceName}}Timeout connects to an {{.Prefix}}{{.ServiceName}} at the specified network address.
func {{.Prefix}}Dial{{.ServiceName}}Timeout(network, addr string, timeout time.Duration) (*{{.Prefix}}{{.ServiceName}}Client, error) {
	c, err := protorpc.DialTimeout(network, addr, timeout)
	if err != nil {
		return nil, err
	}
	return &{{.Prefix}}{{.ServiceName}}Client{c}, nil
}
`
	const clientMethodTmpl = `
func (c *{{.Prefix}}{{.ServiceName}}Client) {{.MethodName}}(in *{{.ArgsType}}) (out *{{.ReplyType}}, err error) {
	if in == nil {
		in = new({{.ArgsType}})
	}

	type Validator interface {
		Validate() error
	}
	if x, ok := proto.Message(in).(Validator); ok {
		if err := x.Validate(); err != nil {
			return nil, err
		}
	}

	out = new({{.ReplyType}})
	if err = c.Call("{{.ServiceRegisterName}}.{{.MethodName}}", in, out); err != nil {
		return nil, err
	}

	if x, ok := proto.Message(out).(Validator); ok {
		if err := x.Validate(); err != nil {
			return out, err
		}
	}

	return out, nil
}

func (c *{{.Prefix}}{{.ServiceName}}Client) Async{{.MethodName}}(in *{{.ArgsType}}, out *{{.ReplyType}}, done chan *rpc.Call) *rpc.Call {
	if in == nil {
		in = new({{.ArgsType}})
	}
	return c.Go(
		"{{.ServiceRegisterName}}.{{.MethodName}}",
		in, out,
		done,
	)
}
`

	// gen client method list
	var methodList string
	for _, m := range svc.Method {
		out := bytes.NewBuffer([]byte{})
		t := template.Must(template.New("").Parse(clientMethodTmpl))
		t.Execute(out, &struct {
			Prefix              string
			ServiceName         string
			ServiceRegisterName string
			MethodName          string
			ArgsType            string
			ReplyType           string
		}{
			Prefix:      flagPrefix,
			ServiceName: generator.CamelCase(svc.GetName()),
			ServiceRegisterName: p.makeServiceRegisterName(
				file, file.GetPackage(), generator.CamelCase(svc.GetName()),
			),
			MethodName: generator.CamelCase(m.GetName()),
			ArgsType:   g.TypeName(g.ObjectNamed(m.GetInputType())),
			ReplyType:  g.TypeName(g.ObjectNamed(m.GetOutputType())),
		})
		methodList += out.String()
	}

	// gen all client code
	{
		out := bytes.NewBuffer([]byte{})
		t := template.Must(template.New("").Parse(clientHelperFuncTmpl))
		t.Execute(out, &struct {
			Prefix      string
			PackageName string
			ServiceName string
			MethodList  string
		}{
			Prefix:      flagPrefix,
			PackageName: file.GetPackage(),
			ServiceName: generator.CamelCase(svc.GetName()),
			MethodList:  methodList,
		})

		return out.String()
	}
}

func (p *protorpcPlugin) makeServiceRegisterName(
	file *generator.FileDescriptor,
	packageName, serviceName string,
) string {
	// return packageName + "." + serviceName
	return serviceName
}
