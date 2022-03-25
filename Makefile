all:

buildTpl:
	go-bindata -nometadata -o tpl/bindata.go -ignore bindata.go -pkg tpl tpl

genexample:
	go install
	ezorm gen -i ./example/mssql_people/people_mssql.yaml -o ./example/mssql_people -p people --goPackage test
	ezorm gen -i ./example/blog/blog.yaml -o ./example/blog -p blog --goPackage test
	ezorm gen -i ./example/mysql_people/people.yaml -o ./example/mysql_people -p people --goPackage test
	ezorm gen -i ./example/redis_people/people.yaml -o ./example/redis_people -p people --goPackage test
	ezorm gen -i ./example/mysql_user/user.yaml -o ./example/mysql_user/ -p user --goPackage test

clean:
	rm ./example/mssql_people/gen_*.go
	rm ./example/blog/gen_*.go
	rm ./example/mysql_people/gen_*.go
	rm ./example/redis_people/gen_*.go

genmongo:
	rm example/blog/gen_*.go
	ezorm gen -i example/blog/blog.yaml -o example/blog -p blog --goPackage test

test: genexample testmssql testmongo testmysql

testmssql:
	go test -v ./example/mssql_people/...

testmongo:
	go test -v ./example/blog/...

testmysql:
	go test -v ./example/mysql_people/...

testredis:
	go test -v ./example/redis_people/...
