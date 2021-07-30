package model

import "time"

var _ time.Time

type RoleUser struct {
	UserId int64 `db:"user_id"`
	RoleId int64 `db:"role_id"`
	isNew  bool
}

const ()

func (p *RoleUser) GetNameSpace() string {
	return "model"
}

func (p *RoleUser) GetClassName() string {
	return "RoleUser"
}

type _RoleUserMgr struct {
}

var RoleUserMgr *_RoleUserMgr

// Get_RoleUserMgr returns the orm manager in case of its name starts with lower letter
func Get_RoleUserMgr() *_RoleUserMgr { return RoleUserMgr }

func (m *_RoleUserMgr) NewRoleUser() *RoleUser {
	rval := new(RoleUser)
	return rval
}
