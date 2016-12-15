package test

import (
	"github.com/ezbuy/ezorm/db"
)

var (
	_store *db.RedisStore
)

type RedisConfig struct {
	Host string
	Port int
	DB   int
}

func RedisSetUp(cf *RedisConfig) {
	store, err := db.NewRedisStore(cf.Host, cf.Port, cf.DB)
	if err != nil {
		panic(err)
	}
	_store = store
}

func redisSetObject(obj db.Object) error {
	return _store.SetObject(obj)
}

func redisGetObject(obj db.Object) error {
	return _store.GetObject(obj)
}

func redisDelObject(obj db.Object) error {
	return _store.DelObject(obj)
}
