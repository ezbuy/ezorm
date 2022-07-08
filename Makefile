all:

build:
	go build -ldflags "-X main.Commit=$(shell git rev-parse --short HEAD)" -o bin/ezorm

build-plugin:
	go build -o bin/ezorm-gen-hello-driver ./e2e/plugins/hello-driver
	mv bin/ezorm-gen-hello-driver /go/bin

cleane2e:
	rm -rf e2e/mongo/user/gen_*.go
	rm -rf e2e/mongo/nested/gen_*.go
	rm -rf e2e/mysql/gen_*.go
	rm -rf e2e/mysqlr/gen*.go

regene2e: cleane2e gene2e

gene2e:
	bin/ezorm gen -i ./e2e/mongo/user/user.yaml -o  ./e2e/mongo/user --goPackage user --namespace mongo_e2e
	bin/ezorm gen -i ./e2e/mongo/nested/nested.yaml -o  ./e2e/mongo/nested --goPackage nested --namespace mongo_e2e
	bin/ezorm gen -i ./e2e/mysql -o ./e2e/mysql --goPackage mysql --namespace mysql_e2e
	bin/ezorm gen -i ./e2e/mysqlr -o ./e2e/mysqlr --goPackage mysqlr --namespace mysqlr_e2e

gen-plugin-e2e:
	bin/ezorm gen -i ./e2e/plugins/hello-driver -o ./e2e/plugins/hello-driver --goPackage hello-driver --plugin hello-driver

test: build build-plugin gene2e test-mysql test-mysqlr test-mongo-go-driver

.PHONY: test-mongo-go-driver
test-mongo-go-driver:
	go test -v ./e2e/mongo/user/...

.PHONY: test-mysql
test-mysql:
	go test -v ./e2e/mysql/...


.PHONY: test-mysqlr
test-mysqlr:
	go test -v ./e2e/mysqlr/...

.PHONY: test-plugin
test-plugin:
	go test -v ./e2e/plugins/hello-driver/...
