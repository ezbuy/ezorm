package test

import "time"

var _ time.Time

type Blog struct {
	BlogId      int32     `db:"BlogId"`
	Title       string    `db:"Title"`
	Hits        int32     `db:"Hits"`
	Slug        string    `db:"Slug"`
	Body        string    `db:"Body"`
	User        int32     `db:"User"`
	IsPublished bool      `db:"IsPublished"`
	Create      time.Time `db:"Create"`
	Update      time.Time `db:"Update"`
	isNew       bool
}

func (p *Blog) GetNameSpace() string {
	return "people"
}

func (p *Blog) GetClassName() string {
	return "Blog"
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
