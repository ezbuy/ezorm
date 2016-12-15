package test

func (m *_SortUserBlogsMgr) SetSortUserBlogs(obj *SortUserBlogs) error {
	return redisSetObject(obj)
}

func (m *_SortUserBlogsMgr) GetSortUserBlogs(obj *SortUserBlogs) error {
	return redisGetObject(obj)
}

func (m *_SortUserBlogsMgr) DelSortUserBlogs(obj *SortUserBlogs) error {
	return redisDelObject(obj)
}
