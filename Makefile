SHELL := /bin/bash
export GOPATH := $(shell pwd)/../../../..
export PATH := ${PATH}:${GOPATH}/bin
export GOBIN := ${GOPATH}/bin

all:

init:
	go get gopkg.in/yaml.v2
	go get github.com/spf13/cobra/cobra
	go get -u github.com/jteeuwen/go-bindata/...
	go get github.com/denisenkom/go-mssqldb
	go get github.com/jmoiron/sqlx
	go get gopkg.in/mgo.v2

debugTpl:
	go-bindata -o tpl/bindata.go -ignore bindata.go -pkg tpl -debug tpl

buildTpl:
	go-bindata -o tpl/bindata.go -ignore bindata.go -pkg tpl tpl

test:
	go install github.com/ezbuy/ezorm
	ezorm gen -i example/example.yaml -o example
	go test github.com/ezbuy/ezorm/example

testmssql:buildTpl
	go install
	ezorm gen -i ./example/mssql_people/people_mssql.yaml -o ./example/mssql_people -p people

