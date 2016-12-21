package test

import (
	"time"
)

var (
	_ time.Time
)

///////////// ZSET /////////////////////////////////////////////////////
func (m *_SortUserBlogMgr) ZRelationSet(key string, score float64, obj *SortUserBlog) error {
	return redisZSetAdd(obj, key, score, obj.Value)
}

func (m *_SortUserBlogMgr) ZRelationRangeByScore(key string, min, max int64) ([]*SortUserBlog, error) {
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

func (m *_SortUserBlogMgr) ZRelationRem(key string, obj *SortUserBlog) error {
	return redisZSetRem(obj, key, obj.Value)
}

func (m *_SortUserBlogMgr) ZRelationDel(key string) error {
	return redisZSetDel(m.NewSortUserBlog(), key)
}
