package model

import "time"

var _ time.Time

type UserDetail struct {
	Id           int64  `db:"id"`
	UserId       int64  `db:"user_id"`
	Email        string `db:"email"`
	Introduction string `db:"introduction"`
	Age          int32  `db:"age"`
	Avatar       string `db:"avatar"`
	isNew        bool
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
