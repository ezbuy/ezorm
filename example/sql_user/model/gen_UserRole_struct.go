package model

import "time"

var _ time.Time

type UserRole struct {
	Id     int64 `db:"id"`
	UserId int64 `db:"user_id"`
	RoleId int64 `db:"role_id"`
	isNew  bool
}

const ()

func (p *UserRole) GetNameSpace() string {
	return "model"
}

func (p *UserRole) GetClassName() string {
	return "UserRole"
}

type _UserRoleMgr struct {
}

var UserRoleMgr *_UserRoleMgr

// Get_UserRoleMgr returns the orm manager in case of its name starts with lower letter
func Get_UserRoleMgr() *_UserRoleMgr { return UserRoleMgr }

func (m *_UserRoleMgr) NewUserRole() *UserRole {
	rval := new(UserRole)
	return rval
}
