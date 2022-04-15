all:

build:
	go build -ldflags "-X main.Version=$(shell git rev-parse --short HEAD)" -o bin/ezorm

buildTpl:
	go-bindata -nometadata -o tpl/bindata.go -ignore bindata.go -pkg tpl tpl

genexample:
	bin/ezorm gen -i ./example/mssql_people/people_mssql.yaml -o ./example/mssql_people -p people --goPackage test
	bin/ezorm gen -i ./example/blog/blog.yaml -o ./example/blog -p blog --goPackage test
	bin/ezorm gen -i ./e2e/mysql -o ./e2e/mysql --goPackage mysql
	bin/ezorm gen -i ./example/redis_people/people.yaml -o ./example/redis_people -p people --goPackage test

genmongo:
	rm example/blog/gen_*.go
	bin/ezorm gen -i example/blog/blog.yaml -o example/blog -p blog --goPackage test

test: build genexample testmssql testmongo testmysql

testmssql:
	go test -v ./example/mssql_people/...

testmongo:
	go test -v ./example/blog/...

testmysql:
	go test -v ./e2e/mysql/...

testredis:
	go test -v ./example/redis_people/...
