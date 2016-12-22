package test

import (
	"time"
)

var (
	_ time.Time
)

func (m *_SortUserBlogMgr) ZAddBySQLs(sqls ...string) error {
	querys := []string{}
	if len(sqls) > 0 {
		querys = append(querys, sqls...)
	} else {
		querys = append(querys, "SELECT CONCAT('UserId:', UserId) AS k, SCORE, ID AS v FROM BLOGS")
	}
	for _, sql := range querys {
		objs, err := m.Query(sql)
		if err != nil {
			return err
		}

		for _, obj := range objs {
			if err := m.ZAdd(obj.Key, obj); err != nil {
				return err
			}
		}
	}
	return nil
}

///////////// ZSET /////////////////////////////////////////////////////
func (m *_SortUserBlogMgr) ZAdd(key string, obj *SortUserBlog) error {
	return redisZSetAdd(obj, key, obj.Score, obj.Value)
}

func (m *_SortUserBlogMgr) ZRangeByScore(key string, min, max int64) ([]*SortUserBlog, error) {
	strs, err := redisZSetRangeByScore(m.NewSortUserBlog(), key, min, max)
	if err != nil {
		return nil, err
	}

	objs := []*SortUserBlog{}
	for _, str := range strs {
		obj := m.NewSortUserBlog()
		obj.Key = key
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
