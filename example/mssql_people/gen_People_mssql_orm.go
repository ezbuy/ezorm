package mssql_people

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/ezbuy/ezorm/db"
)

func (m *_PeopleMgr) Save(obj *People) (sql.Result, error) {
	if obj.Id == 0 {
		return m.saveInsert(obj)
	}
	return m.saveUpdate(obj)
}

func (m *_PeopleMgr) saveInsert(obj *People) (sql.Result, error) {
	var fieldNames []string
	var placeHolders []string
	var fieldValues []interface{}
	t := reflect.TypeOf(obj).Elem()
	v := reflect.ValueOf(obj).Elem()
	nf := t.NumField()
	for i := 0; i < nf; i++ {
		fieldName := t.Field(i).Name
		if fieldName == idFieldName {
			continue
		}
		fieldNames = append(fieldNames, fieldName)
		placeHolders = append(placeHolders, "?")
		fieldValues = append(fieldValues, v.FieldByName(fieldName).Interface())
	}

	query := fmt.Sprintf("insert into dbo.[People] (%s) values (%s)",
		strings.Join(fieldNames, ","),
		strings.Join(placeHolders, ","))
	server := db.GetSqlServer()
	result, err := server.Exec(query, fieldValues...)
	if err != nil {
		return result, err
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return result, err
	}

	obj.Id = int32(lastInsertId)

	return result, err
}

func (m *_PeopleMgr) saveUpdate(obj *People) (sql.Result, error) {
	var fieldSets []string
	var fieldValues []interface{}
	t := reflect.TypeOf(obj).Elem()
	v := reflect.ValueOf(obj).Elem()
	nf := t.NumField()
	for i := 0; i < nf; i++ {
		fieldName := t.Field(i).Name
		if fieldName == idFieldName {
			continue
		}
		fieldSets = append(fieldSets, fieldName+"=?")
		fieldValues = append(fieldValues, v.FieldByName(fieldName).Interface())
	}

	idField, ok := t.FieldByName(idFieldName)
	if !ok {
		return nil, errors.New("no Id field")
	}

	idFieldStr := idField.Tag.Get("db")
	if idFieldStr != "" {
		idFieldStr = idFieldName
	}

	query := fmt.Sprintf("update dbo.[People] set %s where %s=%d",
		strings.Join(fieldSets, ","), idFieldStr, obj.Id)
	server := db.GetSqlServer()
	return server.Exec(query, fieldValues...)
}

func (m *_PeopleMgr) FindOne(where string, args ...interface{}) (*People, error) {
	query := getQuerysql(true, where)
	server := db.GetSqlServer()
	var obj People
	err := server.Query(&obj, query, args...)
	return &obj, err
}

func (m *_PeopleMgr) Find(where string, args ...interface{}) (results []*People, err error) {
	query := getQuerysql(false, where)
	server := db.GetSqlServer()
	err = server.Query(&results, query, args...)
	return
}

func (m *_PeopleMgr) FindAll() (results []*People, err error) {
	return m.Find("")
}

func (m *_PeopleMgr) FindWithOffset(where string, offset int, limit int, args ...interface{}) (results []*People, err error) {
	query := getQuerysql(false, where)

	if !strings.Contains(strings.ToLower(where), "ORDER BY") {
		where = " ORDER BY Name"
	}
	query = query + where + " OFFSET ? Rows FETCH NEXT ? Rows ONLY"
	args = append(args, offset)
	args = append(args, limit)

	server := db.GetSqlServer()
	err = server.Query(&results, query, args...)
	return
}

func getQuerysql(topOne bool, where string) string {
	query := `SELECT `
	if topOne {
		query = query + ` TOP 1 `
	}
	query = query + ` * FROM dbo.[People] WITH(NOLOCK) `

	if where != "" {
		if strings.Index(strings.Trim(where, " "), "WHERE") == -1 {
			where = " WHERE " + where
		}
		query = query + where
	}
	return query
}

func (m *_PeopleMgr) Del(where string, params ...interface{}) (sql.Result, error) {
	query := "delete from People"
	if where != "" {
		query = fmt.Sprintf("delete from People where " + where)
	}
	server := db.GetSqlServer()
	return server.Exec(query, params...)
}

// argument example:
// set:"a=?, b=?"
// where:"c=? and d=?"
// params:[]interface{}{"a", "b", "c", "d"}...
func (m *_PeopleMgr) Update(set, where string, params ...interface{}) (sql.Result, error) {
	query := fmt.Sprintf("update People set %s", set)
	if where != "" {
		query = fmt.Sprintf("update People set %s where %s", set, where)
	}
	server := db.GetSqlServer()
	return server.Exec(query, params...)
}
