package test

import "time"

var _ time.Time

type Blog struct {
	BlogId      int32     `db:"blog_id" json:"blog_id"`
	Title       string    `db:"title" json:"title"`
	Hits        int32     `db:"hits" json:"hits"`
	Slug        string    `db:"slug" json:"slug"`
	Body        string    `db:"body" json:"body"`
	User        int32     `db:"user" json:"user"`
	IsPublished bool      `db:"is_published" json:"is_published"`
	Create      time.Time `db:"create" json:"create"`
	Update      time.Time `db:"update" json:"update"`
	isNew       bool
}

func (p *Blog) GetNameSpace() string {
	return "people"
}

func (p *Blog) GetClassName() string {
	return "Blog"
}
func (p *Blog) GetStoreType() string {
	return "hash"
}

func (p *Blog) GetPrimaryKey() string {
	return "BlogId"
}

func (p *Blog) GetIndexes() []string {
	idx := []string{
		"Slug",
		"User",
		"IsPublished",
		"Create",
		"Update",
	}
	return idx
}

type _BlogMgr struct {
}

var BlogMgr *_BlogMgr

func (m *_BlogMgr) NewBlog() *Blog {
	rval := new(Blog)
	return rval
}
