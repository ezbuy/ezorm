package test

import "gopkg.in/mgo.v2/bson"

import "time"

var _ time.Time

type Blog struct {
	ID bson.ObjectId `json:"id" bson:"_id,omitempty"`
	// 标题
	Title string `bson:"title" json:"title"`
	// 提示
	Hits        int32     `bson:"Hits" json:"Hits"`
	Slug        string    `bson:"Slug" json:"Slug"`
	Body        string    `bson:"Body" json:"Body"`
	CreateDate  time.Time `bson:"createDate" json:"createDate"`
	User        int32     `bson:"User" json:"User"`
	IsPublished bool      `bson:"IsPublished" json:"IsPublished"`
	isNew       bool
}

const (
	BlogMgoFieldID          = "_id"
	BlogMgoFieldTitle       = "title"
	BlogMgoFieldHits        = "Hits"
	BlogMgoFieldSlug        = "Slug"
	BlogMgoFieldBody        = "Body"
	BlogMgoFieldCreateDate  = "createDate"
	BlogMgoFieldUser        = "User"
	BlogMgoFieldIsPublished = "IsPublished"
)
const (
	BlogMgoSortFieldIDAsc          = "_id"
	BlogMgoSortFieldIDDesc         = "-_id"
	BlogMgoSortFieldCreateDateAsc  = "createDate"
	BlogMgoSortFieldCreateDateDesc = "-createDate"
)

func (p *Blog) GetNameSpace() string {
	return "blog"
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
	rval.isNew = true
	rval.ID = bson.NewObjectId()

	return rval
}
