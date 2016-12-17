package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	redis "gopkg.in/redis.v5"
)

const (
	JSON = "json"
	HASH = "hash"
	SET  = "set"
	ZSET = "zset"
	GEO  = "geo"
)

type Object interface {
	GetClassName() string
	GetStoreType() string
	GetPrimaryKey() string
	GetIndexes() []string
}

func KeyOfObject(obj Object) (string, error) {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}
	if t.Kind() != reflect.Struct {
		return "", errors.New("object unsupport type")
	}

	return fmt.Sprintf("%s:%s:%s:%v", obj.GetStoreType(), obj.GetClassName(), obj.GetPrimaryKey(), v.FieldByName(obj.GetPrimaryKey()).Interface()), nil
}

func KeyOfObjectById(obj Object, id string) (string, error) {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}
	if t.Kind() != reflect.Struct {
		return "", errors.New("object unsupport type")
	}

	return fmt.Sprintf("%s:%s:%s:%s", obj.GetStoreType(), obj.GetClassName(), obj.GetPrimaryKey(), id), nil
}

func KeyOfClass(obj Object) (string, error) {
	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return "", errors.New("object unsupport type")
	}
	return fmt.Sprintf("%s:%s:%s", ZSET, obj.GetClassName(), obj.GetPrimaryKey()), nil
}

func KeyOfIndexByClass(class string, indexName string, indexValue interface{}) (string, error) {
	return fmt.Sprintf("%s:%s:%s:%v", SET, class, indexName, indexValue), nil
}

func KeyOfIndexByObject(obj Object, index string) (string, error) {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}
	if t.Kind() != reflect.Struct {
		return "", errors.New("object unsupport type")
	}

	val, err := FieldEncode(v.FieldByName(index))
	if err != nil {
		return "", err
	}
	return KeyOfIndexByClass(obj.GetClassName(), index, val)
}

func FieldEncode(fv reflect.Value) (interface{}, error) {
	if fv.CanInterface() {
		val := fv.Interface()
		if fv.Kind() == reflect.Ptr || fv.Kind() == reflect.Struct {
			if fv.Type() == reflect.TypeOf((*time.Time)(nil)) {
				val = val.(*time.Time).Unix()
				return val, nil
			}
			if fv.Type() == reflect.TypeOf(time.Now()) {
				val = val.(time.Time).Unix()
				return val, nil
			}
			return nil, errors.New("field unsupport complex type")
		}
		return val, nil
	}
	return nil, errors.New("field unexport")
}

func FieldDecode(fv reflect.Value, val interface{}) (interface{}, error) {
	if fv.Kind() == reflect.Ptr || fv.Kind() == reflect.Struct {
		if fv.Type() == reflect.TypeOf((*time.Time)(nil)) {
			tm := time.Unix(val.(int64), 0)
			return &tm, nil
		}
		if fv.Type() == reflect.TypeOf(time.Now()) {
			tm := time.Unix(val.(int64), 0)
			return tm, nil
		}
		return nil, errors.New("field unsupport complex type")
	}
	return val, nil
}

func FieldNum(v reflect.Value) int {
	num := 0
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).CanInterface() {
			num++
		}
	}
	return num
}

func (r *RedisStore) setFieldValue(field reflect.Value, stringCmd *redis.StringCmd) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(stringCmd.Val())
	case reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8, reflect.Uint, reflect.Uintptr:
		v, err := stringCmd.Uint64()
		if err != nil {
			return err
		}
		field.SetUint(v)
	case reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8, reflect.Int:
		v, err := stringCmd.Int64()
		if err != nil {
			return err
		}
		field.SetInt(v)
	case reflect.Float32, reflect.Float64:
		v, err := stringCmd.Float64()
		if err != nil {
			return err
		}
		field.SetFloat(v)
	case reflect.Bool:
		var b bool
		if err := stringCmd.Scan(&b); err != nil {
			return err
		}
		field.SetBool(b)
	case reflect.Struct, reflect.Ptr:
		v, err := stringCmd.Int64()
		if err != nil {
			return err
		}
		valuleDecode, err := FieldDecode(field, v)
		if err != nil {
			return err
		}
		field.Set(reflect.ValueOf(valuleDecode))
	default:
		return errors.New("Unsupported Type")
	}
	return nil
}

