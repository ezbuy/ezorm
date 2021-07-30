package sqlm

import "database/sql"

type Execable interface {
	Exec(sql string, args ...interface{}) (sql.Result, error)
}

type Queryable interface {
	Query(sql string, args ...interface{}) (*sql.Rows, error)
}

func ExecAffected(ea Execable, sql string, args []interface{}) (int64, error) {
	r, err := ea.Exec(sql, args...)
	if err != nil {
		return 0, err
	}
	return r.RowsAffected()
}

func ExecLastId(ea Execable, sql string, args []interface{}) (int64, error) {
	r, err := ea.Exec(sql, args...)
	if err != nil {
		return 0, err
	}
	return r.LastInsertId()
}

func QueryOne(qa Queryable, sql string, args []interface{}, scan func(rows *sql.Rows) error) error {
	rows, err := qa.Query(sql, args...)
	if err != nil {
		return err
	}
	var scanErr error
	if rows.Next() {
		scanErr = scan(rows)
	}
	if err := rows.Close(); err != nil {
		return err
	}
	return scanErr
}

func QueryMany(qa Queryable, sql string, args []interface{}, scan func(rows *sql.Rows) error) error {
	rows, err := qa.Query(sql, args...)
	if err != nil {
		return err
	}
	var scanErr error
	for rows.Next() {
		scanErr = scan(rows)
		if scanErr != nil {
			break
		}
	}
	if err := rows.Close(); err != nil {
		return err
	}
	return scanErr
}
