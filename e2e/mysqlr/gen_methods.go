// Package mysqlr is generated from e2e/mysqlr/sqls directory
// by github.com/ezbuy/ezorm/v2 , DO NOT EDIT!
package mysqlr

import (
	"context"
	"fmt"
	"strings"
	"time"

	sql_driver "database/sql"

	"github.com/ezbuy/ezorm/v2/pkg/db"
	"github.com/ezbuy/ezorm/v2/pkg/sql"
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
	Id         int64 `sql:"id"`
	TitleCount any   `sql:"title_count"`
	Status     int32 `sql:"status"`
}

type BlogReq struct {
	Id int64 `sql:"id"`

	Offset int32 `sql:"offset"`

	Limit int32 `sql:"limit"`
}

func (req *BlogReq) Params() []any {
	var params []any

	if req.Id != 0 {
		params = append(params, req.Id)
	}

	params = append(params, req.Offset)

	params = append(params, req.Limit)

	return params
}

func (req *BlogReq) Condition() string {
	var conditions []string
	if req.Id != 0 {
		conditions = append(conditions, "id = ?")
	}
	var query string
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	query += " LIMIT ?, ?"
	return query
}

const _BlogSQL = "SELECT `Id`,SUM(`title`) AS `title_count`,`status` FROM `blogs` %s"

// Blog is a raw query handler generated function for `e2e/mysqlr/sqls/blog.sql`.
func (m *sqlMethods) Blog(ctx context.Context, req *BlogReq, opts ...RawQueryOptionHandler) ([]*BlogResp, error) {

	rawQueryOption := &RawQueryOption{}

	for _, o := range opts {
		o(rawQueryOption)
	}

	query := fmt.Sprintf(_BlogSQL, req.Condition())

	rows, err := db.GetMysql(db.WithDB(rawQueryOption.db)).QueryContext(ctx, query, req.Params()...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*BlogResp
	for rows.Next() {
		var o BlogResp
		err = rows.Scan(&o.Id, &o.TitleCount, &o.Status)
		if err != nil {
			return nil, err
		}
		results = append(results, &o)
	}
	return results, nil
}
