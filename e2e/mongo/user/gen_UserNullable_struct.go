package user

import (
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var _ time.Time
var _ json.Marshaler

// UserNullable user model for nullable field tests
type UserNullable struct {
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`

	// the identity of each user
	UserId uint64 `bson:"uid" json:"uid"`

	// name of user used to login
	Username     string     `bson:"Username" json:"Username"`
	Age          *int32     `bson:"Age,omitempty" json:"Age,omitempty"`
	Nickname     *string    `bson:"Nickname,omitempty" json:"Nickname,omitempty"`
	RegisterDate *time.Time `bson:"registerDate,omitempty" json:"registerDate,omitempty"`
	isNew        bool
}

const (
	UserNullableMgoFieldID           = "_id"
	UserNullableMgoFieldUserId       = "uid"
	UserNullableMgoFieldUsername     = "Username"
	UserNullableMgoFieldAge          = "Age"
	UserNullableMgoFieldNickname     = "Nickname"
	UserNullableMgoFieldRegisterDate = "registerDate"
)

// UserNullableMgoSortField_WRP is a wrapper of UserNullable sort fields e.g.:
// UserNullableMgoSortField_WRP{UserNullableMgoSortField_X_Asc, UserNullableMgoSortField_Y_DESC}
type UserNullableMgoSortField_WRP = primitive.D

var (
	UserNullableMgoSortFieldIDAsc  = primitive.E{Key: "_id", Value: 1}
	UserNullableMgoSortFieldIDDesc = primitive.E{Key: "_id", Value: -1}
)

func (p *UserNullable) GetNameSpace() string {
	return "mongo_e2e"
}

func (p *UserNullable) GetClassName() string {
	return "UserNullable"
}

type _UserNullableMgr struct {
}

var UserNullableMgr *_UserNullableMgr

// Get_UserNullableMgr returns the orm manager in case of its name starts with lower letter
func Get_UserNullableMgr() *_UserNullableMgr { return UserNullableMgr }

func (m *_UserNullableMgr) NewUserNullable() *UserNullable {
	rval := new(UserNullable)
	rval.isNew = true
	rval.ID = primitive.NewObjectID()

	return rval
}
