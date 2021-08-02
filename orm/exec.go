package orm

import "database/sql"

type Execable interface {
	Exec(sql string, args ...interface{}) (sql.Result, error)
}

type Queryable interface {
	Query(sql string, args ...interface{}) (*sql.Rows, error)
}

func ExecLastId(db Execable, sql string, args []interface{}) (int64, error) {
	r, err := db.Exec(sql, args...)
	if err != nil {
		return 0, err
	}
	return r.LastInsertId()
}

func ExecAffected(db Execable, sql string, args []interface{}) (int64, error) {
	r, err := db.Exec(sql, args...)
	if err != nil {
		return 0, err
	}
	return r.RowsAffected()
}

func ExecQuery(db Queryable, sql string, args []interface{}, fn func(rows *sql.Rows) error) error {
	rows, err := db.Query(sql, args...)
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
