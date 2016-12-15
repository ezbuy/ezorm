package test

import "time"

var _ time.Time

type Blog struct {
	BlogId      int32     `db:"blog_id"`
	Title       string    `db:"title"`
	Hits        int32     `db:"hits"`
	Slug        string    `db:"slug"`
	Body        string    `db:"body"`
	User        int32     `db:"user"`
	IsPublished bool      `db:"is_published"`
	Create      time.Time `db:"create"`
	Update      time.Time `db:"update"`
	isNew       bool
}

func (p *Blog) GetNameSpace() string {
	return "people"
}

func (p *Blog) GetClassName() string {
	return "Blog"
}

func (p *Blog) GetIndexes() []string {
	idx := []string{}
	idx = append(idx, "Slug")
	idx = append(idx, "User")
	idx = append(idx, "IsPublished")
	idx = append(idx, "Create")
	idx = append(idx, "Update")
	return idx
}

type _BlogMgr struct {
}

var BlogMgr *_BlogMgr

func (m *_BlogMgr) NewBlog() *Blog {
	rval := new(Blog)
	return rval
}
