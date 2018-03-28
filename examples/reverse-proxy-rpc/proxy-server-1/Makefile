# Copyright 2018 <chaishushan{AT}gmail.com>. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

PROTO_FILES=$(sort $(wildcard ./*.proto))

generate: $(PROTO_FILES) Makefile
	protoc -I. --go_out=plugins=grpc:. ${PROTO_FILES}

clean:
	@echo "ok"
