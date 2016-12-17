package db

import (
	"fmt"

	redis "gopkg.in/redis.v5"
)

type RedisStore struct {
	conn redis.Cmdable
}

func NewRedisStore(host string, port int, password string, db int) (*RedisStore, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
		DB:       db,
	})

	//! ping the redis-server
	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}

	return &RedisStore{
		conn: client,
	}, nil
}

func ZValue(score float64, member interface{}) redis.Z {
	return redis.Z{
		Score:  score,
		Member: member,
	}
}

func NewGeoLocation(name string, longitude, latitude float64) *redis.GeoLocation {
	return &redis.GeoLocation{
		Name:      name,
		Longitude: longitude,
		Latitude:  latitude,
	}
}
