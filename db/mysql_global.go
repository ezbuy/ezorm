package db

import (
	"context"
	"database/sql"
	"sync"
)

var (
	mysqlInstance *Mysql
	mysqlConnOnce sync.Once
)

func MysqlInitByField(cfg *MysqlFieldConfig) {
	MysqlInit(cfg.Convert())
}

func MysqlInit(cfg *MysqlConfig) {
	mysqlConnOnce.Do(func() {
		var err error
		mysqlInstance, err = NewMysql(cfg)
		if err != nil {
			panic("init mysql: " + err.Error())
		}
		err = mysqlInstance.DB.Ping()
		if err != nil {
			panic("ping mysql: " + err.Error())
		}
	})
}

func getMysqlInstance() *Mysql {
	if mysqlInstance == nil {
		panic("mysql no init, please call MysqlInit first.")
	}
	return mysqlInstance
}

func SetupRawDB(db *sql.DB) {
	var err error
	mysqlInstance, err = NewMysql(nil, WithRawDB(db))
	if err != nil {
		panic("init mysql: " + err.Error())
	}
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