func (r *RedisStore) SetObject(obj Object) error {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}
	if t.Kind() != reflect.Struct {
		return errors.New("object unsupport type")
	}

	key_of_obj, err := KeyOfObject(obj)
	if err != nil {
		return err
	}

	key_of_cls, err := KeyOfClass(obj)
	if err != nil {
		return err
	}

	primary_key_str := fmt.Sprint(v.FieldByName(obj.GetPrimaryKey()).Interface())
	primary_key_num := float64(v.FieldByName(obj.GetPrimaryKey()).Int())

	switch obj.GetStoreType() {
	case JSON:
		bytes, err := json.Marshal(obj)
		if err != nil {
			return err
		}

		if err := r.conn.Set(key_of_obj, bytes, time.Duration(0)).Err(); err != nil {
			return err
		}

		//! object indexs
		for _, idx := range obj.GetIndexes() {
			if key_of_index, err := KeyOfIndexByObject(obj, idx); err == nil {

				if err := r.conn.SAdd(key_of_index, primary_key_str).Err(); err != nil {
					r.conn.Del(key_of_index)
					return err
				}
			}
		}

		//! object primary key
		if err = r.conn.ZAdd(key_of_cls, ZValue(primary_key_num, primary_key_str)).Err(); err != nil {
			r.conn.Del(key_of_obj)
			return err
		}

		return nil
	case HASH:
		//! object fields
		for i := 0; i < v.NumField(); i++ {
			ft := t.Field(i)
			fv := v.Field(i)

			if fv.CanInterface() {
				val, err := FieldEncode(fv)
				if err != nil {
					return err
				}
				if err := r.conn.HSet(key_of_obj, ft.Name, fmt.Sprint(val)).Err(); err != nil {
					r.conn.Del(key_of_obj)
					return err
				}
			}
		}

		//! object indexs
		for _, idx := range obj.GetIndexes() {
			if key_of_index, err := KeyOfIndexByObject(obj, idx); err == nil {

				if err := r.conn.SAdd(key_of_index, primary_key_str).Err(); err != nil {
					r.conn.Del(key_of_index)
					return err
				}
			}
		}

		//! object primary key
		if err = r.conn.ZAdd(key_of_cls, ZValue(primary_key_num, primary_key_str)).Err(); err != nil {
			r.conn.Del(key_of_obj)
			return err
		}
		return nil
	case SET:
		if FieldNum(v) != 2 {
			return errors.New("set struct only support 2 fields <set-key, value>")
		}
		//field 0 should be primary key
		if t.Field(0).Name != obj.GetPrimaryKey() {
			return errors.New("set struct first field should be primary key")
		}

		f1 := v.Field(1)
		v1, err := FieldEncode(f1)
		if err != nil {
			return err
		}
		return r.conn.SAdd(key_of_obj, v1).Err()
	case ZSET:
		if FieldNum(v) != 3 {
			return errors.New("zset struct only support 3 fields <zset-key, score, value>")
		}
		//field 0 should be primary key
		if t.Field(0).Name != obj.GetPrimaryKey() {
			return errors.New("zset struct first field should be primary key")
		}
		//field 1 should be score key
		if strings.ToLower(t.Field(1).Name) != "score" {
			return errors.New("zset struct second field should be score")
		}
		switch v.Field(1).Kind() {
		case reflect.Float32, reflect.Float64:
		default:
			return errors.New("zset struct score field should be float type")
		}

		val, err := FieldEncode(v.Field(2))
		if err != nil {
			return err
		}
		return r.conn.ZAdd(key_of_obj, ZValue(v.Field(1).Float(), val)).Err()
	case GEO:
		if FieldNum(v) != 3 {
			return errors.New("geo struct only support 4 fields <geo-key, longitude, latitude, value>")
		}
		//field 0 should be primary key
		if t.Field(0).Name != obj.GetPrimaryKey() {
			return errors.New("zset struct first field should be primary key")
		}
		v0, err := FieldEncode(v.Field(0))
		if err != nil {
			return err
		}

		//field 1 should be longitude key
		if strings.ToLower(t.Field(1).Name) != "longitude" {
			return errors.New("zset struct second field should be longitude")
		}
		switch v.Field(1).Kind() {
		case reflect.Float32, reflect.Float64:
		default:
			return errors.New("zset struct longitude field should be float type")
		}

		//field 2 should be latitude key
		if strings.ToLower(t.Field(2).Name) != "latitude" {
			return errors.New("zset struct second field should be latitude")
		}
		switch v.Field(2).Kind() {
		case reflect.Float32, reflect.Float64:
		default:
			return errors.New("zset struct latitude field should be float type")
		}

		return r.conn.GeoAdd(key_of_cls, NewGeoLocation(fmt.Sprint(v0), v.Field(1).Float(), v.Field(2).Float())).Err()
	}
	return errors.New("store unsupport type")
}

func (r *RedisStore) GetObjectById(obj Object, id string) error {
	t := reflect.TypeOf(obj)
	if t.Kind() != reflect.Ptr {
		return errors.New("object unsupport type")
	}
	t = t.Elem()
	v := reflect.ValueOf(obj).Elem()
	if t.Kind() != reflect.Struct {
		return errors.New("object unsupport type")
	}

	key_of_obj, err := KeyOfObjectById(obj, id)
	if err != nil {
		return err
	}
	switch obj.GetStoreType() {
	case JSON:
		data, err := r.conn.Get(key_of_obj).Bytes()
		if err != nil {
			return err
		}
		if err := json.Unmarshal(data, obj); err != nil {
			return err
		}
	case HASH:
		for i := 0; i < v.NumField(); i++ {
			fv := v.Field(i)
			if fv.CanInterface() {
				cmd := r.conn.HGet(key_of_obj, t.Field(i).Name)

				if err := r.setFieldValue(fv, cmd); err != nil {
					return err
				}
			}
		}
	case SET, ZSET, GEO:
		return errors.New("set,zset,geo struct not support get object")
	}
	return nil
}

func (r *RedisStore) GetObject(obj Object) error {
	t := reflect.TypeOf(obj)
	if t.Kind() != reflect.Ptr {
		return errors.New("object unsupport type")
	}
	v := reflect.ValueOf(obj).Elem()

	return r.GetObjectById(obj, fmt.Sprint(v.FieldByName(obj.GetPrimaryKey()).Interface()))
}

func (r *RedisStore) DelObject(obj Object) error {
	key_of_obj, err := KeyOfObject(obj)
	if err != nil {
		return err
	}
	return r.conn.Del(key_of_obj).Err()
}

//redisSMEMBERSInts(key string) ([]int, error) {
func (r *RedisStore) SMembersIds(key string) ([]string, error) {
	return r.conn.SMembers(key).Result()
}

func (r *RedisStore) SInterIds(keys ...string) ([]string, error) {
	return r.conn.SInter(keys...).Result()
}
