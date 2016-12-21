package test

import (
	"time"
)

var (
	_ time.Time
)

///////////// SET /////////////////////////////////////////////////////
func (m *_UserBlogMgr) RelationSet(key string, obj *UserBlog) error {
	return redisSetAdd(obj, key, obj.Value)
}

func (m *_UserBlogMgr) RelationGet(key string) ([]*UserBlog, error) {
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

func (m *_UserBlogMgr) RelationRem(key string, obj *UserBlog) error {
	return redisSetRem(obj, key, obj.Value)
}

func (m *_UserBlogMgr) RelationDel(key string) error {
	return redisSetDel(m.NewUserBlog(), key)
}
