package test

func (m *_UserLocationMgr) SetUserLocation(obj *UserLocation) error {
	return redisSetObject(obj)
}

func (m *_UserLocationMgr) GetUserLocation(obj *UserLocation) error {
	return redisGetObject(obj)
}

func (m *_UserLocationMgr) DelUserLocation(obj *UserLocation) error {
	return redisDelObject(obj)
}
