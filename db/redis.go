package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	redis "gopkg.in/redis.v5"
)

const (
	JSON = "json"
	HASH = "hash"
	SET  = "set"
	ZSET = "zset"
	GEO  = "geo"
	LIST = "list"
)

type Object interface {
	GetClassName() string
	GetStoreType() string
	GetPrimaryKey() string
	GetIndexes() []string
}

type RedisStore struct {
	conn redis.Cmdable
}

func NewRedisStore(host string, port int, password string, db int) (*RedisStore, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
		DB:       db,
	})

	if err := client.Ping().Err(); err != nil {
		return nil, err
	}

	return &RedisStore{
		conn: client,
	}, nil
}

//! types

//! common
func keyOfObject(obj Object) string {
	objv := reflect.ValueOf(obj).Elem()
	return fmt.Sprintf("%s:%s:%v", obj.GetStoreType(), obj.GetClassName(),
		objv.FieldByName(obj.GetPrimaryKey()))
}

func indexOfObject(obj Object, index string) string {
	return fmt.Sprintf("%s:%s:%s", SET, obj.GetClassName(), index)
}

func listOfObject(obj Object, list string) string {
	return fmt.Sprintf("%s:%s:%s", LIST, obj.GetClassName(), list)
}

func keyOfClass(obj Object, keys ...string) string {
	if len(keys) > 0 {
		return fmt.Sprintf("%s:%s:%s", obj.GetStoreType(), obj.GetClassName(), strings.Join(keys, ":"))
	}
	return fmt.Sprintf("%s:%s", obj.GetStoreType(), obj.GetClassName())
}

//! functions
func (r *RedisStore) StringScan(str string, val interface{}) error {
	return redis.NewStringResult(str, nil).Scan(val)
}

func (r *RedisStore) Rename(obj Object, oldkey, newkey string) error {
	return r.conn.Rename(keyOfClass(obj, oldkey), keyOfClass(obj, newkey)).Err()
}

func (r *RedisStore) DelObject(obj Object) error {
	return r.conn.Del(keyOfObject(obj)).Err()
}

func (r *RedisStore) DelIndex(obj Object, index string) error {
	return r.conn.Del(indexOfObject(obj, index)).Err()
}

func (r *RedisStore) DelList(obj Object, list string) error {
	return r.conn.Del(listOfObject(obj, list)).Err()
}

func (r *RedisStore) DelKey(obj Object, key string) error {
	return r.conn.Del(keyOfClass(obj, key)).Err()
}

func (r *RedisStore) JsonSet(obj Object) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	return r.conn.Set(keyOfObject(obj), data, 0).Err()
}

func (r *RedisStore) JsonGet(obj Object) error {
	data, err := r.conn.Get(keyOfObject(obj)).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, obj)
}

func (r *RedisStore) FieldSet(obj Object, field string, value interface{}) error {
	return r.conn.HSet(keyOfObject(obj), field, fmt.Sprint(value)).Err()
}

func (r *RedisStore) FieldGet(obj Object, field string, value interface{}) error {
	v := r.conn.HGet(keyOfObject(obj), field)
	if v.Err() != nil {
		return v.Err()
	}

	return v.Scan(value)
}

func (r *RedisStore) IndexAdd(obj Object, indexName string, indexValue interface{}, value interface{}) error {
	return r.conn.SAdd(indexOfObject(obj, fmt.Sprintf("%s:%v", indexName, indexValue)), fmt.Sprint(value)).Err()
}

func (r *RedisStore) IndexGet(obj Object, indexName string, indexValue interface{}) ([]string, error) {
	return r.conn.SMembers(indexOfObject(obj, fmt.Sprintf("%s:%v", indexName, indexValue))).Result()
}

func (r *RedisStore) IndexDel(obj Object, indexName string, indexValue interface{}) error {
	return r.conn.Del(indexOfObject(obj, fmt.Sprintf("%s:%v", indexName, indexValue))).Err()
}

func (r *RedisStore) IndexRem(obj Object, indexName string, indexValue interface{}, value interface{}) error {
	return r.conn.SRem(indexOfObject(obj, fmt.Sprintf("%s:%v", indexName, indexValue)), value).Err()
}

func (r *RedisStore) MultiIndexesGet(obj Object, indexes ...string) ([]string, error) {
	keys := []string{}
	for _, idx := range indexes {
		keys = append(keys, indexOfObject(obj, idx))
	}
	return r.conn.SInter(keys...).Result()
}

