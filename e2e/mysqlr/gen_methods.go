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
	Id     int64 `sql:"id"`
	Status int32 `sql:"status"`
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

const _BlogSQL = "SELECT `id`,`status` FROM `blogs` %s"

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
		err = rows.Scan(&o.Id, &o.Status)
		if err != nil {
			return nil, err
		}
		results = append(results, &o)
	}
	return results, nil
}

type BlogAggrResp struct {
	Count any `sql:"count"`
}

type BlogAggrReq struct {
	Id int64 `sql:"id"`
}

func (req *BlogAggrReq) Params() []any {
	var params []any

	if req.Id != 0 {
		params = append(params, req.Id)
	}

	return params
}

func (req *BlogAggrReq) Condition() string {
	var conditions []string
	if req.Id != 0 {
		conditions = append(conditions, "id = ?")
	}
	var query string
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	return query
}

const _BlogAggrSQL = "SELECT COUNT(`id`) AS `count` FROM `blogs` %s"

// BlogAggr is a raw query handler generated function for `e2e/mysqlr/sqls/blog_aggr.sql`.
func (m *sqlMethods) BlogAggr(ctx context.Context, req *BlogAggrReq, opts ...RawQueryOptionHandler) ([]*BlogAggrResp, error) {

	rawQueryOption := &RawQueryOption{}

	for _, o := range opts {
		o(rawQueryOption)
	}

	query := fmt.Sprintf(_BlogAggrSQL, req.Condition())

	rows, err := db.GetMysql(db.WithDB(rawQueryOption.db)).QueryContext(ctx, query, req.Params()...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*BlogAggrResp
	for rows.Next() {
		var o BlogAggrResp
		err = rows.Scan(&o.Count)
		if err != nil {
			return nil, err
		}
		results = append(results, &o)
	}
	return results, nil
}

type BlogFuncResp struct {
	UTitle   any `sql:"u_title"`
	LenTitle any `sql:"len_title"`
}

type BlogFuncReq struct {
	Id int64 `sql:"id"`
}

func (req *BlogFuncReq) Params() []any {
	var params []any

	if req.Id != 0 {
		params = append(params, req.Id)
	}

	return params
}

func (req *BlogFuncReq) Condition() string {
	var conditions []string
	if req.Id != 0 {
		conditions = append(conditions, "id = ?")
	}
	var query string
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	return query
}

const _BlogFuncSQL = "SELECT UPPER(`title`) AS `u_title`,LENGTH(`title`) AS `len_title` FROM `blogs` %s"

// BlogFunc is a raw query handler generated function for `e2e/mysqlr/sqls/blog_func.sql`.
func (m *sqlMethods) BlogFunc(ctx context.Context, req *BlogFuncReq, opts ...RawQueryOptionHandler) ([]*BlogFuncResp, error) {

	rawQueryOption := &RawQueryOption{}

	for _, o := range opts {
		o(rawQueryOption)
	}

	query := fmt.Sprintf(_BlogFuncSQL, req.Condition())

	rows, err := db.GetMysql(db.WithDB(rawQueryOption.db)).QueryContext(ctx, query, req.Params()...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*BlogFuncResp
	for rows.Next() {
		var o BlogFuncResp
		err = rows.Scan(&o.UTitle, &o.LenTitle)
		if err != nil {
			return nil, err
		}
		results = append(results, &o)
	}
	return results, nil
}
