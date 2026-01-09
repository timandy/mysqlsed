@echo off

set CGO_ENABLED=0
set GOOS=linux
set GOARCH=amd64
go build -o ./out/mysqlsed_linux_amd64.exe .
