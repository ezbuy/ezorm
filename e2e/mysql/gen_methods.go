// Package mysql is generated from e2e/mysql/sqls directory
// by github.com/ezbuy/ezorm/v2 , DO NOT EDIT!
package mysql

import (
	"context"
	sql_driver "database/sql"
	"fmt"
	"time"

	"github.com/ezbuy/ezorm/v2/db"
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

// GetUser is a raw query handler generated function for `e2e/mysql/sqls/get_user.sql`.
func (m *sqlMethods) GetUser(ctx context.Context, req *GetUserReq, opts ...RawQueryOptionHandler) ([]*GetUserResp, error) {

	rawQueryOption := &RawQueryOption{}

	for _, o := range opts {
		o(rawQueryOption)
	}

	query := _GetUserSQL

	rows, err := db.GetMysql(db.WithDB(rawQueryOption.db)).QueryContext(ctx, query, req.Params()...)
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

type GetUserInResp struct {
	UserId int32 `sql:"user_id"`
}

type GetUserInReq struct {
	Name []string `sql:"name"`
}

func (req *GetUserInReq) Params() []any {
	var params []any

	for _, v := range req.Name {
		params = append(params, v)
	}

	return params
}

func (req *GetUserInReq) QueryIn() []any {
	var qs []any

	qs = append(qs, sql.NewIn(len(req.Name)).String())
	return qs
}

const _GetUserInSQL = "SELECT `user_id` FROM `test_user` WHERE `name` IN %s"

// GetUserIn is a raw query handler generated function for `e2e/mysql/sqls/get_user_in.sql`.
func (m *sqlMethods) GetUserIn(ctx context.Context, req *GetUserInReq, opts ...RawQueryOptionHandler) ([]*GetUserInResp, error) {

	rawQueryOption := &RawQueryOption{}

	for _, o := range opts {
		o(rawQueryOption)
	}

	query := fmt.Sprintf(_GetUserInSQL, req.QueryIn()...)

	rows, err := db.GetMysql(db.WithDB(rawQueryOption.db)).QueryContext(ctx, query, req.Params()...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*GetUserInResp
	for rows.Next() {
		var o GetUserInResp
		err = rows.Scan(&o.UserId)
		if err != nil {
			return nil, err
		}
		results = append(results, &o)
	}
	return results, nil
}

type UserJoinBlogResp struct {
	BlogId int32 `sql:"blog_id"`
	UserId int32 `sql:"user_id"`
}

type UserJoinBlogReq struct {
	Name string `sql:"name"`
}

func (req *UserJoinBlogReq) Params() []any {
	var params []any

	params = append(params, req.Name)

	return params
}

const _UserJoinBlogSQL = "SELECT `u`.`user_id`,`b`.`blog_id` FROM `test_user` AS `u` JOIN `blog` AS `b` ON `u`.`user_id`=`b`.`user` WHERE `u`.`name`=?"

// UserJoinBlog is a raw query handler generated function for `e2e/mysql/sqls/user_join_blog.sql`.
func (m *sqlMethods) UserJoinBlog(ctx context.Context, req *UserJoinBlogReq, opts ...RawQueryOptionHandler) ([]*UserJoinBlogResp, error) {

	rawQueryOption := &RawQueryOption{}

	for _, o := range opts {
		o(rawQueryOption)
	}

	query := _UserJoinBlogSQL

	rows, err := db.GetMysql(db.WithDB(rawQueryOption.db)).QueryContext(ctx, query, req.Params()...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*UserJoinBlogResp
	for rows.Next() {
		var o UserJoinBlogResp
		err = rows.Scan(&o.BlogId, &o.UserId)
		if err != nil {
			return nil, err
		}
		results = append(results, &o)
	}
	return results, nil
}
