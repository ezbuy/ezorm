all:

build:
	go build -ldflags "-X main.Version=$(shell git rev-parse --short HEAD)" -o bin/ezorm

buildTpl:
	go-bindata -nometadata -o tpl/bindata.go -ignore bindata.go -pkg tpl tpl

genexample:
	bin/ezorm gen -i ./example/mssql_people/people_mssql.yaml -o ./example/mssql_people -p people --goPackage test
	bin/ezorm gen -i ./example/blog/blog.yaml -o ./example/blog -p blog --goPackage test
	bin/ezorm gen -i ./example/mysql_people/people.yaml -o ./example/mysql_people -p people --goPackage test
	bin/ezorm gen -i ./example/redis_people/people.yaml -o ./example/redis_people -p people --goPackage test
	bin/ezorm gen -i ./example/mysql_user/user.yaml -o ./example/mysql_user/ -p user --goPackage test

genmongo:
	rm example/blog/gen_*.go
	bin/ezorm gen -i example/blog/blog.yaml -o example/blog -p blog --goPackage test

test: build genexample testmssql testmongo testmysql

testmssql:
	go test -v ./example/mssql_people/...

testmongo:
	go test -v ./example/blog/...

testmysql:
	go test -v ./example/mysql_people/...

testredis:
	go test -v ./example/redis_people/...
