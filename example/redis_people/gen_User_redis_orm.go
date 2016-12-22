package test

///////////// JSON & HASH /////////////////////////////////////////////////////
import (
	"fmt"
	"github.com/ezbuy/ezorm/db"
	"time"
)

var (
	_ time.Time
)

func (m *_UserMgr) SetUser(obj *User) error {
	//! object field set
	if err := redisFieldSet(obj, "UserId", obj.UserId); err != nil {
		return err
	}
	if err := redisFieldSet(obj, "UserNumber", obj.UserNumber); err != nil {
		return err
	}
	if err := redisFieldSet(obj, "Name", obj.Name); err != nil {
		return err
	}
	transformed_Create_field := db.TimeFormat(obj.Create)
	if err := redisFieldSet(obj, "Create", transformed_Create_field); err != nil {
		return err
	}
	transformed_Update_field := db.TimeToLocalTime(obj.Update)
	if err := redisFieldSet(obj, "Update", transformed_Update_field); err != nil {
		return err
	}
	//! object index set
	if err := redisIndexSet(obj, "UserNumber", obj.UserNumber, obj.UserId); err != nil {
		return err
	}
	if err := redisIndexSet(obj, "Name", obj.Name, obj.UserId); err != nil {
		return err
	}
	if err := redisIndexSet(obj, "Create", obj.Create, obj.UserId); err != nil {
		return err
	}
	if err := redisIndexSet(obj, "Update", obj.Update, obj.UserId); err != nil {
		return err
	}
	//! object primary key set
	_, err := redisListLPush(obj, "UserId", obj.UserId)
	return err
}

func (m *_UserMgr) DelUser(obj *User) error {
	if err := redisDelObject(obj); err != nil {
		return err
	}
	if err := redisIndexRemove(obj, "UserNumber", obj.UserNumber, obj.UserId); err != nil {
		return err
	}
	if err := redisIndexRemove(obj, "Name", obj.Name, obj.UserId); err != nil {
		return err
	}
	if err := redisIndexRemove(obj, "Create", obj.Create, obj.UserId); err != nil {
		return err
	}
	if err := redisIndexRemove(obj, "Update", obj.Update, obj.UserId); err != nil {
		return err
	}
	return redisListRemove(obj, "UserId", obj.UserId)
}

func (m *_UserMgr) GetUser(obj *User) error {
	//! object field get
	if err := redisFieldGet(obj, "UserId", &obj.UserId); err != nil {
		return err
	}
	if err := redisFieldGet(obj, "UserNumber", &obj.UserNumber); err != nil {
		return err
	}
	if err := redisFieldGet(obj, "Name", &obj.Name); err != nil {
		return err
	}
	var Create string
	if err := redisFieldGet(obj, "Create", &Create); err != nil {
		return err
	}
	obj.Create = db.TimeParse(Create)
	var Update string
	if err := redisFieldGet(obj, "Update", &Update); err != nil {
		return err
	}
	obj.Update = db.TimeParseLocalTime(Update)
	return nil
}

func (m *_UserMgr) GetUserById(id int32) (*User, error) {
	obj := m.NewUser()
	obj.UserId = id
	if err := m.GetUser(obj); err != nil {
		return nil, err
	}
	return obj, nil
}

func (m *_UserMgr) GetUsersByIds(ids []int32) ([]*User, error) {
	objs := []*User{}
	for _, id := range ids {
		obj, err := m.GetUserById(id)
		if err != nil {
			return objs, err
		}
		objs = append(objs, obj)
	}
	return objs, nil
}

func (m *_UserMgr) GetUsersByUserNumber(val int32) ([]*User, error) {
	strs, err := redisIndexGet(m.NewUser(), "UserNumber", val)
	if err != nil {
		return nil, err
	}
	objs := []*User{}
	for _, str := range strs {
		var id int32
		if err := redisStringScan(str, &id); err != nil {
			return nil, err
		}

		obj, err := m.GetUserById(id)
		if err != nil {
			return nil, err
		}

		objs = append(objs, obj)
	}
	return objs, nil
}

func (m *_UserMgr) GetUsersByName(val string) ([]*User, error) {
	strs, err := redisIndexGet(m.NewUser(), "Name", val)
	if err != nil {
		return nil, err
	}
	objs := []*User{}
	for _, str := range strs {
		var id int32
		if err := redisStringScan(str, &id); err != nil {
			return nil, err
		}

		obj, err := m.GetUserById(id)
		if err != nil {
			return nil, err
		}

		objs = append(objs, obj)
	}
	return objs, nil
}

func (m *_UserMgr) GetUsersByCreate(val time.Time) ([]*User, error) {
	strs, err := redisIndexGet(m.NewUser(), "Create", val)
	if err != nil {
		return nil, err
	}
	objs := []*User{}
	for _, str := range strs {
		var id int32
		if err := redisStringScan(str, &id); err != nil {
			return nil, err
		}

		obj, err := m.GetUserById(id)
		if err != nil {
			return nil, err
		}

		objs = append(objs, obj)
	}
	return objs, nil
}

func (m *_UserMgr) GetUsersByUpdate(val time.Time) ([]*User, error) {
	strs, err := redisIndexGet(m.NewUser(), "Update", val)
	if err != nil {
		return nil, err
	}
	objs := []*User{}
	for _, str := range strs {
		var id int32
		if err := redisStringScan(str, &id); err != nil {
			return nil, err
		}

		obj, err := m.GetUserById(id)
		if err != nil {
			return nil, err
		}

		objs = append(objs, obj)
	}
	return objs, nil
}

func (m *_UserMgr) GetUsersByIndexes(indexes map[string]interface{}) ([]*User, error) {
	index_keys := []string{}
	if val, ok := indexes["UserNumber"]; ok {
		index_keys = append(index_keys, fmt.Sprintf("UserNumber:%v", val))
	}
	if val, ok := indexes["Name"]; ok {
		index_keys = append(index_keys, fmt.Sprintf("Name:%v", val))
	}
	if val, ok := indexes["Create"]; ok {
		transformed_Create_field := db.TimeFormat(val.(time.Time))
		index_keys = append(index_keys, fmt.Sprintf("Create:%v", transformed_Create_field))
	}
	if val, ok := indexes["Update"]; ok {
		transformed_Update_field := db.TimeToLocalTime(val.(time.Time))
		index_keys = append(index_keys, fmt.Sprintf("Update:%v", transformed_Update_field))
	}

	strs, err := redisMultiIndexesGet(m.NewUser(), index_keys...)
	if err != nil {
		return nil, err
	}

	objs := []*User{}
	for _, str := range strs {
		var id int32
		if err := redisStringScan(str, &id); err != nil {
			return nil, err
		}

		obj, err := m.GetUserById(id)
		if err != nil {
			return nil, err
		}

		objs = append(objs, obj)
	}
	return objs, nil
}

func (m *_UserMgr) ListRange(start, stop int64) ([]*User, error) {
	strs, err := redisListRange(m.NewUser(), "UserId", start, stop)
	if err != nil {
		return nil, err
	}

	objs := []*User{}
	for _, str := range strs {
		var id int32
		if err := redisStringScan(str, &id); err != nil {
			return nil, err
		}

		obj, err := m.GetUserById(id)
		if err != nil {
			return nil, err
		}

		objs = append(objs, obj)
	}
	return objs, nil
}

func (m *_UserMgr) ListCount() (int64, error) {
	return redisListCount(m.NewUser(), "UserId")
}
