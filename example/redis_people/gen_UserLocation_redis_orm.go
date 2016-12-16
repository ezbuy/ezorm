package test

func (m *_UserLocationMgr) SetUserLocation(obj *UserLocation) error {
	return redisSetObject(obj)
}

func (m *_UserLocationMgr) DelUserLocation(obj *UserLocation) error {
	return redisDelObject(obj)
}

///////////// GEO /////////////////////////////////////////////////////
