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
func (m *_UserMgr) Import() error {
	return m.AddBySQL("SELECT `id`,`name`,`mailbox`,`sex`,`longitude`,`latitude`,`description`,`password`,`head_url`,`status`,`created_at`, `updated_at` FROM users")
}

func (m *_UserMgr) Set(obj *User) error {
	//! object field set
	pipeline := redisPipeline()
	if err := pipeline.FieldSet(obj, "Id", obj.Id); err != nil {
		return err
	}
	if err := pipeline.FieldSet(obj, "Name", obj.Name); err != nil {
		return err
	}
	if err := pipeline.FieldSet(obj, "Mailbox", obj.Mailbox); err != nil {
		return err
	}
	if err := pipeline.FieldSet(obj, "Sex", obj.Sex); err != nil {
		return err
	}
	if err := pipeline.FieldSet(obj, "Longitude", obj.Longitude); err != nil {
		return err
	}
	if err := pipeline.FieldSet(obj, "Latitude", obj.Latitude); err != nil {
		return err
	}
	if err := pipeline.FieldSet(obj, "Description", obj.Description); err != nil {
		return err
	}
	if err := pipeline.FieldSet(obj, "Password", obj.Password); err != nil {
		return err
	}
	if err := pipeline.FieldSet(obj, "HeadUrl", obj.HeadUrl); err != nil {
		return err
	}
	if err := pipeline.FieldSet(obj, "Status", obj.Status); err != nil {
		return err
	}
	transformed_CreatedAt_field := db.TimeFormat(obj.CreatedAt)
	if err := pipeline.FieldSet(obj, "CreatedAt", transformed_CreatedAt_field); err != nil {
		return err
	}
	transformed_UpdatedAt_field := db.TimeFormat(obj.UpdatedAt)
	if err := pipeline.FieldSet(obj, "UpdatedAt", transformed_UpdatedAt_field); err != nil {
		return err
	}
	//! object index set
	if err := pipeline.IndexSet(obj, "Name", obj.Name, obj.Id); err != nil {
		return err
	}
	//! object primary key set
	if err := pipeline.ListLPush(obj, "Id", obj.Id); err != nil {
		return err
	}
	_, err := pipeline.Exec()
	return err
}

func (m *_UserMgr) Remove(obj *User) error {
	if err := redisDelObject(obj); err != nil {
		return err
	}
	if err := redisIndexRemove(obj, "Name", obj.Name, obj.Id); err != nil {
		return err
	}
	return redisListRemove(obj, "Id", obj.Id)
}

func (m *_UserMgr) Clear() error {
	return redisDrop(m.NewUser())
}

func (m *_UserMgr) Get(obj *User) error {
	//! object field get
	if err := redisFieldGet(obj, "Id", &obj.Id); err != nil {
		return err
	}
	if err := redisFieldGet(obj, "Name", &obj.Name); err != nil {
		return err
	}
	if err := redisFieldGet(obj, "Mailbox", &obj.Mailbox); err != nil {
		return err
	}
	if err := redisFieldGet(obj, "Sex", &obj.Sex); err != nil {
		return err
	}
	if err := redisFieldGet(obj, "Longitude", &obj.Longitude); err != nil {
		return err
	}
	if err := redisFieldGet(obj, "Latitude", &obj.Latitude); err != nil {
		return err
	}
	if err := redisFieldGet(obj, "Description", &obj.Description); err != nil {
		return err
	}
	if err := redisFieldGet(obj, "Password", &obj.Password); err != nil {
		return err
	}
	if err := redisFieldGet(obj, "HeadUrl", &obj.HeadUrl); err != nil {
		return err
	}
	if err := redisFieldGet(obj, "Status", &obj.Status); err != nil {
		return err
	}
	var CreatedAt string
	if err := redisFieldGet(obj, "CreatedAt", &CreatedAt); err != nil {
		return err
	}
	obj.CreatedAt = db.TimeParse(CreatedAt)
	var UpdatedAt string
	if err := redisFieldGet(obj, "UpdatedAt", &UpdatedAt); err != nil {
		return err
	}
	obj.UpdatedAt = db.TimeParse(UpdatedAt)
	return nil
}

func (m *_UserMgr) GetById(id int32) (*User, error) {

	obj := m.NewUser()
	obj.Id = id
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

func (m *_UserMgr) GetByIndexes(indexes map[string]interface{}) ([]*User, error) {

	index_keys := []string{}
	if val, ok := indexes["Name"]; ok {
		index_keys = append(index_keys, fmt.Sprintf("Name:%v", val))
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

	strs, err := redisListRange(m.NewUser(), "Id", start, stop)
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
	return redisListCount(m.NewUser(), "Id")
}
