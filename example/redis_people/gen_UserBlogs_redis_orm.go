package test

func (m *_UserBlogsMgr) SetUserBlogs(obj *UserBlogs) error {
	return redisSetObject(obj)
}

func (m *_UserBlogsMgr) DelUserBlogs(obj *UserBlogs) error {
	return redisDelObject(obj)
}

///////////// SET /////////////////////////////////////////////////////
