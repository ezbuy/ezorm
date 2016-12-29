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

func (m *_BlogMgr) Set(obj *Blog) error {
	//! object field set
	pipeline := redisPipeline()
	if err := pipeline.FieldSet(obj, "BlogId", obj.BlogId); err != nil {
		return err
	}
	if err := pipeline.FieldSet(obj, "Title", obj.Title); err != nil {
		return err
	}
	if err := pipeline.FieldSet(obj, "Hits", obj.Hits); err != nil {
		return err
	}
	if err := pipeline.FieldSet(obj, "Slug", obj.Slug); err != nil {
		return err
	}
	if err := pipeline.FieldSet(obj, "Body", obj.Body); err != nil {
		return err
	}
	if err := pipeline.FieldSet(obj, "User", obj.User); err != nil {
		return err
	}
	if err := pipeline.FieldSet(obj, "IsPublished", obj.IsPublished); err != nil {
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
	if err := pipeline.IndexSet(obj, "Slug", obj.Slug, obj.BlogId); err != nil {
		return err
	}
	if err := pipeline.IndexSet(obj, "User", obj.User, obj.BlogId); err != nil {
		return err
	}
	if err := pipeline.IndexSet(obj, "IsPublished", obj.IsPublished, obj.BlogId); err != nil {
		return err
	}
	if err := pipeline.IndexSet(obj, "Create", obj.Create, obj.BlogId); err != nil {
		return err
	}
	if err := pipeline.IndexSet(obj, "Update", obj.Update, obj.BlogId); err != nil {
		return err
	}
	//! object primary key set
	if err := pipeline.ListLPush(obj, "BlogId", obj.BlogId); err != nil {
		return err
	}
	_, err := pipeline.Exec()
	return err
}

func (m *_BlogMgr) Remove(obj *Blog) error {
	if err := redisDelObject(obj); err != nil {
		return err
	}
	if err := redisIndexRemove(obj, "Slug", obj.Slug, obj.BlogId); err != nil {
		return err
	}
	if err := redisIndexRemove(obj, "User", obj.User, obj.BlogId); err != nil {
		return err
	}
	if err := redisIndexRemove(obj, "IsPublished", obj.IsPublished, obj.BlogId); err != nil {
		return err
	}
	if err := redisIndexRemove(obj, "Create", obj.Create, obj.BlogId); err != nil {
		return err
	}
	if err := redisIndexRemove(obj, "Update", obj.Update, obj.BlogId); err != nil {
		return err
	}
	return redisListRemove(obj, "BlogId", obj.BlogId)
}

func (m *_BlogMgr) Get(obj *Blog) error {
	//! object field get
	if err := redisFieldGet(obj, "BlogId", &obj.BlogId); err != nil {
		return err
	}
	if err := redisFieldGet(obj, "Title", &obj.Title); err != nil {
		return err
	}
	if err := redisFieldGet(obj, "Hits", &obj.Hits); err != nil {
		return err
	}
	if err := redisFieldGet(obj, "Slug", &obj.Slug); err != nil {
		return err
	}
	if err := redisFieldGet(obj, "Body", &obj.Body); err != nil {
		return err
	}
	if err := redisFieldGet(obj, "User", &obj.User); err != nil {
		return err
	}
	if err := redisFieldGet(obj, "IsPublished", &obj.IsPublished); err != nil {
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

func (m *_BlogMgr) GetById(id int32) (*Blog, error) {
	obj := m.NewBlog()
	obj.BlogId = id
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

func (m *_BlogMgr) GetBySlug(val string) ([]*Blog, error) {
	strs, err := redisIndexGet(m.NewBlog(), "Slug", val)
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

func (m *_BlogMgr) GetByUser(val int32) ([]*Blog, error) {
	strs, err := redisIndexGet(m.NewBlog(), "User", val)
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

func (m *_BlogMgr) GetByIsPublished(val bool) ([]*Blog, error) {
	strs, err := redisIndexGet(m.NewBlog(), "IsPublished", val)
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

func (m *_BlogMgr) GetByCreate(val time.Time) ([]*Blog, error) {
	strs, err := redisIndexGet(m.NewBlog(), "Create", val)
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

func (m *_BlogMgr) GetByUpdate(val time.Time) ([]*Blog, error) {
	strs, err := redisIndexGet(m.NewBlog(), "Update", val)
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
	if val, ok := indexes["Slug"]; ok {
		index_keys = append(index_keys, fmt.Sprintf("Slug:%v", val))
	}
	if val, ok := indexes["User"]; ok {
		index_keys = append(index_keys, fmt.Sprintf("User:%v", val))
	}
	if val, ok := indexes["IsPublished"]; ok {
		index_keys = append(index_keys, fmt.Sprintf("IsPublished:%v", val))
	}
	if val, ok := indexes["Create"]; ok {
		transformed_Create_field := db.TimeFormat(val.(time.Time))
		index_keys = append(index_keys, fmt.Sprintf("Create:%v", transformed_Create_field))
	}
	if val, ok := indexes["Update"]; ok {
		transformed_Update_field := db.TimeToLocalTime(val.(time.Time))
		index_keys = append(index_keys, fmt.Sprintf("Update:%v", transformed_Update_field))
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
	strs, err := redisListRange(m.NewBlog(), "BlogId", start, stop)
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
	return redisListCount(m.NewBlog(), "BlogId")
}
