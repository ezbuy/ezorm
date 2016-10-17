package people

import "time"

var _ time.Time

type User struct {
	UserId     int32  `db:"UserId"`
	UserNumber int32  `db:"UserNumber"`
	Name       string `db:"Name"`
	isNew      bool
}

func (p *User) GetNameSpace() string {
	return "people"
}

func (p *User) GetClassName() string {
	return "User"
}

type _UserMgr struct {
}

var UserMgr *_UserMgr

func (m *_UserMgr) NewUser() *User {
	rval := new(User)
	return rval
}
