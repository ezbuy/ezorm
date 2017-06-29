package test

import "time"

var _ time.Time

type BlogTemp struct {
	BlogId      int32  `db:"blog_id"`
	Title       string `db:"title"`
	Hits        int32  `db:"hits"`
	Slug        string `db:"slug"`
	Body        string `db:"body"`
	User        int32  `db:"user"`
	GhostNumber int32  `db:"ghost_number"`
	isNew       bool
}

func (p *BlogTemp) GetNameSpace() string {
	return "people"
}

func (p *BlogTemp) GetClassName() string {
	return "BlogTemp"
}

type _BlogTempMgr struct {
}

var BlogTempMgr *_BlogTempMgr

func (m *_BlogTempMgr) NewBlogTemp() *BlogTemp {
	rval := new(BlogTemp)
	return rval
}
