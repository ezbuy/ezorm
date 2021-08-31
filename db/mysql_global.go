package db

import (
	"context"
	"database/sql"
)

var mysqlInstance *Mysql

func MysqlInitByField(cfg *MysqlFieldConfig) {
	MysqlInit(cfg.Convert())
}

func MysqlInit(cfg *MysqlConfig) {
	var err error
	mysqlInstance, err = NewMysql(cfg)
	if err != nil {
		panic(err)
	}
}

func getMysqlInstance() *Mysql {
	if mysqlInstance == nil {
		panic("mysql no init, please call MysqlInit first.")
	}
	return mysqlInstance
}

func GetMysql() *Mysql {
	return getMysqlInstance()
}

func MysqlQuery(query string, args ...interface{}) (*sql.Rows, error) {
	return getMysqlInstance().Query(query, args...)
}

func MysqlQueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return getMysqlInstance().QueryContext(ctx, query, args...)
}

func MysqlExec(query string, args ...interface{}) (sql.Result, error) {
	return getMysqlInstance().Exec(query, args...)
}
