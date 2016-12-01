package example

import "gopkg.in/mgo.v2/bson"

import "time"

var _ time.Time

type Page struct {
	ID bson.ObjectId `json:"id" bson:"_id,omitempty"`

	Title string `bson:"Title" json:"Title"`

	Hits int32 `bson:"Hits" json:"Hits"`

	Slug string `bson:"Slug" json:"Slug"`

	Sections []Section `bson:"Sections" json:"Sections"`

	Meta  map[string][]map[string]int `bson:"Meta" json:"Meta"`
	isNew bool
}

func (p *Page) GetNameSpace() string {
	return "example"
}

func (p *Page) GetClassName() string {
	return "Page"
}

type _PageMgr struct {
}

var PageMgr *_PageMgr

func (m *_PageMgr) NewPage() *Page {
	rval := new(Page)
	rval.isNew = true
	rval.ID = bson.NewObjectId()

	return rval
}
