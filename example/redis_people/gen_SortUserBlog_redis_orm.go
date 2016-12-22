package test

import (
	"time"
)

var (
	_ time.Time
)

///////////// ZSET /////////////////////////////////////////////////////
func (m *_SortUserBlogMgr) ZAdd(key string, score float64, obj *SortUserBlog) error {
	return redisZSetAdd(obj, key, score, obj.Value)
}

func (m *_SortUserBlogMgr) ZRangeByScore(key string, min, max int64) ([]*SortUserBlog, error) {
	strs, err := redisZSetRangeByScore(m.NewSortUserBlog(), key, min, max)
	if err != nil {
		return nil, err
	}

	objs := []*SortUserBlog{}
	for _, str := range strs {
		obj := m.NewSortUserBlog()
		if err := redisStringScan(str, &obj.Value); err != nil {
			return nil, err
		}
		objs = append(objs, obj)
	}
	return objs, nil
}

func (m *_SortUserBlogMgr) ZRem(key string, obj *SortUserBlog) error {
	return redisZSetRem(obj, key, obj.Value)
}

func (m *_SortUserBlogMgr) ZDel(key string) error {
	return redisZSetDel(m.NewSortUserBlog(), key)
}
func (m *_SortUserBlogMgr) ZRangeRelatedBlog(key string, min, max int64) ([]*Blog, error) {
	strs, err := redisZSetRangeByScore(m.NewSortUserBlog(), key, min, max)
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
