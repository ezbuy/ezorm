package db

import (
	"database/sql"
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
	wrappers []QueryWrapper
	*sqlx.DB
}

func (s *SqlServer) Query(dest interface{}, query string, args ...interface{}) error {
	if len(s.wrappers) == 0 {
		return s.query(dest, query, args...)
	}

	queryer := func(query string, args ...interface{}) (interface{}, error) {
		return nil, s.query(dest, query, args...)
	}

	for _, r := range s.wrappers {
		queryer = r(queryer, query, args...)
	}

	_, err := queryer(query, args...)
	return err
}

func (s *SqlServer) Exec(query string, args ...interface{}) (sql.Result, error) {
	if len(s.wrappers) == 0 {
		return s.DB.Exec(query, args...)
	}

	queryer := func(query string, args ...interface{}) (interface{}, error) {
		return s.DB.Exec(query, args...)
	}

	for _, r := range s.wrappers {
		queryer = r(queryer, query, args...)
	}

	resultItf, err := queryer(query, args...)
	if err != nil {
		return nil, err
	}
	result := resultItf.(sql.Result)
	return result, err
}

func (s *SqlServer) AddQueryWrapper(wrapper QueryWrapper) {
	s.wrappers = append(s.wrappers, wrapper)
}

func (s *SqlServer) query(dest interface{}, query string, args ...interface{}) error {
	t := reflect.TypeOf(dest)
	if reflectx.Deref(t).Kind() == reflect.Slice {
		return s.DB.Select(dest, query, args...)
	}

	return s.DB.Get(dest, query, args...)
}

type Queryer func(query string, args ...interface{}) (interface{}, error)

type QueryWrapper func(queryer Queryer, query string, args ...interface{}) Queryer
