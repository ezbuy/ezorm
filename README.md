# ezorm.v2

[![parser](https://github.com/ezbuy/ezorm/workflows/parser/badge.svg)](https://github.com/ezbuy/ezorm/actions/workflows/parser.yml)
[![e2e](https://github.com/ezbuy/ezorm/workflows/e2e/badge.svg)](https://github.com/ezbuy/ezorm/actions/workflows/e2e.yml)
[![CodeQL](https://github.com/ezbuy/ezorm/workflows/CodeQL/badge.svg)](https://github.com/ezbuy/ezorm/actions/workflows/codeql.yml)


## Why another ORM for Go ?

With many years Go development experience in ezbuy , we find that `define db schema first` and share this schema within the project members or DBAs is really an advanced idea , like [Protobuf](https://developers.google.com/protocol-buffers) for API-oriented development.

And the idea seems not alone within the Go community , projects like [ent-go](https://entgo.io/) prove that there is a way to make a best-practice for Go ORM(like) Programming.

ezorm was built with this key idea in mind , but we describe the `database schema` as `YAML` file (or raw `SQL` file) , and the builtin generator(also the compiler) can generate some safe and most-used database operating methods for us , which lets business developers to focus on the business logic , rather than the bored CRUD .

## Support Database

* MySQL
	* driver: `db: mysql`
	* driver: `db: mysqlr`
* MongoDB
* ~~Redis(deprecated since v2)~~
* ~~SQL Server(deprecated since v2)~~
* [3rd customized plugin](./doc/plugin.md)

## Schema

> See full list of schema in our [doc](doc/schema/) .

### YAML Schema

Schema is defined with YAML file like:

```yaml
Blog:
  db: mongo
  fields:
    - Title: string
    - Hits: int32
    - Slug: string
      flags: [unique]
    - Body: string
    - User: int32
    - CreateDate: datetime
      flags: [sort]
    - IsPublished: bool
      flags: [index]
  indexes: [[User, IsPublished]]
```

> Id field will be automatically included for mongo/mysql.

### SQL Schema

SQL Schema is introduced in v2 , and tries to help with raw query , like **table JOIN** , which can not be handled properly by YAML schema.

Inspired by [sqlc](https://github.com/kyleconroy/sqlc)'s great AST analysis , we can extract the AST of SQL like:

```SQL
SELECT
  name
FROM test_user
WHERE name = "me";
```

and generates the following Go Code :

```go
type GetUserResp struct {
	Name string `sql:"name"`
}

type GetUserReq struct {
	Name string `sql:"name"`
}

func (req *GetUserReq) Params() []any {
	var params []any

	params = append(params, req.Name)

	return params
}

const _GetUserSQL = "SELECT `name` FROM `test_user` WHERE `name`=?"

// GetUser is a raw query handler generated function for `example/mysql_people/sqls/get_user.sql`.
func (*sqlMethods) GetUser(ctx context.Context, req *GetUserReq) ([]*GetUserResp, error) {

	query := _GetUserSQL

	rows, err := db.MysqlQuery(query, req.Params()...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*GetUserResp
	for rows.Next() {
		var o GetUserResp
		err = rows.Scan(&o.Name)
		if err != nil {
			return nil, err
		}
		results = append(results, &o)
	}
	return results, nil
}
```

## Usage

> ezorm requires **Go 1.18** or later.

```shell
	$ go install github.com/ezbuy/ezorm/v2
	$ ezorm gen -i blog.yaml -o .
```

To generate codes, for model like `Blog`, a blog manager will be generated, supporting ActiveRecord like:

```go
p := blog.BlogMgr.NewBlog()
p.Title = "I like ezorm"
p.Slug = "ezorm"
p.Save()

p, err := blog.BlogMgr.FindBySlug("ezorm")
if err != nil {
  // handle error
}
fmt.Println("%v", p)
page.PageMgr.RemoveByID(p.Id())

_, err = blog.BlogMgr.FindBySlug("ezorm")
if err == nil {
  // handle error
}
  ```

> use `ezorm -h` for more help

