{{define "conf.mysql"}}package {{.GoPackage}}
// Package {{.GoPackage}} defines the mysql config
import (
    "errors"
	"sync"
	"time"

	"github.com/ezbuy/redis-orm/orm"
	"github.com/ezbuy/wrapper/database"
)

var (
	mysqlStore *orm.DBStore
	mysqlCfg   MySQLConfig
	mysqlDSN   string

	mysqlDSNs       = map[string]string{}
	mysqlMultiStore = map[string]*orm.DBStore{}

    mysqlMultiOnce  sync.Once
	mysqlOnce  sync.Once
)

type MySQLConfig struct {
	Host            string
	Port            int
	UserName        string
	Password        string
	Database        string
	PoolSize        int
	ConnMaxLifeTime time.Duration
}

func MySQLSetup(cf *MySQLConfig) {
	mysqlCfg = *cf
}

func MySQLDSNSetup(dsn string) {
	mysqlDSN = dsn
}

func MySQLMultiDSNSetup(key, dsn string) {
    if _, ok := mysqlDSNs[key]; ok {
		panic(errors.New("ezorm: setup: "key + " exists"))
	}

	mysqlDSNs[key] = dsn
}

func MySQLInstance(key string) *orm.DBStore {
    mysqlMultiOnce.Do(func() {
		for key, dsn := range mysqlDSNs {
			s, err := orm.NewDBDSNStore("mysql", dsn)
			if err != nil {
				panic(err)
			}

			mysqlMultiStore[key] = s
		}
	})

	s, ok := mysqlMultiStore[key]
	if !ok {
		panic(errors.New("ezorm: getMySQLInstance: "key + " not found"))
	}
	return s
}

func MySQL() *orm.DBStore {
	var err error
	mysqlOnce.Do(func() {
	    if mysqlDSN != "" {
	        mysqlStore, err = orm.NewDBDSNStore("mysql", mysqlDSN)
            if err != nil {
        	    panic(err)
            }
            return
        }
	    mysqlStore, err = orm.NewDBStore("mysql",
            mysqlCfg.Host,
            mysqlCfg.Port,
            mysqlCfg.Database,
            mysqlCfg.UserName,
            mysqlCfg.Password,
        )
        if err != nil {
          panic(errors.New("ezorm: "+ err.Error()))
        }
        mysqlStore.SetConnMaxLifetime(time.Hour)
        if mysqlCfg.ConnMaxLifeTime > 0 {
            mysqlStore.SetConnMaxLifetime(mysqlCfg.ConnMaxLifeTime)
        }
        mysqlStore.SetMaxIdleConns(mysqlCfg.PoolSize)
        mysqlStore.SetMaxOpenConns(mysqlCfg.PoolSize)
	    mysqlStore.AddWrappers(
		    database.NewMySQLTracerWrapper(),
		)
	})
	return mysqlStore
}

{{end}}
