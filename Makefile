all:

build:
	go build -ldflags "-X main.Version=$(shell git rev-parse --short HEAD)" -o bin/ezorm

gene2e:
	bin/ezorm gen -i ./e2e/mongo/user/user.yaml -o  ./e2e/mongo/user --goPackage user
	bin/ezorm gen -i ./e2e/mysql -o ./e2e/mysql --goPackage mysql

test: build gene2e test-mysql test-mongo-go-driver

.PHONY: test-mongo-go-driver
test-mongo-go-driver:
	go test -v ./e2e/mongo/user/...

.PHONY: test-mysql
test-mysql:
	go test -v ./e2e/mysql/...

