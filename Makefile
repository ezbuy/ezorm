all:

debugTpl:
	go-bindata -o tpl/bindata.go -ignore bindata.go -pkg tpl -debug tpl

buildTpl:
	go-bindata -o tpl/bindata.go -ignore bindata.go -pkg tpl tpl

test: testmssql testmongo testmysql

testmssql:
	go install
	ezorm gen -i ./example/mssql_people/people_mssql.yaml -o ./example/mssql_people -p people
	go test -v ./example/mssql_people/...

testmongo:
	go install
	ezorm gen -i ./example/blog/blog.yaml -o ./example/blog -p blog
	go test -v ./example/blog/...

testmysql:
	go install
	ezorm gen -i ./example/mysql_people/people.yaml -o ./example/mysql_people -p people
	go test -v ./example/mysql_people/...
