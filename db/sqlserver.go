package db

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/opentracing/opentracing-go"
	tags "github.com/opentracing/opentracing-go/ext"

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
	wrappers        []QueryWrapper
	contextWrappers []QueryContextWrapper
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

func (s *SqlServer) QueryContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	if len(s.contextWrappers) == 0 {
		return s.queryContext(ctx, dest, query, args...)
	}

	queryer := func(ctx context.Context, query string, args ...interface{}) (interface{}, error) {
		return nil, s.queryContext(ctx, dest, query, args...)
	}

	for _, r := range s.contextWrappers {
		queryer = r(ctx, queryer, query, args...)
	}

	_, err := queryer(ctx, query, args...)
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

func (s *SqlServer) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	if len(s.contextWrappers) == 0 {
		return s.DB.ExecContext(ctx, query, args...)
	}

	queryer := func(ctx context.Context, query string, args ...interface{}) (interface{}, error) {
		return s.DB.ExecContext(ctx, query, args...)
	}

	for _, r := range s.contextWrappers {
		queryer = r(ctx, queryer, query, args...)
	}

	resultItf, err := queryer(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	result := resultItf.(sql.Result)
	return result, err
}

func (s *SqlServer) AddQueryWrapper(wrapper QueryWrapper) {
	s.wrappers = append(s.wrappers, wrapper)
}

func (s *SqlServer) AddQueryContextWrapper(wrapper QueryContextWrapper) {
	s.contextWrappers = append(s.contextWrappers, wrapper)
}

func (s *SqlServer) GetContextWrappers() []QueryContextWrapper {
	return s.contextWrappers
}

func (s *SqlServer) GetWrappers() []QueryWrapper {
	return s.wrappers
}

func (s *SqlServer) query(dest interface{}, query string, args ...interface{}) error {
	t := reflect.TypeOf(dest)
	if reflectx.Deref(t).Kind() == reflect.Slice {
		return s.DB.Select(dest, query, args...)
	}
	return s.DB.Get(dest, query, args...)
}

func (s *SqlServer) queryContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	t := reflect.TypeOf(dest)
	if reflectx.Deref(t).Kind() == reflect.Slice {
		return s.DB.SelectContext(ctx, dest, query, args...)
	}
	return s.DB.GetContext(ctx, dest, query, args...)
}

type Queryer func(query string, args ...interface{}) (interface{}, error)

type QueryWrapper func(queryer Queryer, query string, args ...interface{}) Queryer

type ContextQueryer func(ctx context.Context, query string, args ...interface{}) (interface{}, error)

type QueryContextWrapper func(ctx context.Context, queryer ContextQueryer, query string, args ...interface{}) ContextQueryer

func SQLServerTracerWrapper(ctx context.Context, queryer ContextQueryer, query string, args ...interface{}) ContextQueryer {
	return func(ctx context.Context, query string, args ...interface{}) (interface{}, error) {
		tracer := opentracing.GlobalTracer()
		if ctx == nil {
			ctx = context.TODO()
		}
		span := opentracing.SpanFromContext(ctx)
		if span == nil {
			return queryer(ctx, query, args...)
		}
		span = tracer.StartSpan("mssql", opentracing.ChildOf(span.Context()))
		tags.DBInstance.Set(span, "")
		tags.DBStatement.Set(span, hackQueryBuilder(query))
		tags.DBType.Set(span, "mssql")
		tags.DBUser.Set(span, "")
		ctx = opentracing.ContextWithSpan(ctx, span)
		defer span.Finish()
		return queryer(ctx, query, args...)
	}
}

func hackQueryBuilder(query string) string {
	query = strings.Replace(query, "select", "SELECT", -1)
	query = strings.Replace(query, "from", "FROM", -1)
	r := regexp.MustCompile("(?s)SELECT (.*) FROM")
	query = r.ReplaceAllString(query, "SELECT ... FROM")
	return query
}
