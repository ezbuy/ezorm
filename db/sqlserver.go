package db

import (
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
)

var _config *SqlDbConfig

func GetSqlServer() *sql.DB {
	sqlServer, err := sql.Open("mssql", _config.SqlConnStr)
	if err != nil {
		fmt.Printf("[db.GetSqlServer] open sql fail:%s", err.Error())
	}
	return sqlServer
}

func SetDBConfig(conf *SqlDbConfig) {
	_config = conf
}
