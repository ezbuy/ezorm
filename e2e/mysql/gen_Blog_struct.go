package mysql

import (
	"encoding/json"
	"time"
)

var _ time.Time
var _ json.Marshaler

type Blog struct {
	BlogId      int32     `db:"blog_id"`
	Title       string    `db:"title"`
	Hits        int32     `db:"hits"`
	Slug        string    `db:"slug"`
	Body        string    `db:"body"`
	User        int32     `db:"user"`
	IsPublished bool      `db:"is_published"`
	GroupId     int64     `db:"group_id"`
	Create      time.Time `db:"create"`
	Update      time.Time `db:"update"`
	isNew       bool
}

const (
	BlogMysqlFieldBlogId      = "blog_id"
	BlogMysqlFieldTitle       = "title"
	BlogMysqlFieldHits        = "hits"
	BlogMysqlFieldSlug        = "slug"
	BlogMysqlFieldBody        = "body"
	BlogMysqlFieldUser        = "user"
	BlogMysqlFieldIsPublished = "is_published"
	BlogMysqlFieldGroupId     = "group_id"
	BlogMysqlFieldCreate      = "create"
	BlogMysqlFieldUpdate      = "update"
)

func (p *Blog) GetNameSpace() string {
	return "mysql_e2e"
}

func (p *Blog) GetClassName() string {
	return "Blog"
}

type _BlogMgr struct {
}

var BlogMgr *_BlogMgr

// Get_BlogMgr returns the orm manager in case of its name starts with lower letter
func Get_BlogMgr() *_BlogMgr { return BlogMgr }

func (m *_BlogMgr) NewBlog() *Blog {
	rval := new(Blog)
	return rval
}
