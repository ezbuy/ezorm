package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"io/ioutil"
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
