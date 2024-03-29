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
	Id     int64 `sql:"id"`
	Offset int32 `sql:"offset"`
	Limit  int32 `sql:"limit"`
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
		conditions = append(conditions, "`b`.`id` = ?")
	}
	var query string
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	query += " LIMIT ?, ?"

	return query
}

const _BlogSQL = "SELECT `b`.`id`,`b`.`status` FROM `blogs` AS `b` %s"

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
		conditions = append(conditions, "`id` = ?")
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
	LenTitle any `sql:"len_title"`
	UTitle   any `sql:"u_title"`
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
		conditions = append(conditions, "`id` = ?")
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

type BlogJoinOrderByResp struct {
	Bid int64 `sql:"bid"`
	Uid int64 `sql:"uid"`
}

type BlogJoinOrderByReq struct {
	Id int64 `sql:"id"`
}

func (req *BlogJoinOrderByReq) Params() []any {
	var params []any

	if req.Id != 0 {
		params = append(params, req.Id)
	}

	return params
}

func (req *BlogJoinOrderByReq) Condition() string {
	var conditions []string
	if req.Id != 0 {
		conditions = append(conditions, "`u`.`id` = ?")
	}
	var query string
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	query += " ORDER BY `u`.`id` DESC"

	return query
}

const _BlogJoinOrderBySQL = "SELECT `u`.`id` AS `uid`,`b`.`id` AS `bid` FROM `users` AS `u` JOIN `blogs` AS `b` ON `u`.`id`=`b`.`user_id` %s"

// BlogJoinOrderBy is a raw query handler generated function for `e2e/mysqlr/sqls/blog_join_order_by.sql`.
func (m *sqlMethods) BlogJoinOrderBy(ctx context.Context, req *BlogJoinOrderByReq, opts ...RawQueryOptionHandler) ([]*BlogJoinOrderByResp, error) {

	rawQueryOption := &RawQueryOption{}

	for _, o := range opts {
		o(rawQueryOption)
	}

	query := fmt.Sprintf(_BlogJoinOrderBySQL, req.Condition())

	rows, err := db.GetMysql(db.WithDB(rawQueryOption.db)).QueryContext(ctx, query, req.Params()...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*BlogJoinOrderByResp
	for rows.Next() {
		var o BlogJoinOrderByResp
		err = rows.Scan(&o.Uid, &o.Bid)
		if err != nil {
			return nil, err
		}
		results = append(results, &o)
	}
	return results, nil
}

type BlogLikeResp struct {
	Id     int64 `sql:"id"`
	Status int32 `sql:"status"`
}

type BlogLikeReq struct {
	Title string `sql:"title"`
}

func (req *BlogLikeReq) Params() []any {
	var params []any

	if req.Title != "" {
		params = append(params, req.Title)
	}

	return params
}

func (req *BlogLikeReq) Condition() string {
	var conditions []string
	if req.Title != "" {
		conditions = append(conditions, "`b`.`title` LIKE ?")
	}
	var query string
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	return query
}

const _BlogLikeSQL = "SELECT `b`.`id`,`b`.`status` FROM `blogs` AS `b` %s"

// BlogLike is a raw query handler generated function for `e2e/mysqlr/sqls/blog_like.sql`.
func (m *sqlMethods) BlogLike(ctx context.Context, req *BlogLikeReq, opts ...RawQueryOptionHandler) ([]*BlogLikeResp, error) {

	rawQueryOption := &RawQueryOption{}

	for _, o := range opts {
		o(rawQueryOption)
	}

	query := fmt.Sprintf(_BlogLikeSQL, req.Condition())

	rows, err := db.GetMysql(db.WithDB(rawQueryOption.db)).QueryContext(ctx, query, req.Params()...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*BlogLikeResp
	for rows.Next() {
		var o BlogLikeResp
		err = rows.Scan(&o.Id, &o.Status)
		if err != nil {
			return nil, err
		}
		results = append(results, &o)
	}
	return results, nil
}

type BlogOrderByResp struct {
	Status int32 `sql:"status"`
	UserId int32 `sql:"user_id"`
}

type BlogOrderByReq struct {
	UserId int64 `sql:"user_id"`
	Offset int32 `sql:"offset"`
	Limit  int32 `sql:"limit"`
}

func (req *BlogOrderByReq) Params() []any {
	var params []any

	if req.UserId != 0 {
		params = append(params, req.UserId)
	}

	params = append(params, req.Offset)

	params = append(params, req.Limit)

	return params
}

func (req *BlogOrderByReq) Condition() string {
	var conditions []string
	if req.UserId != 0 {
		conditions = append(conditions, "`b`.`user_id` = ?")
	}
	var query string
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	query += " ORDER BY `b`.`id` DESC"
	query += " LIMIT ?, ?"

	return query
}

const _BlogOrderBySQL = "SELECT `b`.`user_id`,`b`.`status` FROM `blogs` AS `b` %s"

// BlogOrderBy is a raw query handler generated function for `e2e/mysqlr/sqls/blog_order_by.sql`.
func (m *sqlMethods) BlogOrderBy(ctx context.Context, req *BlogOrderByReq, opts ...RawQueryOptionHandler) ([]*BlogOrderByResp, error) {

	rawQueryOption := &RawQueryOption{}

	for _, o := range opts {
		o(rawQueryOption)
	}

	query := fmt.Sprintf(_BlogOrderBySQL, req.Condition())

	rows, err := db.GetMysql(db.WithDB(rawQueryOption.db)).QueryContext(ctx, query, req.Params()...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*BlogOrderByResp
	for rows.Next() {
		var o BlogOrderByResp
		err = rows.Scan(&o.UserId, &o.Status)
		if err != nil {
			return nil, err
		}
		results = append(results, &o)
	}
	return results, nil
}
