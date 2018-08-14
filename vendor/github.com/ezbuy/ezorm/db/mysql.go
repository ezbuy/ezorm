package db

import (
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

func NewMysql(cfg *MysqlConfig) (*Mysql, error) {
	if cfg == nil {
		cfg = new(MysqlConfig)
	}
	cfg.init()

	db, err := sql.Open("mysql", cfg.DataSource)
	if err != nil {
		return nil, fmt.Errorf("sql.Open:", err)
	}
	db.SetConnMaxLifetime(time.Hour)
	db.SetMaxIdleConns(cfg.PoolSize)
	db.SetMaxOpenConns(cfg.PoolSize)

	return &Mysql{
		cfg: cfg,
		DB:  db,
	}, nil
}
