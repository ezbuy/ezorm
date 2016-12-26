package test

import "time"

var _ time.Time

type User struct {
	UserId     int32     `db:"UserId"	json:"UserId"`
	UserNumber int32     `db:"UserNumber"	json:"UserNumber"`
	Name       string    `db:"Name"	json:"Name"`
	Create     time.Time `db:"Create"	json:"Create"`
	Update     time.Time `db:"Update"	json:"Update"`
	isNew      bool
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
	idx := []string{
		"UserNumber",
		"Name",
		"Create",
		"Update",
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
