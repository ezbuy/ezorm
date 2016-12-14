package test

///////redis_orm//////

func (m *_BlogMgr) SetBlog(obj *Blog) error {
	return redisSetObject(obj)
}

func (m *_BlogMgr) GetBlog(obj *Blog) error {
	return redisGetObject(obj)
}

func (m *_BlogMgr) DelBlog(obj *Blog) error {
	return redisDelObject(obj)
}
