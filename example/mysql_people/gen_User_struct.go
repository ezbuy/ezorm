package test

import "time"

var _ time.Time

type User struct {
	UserId     int32  `db:"user_id" json:"user_id"`
	UserNumber int32  `db:"user_number" json:"user_number"`
	Name       string `db:"name" json:"name"`
	isNew      bool
}

func (p *User) GetNameSpace() string {
	return "people"
}

func (p *User) GetClassName() string {
	return "User"
}

func (p *User) GetIndexes() []string {
	idx := []string{
		"UserNumber",
		"Name",
	}
	return idx
}

type _UserMgr struct {
}

var UserMgr *_UserMgr

func (m *_UserMgr) NewUser() *User {
	rval := new(User)
	return rval
}
