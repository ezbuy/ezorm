package test

import (
	"time"

	"github.com/ezbuy/ezorm/db"
)

func (m *_UserMgr) SetUser(obj *User) error {
	return redisSetObject(obj)
}

func (m *_UserMgr) DelUser(obj *User) error {
	return redisDelObject(obj)
}

///////////// JSON & HASH /////////////////////////////////////////////////////
func (m *_UserMgr) GetUser(obj *User) error {
	return redisGetObject(obj)
}

func (m *_UserMgr) GetUsersByIds(ids []int32) ([]*User, error) {
	objs := []*User{}
	for _, id := range ids {
		obj := m.NewUser()
		obj.UserId = id
		if err := redisGetObject(obj); err != nil {
			return objs, err
		}
		objs = append(objs, obj)
	}
	return objs, nil
}

func (m *_UserMgr) GetUsersByUserNumber(val int32) ([]*User, error) {
	obj := m.NewUser()
	obj.UserNumber = val

	key_of_index, err := db.KeyOfIndexByObject(obj, "UserNumber")
	if err != nil {
		return nil, err
	}

	ids, err := redisSMEMBERSInts(key_of_index)
	if err != nil {
		return nil, err
	}

	keys := []int32{}
	for _, id := range ids {
		keys = append(keys, int32(id))
	}

	return m.GetUsersByIds(keys)
}

func (m *_UserMgr) GetUsersByName(val string) ([]*User, error) {
	obj := m.NewUser()
	obj.Name = val

	key_of_index, err := db.KeyOfIndexByObject(obj, "Name")
	if err != nil {
		return nil, err
	}

	ids, err := redisSMEMBERSInts(key_of_index)
	if err != nil {
		return nil, err
	}

	keys := []int32{}
	for _, id := range ids {
		keys = append(keys, int32(id))
	}

	return m.GetUsersByIds(keys)
}

func (m *_UserMgr) GetUsersByCreate(val time.Time) ([]*User, error) {
	obj := m.NewUser()
	obj.Create = val

	key_of_index, err := db.KeyOfIndexByObject(obj, "Create")
	if err != nil {
		return nil, err
	}

	ids, err := redisSMEMBERSInts(key_of_index)
	if err != nil {
		return nil, err
	}

	keys := []int32{}
	for _, id := range ids {
		keys = append(keys, int32(id))
	}

	return m.GetUsersByIds(keys)
}

func (m *_UserMgr) GetUsersByUpdate(val time.Time) ([]*User, error) {
	obj := m.NewUser()
	obj.Update = val

	key_of_index, err := db.KeyOfIndexByObject(obj, "Update")
	if err != nil {
		return nil, err
	}

	ids, err := redisSMEMBERSInts(key_of_index)
	if err != nil {
		return nil, err
	}

	keys := []int32{}
	for _, id := range ids {
		keys = append(keys, int32(id))
	}

	return m.GetUsersByIds(keys)
}

func (m *_UserMgr) GetUsersByIndexes(indexes map[string]interface{}) ([]*User, error) {
	obj := m.NewUser()

	index_keys := []interface{}{}
	for k, v := range indexes {
		if idx, err := db.KeyOfIndexByClass(obj.GetClassName(), k, v); err == nil {
			index_keys = append(index_keys, idx)
		}
	}

	ids, err := redisSINTERInts(index_keys...)
	if err != nil {
		return nil, err
	}

	keys := []int32{}
	for _, id := range ids {
		keys = append(keys, int32(id))
	}
	return m.GetUsersByIds(keys)
}
