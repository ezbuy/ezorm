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

func (m *_UserMgr) AddBySQL(sql string, args ...interface{}) error {
	objs, err := m.Query(sql)
	if err != nil {
		return err
	}

	for _, obj := range objs {
		if err := m.Set(obj); err != nil {
			return err
		}
	}
	return nil
}

func (m *_UserMgr) DelBySQL(sql string, args ...interface{}) error {
	objs, err := m.Query(sql)
	if err != nil {
		return err
	}

	for _, obj := range objs {
		if err := m.Remove(obj); err != nil {
			return err
		}
	}
	return nil
}

func (m *_UserMgr) Set(obj *User) error {
	//! object field set
	pipeline := redisPipeline()
	if err := pipeline.FieldSet(obj, "UserId", obj.UserId); err != nil {
		return err
	}
	if err := pipeline.FieldSet(obj, "UserNumber", obj.UserNumber); err != nil {
		return err
	}
	if err := pipeline.FieldSet(obj, "Name", obj.Name); err != nil {
		return err
	}
	transformed_Create_field := db.TimeFormat(obj.Create)
	if err := pipeline.FieldSet(obj, "Create", transformed_Create_field); err != nil {
		return err
	}
	transformed_Update_field := db.TimeToLocalTime(obj.Update)
	if err := pipeline.FieldSet(obj, "Update", transformed_Update_field); err != nil {
		return err
	}
	//! object index set
	if err := pipeline.IndexSet(obj, "UserNumber", obj.UserNumber, obj.UserId); err != nil {
		return err
	}
	if err := pipeline.IndexSet(obj, "Name", obj.Name, obj.UserId); err != nil {
		return err
	}
	if err := pipeline.IndexSet(obj, "Create", obj.Create, obj.UserId); err != nil {
		return err
	}
	if err := pipeline.IndexSet(obj, "Update", obj.Update, obj.UserId); err != nil {
		return err
	}
	//! object primary key set
	if err := pipeline.ListLPush(obj, "UserId", obj.UserId); err != nil {
		return err
	}
	_, err := pipeline.Exec()
	return err
}

func (m *_UserMgr) Remove(obj *User) error {
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

func (m *_UserMgr) Get(obj *User) error {
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

func (m *_UserMgr) GetById(id int32) (*User, error) {
	obj := m.NewUser()
	obj.UserId = id
	if err := m.Get(obj); err != nil {
		return nil, err
	}
	return obj, nil
}

func (m *_UserMgr) GetByIds(ids []int32) ([]*User, error) {
	objs := make([]*User, 0, len(ids))
	for _, id := range ids {
		obj, err := m.GetById(id)
		if err != nil {
			return objs, err
		}
		objs = append(objs, obj)
	}
	return objs, nil
}

func (m *_UserMgr) GetByUserNumber(val int32) ([]*User, error) {
	strs, err := redisIndexGet(m.NewUser(), "UserNumber", val)
	if err != nil {
		return nil, err
	}
	objs := make([]*User, 0, len(strs))
	for _, str := range strs {
		var id int32
		if err := redisStringScan(str, &id); err != nil {
			return nil, err
		}

		obj, err := m.GetById(id)
		if err != nil {
			return nil, err
		}

		objs = append(objs, obj)
	}
	return objs, nil
}

func (m *_UserMgr) GetByName(val string) ([]*User, error) {
	strs, err := redisIndexGet(m.NewUser(), "Name", val)
	if err != nil {
		return nil, err
	}
	objs := make([]*User, 0, len(strs))
	for _, str := range strs {
		var id int32
		if err := redisStringScan(str, &id); err != nil {
			return nil, err
		}

		obj, err := m.GetById(id)
		if err != nil {
			return nil, err
		}

		objs = append(objs, obj)
	}
	return objs, nil
}

func (m *_UserMgr) GetByCreate(val time.Time) ([]*User, error) {
	strs, err := redisIndexGet(m.NewUser(), "Create", val)
	if err != nil {
		return nil, err
	}
	objs := make([]*User, 0, len(strs))
	for _, str := range strs {
		var id int32
		if err := redisStringScan(str, &id); err != nil {
			return nil, err
		}

		obj, err := m.GetById(id)
		if err != nil {
			return nil, err
		}

		objs = append(objs, obj)
	}
	return objs, nil
}

func (m *_UserMgr) GetByUpdate(val time.Time) ([]*User, error) {
	strs, err := redisIndexGet(m.NewUser(), "Update", val)
	if err != nil {
		return nil, err
	}
	objs := make([]*User, 0, len(strs))
	for _, str := range strs {
		var id int32
		if err := redisStringScan(str, &id); err != nil {
			return nil, err
		}

		obj, err := m.GetById(id)
		if err != nil {
			return nil, err
		}

		objs = append(objs, obj)
	}
	return objs, nil
}

func (m *_UserMgr) GetByIndexes(indexes map[string]interface{}) ([]*User, error) {
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

	objs := make([]*User, 0, len(strs))
	for _, str := range strs {
		var id int32
		if err := redisStringScan(str, &id); err != nil {
			return nil, err
		}

		obj, err := m.GetById(id)
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

		obj, err := m.GetById(id)
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
