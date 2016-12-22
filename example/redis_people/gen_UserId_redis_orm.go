package test

import (
	"time"
)

var (
	_ time.Time
)

func (m *_UserIdMgr) ListLPushBySQLs(sqls ...string) error {
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
			if _, err := m.ListLPush(obj.Key, obj); err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *_UserIdMgr) ListRPushBySQLs(sqls ...string) error {
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
			if _, err := m.ListRPush(obj.Key, obj); err != nil {
				return err
			}
		}
	}
	return nil
}

///////////// LIST /////////////////////////////////////////////////////
func (m *_UserIdMgr) ListLPush(key string, obj *UserId) (int64, error) {
	return redisListLPush(obj, key, obj.Value)
}

func (m *_UserIdMgr) ListRPush(key string, obj *UserId) (int64, error) {
	return redisListRPush(obj, key, obj.Value)
}

func (m *_UserIdMgr) ListLPop(key string, obj *UserId) error {
	obj.Key = key
	return redisListLPop(obj, key, &obj.Value)
}

func (m *_UserIdMgr) ListRPop(key string, obj *UserId) error {
	obj.Key = key
	return redisListRPop(obj, key, &obj.Value)
}

func (m *_UserIdMgr) ListInsertBefore(key string, pivot, obj *UserId) (int64, error) {
	return redisListInsertBefore(obj, key, pivot.Value, obj.Value)
}

func (m *_UserIdMgr) ListInsertAfter(key string, pivot, obj *UserId) (int64, error) {
	return redisListInsertAfter(obj, key, pivot.Value, obj.Value)
}

func (m *_UserIdMgr) ListRange(key string, start, stop int64) ([]*UserId, error) {
	strs, err := redisListRange(m.NewUserId(), key, start, stop)
	if err != nil {
		return nil, err
	}

	objs := []*UserId{}
	for _, str := range strs {
		obj := m.NewUserId()
		obj.Key = key
		if err := redisStringScan(str, &obj.Value); err != nil {
			return nil, err
		}
		objs = append(objs, obj)
	}
	return objs, nil
}
func (m *_UserIdMgr) ListRangeRelatedUsers(key string, start, stop int64) ([]*User, error) {
	strs, err := redisListRange(m.NewUserId(), key, start, stop)
	if err != nil {
		return nil, err
	}

	objs := []*User{}
	for _, str := range strs {
		var val int32
		if err := redisStringScan(str, &val); err != nil {
			return nil, err
		}
		if obj, err := UserMgr.GetUserById(val); err == nil {
			objs = append(objs, obj)
		}
	}
	return objs, nil
}

func (m *_UserIdMgr) ListCount(key string) (int64, error) {
	return redisListCount(m.NewUserId(), key)
}

func (m *_UserIdMgr) ListDel(key string) error {
	return redisListDel(m.NewUserId(), key)
}
