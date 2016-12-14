package test

///////redis_orm//////

func (m *_UserMgr) SetUser(obj *User) error {
	return redisSetObject(obj)
}

func (m *_UserMgr) GetUser(obj *User) error {
	return redisGetObject(obj)
}

func (m *_UserMgr) DelUser(obj *User) error {
	return redisDelObject(obj)
}
