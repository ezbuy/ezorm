package test

import "time"

var _ time.Time

type User struct {
	UserId     int32  `bson:"UserId" json:"UserId"`
	UserNumber int32  `bson:"UserNumber" json:"UserNumber"`
	Name       string `bson:"Name" json:"Name"`
}

func (p *User) GetNameSpace() string {
	return "people"
}

func (p *User) GetClassName() string {
	return "User"
}
func (p *User) GetStoreType() string {
	return "hash"
}

func (p *User) GetPrimaryKey() string {
	return "UserId"
}

func (p *User) GetIndexes() []string {
	idx := []string{}
	idx = append(idx, "UserNumber")
	idx = append(idx, "Name")
	return idx
}

type _UserMgr struct {
}

var UserMgr *_UserMgr

func (m *_UserMgr) NewUser() *User {
	rval := new(User)
	return rval
}
