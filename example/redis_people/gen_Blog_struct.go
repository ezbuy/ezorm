package test

import "time"

var _ time.Time

type Blog struct {
	BlogId      int32     `db:"BlogId"	json:"BlogId"`
	Title       string    `db:"Title"	json:"Title"`
	Hits        int32     `db:"Hits"	json:"Hits"`
	Slug        string    `db:"Slug"	json:"Slug"`
	Body        string    `db:"Body"	json:"Body"`
	User        int32     `db:"User"	json:"User"`
	IsPublished bool      `db:"IsPublished"	json:"IsPublished"`
	Create      time.Time `db:"Create"	json:"Create"`
	Update      time.Time `db:"Update"	json:"Update"`
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
