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

// MysqlFieldConfig uses fields to config mysql, it can be converted
// to DSN style.
type MysqlFieldConfig struct {
	Host            string
	Port            int
	UserName        string
	Password        string
	Database        string
	PoolSize        int
	ConnMaxLifeTime time.Duration
}

func (cfg *MysqlFieldConfig) Convert() *MysqlConfig {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		cfg.UserName,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database)
	return &MysqlConfig{
		DataSource:      dsn,
		PoolSize:        cfg.PoolSize,
		ConnMaxLifeTime: cfg.ConnMaxLifeTime,
	}
}

func NewMysql(cfg *MysqlConfig) (*Mysql, error) {
	if cfg == nil {
		cfg = new(MysqlConfig)
	}
	cfg.init()

	db, err := sql.Open("mysql", cfg.DataSource)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %v", err)
	}
	db.SetConnMaxLifetime(time.Hour)
	db.SetMaxIdleConns(cfg.PoolSize)
	db.SetMaxOpenConns(cfg.PoolSize)

	return &Mysql{
		cfg: cfg,
		DB:  db,
	}, nil
}
