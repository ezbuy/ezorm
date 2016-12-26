package test

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ezbuy/ezorm/db"
	redis "gopkg.in/redis.v5"
	"reflect"
	"strings"
)

var (
	_redis_store *db.RedisStore
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
	_redis_store = store
}

//! util functions
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

func redisStringScan(str string, val interface{}) error {
	return _redis_store.StringScan(str, val)
}

/////////////
func redisRename(obj Object, oldkey, newkey string) error {
	return _redis_store.Rename(keyOfClass(obj, oldkey), keyOfClass(obj, newkey)).Err()
}

func redisDelObject(obj Object) error {
	return _redis_store.Del(keyOfObject(obj)).Err()
}

func redisDelIndex(obj Object, index string) error {
	return _redis_store.Del(indexOfObject(obj, index)).Err()
}

func redisDelList(obj Object, list string) error {
	return _redis_store.Del(listOfObject(obj, list)).Err()
}

func redisDelKey(obj Object, key string) error {
	return _redis_store.Del(keyOfClass(obj, key)).Err()
}

func redisJsonSet(obj Object) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	return _redis_store.Set(keyOfObject(obj), data, 0).Err()
}

func redisJsonGet(obj Object) error {
	data, err := _redis_store.Get(keyOfObject(obj)).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, obj)
}

func redisFieldSet(obj Object, field string, value interface{}) error {
	return _redis_store.HSet(keyOfObject(obj), field, fmt.Sprint(value)).Err()
}

func redisFieldGet(obj Object, field string, value interface{}) error {
	v := _redis_store.HGet(keyOfObject(obj), field)
	if v.Err() != nil {
		return v.Err()
	}

	return v.Scan(value)
}

func redisIndexSet(obj Object, indexName string, indexValue interface{}, value interface{}) error {
	return _redis_store.SAdd(indexOfObject(obj, fmt.Sprintf("%s:%v", indexName, indexValue)), fmt.Sprint(value)).Err()
}

func redisIndexGet(obj Object, indexName string, indexValue interface{}) ([]string, error) {
	return _redis_store.SMembers(indexOfObject(obj, fmt.Sprintf("%s:%v", indexName, indexValue))).Result()
}

func redisIndexDel(obj Object, indexName string, indexValue interface{}) error {
	return _redis_store.Del(indexOfObject(obj, fmt.Sprintf("%s:%v", indexName, indexValue))).Err()
}

func redisIndexRemove(obj Object, indexName string, indexValue interface{}, value interface{}) error {
	return _redis_store.SRem(indexOfObject(obj, fmt.Sprintf("%s:%v", indexName, indexValue)), value).Err()
}

func redisMultiIndexesGet(obj Object, indexes ...string) ([]string, error) {
	keys := []string{}
	for _, idx := range indexes {
		keys = append(keys, indexOfObject(obj, idx))
	}
	return _redis_store.SInter(keys...).Result()
}

func redisSetAdd(obj Object, key string, value interface{}) error {
	return _redis_store.SAdd(keyOfClass(obj, key), fmt.Sprint(value)).Err()
}

func redisSetGet(obj Object, key string) ([]string, error) {
	return _redis_store.SMembers(keyOfClass(obj, key)).Result()
}

func redisSetDel(obj Object, key string) error {
	return _redis_store.Del(keyOfClass(obj, key)).Err()
}

func redisSetRem(obj Object, key string, value interface{}) error {
	return _redis_store.SRem(keyOfClass(obj, key), value).Err()
}

func redisZSetAdd(obj Object, key string, score float64, value interface{}) error {
	return _redis_store.ZAdd(keyOfClass(obj, key), redis.Z{Score: score, Member: value}).Err()
}

func redisZSetRangeByScore(obj Object, key string, min, max int64) ([]string, error) {
	return _redis_store.ZRangeByScore(keyOfClass(obj, key), redis.ZRangeBy{
		Min: fmt.Sprint(min),
		Max: fmt.Sprint(max),
	}).Result()
}

func redisZSetDel(obj Object, key string) error {
	return _redis_store.Del(keyOfClass(obj, key)).Err()
}

func redisZSetRem(obj Object, key string, value interface{}) error {
	return _redis_store.ZRem(keyOfClass(obj, key), value).Err()
}

func redisListLPush(obj Object, listName string, value interface{}) (int64, error) {
	return _redis_store.LPush(listOfObject(obj, listName), value).Result()
}

func redisListRPush(obj Object, listName string, value interface{}) (int64, error) {
	return _redis_store.RPush(listOfObject(obj, listName), value).Result()
}

func redisListLPop(obj Object, listName string, value interface{}) error {
	if !reflect.ValueOf(value).CanSet() {
		return errors.New("value can't be set")
	}

	v := _redis_store.LPop(listOfObject(obj, listName))
	if v.Err() != nil {
		return v.Err()
	}

	return v.Scan(value)
}

func redisListRPop(obj Object, listName string, value interface{}) error {
	v := _redis_store.RPop(listOfObject(obj, listName))
	if v.Err() != nil {
		return v.Err()
	}

	return v.Scan(value)
}

func redisListInsertBefore(obj Object, listName string, pivot, value interface{}) (int64, error) {
	return _redis_store.LInsertBefore(listOfObject(obj, listName), pivot, value).Result()
}

func redisListInsertAfter(obj Object, listName string, pivot, value interface{}) (int64, error) {
	return _redis_store.LInsertAfter(listOfObject(obj, listName), pivot, value).Result()
}

func redisListLength(obj Object, listName string) (int64, error) {
	return _redis_store.LLen(listOfObject(obj, listName)).Result()
}

func redisListRange(obj Object, listName string, start, stop int64) ([]string, error) {
	return _redis_store.LRange(listOfObject(obj, listName), start, stop).Result()
}

func redisListRemove(obj Object, listName string, value interface{}) error {
	return _redis_store.LRem(listOfObject(obj, listName), 0, fmt.Sprint(value)).Err()
}

func redisListCount(obj Object, listName string) (int64, error) {
	return _redis_store.LLen(listOfObject(obj, listName)).Result()
}

func redisListDel(obj Object, listName string) error {
	return _redis_store.Del(listOfObject(obj, listName)).Err()
}

func redisGeoAdd(obj Object, key string, longitude float64, latitude float64, value interface{}) error {
	return _redis_store.GeoAdd(keyOfClass(obj, key), &redis.GeoLocation{
		Longitude: longitude,
		Latitude:  latitude,
		Name:      fmt.Sprint(value),
	}).Err()
}

func redisGeoRadius(obj Object, key string, longitude float64, latitude float64, query *redis.GeoRadiusQuery) ([]string, error) {
	locations, err := _redis_store.GeoRadius(keyOfClass(obj, key), longitude, latitude, query).Result()
	if err != nil {
		return nil, err
	}

	strs := []string{}
	for _, loc := range locations {
		strs = append(strs, loc.Name)
	}
	return strs, nil
}

func redisGeoDel(obj Object, key string) error {
	return _redis_store.Del(keyOfClass(obj, key)).Err()
}
