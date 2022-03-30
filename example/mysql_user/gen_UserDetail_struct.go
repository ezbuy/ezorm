package test

import "time"

var _ time.Time

type UserDetail struct {
	Id      int64  `db:"id"`
	UserId  int64  `db:"user_id"`
	Score   int32  `db:"score"`
	Balance int32  `db:"balance"`
	Text    string `db:"text"`
	isNew   bool
}

const (
	UserDetailMysqlFieldId      = "id"
	UserDetailMysqlFieldUserId  = "user_id"
	UserDetailMysqlFieldScore   = "score"
	UserDetailMysqlFieldBalance = "balance"
	UserDetailMysqlFieldText    = "text"
)

func (p *UserDetail) GetNameSpace() string {
	return "user"
}

func (p *UserDetail) GetClassName() string {
	return "UserDetail"
}

type _UserDetailMgr struct {
}

var UserDetailMgr *_UserDetailMgr

// Get_UserDetailMgr returns the orm manager in case of its name starts with lower letter
func Get_UserDetailMgr() *_UserDetailMgr { return UserDetailMgr }

func (m *_UserDetailMgr) NewUserDetail() *UserDetail {
	rval := new(UserDetail)
	return rval
}
