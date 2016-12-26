package test

import "time"

var _ time.Time

type User struct {
	UserId     int32     `db:"user_id" json:"user_id"`
	UserNumber int32     `db:"user_number" json:"user_number"`
	Name       string    `db:"name" json:"name"`
	Create     time.Time `db:"create" json:"create"`
	Update     time.Time `db:"update" json:"update"`
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
