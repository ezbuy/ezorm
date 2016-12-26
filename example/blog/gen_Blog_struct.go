package test

import "gopkg.in/mgo.v2/bson"

import "time"

var _ time.Time

type Blog struct {
	ID          bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Title       string        `bson:"title"	json:"title"`
	Hits        int32         `bson:"hits"	json:"hits"`
	Slug        string        `bson:"slug"	json:"slug"`
	Body        string        `bson:"body"	json:"body"`
	User        int32         `bson:"user"	json:"user"`
	IsPublished bool          `bson:"is_published"	json:"is_published"`
	isNew       bool
}

func (p *Blog) GetNameSpace() string {
	return "blog"
}

func (p *Blog) GetClassName() string {
	return "Blog"
}

func (p *Blog) GetIndexes() []string {
	idx := []string{
		"Slug",
		"IsPublished",
	}
	return idx
}

type _BlogMgr struct {
}

var BlogMgr *_BlogMgr

func (m *_BlogMgr) NewBlog() *Blog {
	rval := new(Blog)
	rval.isNew = true
	rval.ID = bson.NewObjectId()

	return rval
}
