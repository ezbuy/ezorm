package db

import (
	"database/sql"
	"fmt"
	"reflect"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
)

var _config *SqlDbConfig

func GetSqlServer() *SqlServer {
	db, err := sqlx.Connect("mssql", _config.SqlConnStr)
	if err != nil {
		fmt.Printf("[db.GetSqlServer] open sql fail:%s", err.Error())
	}
	return &SqlServer{DB: db}
}

func SetDBConfig(conf *SqlDbConfig) {
	_config = conf
}

type SqlServer struct {
	*sqlx.DB
}

func (s *SqlServer) Query(dest interface{}, query string, args ...interface{}) error {
	t := reflect.TypeOf(dest)
	if reflectx.Deref(t).Kind() == reflect.Slice {
		return s.DB.Select(dest, query, args...)
	}

	return s.DB.Get(dest, query, args...)
}

func Query(dest interface{}, query string, args ...interface{}) error {
	return GetSqlServer().Query(dest, query, args...)
}

func Exec(query string, args ...interface{}) (sql.Result, error) {
	return GetSqlServer().Exec(query, args...)
}
