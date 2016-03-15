package db

import (
	"database/sql"
	"fmt"
	"reflect"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
)

var (
	_config *SqlDbConfig
	_db     *SqlServer
)

func GetSqlServer() *SqlServer {
	return _db
}

func SetDBConfig(conf *SqlDbConfig) {
	_config = conf
	db, err := sqlx.Connect("mssql", _config.SqlConnStr)
	if err != nil {
		fmt.Printf("[db.GetSqlServer] open sql fail:%s", err.Error())
	}

	_db = &SqlServer{DB: db}
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
	return _db.Query(dest, query, args...)
}

func Exec(query string, args ...interface{}) (result sql.Result, err error) {
	return _db.Exec(query, args...)
}
