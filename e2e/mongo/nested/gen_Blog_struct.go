package nested

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var _ time.Time

type Blog struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	BlogId uint64             `bson:"bid" json:"bid"`
	isNew  bool
}

const (
	BlogMgoFieldID     = "_id"
	BlogMgoFieldBlogId = "bid"
)

// BlogMgoSortField_WRP is a wrapper of Blog sort fields e.g.:
// BlogMgoSortField_WRP{BlogMgoSortField_X_Asc, BlogMgoSortField_Y_DESC}
type BlogMgoSortField_WRP = primitive.D

var (
	BlogMgoSortFieldIDAsc  = primitive.E{Key: "_id", Value: 1}
	BlogMgoSortFieldIDDesc = primitive.E{Key: "_id", Value: -1}
)

func (p *Blog) GetNameSpace() string {
	return "nested"
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
	rval.ID = primitive.NewObjectID()

	return rval
}
