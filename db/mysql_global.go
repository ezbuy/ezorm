package db

import (
	"database/sql"
	"sync"
)

var (
	mysqlCfg      MysqlConfig
	mysqlInstance *Mysql
	mysqlConnOnce sync.Once
)

func MysqlInit(cfg *MysqlConfig) {
	mysqlCfg = *cfg
}

func getMysqlInstance() *Mysql {
	var err error
	mysqlConnOnce.Do(func() {
		mysqlInstance, err = NewMysql(&mysqlCfg)
		if err != nil {
			panic(err)
		}
	})

	return mysqlInstance
}

func MysqlQuery(query string, args ...interface{}) (*sql.Rows, error) {
	return getMysqlInstance().Query(query, args...)
}

func MysqlExec(query string, args ...interface{}) (sql.Result, error) {
	return getMysqlInstance().Exec(query, args...)
}
