package model

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

const ()

func (p *Blog) GetNameSpace() string {
	return "model"
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
