all:

build:
	go build -ldflags "-X main.Version=$(shell git rev-parse --short HEAD)" -o bin/ezorm

gene2e:
	bin/ezorm gen -i ./e2e/mongo/blog/blog.yaml -o ./e2e/mongo/blog --goPackage blog
	bin/ezorm gen -i ./e2e/mysql -o ./e2e/mysql --goPackage mysql

test: build gene2e testmongo testmysql

.PHONY: testmongo
testmongo:
	go test -v ./e2e/mongo/blog/...

.PHONY: testmysql
testmysql:
	go test -v ./e2e/mysql/...

