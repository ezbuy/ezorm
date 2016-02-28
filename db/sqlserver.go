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
	sqlServer, err := sql.Open("mssql", _config.DB.SqlServerConn)
	if err != nil {
		fmt.Println("DB error", err)
	}
	return sqlServer
}

func init() {
	_config = new(Config)
	setConfig()
}

func setConfig() {
	confPath := "../conf/default.json"
	cfgbuf, err := ioutil.ReadFile(confPath)
	if err != nil {
		fmt.Println("Read config file failed: %s %s", confPath, err)
	}
	fmt.Println(confPath)
	err = json.Unmarshal(cfgbuf, _config)
	if err != nil {
		fmt.Println("Unmarshal failed:", err.Error())
	}
}
