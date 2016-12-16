package db

import (
	"errors"

	"github.com/garyburd/redigo/redis"
)

func (r *RedisStore) do(cmd string, args ...interface{}) (interface{}, error) {
	conn := r.pool.Get()
	defer conn.Close()

	res, err := conn.Do(cmd, args...)
	if err == redis.ErrNil {
		return nil, ErrNil
	}
	return res, err
}

func (r *RedisStore) GET(key string) (interface{}, error) {
	return r.do("GET", key)
}

func (r *RedisStore) SET(key string, value interface{}) error {
	_, err := r.do("SET", key, value)
	return err
}

func (r *RedisStore) DEL(keys ...string) (int64, error) {
	ks := make([]interface{}, len(keys))
	for i, key := range keys {
		ks[i] = key
	}
	return r.Int64(r.do("DEL", ks...))
}

func (r *RedisStore) HGET(key string, field string) (interface{}, error) {
	return r.do("HGET", key, field)
}

func (r *RedisStore) HLEN(key string) (int64, error) {
	return r.Int64(r.do("HLEN", key))
}

func (r *RedisStore) HSET(key string, field string, val interface{}) error {
	_, err := r.do("HSET", key, field, val)
	return err
}

func (r *RedisStore) HDEL(key string, fields ...string) (int64, error) {
	ks := make([]interface{}, len(fields)+1)
	ks[0] = key
	for i, key := range fields {
		ks[i+1] = key
	}
	return r.Int64(r.do("HDEL", ks...))
}

func (r *RedisStore) HMGET(key string, fields ...string) (interface{}, error) {
	if len(fields) == 0 {
		return nil, ErrNil
	}
	args := make([]interface{}, len(fields)+1)
	args[0] = key
	for i, field := range fields {
		args[i+1] = field
	}
	return r.do("HMGET", args...)

}

func (r *RedisStore) HMSET(key string, kvs ...interface{}) error {
	if len(kvs) == 0 {
		return nil
	}
	if len(kvs)%2 != 0 {
		return ErrWrongArgsNum
	}
	args := make([]interface{}, len(kvs)+1)
	args[0] = key
	for i := 0; i < len(kvs); i += 2 {
		if _, ok := kvs[i].(string); !ok {
			return errors.New("field must be string")
		}
		args[i+1] = kvs[i]
		args[i+2] = kvs[i+1]
	}
	_, err := r.do("HMSET", args...)
	return err
}

func (r *RedisStore) SADD(key string, members ...interface{}) (int64, error) {
	args := make([]interface{}, len(members)+1)
	args[0] = key
	for i, m := range members {
		args[i+1] = m
	}
	return r.Int64(r.do("SADD", args...))
}

func (r *RedisStore) SMEMBERS(key string) (interface{}, error) {
	return r.do("SMEMBERS", key)
}

func (r *RedisStore) SINTER(keys ...interface{}) (interface{}, error) {
	if len(keys) == 0 {
		return nil, errors.New("absent keys")
	}
	return r.do("SINTER", keys...)
}

func (r *RedisStore) ZADD(key string, kvs ...interface{}) (int64, error) {
	if len(kvs) == 0 {
		return 0, nil
	}
	if len(kvs)%2 != 0 {
		return 0, errors.New("args num error")
	}
	args := make([]interface{}, len(kvs)+1)
	args[0] = key
	for i := 0; i < len(kvs); i += 2 {
		args[i+1] = kvs[i]
		args[i+2] = kvs[i+1]
	}
	return r.Int64(r.do("ZADD", args...))
}

func (r *RedisStore) ZCOUNT(key string) (int64, error) {
	return r.Int64(r.do("ZCOUNT", key))
}

func (r *RedisStore) ZREM(key string, members ...string) (int64, error) {
	args := make([]interface{}, len(members)+1)
	args[0] = key
	for i, m := range members {
		args[i+1] = m
	}
	return r.Int64(r.do("ZREM", args...))
}

func (r *RedisStore) ZRANGE(key string, min, max int64, withScores bool) (interface{}, error) {
	if withScores {
		return r.do("ZRANGE", key, min, max, "WITHSCORES")
	} else {
		return r.do("ZRANGE", key, min, max)
	}
}

func (r *RedisStore) GEOADD(key string, longitude float64, latitude float64, value interface{}) (int64, error) {
	return r.Int64(r.do("GEOADD", key, longitude, latitude, value))
}
