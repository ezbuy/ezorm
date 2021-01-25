all:

.PHONY: buildTpl
buildTpl:
	go get -u github.com/jteeuwen/go-bindata/go-bindata@master
	go-bindata -nometadata -o tpl/bindata.go -ignore bindata.go -pkg tpl tpl

.PHONY: install
install:
	go install

.PHONY: genexample
genexample:
	go install
	ezorm gen -i ./example/mssql_people/people_mssql.yaml -o ./example/mssql_people -p people --goPackage test
	ezorm gen -i ./example/blog/blog.yaml -o ./example/blog -p blog --goPackage test
	ezorm gen -i ./example/user/user.yaml -o ./example/user -p user --goPackage test
	ezorm gen -i ./example/mysql_people/people.yaml -o ./example/mysql_people -p people --goPackage test
	ezorm gen -i ./example/redis_people/people.yaml -o ./example/redis_people -p people --goPackage test

.PHONY: clean
clean:
	rm ./example/mssql_people/gen_*.go
	rm ./example/blog/gen_*.go
	rm ./example/mysql_people/gen_*.go
	rm ./example/redis_people/gen_*.go

.PHONY: genmongo
genmongo:
	rm example/blog/gen_*.go
	ezorm gen -i example/blog/blog.yaml -o example/blog -p blog --goPackage test

.PHONY: genmongodriver
genmongodriver:
	rm -f example/user/gen_*.go
	ezorm gen -i example/user/user.yaml -o example/user -p user --goPackage user

.PHONY: test
test: genexample testmssql testmongo testmysql

.PHONY: testmssql
testmssql:
	go test -v ./example/mssql_people/...

.PHONY: testmongo
testmongo:
	go test -v ./example/blog/...

.PHONY: testmysql
testmysql:
	go test -v ./example/mysql_people/...

.PHONY: testredis
testredis:
	go test -v ./example/redis_people/...
