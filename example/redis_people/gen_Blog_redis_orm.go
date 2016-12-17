package test

import (
	"time"

	"github.com/ezbuy/ezorm/db"
)

func (m *_BlogMgr) SetBlog(obj *Blog) error {
	return redisSetObject(obj)
}

func (m *_BlogMgr) DelBlog(obj *Blog) error {
	return redisDelObject(obj)
}

///////////// JSON & HASH /////////////////////////////////////////////////////
func (m *_BlogMgr) GetBlog(obj *Blog) error {
	return redisGetObject(obj)
}

func (m *_BlogMgr) GetBlogById(obj *Blog, id string) error {
	return redisGetObjectById(obj, id)
}

func (m *_BlogMgr) GetBlogsByIds(ids []string) ([]*Blog, error) {
	objs := []*Blog{}
	for _, id := range ids {
		obj := m.NewBlog()
		if err := redisGetObjectById(obj, id); err != nil {
			return objs, err
		}
		objs = append(objs, obj)
	}
	return objs, nil
}

func (m *_BlogMgr) GetBlogsBySlug(val string) ([]*Blog, error) {
	obj := m.NewBlog()
	obj.Slug = val

	key_of_index, err := db.KeyOfIndexByObject(obj, "Slug")
	if err != nil {
		return nil, err
	}

	ids, err := redisSMEMBERIds(key_of_index)
	if err != nil {
		return nil, err
	}
	return m.GetBlogsByIds(ids)
}

func (m *_BlogMgr) GetBlogsByUser(val int32) ([]*Blog, error) {
	obj := m.NewBlog()
	obj.User = val

	key_of_index, err := db.KeyOfIndexByObject(obj, "User")
	if err != nil {
		return nil, err
	}

	ids, err := redisSMEMBERIds(key_of_index)
	if err != nil {
		return nil, err
	}
	return m.GetBlogsByIds(ids)
}

func (m *_BlogMgr) GetBlogsByIsPublished(val bool) ([]*Blog, error) {
	obj := m.NewBlog()
	obj.IsPublished = val

	key_of_index, err := db.KeyOfIndexByObject(obj, "IsPublished")
	if err != nil {
		return nil, err
	}

	ids, err := redisSMEMBERIds(key_of_index)
	if err != nil {
		return nil, err
	}
	return m.GetBlogsByIds(ids)
}

func (m *_BlogMgr) GetBlogsByCreate(val time.Time) ([]*Blog, error) {
	obj := m.NewBlog()
	obj.Create = val

	key_of_index, err := db.KeyOfIndexByObject(obj, "Create")
	if err != nil {
		return nil, err
	}

	ids, err := redisSMEMBERIds(key_of_index)
	if err != nil {
		return nil, err
	}
	return m.GetBlogsByIds(ids)
}

func (m *_BlogMgr) GetBlogsByUpdate(val time.Time) ([]*Blog, error) {
	obj := m.NewBlog()
	obj.Update = val

	key_of_index, err := db.KeyOfIndexByObject(obj, "Update")
	if err != nil {
		return nil, err
	}

	ids, err := redisSMEMBERIds(key_of_index)
	if err != nil {
		return nil, err
	}
	return m.GetBlogsByIds(ids)
}

func (m *_BlogMgr) GetBlogsByIndexes(indexes map[string]interface{}) ([]*Blog, error) {
	obj := m.NewBlog()

	index_keys := []string{}
	for k, v := range indexes {
		if idx, err := db.KeyOfIndexByClass(obj.GetClassName(), k, v); err == nil {
			index_keys = append(index_keys, idx)
		}
	}

	ids, err := redisSINTERIds(index_keys...)
	if err != nil {
		return nil, err
	}
	return m.GetBlogsByIds(ids)
}
