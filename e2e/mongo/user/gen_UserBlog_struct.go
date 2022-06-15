package user

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var _ time.Time

type UserBlog struct {
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`

	// the identity of each user
	UserId uint64 `bson:"uid" json:"uid"`

	// the identity of each blog
	BlogId  uint64 `bson:"bid" json:"bid"`
	Content string `bson:"Content" json:"Content"`
	isNew   bool
}

const (
	UserBlogMgoFieldID      = "_id"
	UserBlogMgoFieldUserId  = "uid"
	UserBlogMgoFieldBlogId  = "bid"
	UserBlogMgoFieldContent = "Content"
)

// UserBlogMgoSortField_WRP is a wrapper of UserBlog sort fields e.g.:
// UserBlogMgoSortField_WRP{UserBlogMgoSortField_X_Asc, UserBlogMgoSortField_Y_DESC}
type UserBlogMgoSortField_WRP = primitive.D

var (
	UserBlogMgoSortFieldIDAsc  = primitive.E{Key: "_id", Value: 1}
	UserBlogMgoSortFieldIDDesc = primitive.E{Key: "_id", Value: -1}
)

func (p *UserBlog) GetNameSpace() string {
	return "user"
}

func (p *UserBlog) GetClassName() string {
	return "UserBlog"
}

type _UserBlogMgr struct {
}

var UserBlogMgr *_UserBlogMgr

// Get_UserBlogMgr returns the orm manager in case of its name starts with lower letter
func Get_UserBlogMgr() *_UserBlogMgr { return UserBlogMgr }

func (m *_UserBlogMgr) NewUserBlog() *UserBlog {
	rval := new(UserBlog)
	rval.isNew = true
	rval.ID = primitive.NewObjectID()

	return rval
}
