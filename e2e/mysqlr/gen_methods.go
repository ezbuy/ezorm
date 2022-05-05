// Package mysqlr is generated from e2e/mysqlr/sqls directory
// by github.com/ezbuy/ezorm/v2 , DO NOT EDIT!
package mysqlr

import (
	"context"
	sql_driver "database/sql"
	"fmt"
	"github.com/ezbuy/ezorm/v2/db"
	"github.com/ezbuy/ezorm/v2/pkg/sql"
	"time"
)

var (
	_ time.Time
	_ context.Context
	_ sql.InBuilder
	_ fmt.Stringer
)

var rawQuery = &sqlMethods{}

type sqlMethods struct{}

type RawQueryOption struct {
	db *sql_driver.DB
}

type RawQueryOptionHandler func(*RawQueryOption)

func GetRawQuery() *sqlMethods {
	return rawQuery
}

func WithDB(db *sql_driver.DB) RawQueryOptionHandler {
	return func(o *RawQueryOption) {
		o.db = db
	}
}

type BlogResp struct {
	Title string `sql:"title"`
}

type BlogReq struct {
	Id int64 `sql:"id"`
}

func (req *BlogReq) Params() []any {
	var params []any

	params = append(params, req.Id)

	return params
}

const _BlogSQL = "SELECT `title` FROM `blogs` WHERE `id`=?"

// Blog is a raw query handler generated function for `e2e/mysqlr/sqls/blog.sql`.
func (m *sqlMethods) Blog(ctx context.Context, req *BlogReq, opts ...RawQueryOptionHandler) ([]*BlogResp, error) {

	rawQueryOption := &RawQueryOption{}

	for _, o := range opts {
		o(rawQueryOption)
	}

	query := _BlogSQL

	rows, err := db.GetMysql(db.WithDB(rawQueryOption.db)).QueryContext(ctx, query, req.Params()...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*BlogResp
	for rows.Next() {
		var o BlogResp
		err = rows.Scan(&o.Title)
		if err != nil {
			return nil, err
		}
		results = append(results, &o)
	}
	return results, nil
}