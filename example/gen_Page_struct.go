package example


import "gopkg.in/mgo.v2/bson"


type Page struct {
	ID         bson.ObjectId `bson:"_id,omitempty"`
	Title  string `bson:"Title"`
	Hits  int32 `bson:"Hits"`
	Slug  string `bson:"Slug"`
	Sections  []Section `bson:"Sections"`
	Meta  map[string][]map[string]int `bson:"Meta"`
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
