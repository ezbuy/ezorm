package mssql_people

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/ezbuy/ezorm/db"
)

func (o *_PeopleMgr) Save(people *People) (sql.Result, error) {
	var fieldNames []string
	var placeHolders []string
	var fieldValues []interface{}
	t := reflect.TypeOf(people).Elem()
	v := reflect.ValueOf(people).Elem()
	nf := t.NumField()
	for i := 0; i < nf; i++ {
		fieldName := t.Field(i).Name
		fieldNames = append(fieldNames, fieldName)
		placeHolders = append(placeHolders, "?")
		fieldValues = append(fieldValues, v.FieldByName(fieldName).Interface())
	}

	query := fmt.Sprintf("insert into dbo.[People] (%s) values (%s)",
		strings.Join(fieldNames, ","),
		strings.Join(placeHolders, ","))
	server := db.GetSqlServer()
	return server.Exec(query, fieldValues...)
}

func (o *_PeopleMgr) FindOne(where string, args ...interface{}) (*People, error) {
	query := getQuerysql(true, where)
	server := db.GetSqlServer()
	var people People
	err := server.Query(&people, query, args...)
	return &people, err
}

func (o *_PeopleMgr) Find(where string, args ...interface{}) (results []*People, err error) {
	query := getQuerysql(false, where)
	server := db.GetSqlServer()
	err = server.Query(&results, query, args...)
	return
}

func (o *_PeopleMgr) FindAll() (results []*People, err error) {
	return o.Find("")
}

func (o *_PeopleMgr) FindWithOffset(where string, offset int, limit int, args ...interface{}) (results []*People, err error) {
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
