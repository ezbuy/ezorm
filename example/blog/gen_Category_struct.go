package test

import "time"

var _ time.Time

type Category struct {
	Cid  int64  `bson:"Cid" json:"Cid"`
	Name string `bson:"Name" json:"Name"`
}

const ()

func (p *Category) GetNameSpace() string {
	return "blog"
}

func (p *Category) GetClassName() string {
	return "Category"
}

type _CategoryMgr struct {
}

var CategoryMgr *_CategoryMgr

func (m *_CategoryMgr) NewCategory() *Category {
	rval := new(Category)
	return rval
}
