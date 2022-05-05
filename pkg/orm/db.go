package orm

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/ezbuy/wrapper/database"
	_ "github.com/go-sql-driver/mysql" // register driver for mysql
)

type DB interface {
	Query(ctx context.Context, sql string, args ...interface{}) (*sql.Rows, error)
	Exec(ctx context.Context, sql string, args ...interface{}) (sql.Result, error)
	SetError(err error)
}

type DBStore struct {
	*sql.DB
	debug    bool
	slowlog  time.Duration
	wrappers []database.Wrapper
}

func NewDBStore(driver, host string, port int, database, username, password string) (*DBStore, error) {
	var dsn string
	switch strings.ToLower(driver) {
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&autocommit=true&parseTime=True",
			username,
			password,
			host,
			port,
			database)
	case "mssql":
		dsn = fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s",
			host, username, password, port, database)
	default:
		return nil, fmt.Errorf("unsupport db driver: %s", driver)
	}
	return NewDBDSNStore(driver, dsn)
}

func NewDBStoreWithRawDB(db *sql.DB) *DBStore {
	wps := []database.Wrapper{
		// insert common wrappers here...
	}
	return &DBStore{db, false, time.Duration(0), wps}
}

func NewDBDSNStore(driver, dsn string) (*DBStore, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	wps := []database.Wrapper{
		// insert common wrappers here...
	}
	return &DBStore{db, false, time.Duration(0), wps}, nil
}

func NewDBStoreCharset(driver, host string, port int, databaseName, username, password, charset string) (*DBStore, error) {
	var dsn string
	switch strings.ToLower(driver) {
	case "mysql":
		if charset == "" {
			charset = "utf8"
		}
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&autocommit=true&parseTime=True",
			username,
			password,
			host,
			port,
			databaseName,
			charset)
	default:
		return nil, fmt.Errorf("do not support db driver: %s", driver)
	}
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	wps := []database.Wrapper{
		// insert common wrappers here...
	}
	return &DBStore{db, false, time.Duration(0), wps}, nil
}

func (store *DBStore) Debug(b bool) {
	store.debug = b
}

func (store *DBStore) SlowLog(duration time.Duration) {
	store.slowlog = duration
}

func (store *DBStore) SetError(err error) {}

func (store *DBStore) Close() error {
	if err := store.DB.Close(); err != nil {
		return err
	}
	store.DB = nil
	return nil
}

func (store *DBStore) AddWrappers(wp ...database.Wrapper) {
	store.wrappers = append(store.wrappers, wp...)
}

func (store *DBStore) Query(ctx context.Context, query string,
	args ...interface{}) (*sql.Rows, error) {
	fn := func(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
		return store.DB.QueryContext(ctx, query, args...)
	}
	for _, wp := range store.wrappers {
		fn = wp.WrapQueryContext(fn, query, args...)
	}
	return fn(ctx, query, args...)
}

func (store *DBStore) Exec(ctx context.Context, query string,
	args ...interface{}) (sql.Result, error) {
	fn := func(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
		return store.DB.ExecContext(ctx, query, args...)
	}
	for _, wp := range store.wrappers {
		fn = wp.WrapExecContext(fn, query, args...)
	}
	return fn(ctx, query, args...)
}

type DBTx struct {
	tx           *sql.Tx
	debug        bool
	slowlog      time.Duration
	err          error
	rowsAffected int64
	wrappers     []database.Wrapper
}

func (store *DBStore) BeginTx() (*DBTx, error) {
	tx, err := store.Begin()
	if err != nil {
		return nil, err
	}

	return &DBTx{
		tx:       tx,
		debug:    store.debug,
		slowlog:  store.slowlog,
		wrappers: store.wrappers,
	}, nil
}

func (tx *DBTx) Close() error {
	if tx.err != nil {
		return tx.tx.Rollback()
	}
	return tx.tx.Commit()
}

func (tx *DBTx) IsRollback() bool {
	return tx.err != nil
}

func (tx *DBTx) Query(ctx context.Context, query string,
	args ...interface{}) (*sql.Rows, error) {
	fn := func(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
		return tx.tx.QueryContext(ctx, query, args...)
	}
	for _, wp := range tx.wrappers {
		fn = wp.WrapQueryContext(fn, query, args...)
	}
	return fn(ctx, query, args...)
}

func (tx *DBTx) Exec(ctx context.Context, query string,
	args ...interface{}) (sql.Result, error) {
	fn := func(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
		return tx.tx.ExecContext(ctx, query, args...)
	}
	for _, wp := range tx.wrappers {
		fn = wp.WrapExecContext(fn, query, args...)
	}
	return fn(ctx, query, args...)
}

func (tx *DBTx) SetError(err error) {
	tx.err = err
}

func TransactFunc(db *DBStore, txFunc func(*DBTx) error) (err error) {
	tx, err := db.BeginTx()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.SetError(fmt.Errorf("panic: %v", p))
			tx.Close()
			panic(p)
		} else if err != nil {
			tx.SetError(err)
			tx.Close()
		} else {
			err = tx.Close()
		}
	}()

	err = txFunc(tx)
	return err
}

type Transactor interface {
	Transact(tx *DBTx) error
}

func Transact(db *DBStore, t Transactor) error {
	return TransactFunc(db, t.Transact)
}
