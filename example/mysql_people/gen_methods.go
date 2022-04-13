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

func (*sqlMethods) GetUser(ctx context.Context, req *GetUserReq) ([]*GetUserResp, error) {
	rows, err := db.MysqlQuery(_GetUserSQL, req.Params()...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*GetUserResp
	for rows.Next() {
		var o GetUserResp
		err = rows.Scan()
		if err != nil {
			return nil, err
		}
		results = append(results, &o)
	}
	return results, nil
}
