package test

import (
	"time"
)

var (
	_ time.Time
)

///////////// SET /////////////////////////////////////////////////////
func (m *_UserBlogMgr) SetAdd(key string, obj *UserBlog) error {
	return redisSetAdd(obj, key, obj.Value)
}

func (m *_UserBlogMgr) SetGet(key string) ([]*UserBlog, error) {
	strs, err := redisSetGet(m.NewUserBlog(), key)
	if err != nil {
		return nil, err
	}

	objs := []*UserBlog{}
	for _, str := range strs {
		obj := m.NewUserBlog()
		if err := redisStringScan(str, &obj.Value); err != nil {
			return nil, err
		}
		objs = append(objs, obj)
	}
	return objs, nil
}

func (m *_UserBlogMgr) SetRem(key string, obj *UserBlog) error {
	return redisSetRem(obj, key, obj.Value)
}

func (m *_UserBlogMgr) SetDel(key string) error {
	return redisSetDel(m.NewUserBlog(), key)
}
func (m *_UserBlogMgr) RelatedBlogs(key string) ([]*Blog, error) {
	strs, err := redisSetGet(m.NewUserBlog(), key)
	if err != nil {
		return nil, err
	}

	objs := []*Blog{}
	for _, str := range strs {
		var val int32
		if err := redisStringScan(str, &val); err != nil {
			return nil, err
		}
		if obj, err := BlogMgr.GetBlogById(val); err == nil {
			objs = append(objs, obj)
		}
	}
	return objs, nil
}