func (r *RedisStore) SetAdd(obj Object, key string, value interface{}) error {
	return r.conn.SAdd(keyOfClass(obj, key), fmt.Sprint(value)).Err()
}

func (r *RedisStore) SetGet(obj Object, key string) ([]string, error) {
	return r.conn.SMembers(keyOfClass(obj, key)).Result()
}

func (r *RedisStore) SetDel(obj Object, key string) error {
	return r.conn.Del(keyOfClass(obj, key)).Err()
}

func (r *RedisStore) SetRem(obj Object, key string, value interface{}) error {
	return r.conn.SRem(keyOfClass(obj, key), value).Err()
}

func (r *RedisStore) ZSetAdd(obj Object, key string, score float64, value interface{}) error {
	return r.conn.ZAdd(keyOfClass(obj, key), redis.Z{Score: score, Member: value}).Err()
}

func (r *RedisStore) ZSetRangeByScore(obj Object, key string, min, max int64) ([]string, error) {
	return r.conn.ZRangeByScore(keyOfClass(obj, key), redis.ZRangeBy{
		Min: fmt.Sprint(min),
		Max: fmt.Sprint(max),
	}).Result()
}

func (r *RedisStore) ZSetDel(obj Object, key string) error {
	return r.conn.Del(keyOfClass(obj, key)).Err()
}

func (r *RedisStore) ZSetRem(obj Object, key string, value interface{}) error {
	return r.conn.ZRem(keyOfClass(obj, key), value).Err()
}

func (r *RedisStore) ListLPush(obj Object, listName string, value interface{}) (int64, error) {
	return r.conn.LPush(listOfObject(obj, listName), value).Result()
}

func (r *RedisStore) ListRPush(obj Object, listName string, value interface{}) (int64, error) {
	return r.conn.RPush(listOfObject(obj, listName), value).Result()
}

func (r *RedisStore) ListLPop(obj Object, listName string, value interface{}) error {
	if !reflect.ValueOf(value).CanSet() {
		return errors.New("value can't be set")
	}

	v := r.conn.LPop(listOfObject(obj, listName))
	if v.Err() != nil {
		return v.Err()
	}

	return v.Scan(value)
}

func (r *RedisStore) ListRPop(obj Object, listName string, value interface{}) error {
	v := r.conn.RPop(listOfObject(obj, listName))
	if v.Err() != nil {
		return v.Err()
	}

	return v.Scan(value)
}

func (r *RedisStore) ListInsertBefore(obj Object, listName string, pivot, value interface{}) (int64, error) {
	return r.conn.LInsertBefore(listOfObject(obj, listName), pivot, value).Result()
}

func (r *RedisStore) ListInsertAfter(obj Object, listName string, pivot, value interface{}) (int64, error) {
	return r.conn.LInsertAfter(listOfObject(obj, listName), pivot, value).Result()
}

func (r *RedisStore) ListLength(obj Object, listName string) (int64, error) {
	return r.conn.LLen(listOfObject(obj, listName)).Result()
}

func (r *RedisStore) ListRange(obj Object, listName string, start, stop int64) ([]string, error) {
	return r.conn.LRange(listOfObject(obj, listName), start, stop).Result()
}

func (r *RedisStore) ListRem(obj Object, listName string, value interface{}) error {
	return r.conn.LRem(listOfObject(obj, listName), 0, fmt.Sprint(value)).Err()
}

func (r *RedisStore) ListCount(obj Object, listName string) (int64, error) {
	return r.conn.LLen(listOfObject(obj, listName)).Result()
}

func (r *RedisStore) ListDel(obj Object, listName string) error {
	return r.conn.Del(listOfObject(obj, listName)).Err()
}

func (r *RedisStore) GeoAdd(obj Object, key string, longitude float64, latitude float64, value interface{}) error {
	return r.conn.GeoAdd(keyOfClass(obj, key), &redis.GeoLocation{
		Longitude: longitude,
		Latitude:  latitude,
		Name:      fmt.Sprint(value),
	}).Err()
}

func (r *RedisStore) GeoRadius(obj Object, key string, longitude float64, latitude float64, query *redis.GeoRadiusQuery) ([]string, error) {
	locations, err := r.conn.GeoRadius(keyOfClass(obj, key), longitude, latitude, query).Result()
	if err != nil {
		return nil, err
	}

	strs := []string{}
	for _, loc := range locations {
		strs = append(strs, loc.Name)
	}
	return strs, nil
}

func (r *RedisStore) GeoDel(obj Object, key string) error {
	return r.conn.Del(keyOfClass(obj, key)).Err()
}
