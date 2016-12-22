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

func (m *_BlogMgr) SetBlog(obj *Blog) error {
	//! object field set
	if err := redisFieldSet(obj, "BlogId", obj.BlogId); err != nil {
		return err
	}
	if err := redisFieldSet(obj, "Title", obj.Title); err != nil {
		return err
	}
	if err := redisFieldSet(obj, "Hits", obj.Hits); err != nil {
		return err
	}
	if err := redisFieldSet(obj, "Slug", obj.Slug); err != nil {
		return err
	}
	if err := redisFieldSet(obj, "Body", obj.Body); err != nil {
		return err
	}
	if err := redisFieldSet(obj, "User", obj.User); err != nil {
		return err
	}
	if err := redisFieldSet(obj, "IsPublished", obj.IsPublished); err != nil {
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
	if err := redisIndexSet(obj, "Slug", obj.Slug, obj.BlogId); err != nil {
		return err
	}
	if err := redisIndexSet(obj, "User", obj.User, obj.BlogId); err != nil {
		return err
	}
	if err := redisIndexSet(obj, "IsPublished", obj.IsPublished, obj.BlogId); err != nil {
		return err
	}
	if err := redisIndexSet(obj, "Create", obj.Create, obj.BlogId); err != nil {
		return err
	}
	if err := redisIndexSet(obj, "Update", obj.Update, obj.BlogId); err != nil {
		return err
	}
	//! object primary key set
	_, err := redisListLPush(obj, "BlogId", obj.BlogId)
	return err
}

func (m *_BlogMgr) DelBlog(obj *Blog) error {
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

func (m *_BlogMgr) GetBlog(obj *Blog) error {
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

func (m *_BlogMgr) GetBlogById(id int32) (*Blog, error) {
	obj := m.NewBlog()
	obj.BlogId = id
	if err := m.GetBlog(obj); err != nil {
		return nil, err
	}
	return obj, nil
}

func (m *_BlogMgr) GetBlogsByIds(ids []int32) ([]*Blog, error) {
	objs := []*Blog{}
	for _, id := range ids {
		obj, err := m.GetBlogById(id)
		if err != nil {
			return objs, err
		}
		objs = append(objs, obj)
	}
	return objs, nil
}

func (m *_BlogMgr) GetBlogsBySlug(val string) ([]*Blog, error) {
	strs, err := redisIndexGet(m.NewBlog(), "Slug", val)
	if err != nil {
		return nil, err
	}
	objs := []*Blog{}
	for _, str := range strs {
		var id int32
		if err := redisStringScan(str, &id); err != nil {
			return nil, err
		}

		obj, err := m.GetBlogById(id)
		if err != nil {
			return nil, err
		}

		objs = append(objs, obj)
	}
	return objs, nil
}

func (m *_BlogMgr) GetBlogsByUser(val int32) ([]*Blog, error) {
	strs, err := redisIndexGet(m.NewBlog(), "User", val)
	if err != nil {
		return nil, err
	}
	objs := []*Blog{}
	for _, str := range strs {
		var id int32
		if err := redisStringScan(str, &id); err != nil {
			return nil, err
		}

		obj, err := m.GetBlogById(id)
		if err != nil {
			return nil, err
		}

		objs = append(objs, obj)
	}
	return objs, nil
}

func (m *_BlogMgr) GetBlogsByIsPublished(val bool) ([]*Blog, error) {
	strs, err := redisIndexGet(m.NewBlog(), "IsPublished", val)
	if err != nil {
		return nil, err
	}
	objs := []*Blog{}
	for _, str := range strs {
		var id int32
		if err := redisStringScan(str, &id); err != nil {
			return nil, err
		}

		obj, err := m.GetBlogById(id)
		if err != nil {
			return nil, err
		}

		objs = append(objs, obj)
	}
	return objs, nil
}

func (m *_BlogMgr) GetBlogsByCreate(val time.Time) ([]*Blog, error) {
	strs, err := redisIndexGet(m.NewBlog(), "Create", val)
	if err != nil {
		return nil, err
	}
	objs := []*Blog{}
	for _, str := range strs {
		var id int32
		if err := redisStringScan(str, &id); err != nil {
			return nil, err
		}

		obj, err := m.GetBlogById(id)
		if err != nil {
			return nil, err
		}

		objs = append(objs, obj)
	}
	return objs, nil
}

func (m *_BlogMgr) GetBlogsByUpdate(val time.Time) ([]*Blog, error) {
	strs, err := redisIndexGet(m.NewBlog(), "Update", val)
	if err != nil {
		return nil, err
	}
	objs := []*Blog{}
	for _, str := range strs {
		var id int32
		if err := redisStringScan(str, &id); err != nil {
			return nil, err
		}

		obj, err := m.GetBlogById(id)
		if err != nil {
			return nil, err
		}

		objs = append(objs, obj)
	}
	return objs, nil
}

func (m *_BlogMgr) GetBlogsByIndexes(indexes map[string]interface{}) ([]*Blog, error) {
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

	objs := []*Blog{}
	for _, str := range strs {
		var id int32
		if err := redisStringScan(str, &id); err != nil {
			return nil, err
		}

		obj, err := m.GetBlogById(id)
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

		obj, err := m.GetBlogById(id)
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
