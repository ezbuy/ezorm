package test

import "github.com/ezbuy/ezorm/db"

var (
	_store *db.RedisStore
)

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

func RedisSetUp(cf *RedisConfig) {
	store, err := db.NewRedisStore(cf.Host, cf.Port, cf.Password, cf.DB)
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

func redisGetObjectById(obj db.Object, id string) error {
	return _store.GetObjectById(obj, id)
}

func redisDelObject(obj db.Object) error {
	return _store.DelObject(obj)
}

func redisSMEMBERIds(key string) ([]string, error) {
	return _store.SMembersIds(key)
}

func redisSINTERIds(keys ...string) ([]string, error) {
	return _store.SInterIds(keys...)
}
