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

func (m *_BlogMgr) AddBySQL(sql string, args ...interface{}) error {
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

func (m *_BlogMgr) DelBySQL(sql string, args ...interface{}) error {
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
func (m *_BlogMgr) Import() error {
	return m.AddBySQL("SELECT `id`,`user_id`,`title`,`content`,`status`,`readed`, `created_at`, `updated_at` FROM blogs")
}

func (m *_BlogMgr) Set(obj *Blog) error {
	//! object field set
	pipeline := redisPipeline()
	if err := pipeline.FieldSet(obj, "Id", obj.Id); err != nil {
		return err
	}
	if err := pipeline.FieldSet(obj, "UserId", obj.UserId); err != nil {
		return err
	}
	if err := pipeline.FieldSet(obj, "Title", obj.Title); err != nil {
		return err
	}
	if err := pipeline.FieldSet(obj, "Content", obj.Content); err != nil {
		return err
	}
	if err := pipeline.FieldSet(obj, "Status", obj.Status); err != nil {
		return err
	}
	if err := pipeline.FieldSet(obj, "Readed", obj.Readed); err != nil {
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
	if err := pipeline.IndexSet(obj, "UserId", obj.UserId, obj.Id); err != nil {
		return err
	}
	//! object primary key set
	if err := pipeline.ListLPush(obj, "Id", obj.Id); err != nil {
		return err
	}
	_, err := pipeline.Exec()
	return err
}

func (m *_BlogMgr) Remove(obj *Blog) error {
	if err := redisDelObject(obj); err != nil {
		return err
	}
	if err := redisIndexRemove(obj, "UserId", obj.UserId, obj.Id); err != nil {
		return err
	}
	return redisListRemove(obj, "Id", obj.Id)
}

func (m *_BlogMgr) Clear() error {
	return redisDrop(m.NewBlog())
}

func (m *_BlogMgr) Get(obj *Blog) error {
	//! object field get
	if err := redisFieldGet(obj, "Id", &obj.Id); err != nil {
		return err
	}
	if err := redisFieldGet(obj, "UserId", &obj.UserId); err != nil {
		return err
	}
	if err := redisFieldGet(obj, "Title", &obj.Title); err != nil {
		return err
	}
	if err := redisFieldGet(obj, "Content", &obj.Content); err != nil {
		return err
	}
	if err := redisFieldGet(obj, "Status", &obj.Status); err != nil {
		return err
	}
	if err := redisFieldGet(obj, "Readed", &obj.Readed); err != nil {
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

func (m *_BlogMgr) GetById(id int32) (*Blog, error) {

	obj := m.NewBlog()
	obj.Id = id
	if err := m.Get(obj); err != nil {
		return nil, err
	}
	return obj, nil
}

func (m *_BlogMgr) GetByIds(ids []int32) ([]*Blog, error) {

	objs := make([]*Blog, 0, len(ids))
	for _, id := range ids {
		obj, err := m.GetById(id)
		if err != nil {
			return objs, err
		}
		objs = append(objs, obj)
	}
	return objs, nil
}

func (m *_BlogMgr) GetByUserId(val int32) ([]*Blog, error) {

	strs, err := redisIndexGet(m.NewBlog(), "UserId", val)
	if err != nil {
		return nil, err
	}
	objs := make([]*Blog, 0, len(strs))
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

func (m *_BlogMgr) GetByIndexes(indexes map[string]interface{}) ([]*Blog, error) {

	index_keys := []string{}
	if val, ok := indexes["UserId"]; ok {
		index_keys = append(index_keys, fmt.Sprintf("UserId:%v", val))
	}

	strs, err := redisMultiIndexesGet(m.NewBlog(), index_keys...)
	if err != nil {
		return nil, err
	}

	objs := make([]*Blog, 0, len(strs))
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

func (m *_BlogMgr) ListRange(start, stop int64) ([]*Blog, error) {

	strs, err := redisListRange(m.NewBlog(), "Id", start, stop)
	if err != nil {
		return nil, err
	}

	objs := []*Blog{}
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

func (m *_BlogMgr) ListCount() (int64, error) {
	return redisListCount(m.NewBlog(), "Id")
}
