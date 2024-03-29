package mysql

import "time"

var _ time.Time

type User struct {
	UserId     int32  `db:"user_id"`
	UserNumber int32  `db:"user_number"`
	Name       string `db:"name"`
	isNew      bool
}

const (
	UserMysqlFieldUserId     = "user_id"
	UserMysqlFieldUserNumber = "user_number"
	UserMysqlFieldName       = "name"
)

func (p *User) GetNameSpace() string {
	return "mysql_e2e"
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
	return rval
}
