all:

build:
	go build -ldflags "-X main.Version=$(shell git rev-parse --short HEAD)" -o bin/ezorm

buildTpl:
	go-bindata -nometadata -o tpl/bindata.go -ignore bindata.go -pkg tpl tpl

gene2e:
	bin/ezorm gen -i ./e2e/mongo/blog/blog.yaml -o ./e2e/mongo/blog --goPackage blog
	bin/ezorm gen -i ./e2e/mysql -o ./e2e/mysql --goPackage mysql

test: build gene2e testmongo testmysql

testmongo:
	go test -v ./e2e/mongo/blog/...

testmysql:
	go test -v ./e2e/mysql/...

