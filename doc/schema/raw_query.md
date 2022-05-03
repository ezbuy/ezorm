# Raw Query Support

> See more discussion and context in [#127](https://github.com/ezbuy/ezorm/issues/127) and [#139](https://github.com/ezbuy/ezorm/pull/139)

## Getting Start

### Table Schema

Table Schema follows the v1 mechanism and defined as the YAML file .

Checkout the schema doc for more details.

### Raw Query Declaration

To use the raw query mechanism, you first need to declare the raw query in the the `sqls` folder, like `example/mysql_people/sqls`.

And the declaration is nothing special , just follows the SQL syntax, but the compiler(generator) will check both the field name and the type between the declaration(SQL File) and the table schema(YAML).

That is , if we have a field named `name1` in the table schema, but the SQL file declares a field named `name2` and can not find `name2` in table schema , the compiler(generator) will throw an error and stop the process. Another case, if we have `WHERE name = 1` in our SQL declaration , and the table schema should be equal to type `int`, otherwise, the compiler(generator) will also throw an error and stop the process.

After that , things are almost done , `ezorm gen` will do the rest for you.

For example , lets write a simple SQL Query like:

> you can find more sql query [here](https://github.com/ezbuy/ezorm/tree/main/e2e/mysql/sqls)

```sql

SELECT
	name
FROM
	test_user
WHERE
	name = "me"
```

After run the `ezorm gen -i ./example/mysql_people/people.yaml -o ./example/mysql_people -p people --goPackage test` , the Go code helper will be auto-generated as the following:

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

## Pros and Cons

Our original idea is inspired by [sqlc](https://github.com/kyleconroy/sqlc) , but the query declaration is bad for the SQL highlighting by the SQL lint extension.

### Pros

* Native SQL Query declaration, which can be better highlighting and easy debugging.

### Cons

* The SQL AST parser is at the early stage , and can not support some SQL keywords.


## Support and Limitation

|Keyword|Support|Remark|
|--|--|--|
|`IN`| YES | also support the condition is a sub-query|
|`JOIN`| YES | All kinds of `JOIN` type|
|`SQL Function`| NO | not planned|
|`LIKE`| NO | planned |
