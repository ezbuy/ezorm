// Package test is generated from example/mysql_people/sqls directory
// by github.com/ezbuy/ezorm/v2 , DO NOT EDIT!
package test

import (
	"context"
	"time"

	"github.com/ezbuy/ezorm/v2/db"
)

var (
	_ time.Time
	_ context.Context
)

type sqlMethods struct{}

var SQL = &sqlMethods{}

type GetUserResp struct {
	Name string `sql:"name"`
}

type GetUserReq struct {
	Name string `sql:"name"`
}

func (req *GetUserReq) Params() []any {
	return []any{
		req.Name}
}

const _GetUserSQL = "SELECT `name` FROM `test_user` WHERE `name`=?"

// GetUser is a raw query handler generated function for `example/mysql_people/sqls/get_user.sql`.
func (*sqlMethods) GetUser(ctx context.Context, req *GetUserReq) ([]*GetUserResp, error) {
	rows, err := db.MysqlQuery(_GetUserSQL, req.Params()...)
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
