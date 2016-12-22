package test

import (
	"time"
)

var (
	_ time.Time
)

///////////// LIST /////////////////////////////////////////////////////
func (m *_UserIdMgr) ListLPush(key string, obj *UserId) (int64, error) {
	return redisListLPush(obj, key, obj.Value)
}

func (m *_UserIdMgr) ListRPush(key string, obj *UserId) (int64, error) {
	return redisListRPush(obj, key, obj.Value)
}

func (m *_UserIdMgr) ListLPop(key string, obj *UserId) error {
	return redisListLPop(obj, key, &obj.Value)
}

func (m *_UserIdMgr) ListRPop(key string, obj *UserId) error {
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
