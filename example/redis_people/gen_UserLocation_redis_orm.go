package test

import (
	redis "gopkg.in/redis.v5"
	"time"
)

var (
	_ time.Time
)

///////////// GEO /////////////////////////////////////////////////////
func (m *_UserLocationMgr) GeoAdd(key string, longitude float64, latitude float64, obj *UserLocation) error {
	return redisGeoAdd(obj, key, longitude, latitude, obj.Value)
}

func (m *_UserLocationMgr) GeoRadius(key string, longitude float64, latitude float64, query *redis.GeoRadiusQuery) ([]*UserLocation, error) {
	strs, err := redisGeoRadius(m.NewUserLocation(), key, longitude, latitude, query)
	if err != nil {
		return nil, err
	}

	objs := []*UserLocation{}
	for _, str := range strs {
		obj := m.NewUserLocation()
		if err := redisStringScan(str, &obj.Value); err != nil {
			return nil, err
		}
		objs = append(objs, obj)
	}
	return objs, nil
}

func (m *_UserLocationMgr) GeoDel(key string) error {
	return redisGeoDel(m.NewUserLocation(), key)
}
