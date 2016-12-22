package test

import (
	redis "gopkg.in/redis.v5"
	"time"
)

var (
	_ time.Time
)

func (m *_UserLocationMgr) GeoAddBySQLs(sqls ...string) error {
	querys := []string{}
	if len(sqls) > 0 {
		querys = append(querys, sqls...)
	} else {
		querys = append(querys, "SELECT CONCAT('UserId:', UserId) AS k, Longitude, Latitude, ID AS v FROM BLOGS")
	}
	for _, sql := range querys {
		objs, err := m.Query(sql)
		if err != nil {
			return err
		}

		for _, obj := range objs {
			if err := m.GeoAdd(obj.Key, obj); err != nil {
				return err
			}
		}
	}
	return nil
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

func (m *_UserLocationMgr) GeoDel(key string) error {
	return redisGeoDel(m.NewUserLocation(), key)
}
