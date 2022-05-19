package user

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var _ time.Time

// User all registered user use our systems
type User struct {
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`

	// the identity of each user
	UserId uint64 `bson:"uid" json:"uid"`

	// name of user used to login
	Username     string    `bson:"Username" json:"Username"`
	Age          int32     `bson:"Age" json:"Age"`
	RegisterDate time.Time `bson:"registerDate" json:"registerDate"`
	isNew        bool
}

const (
	UserMgoFieldID           = "_id"
	UserMgoFieldUserId       = "uid"
	UserMgoFieldUsername     = "Username"
	UserMgoFieldAge          = "Age"
	UserMgoFieldRegisterDate = "registerDate"
)

// UserMgoSortField_WRP is a wrapper of User sort fields e.g.:
// UserMgoSortField_WRP{UserMgoSortField_X_Asc, UserMgoSortField_Y_DESC}
type UserMgoSortField_WRP = primitive.D

var (
	UserMgoSortFieldIDAsc   = primitive.E{Key: "_id", Value: 1}
	UserMgoSortFieldIDDesc  = primitive.E{Key: "_id", Value: -1}
	UserMgoSortFieldAgeAsc  = primitive.E{Key: "Age", Value: 1}
	UserMgoSortFieldAgeDesc = primitive.E{Key: "Age", Value: -1}
)

func (p *User) GetNameSpace() string {
	return "user"
}

func (p *User) GetClassName() string {
	return "User"
}

type _UserMgr struct {
}

var UserMgr *_UserMgr

// Get_UserMgr returns the orm manager in case of its name starts with lower letter
func Get_UserMgr() *_UserMgr { return UserMgr }

func (m *_UserMgr) NewUser() *User {
	rval := new(User)
	rval.isNew = true
	rval.ID = primitive.NewObjectID()

	return rval
}
