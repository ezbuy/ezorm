package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
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

	return fmt.Sprintf("%s:%s:%s:%d", obj.GetStoreType(), obj.GetClassName(), obj.GetPrimaryKey(), v.FieldByName(obj.GetPrimaryKey()).Int()), nil
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

func KeyOfObjectIndex(obj Object, index string) (string, error) {
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
	if val != nil {
		return fmt.Sprintf("%s:%s:%s:%v", SET, obj.GetClassName(), index, val), nil
	}
	return "", errors.New("object index invalid")
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
	return nil, nil
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

func (r *RedisStore) setFieldValue(field reflect.Value, val interface{}) error {
	switch field.Kind() {
	case reflect.String:
		valueAsString, err := r.String(val, nil)
		if err != nil {
			return err
		}
		field.SetString(valueAsString)
	case reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8, reflect.Uint, reflect.Uintptr:
		valueAsUint, err := r.Uint64(val, nil)
		if err != nil {
			return err
		}
		field.SetUint(valueAsUint)
	case reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8, reflect.Int:
		valueAsInt, err := r.Int64(val, nil)
		if err != nil {
			return err
		}
		field.SetInt(valueAsInt)
	case reflect.Float32, reflect.Float64:
		valueAsFloat, err := r.Float64(val, nil)
		if err != nil {
			return err
		}
		field.SetFloat(valueAsFloat)
	case reflect.Bool:
		boolValue, err := r.Bool(val, nil)
		if err != nil {
			return err
		}
		field.SetBool(boolValue)
	case reflect.Struct, reflect.Ptr:
		valueAsInt, err := r.Int64(val, nil)
		if err != nil {
			return err
		}
		valuleDecode, err := FieldDecode(field, valueAsInt)
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

	primary_key := v.FieldByName(obj.GetPrimaryKey()).Int()

	switch obj.GetStoreType() {
	case JSON:
		bytes, err := json.Marshal(obj)
		if err != nil {
			return err
		}
		return r.SET(key_of_obj, bytes)
	case HASH:
		//! object fields
		for i := 0; i < v.NumField(); i++ {
			ft := t.Field(i)
			fv := v.Field(i)

			if fv.CanInterface() {
				val, _ := FieldEncode(fv)
				if err := r.HSET(key_of_obj, ft.Name, val); err != nil {
					r.DEL(key_of_obj)
					return err
				}
			}
		}

		//! object indexs
		for _, idx := range obj.GetIndexes() {
			if key_of_index, err := KeyOfObjectIndex(obj, idx); err == nil {
				_, err := r.SADD(key_of_index, primary_key)
				if err != nil {
					r.DEL(key_of_obj)
					return err
				}
			}
		}

		//! object primary key
		_, err := r.ZADD(key_of_cls, primary_key, primary_key)
		if err != nil {
			r.DEL(key_of_obj)
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
		_, err = r.SADD(key_of_obj, v1)
		return err
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
		case reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8, reflect.Int:
		default:
			return errors.New("zset struct score field should be number type")
		}

		val, err := FieldEncode(v.Field(2))
		if err != nil {
			return err
		}

		_, err = r.ZADD(key_of_obj, v.Field(1).Int(), val)
		return err
	case GEO:
		if FieldNum(v) != 4 {
			return errors.New("geo struct only support 4 fields <geo-key, longitude, latitude, value>")
		}
		//field 0 should be primary key
		if t.Field(0).Name != obj.GetPrimaryKey() {
			return errors.New("zset struct first field should be primary key")
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
		val, err := FieldEncode(v.Field(3))
		if err != nil {
			return err
		}

		_, err = r.GEOADD(key_of_obj, v.Field(1).Float(), v.Field(2).Float(), val)
		return err
	}
	return errors.New("store unsupport type")
}

func (r *RedisStore) GetObject(obj Object) error {
	t := reflect.TypeOf(obj)
	if t.Kind() != reflect.Ptr {
		return errors.New("object unsupport type")
	}
	t = t.Elem()
	v := reflect.ValueOf(obj).Elem()
	if t.Kind() != reflect.Struct {
		return errors.New("object unsupport type")
	}

	key_of_obj, err := KeyOfObject(obj)
	if err != nil {
		return err
	}
	switch obj.GetStoreType() {
	case JSON:
		data, err := r.Bytes(r.GET(key_of_obj))
		if err != nil {
			return err
		}
		if err := json.Unmarshal(data, obj); err != nil {
			return err
		}
		return nil
	case HASH:
		for i := 0; i < v.NumField(); i++ {
			if v.Field(i).CanInterface() {
				val, err := r.HGET(key_of_obj, t.Field(i).Name)
				if err != nil {
					return err
				}
				if err := r.setFieldValue(v.Field(i), val); err != nil {
					return err
				}
			}
		}
	case SET, ZSET, GEO:
		return errors.New("set,zset,geo struct not support get object")
	}
	return nil
}

func (r *RedisStore) DelObject(obj Object) error {
	key_of_obj, err := KeyOfObject(obj)
	if err != nil {
		return err
	}
	_, err = r.DEL(key_of_obj)
	return err
}
