package test

import "time"

var _ time.Time

type User struct {
	UserId     int64   `db:"user_id"`
	Name       string  `db:"name"`
	Phone      string  `db:"phone"`
	Age        int32   `db:"age"`
	Balance    float64 `db:"balance"`
	Text       string  `db:"text"`
	CreateDate int64   `db:"create_date"`
	isNew      bool
}

const (
	UserMysqlFieldUserId     = "user_id"
	UserMysqlFieldName       = "name"
	UserMysqlFieldPhone      = "phone"
	UserMysqlFieldAge        = "age"
	UserMysqlFieldBalance    = "balance"
	UserMysqlFieldText       = "text"
	UserMysqlFieldCreateDate = "create_date"
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
	return rval
}
