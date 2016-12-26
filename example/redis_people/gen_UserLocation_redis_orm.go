package test

import (
	redis "gopkg.in/redis.v5"
	"time"
)

var (
	_ time.Time
)

func (m *_UserLocationMgr) AddBySQL(sql string, args ...interface{}) error {
	objs, err := m.Query(sql)
	if err != nil {
		return err
	}

	for _, obj := range objs {
		if err := m.GeoAdd(obj.Key, obj); err != nil {
			return err
		}
	}
	return nil
}

func (m *_UserLocationMgr) DelBySQL(sql string, args ...interface{}) error {
	objs, err := m.Query(sql)
	if err != nil {
		return err
	}

	for _, obj := range objs {
		if err := m.GeoRem(obj.Key, obj); err != nil {
			return err
		}
	}
	return nil
}
func (m *_UserLocationMgr) Import() error {
	return m.AddBySQL("SELECT CONCAT('UserId:', UserId) AS k, Longitude, Latitude, ID AS v FROM BLOGS")
}

///////////// GEO /////////////////////////////////////////////////////
func (m *_UserLocationMgr) GeoAdd(key string, obj *UserLocation) error {
	return redisGeoAdd(obj, key, obj.Longitude, obj.Latitude, obj.Value)
}

func (m *_UserLocationMgr) GeoRadius(key string, longitude float64, latitude float64, query *redis.GeoRadiusQuery) ([]*UserLocation, error) {
	strs, err := redisGeoRadius(m.NewUserLocation(), key, longitude, latitude, query)
	if err != nil {
		return nil, err
	}

	objs := []*UserLocation{}
	for _, str := range strs {
		obj := m.NewUserLocation()
		obj.Key = key
		if err := redisStringScan(str, &obj.Value); err != nil {
			return nil, err
		}
		objs = append(objs, obj)
	}
	return objs, nil
}
func (m *_UserLocationMgr) GeoRadiusRelatedUsers(key string, longitude float64, latitude float64, query *redis.GeoRadiusQuery) ([]*User, error) {
	strs, err := redisGeoRadius(m.NewUserLocation(), key, longitude, latitude, query)
	if err != nil {
		return nil, err
	}

	objs := []*User{}
	for _, str := range strs {
		var val int32
		if err := redisStringScan(str, &val); err != nil {
			return nil, err
		}
		if obj, err := UserMgr.GetUserById(val); err == nil {
			objs = append(objs, obj)
		}
	}
	return objs, nil
}

func (m *_UserLocationMgr) GeoRem(key string, obj *UserLocation) error {
	return redisZSetRem(m.NewUserLocation(), key, obj)
}

func (m *_UserLocationMgr) GeoDel(key string) error {
	return redisGeoDel(m.NewUserLocation(), key)
}