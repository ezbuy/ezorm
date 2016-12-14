package db

import (
	"errors"
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
)

var ErrNil = errors.New("nil return")
var ErrWrongType = errors.New("wrong type")
var ErrWrongArgsNum = errors.New("args num error")

const redisMaxIdleConn = 64
const redisMaxActive = 128

type RedisStore struct {
	pool        *redis.Pool
	host        string
	port        int
	db          int
	withTimeout bool
}

func newRedisStore(host string, port int, db int, withTimeout bool) (*RedisStore, error) {
	f := func() (redis.Conn, error) {
		var c redis.Conn
		var err error
		if withTimeout {
			c, err = redis.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), time.Second*10, time.Second*3, time.Second*3)
		} else {
			c, err = redis.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
		}
		if err != nil {
			return nil, err
		}
		if _, err := c.Do("SELECT", db); err != nil {
			return nil, err
		}
		return c, err
	}
	pool := redis.NewPool(f, redisMaxIdleConn)
	pool.MaxActive = redisMaxActive
	pool.Wait = true

	store := &RedisStore{pool: pool, host: host, port: port, db: db, withTimeout: withTimeout}
	return store, nil
}

func NewRedisStore(host string, port int, db int) (*RedisStore, error) {
	return newRedisStore(host, port, db, true)
}

func NewRedisStoreWithoutTimeout(host string, port int, db int) (*RedisStore, error) {
	return newRedisStore(host, port, db, false)
}

func (r *RedisStore) SetMaxIdle(maxIdle int) {
	r.pool.MaxIdle = maxIdle
}

func (r *RedisStore) SetMaxActive(maxActive int) {
	r.pool.MaxActive = maxActive
}

func (r *RedisStore) GetPool() *redis.Pool {
	return r.pool
}
