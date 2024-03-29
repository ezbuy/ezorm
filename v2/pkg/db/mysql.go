package db

import (
	"bytes"
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Mysql struct {
	cfg  *MysqlConfig
	init sync.Once
	*sql.DB
}

type MysqlConfig struct {
	DataSource      string
	PoolSize        int
	ConnMaxLifeTime time.Duration
}

func (cfg *MysqlConfig) init() {
	if cfg.PoolSize == 0 {
		cfg.PoolSize = 2
	}
	if cfg.ConnMaxLifeTime == 0 {
		cfg.ConnMaxLifeTime = time.Hour
	}
}

// MysqlFieldConfig uses fields to config mysql, it can be converted
// to DSN style.
type MysqlFieldConfig struct {
	Addr     string
	UserName string
	Password string
	Database string

	PoolSize        int
	ConnMaxLifeTime time.Duration

	Options map[string]string
}

func (cfg *MysqlFieldConfig) Convert() *MysqlConfig {
	var userDSN string
	if cfg.Password == "" {
		userDSN = cfg.UserName
	} else if cfg.UserName != "" {
		userDSN = fmt.Sprintf("%s:%s", cfg.UserName, cfg.Password)
	}
	if userDSN != "" {
		userDSN += "@"
	}
	var buf bytes.Buffer
	for key, val := range cfg.Options {
		param := fmt.Sprintf("&%s=%s", key, val)
		buf.WriteString(param)
	}
	dsn := fmt.Sprintf("%stcp(%s)/%s?charset=utf8mb4%s",
		userDSN,
		cfg.Addr,
		cfg.Database,
		buf.String())
	mysqlCfg := &MysqlConfig{
		DataSource:      dsn,
		PoolSize:        cfg.PoolSize,
		ConnMaxLifeTime: cfg.ConnMaxLifeTime,
	}
	mysqlCfg.init()
	return mysqlCfg
}

type NewMySQLOption struct {
	db *sql.DB
}

type NewMySQLOptionFunc func(*NewMySQLOption)

func WithRawDB(db *sql.DB) NewMySQLOptionFunc {
	return func(opt *NewMySQLOption) {
		opt.db = db
	}
}

func NewMysql(cfg *MysqlConfig, opts ...NewMySQLOptionFunc) (*Mysql, error) {
	newOptions := &NewMySQLOption{}
	for _, opt := range opts {
		opt(newOptions)
	}
	if cfg == nil {
		cfg = new(MysqlConfig)
	}
	cfg.init()
	db := newOptions.db
	if db == nil {
		var err error
		db, err = sql.Open("mysql", cfg.DataSource)
		if err != nil {
			return nil, fmt.Errorf("sql.Open: %v", err)
		}
		db.SetConnMaxLifetime(time.Hour)
		db.SetMaxIdleConns(cfg.PoolSize)
		db.SetMaxOpenConns(cfg.PoolSize)
	}
	return &Mysql{
		cfg: cfg,
		DB:  db,
	}, nil
}
