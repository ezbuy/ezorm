// Package mysqlr is generated from e2e/mysqlr/sqls directory
// by github.com/ezbuy/ezorm/v2 , DO NOT EDIT!
package mysqlr

import (
	"context"
	sql_driver "database/sql"
	"fmt"
	"time"

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
	Blog any
}

type BlogReq struct {
	Id int64 `sql:"id"`
}

func (req *BlogReq) Params() []any {
	var params []any

	params = append(params, req.Id)

	return params
}

const _BlogSQL = "SELECT SUM(`title`) AS `title_count` FROM `blogs` WHERE `id`=?"

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
		err = rows.Scan(&o.Blog)
		if err != nil {
			return nil, err
		}
		results = append(results, &o)
	}
	return results, nil
}

type TestResp struct {
	Iid           int64  `sql:"iid"`
	WarehouseId   int64  `sql:"warehouse_id"`
	SkuCode       string `sql:"sku_code"`
	Barcode       string `sql:"barcode"`
	QuantityTotal int64  `sql:"quantity_total"`
	Bid           int64  `sql:"bid"`
	AreaId        int64  `sql:"area_id"`
	Code          string `sql:"code"`
}

type TestReq struct {
	SkuCode string `sql:"sku_code"`
}

func (req *TestReq) Params() []any {
	var params []any

	params = append(params, req.SkuCode)

	return params
}

const _TestSQL = "SELECT `oper_inventory`.`id` AS `iid`,`oper_inventory`.`warehouse_id`,`oper_inventory`.`sku_code`,`oper_inventory`.`barcode`,`oper_inventory`.`quantity_total`,`oper_storage_bin`.`id` AS `bid`,`oper_storage_bin`.`area_id`,`oper_storage_bin`.`code` FROM `oper_inventory` AS `i` JOIN `oper_storage_bin` AS `b` ON `i`.`bin_id`=`b`.`id` WHERE `oper_inventory`.`sku_code`=?"

// Test is a raw query handler generated function for `e2e/mysqlr/sqls/test.sql`.
func (m *sqlMethods) Test(ctx context.Context, req *TestReq, opts ...RawQueryOptionHandler) ([]*TestResp, error) {

	rawQueryOption := &RawQueryOption{}

	for _, o := range opts {
		o(rawQueryOption)
	}

	query := _TestSQL

	rows, err := db.GetMysql(db.WithDB(rawQueryOption.db)).QueryContext(ctx, query, req.Params()...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*TestResp
	for rows.Next() {
		var o TestResp
		err = rows.Scan(&o.Iid, &o.WarehouseId, &o.SkuCode, &o.Barcode, &o.QuantityTotal, &o.Bid, &o.AreaId, &o.Code)
		if err != nil {
			return nil, err
		}
		results = append(results, &o)
	}
	return results, nil
}
