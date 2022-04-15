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

type GetMySQLOption struct {
	db *sql.DB
}

type GetMySQLOptionFunc func(*GetMySQLOption)

func WithDB(db *sql.DB) GetMySQLOptionFunc {
	return func(opt *GetMySQLOption) {
		opt.db = db
	}
}

func GetMysql(opts ...GetMySQLOptionFunc) *Mysql {
	getOption := &GetMySQLOption{}
	for _, opt := range opts {
		opt(getOption)
	}
	if getOption.db == nil {
		return getMysqlInstance()
	}
	s, err := NewMysql(nil, WithRawDB(getOption.db))
	if err != nil {
		panic(err)
	}
	return s
}

func MysqlQuery(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return getMysqlInstance().Query(query, args...)
}

func MysqlExec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return getMysqlInstance().ExecContext(ctx, query, args...)
}
