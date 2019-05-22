#!/bin/sh

#windows
go build -o sq.exe main.go access.go circularAlg.go  jwt.go redis.go 
set SHA1-PWD=GuOqAxqkaL7E2Hr1LUb8PjLX7dE=
set SQ-SECRET=123
#192.168.1.109
#./sq.exe 8080
#lx
#go build main.go access.go circularAlg.go  jwt.go redis.go -o sq