package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"io/ioutil"
)

var _config *Config

func GetSqlServer() *sql.DB {
	sqlServer, err := sql.Open("mssql", _config.SqlConnStr)
	if err != nil {
		fmt.Printf("[db.GetSqlServer] open sql fail:%s", err.Error())
	}
	return sqlServer
}

func SetDBConfig(conf *Config) {
	_config = conf
}
