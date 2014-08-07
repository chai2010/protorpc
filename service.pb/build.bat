:: Copyright 2013 <chaishushan{AT}gmail.com>. All rights reserved.
:: Use of this source code is governed by a BSD-style
:: license that can be found in the LICENSE file.

setlocal

cd %~dp0

del *.pb.go

protoc --go_out=. arith.proto echo.proto
if not %errorlevel% == 0 (
	pause
)
