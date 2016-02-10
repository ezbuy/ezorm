SHELL := /bin/bash
export GOPATH := $(shell pwd)
export PATH := ${PATH}:${GOPATH}/bin
export GOBIN := ${GOPATH}/bin

all:

init:
	go get gopkg.in/yaml.v2
	go get github.com/spf13/cobra/cobra
	go get -u github.com/jteeuwen/go-bindata/...

debugTpl:
	go-bindata -o tpl/bindata.go -ignore bindata.go -pkg tpl -debug tpl

buildTpl:
	go-bindata -o tpl/bindata.go -ignore bindata.go -pkg tpl tpl
