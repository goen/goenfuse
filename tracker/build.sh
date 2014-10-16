#!/bin/sh

go build
~/Desktop/GOLANG/MYGOPROJECTS/bin/go-bindata -nomemcopy -pkg="main" -o="../bindata.go" tracker
