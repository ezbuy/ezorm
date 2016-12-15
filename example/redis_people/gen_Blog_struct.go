package test

import "time"

var _ time.Time

type Blog struct {
	BlogId      int32     `bson:"BlogId" json:"BlogId"`
	Title       string    `bson:"Title" json:"Title"`
	Hits        int32     `bson:"Hits" json:"Hits"`
	Slug        string    `bson:"Slug" json:"Slug"`
	Body        string    `bson:"Body" json:"Body"`
	User        int32     `bson:"User" json:"User"`
	IsPublished bool      `bson:"IsPublished" json:"IsPublished"`
	Create      time.Time `bson:"Create" json:"Create"`
	Update      time.Time `bson:"Update" json:"Update"`
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
