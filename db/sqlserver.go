package db

import (
	"fmt"
	"reflect"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
)

func GetSqlServer(dataSourceName string) *SqlServer {
	db, err := sqlx.Connect("mssql", dataSourceName)
	if err != nil {
		fmt.Printf("[db.GetSqlServer] open sql fail:%s", err.Error())
	}

	return &SqlServer{DB: db}
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
