package db

import (
	"github.com/garyburd/redigo/redis"
)

func (r *RedisStore) Int(reply interface{}, err error) (int, error) {
	if err != nil {
		return 0, err
	}
	res, err := redis.Int(reply, err)
	if err == redis.ErrNil {
		return 0, ErrNil
	}
	return res, err
}

func (r *RedisStore) Int64(reply interface{}, err error) (int64, error) {
	if err != nil {
		return 0, err
	}
	res, err := redis.Int64(reply, err)
	if err == redis.ErrNil {
		return 0, ErrNil
	}
	return res, err
}

func (r *RedisStore) Uint64(reply interface{}, err error) (uint64, error) {
	if err != nil {
		return 0, err
	}
	res, err := redis.Uint64(reply, err)
	if err == redis.ErrNil {
		return 0, ErrNil
	}
	return res, err
}

func (r *RedisStore) Float64(reply interface{}, err error) (float64, error) {
	if err != nil {
		return 0, err
	}
	res, err := redis.Float64(reply, err)
	if err == redis.ErrNil {
		return 0, ErrNil
	}
	return res, err
}

func (r *RedisStore) Bool(reply interface{}, err error) (bool, error) {
	if err != nil {
		return false, err
	}
	res, err := redis.Bool(reply, err)
	if err == redis.ErrNil {
		return false, ErrNil
	}
	return res, err
}

func (r *RedisStore) Bytes(reply interface{}, err error) ([]byte, error) {
	if err != nil {
		return nil, err
	}
	res, err := redis.Bytes(reply, err)
	if err == redis.ErrNil {
		return nil, ErrNil
	}
	return res, err
}

func (r *RedisStore) String(reply interface{}, err error) (string, error) {
	if err != nil {
		return "", err
	}
	res, err := redis.String(reply, err)
	if err == redis.ErrNil {
		return "", ErrNil
	}
	return res, err
}

func (r *RedisStore) Strings(reply interface{}, err error) ([]string, error) {
	if err != nil {
		return nil, err
	}
	res, err := redis.Strings(reply, err)
	if err == redis.ErrNil {
		return nil, ErrNil
	}
	return res, err
}

func (r *RedisStore) Values(reply interface{}, err error) ([]interface{}, error) {
	if err != nil {
		return nil, err
	}
	res, err := redis.Values(reply, err)
	if err == redis.ErrNil {
		return nil, ErrNil
	}
	return res, err
}

func (r *RedisStore) Ints(reply interface{}, err error) ([]int, error) {
	if err != nil {
		return nil, err
	}
	res, err := redis.Ints(reply, err)
	if err == redis.ErrNil {
		return nil, ErrNil
	}
	return res, err
}
