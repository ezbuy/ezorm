package model

import "time"

var _ time.Time

type UserDetail struct {
	UserId int32  `db:"user_id"`
	Desc   string `db:"desc"`
	Age    int32  `db:"age"`
	Phone  string `db:"phone"`
	Email  string `db:"email"`
	isNew  bool
}

const ()

func (p *UserDetail) GetNameSpace() string {
	return "model"
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
