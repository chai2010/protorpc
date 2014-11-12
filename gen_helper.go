// Copyright 2014 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	flagRevert = flag.Bool("revert", false, "revert all changes")
)

var (
	targetVersion = GetGoProtobufTipVersion()
	targeFilename = `goprotobuf-` + targetVersion[:12] + `.tar.gz`
	targetURL     = `https://goprotobuf.googlecode.com/archive/` + targetVersion + `.tar.gz`
)

var convertMap = [][2]string{
	[2]string{
		`"code.google.com/p/goprotobuf/proto"`,
		`"github.com/chai2010/protorpc/proto"`,
	},
	[2]string{
		`"code.google.com/p/goprotobuf/protoc-gen-go/descriptor"`,
		`"github.com/chai2010/protorpc/protoc-gen-go/descriptor"`,
	},
	[2]string{
		`"code.google.com/p/goprotobuf/protoc-gen-go/generator"`,
		`"github.com/chai2010/protorpc/protoc-gen-go/generator"`,
	},
	[2]string{
		`"code.google.com/p/goprotobuf/protoc-gen-go/plugin"`,
		`"github.com/chai2010/protorpc/protoc-gen-go/plugin"`,
	},
}

func main() {
	flag.Parse()

	// try download goprotobuf
	if _, err := os.Stat(targeFilename); err != nil {
		downloadFile(targetURL, targeFilename)
	}

	// unpack proto and protoc-gen-go
	unpackSourceCode(targeFilename)

	// fix import path
	fixAllImportPath("proto")
	fixAllImportPath("protoc-gen-go")

	// save proto version
	saveProtoVersion()

	// Done
	fmt.Println("Done")
}

func downloadFile(targetURL, targeFilename string) {
	var (
		f    *os.File
		resp *http.Response
		err  error
	)

	if f, err = os.Create(targeFilename); err != nil {
		log.Fatalf("downloadFile: failed to create %s: %s", targeFilename, err)
	}
	defer f.Close()

	if resp, err = http.Get(targetURL); err != nil {
		log.Fatalf("downloadFile: failed to download %s: %s", targetURL, err)
	}
	defer resp.Body.Close()

	if _, err = io.Copy(f, resp.Body); err != nil {
		log.Fatalf("downloadFile: failed to write %s: %s", targeFilename, err)
	}
}

func unpackSourceCode(filename string) {
	baseName := filename[:len(filename)-len(".tar.gz")]

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("unpackSourceCode: ioutil.ReadFile filed, err = %v", err)
	}

	gzReader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		log.Fatalf("unpackSourceCode: gzip.NewReader filed, err = %v", err)
	}
	defer gzReader.Close()

	trReader := tar.NewReader(gzReader)
	for {
		header, err := trReader.Next()
		if err != nil {
			if err != io.EOF {
				log.Fatalf("unpackSourceCode: trReader.Next filed, err = %v", err)
			}
			break
		}

		// skip other files
		if header.FileInfo().IsDir() {
			continue
		}

		// proto
		if strings.HasPrefix(header.Name, baseName+"/proto/") {
			name := header.Name[len(baseName+"/"):]
			os.MkdirAll(path.Dir(name), 0666)
			fw, err := os.Create(name)
			if err != nil {
				log.Fatalf("unpackSourceCode: os.Create filed, err = %v", err)
			}
			defer fw.Close()

			_, err = io.Copy(fw, trReader)
			if err != nil {
				log.Fatalf("unpackSourceCode: io.Copy filed, err = %v", err)
			}
		}

		// protoc-gen-go
		if strings.HasPrefix(header.Name, baseName+"/protoc-gen-go/") {
			name := header.Name[len(baseName+"/"):]
			os.MkdirAll(path.Dir(name), 0666)
			fw, err := os.Create(name)
			if err != nil {
				log.Fatalf("unpackSourceCode: os.Create filed, err = %v", err)
			}
			defer fw.Close()

			_, err = io.Copy(fw, trReader)
			if err != nil {
				log.Fatalf("unpackSourceCode: io.Copy filed, err = %v", err)
			}
		}
	}
}

func fixAllImportPath(dir string) {
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal("filepath.Walk: ", err)
			return err
		}
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, "gen_helper.go") {
			return nil
		}
		if strings.HasSuffix(path, ".go") {
			fixImportPath(path)
		}
		return nil
	})
}

func fixImportPath(filename string) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal("ioutil.ReadFile: ", err)
	}

	for _, v := range convertMap {
		oldPath, newPath := v[0], v[1]
		if !*flagRevert {
			data = bytes.Replace(data, []byte(oldPath), []byte(newPath), -1)
		} else {
			data = bytes.Replace(data, []byte(newPath), []byte(oldPath), -1)
		}
	}

	if err = ioutil.WriteFile(filename, data, 0666); err != nil {
		log.Fatal("ioutil.WriteFile: ", err)
	}

	if !*flagRevert {
		fmt.Printf("convert %s ok\n", filename)
	} else {
		fmt.Printf("revert %s ok\n", filename)
	}
}

func GetGoProtobufTipVersion() string {
	resp, err := http.Get("https://code.google.com/p/goprotobuf/source/browse/")
	if err != nil {
		log.Fatalf("GetGoProtobufTipTarGzUrl: http.Get failed, %v", err)
	}
	defer resp.Body.Close()

	var buf bytes.Buffer
	if _, err = io.Copy(&buf, resp.Body); err != nil {
		log.Fatalf("GetGoProtobufTipTarGzUrl: io.Copy failed, %v", err)
	}

	// href="//goprotobuf.googlecode.com/archive/26b4b5c8a6ec0e0bc5a8d2807c9f42f7d1bb291a.tar.gz" rel="nofollow">tar.gz</a>
	re := regexp.MustCompile(`href\=\"//goprotobuf.googlecode.com/archive/[0-9a-z]+\.tar\.gz\"`)
	url := re.FindString(string(buf.Bytes()))
	if url == "" {
		log.Fatalf("GetGoProtobufTipTarGzUrl: re.FindString failed")
	}
	return url[len(`href="//goprotobuf.googlecode.com/archive/`) : len(url)-len(`.tar.gz"`)]
}

func saveProtoVersion() {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "goprotobuf: %s\n", targetURL)
	if err := ioutil.WriteFile("goprotobuf-version.txt", buf.Bytes(), 0666); err != nil {
		log.Fatal("ioutil.WriteFile: ", err)
	}
}