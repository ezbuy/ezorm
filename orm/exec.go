package orm

import (
	"context"
	"database/sql"
)

type Execable interface {
	ExecContext(ctx context.Context, sql string, args ...interface{}) (sql.Result, error)
}

type Queryable interface {
	QueryContext(ctx context.Context, sql string, args ...interface{}) (*sql.Rows, error)
}

func ExecLastId(ctx context.Context, db Execable, sql string, args []interface{}) (int64, error) {
	r, err := db.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, err
	}
	return r.LastInsertId()
}

func ExecAffected(ctx context.Context, db Execable, sql string, args []interface{}) (int64, error) {
	r, err := db.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, err
	}
	return r.RowsAffected()
}

func ExecQuery(ctx context.Context, db Queryable, sql string, args []interface{}, fn func(rows *sql.Rows) error) error {
	rows, err := db.QueryContext(ctx, sql, args...)
	if err != nil {
		return err
	}
	var scanErr error
	for rows.Next() {
		scanErr = fn(rows)
		if scanErr != nil {
			break
		}
	}
	if err := rows.Close(); err != nil {
		return err
	}
	return scanErr
}
