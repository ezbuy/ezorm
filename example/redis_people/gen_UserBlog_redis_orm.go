package test

import (
	"time"
)

var (
	_ time.Time
)

func (m *_UserBlogMgr) SetAddBySQLs(sqls ...string) error {
	querys := []string{}
	if len(sqls) > 0 {
		querys = append(querys, sqls...)
	} else {
		querys = append(querys, "SELECT CONCAT('UserId:', UserId) AS k, ID AS v FROM BLOGS")
	}
	for _, sql := range querys {
		objs, err := m.Query(sql)
		if err != nil {
			return err
		}

		for _, obj := range objs {
			if err := m.SetAdd(obj.Key, obj); err != nil {
				return err
			}
		}
	}
	return nil
}

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
		obj.Key = key
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
