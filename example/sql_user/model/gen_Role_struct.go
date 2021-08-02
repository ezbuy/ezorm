package model

import "time"

var _ time.Time

type Role struct {
	Id    int64  `db:"id"`
	Name  string `db:"name"`
	isNew bool
}

const ()

func (p *Role) GetNameSpace() string {
	return "model"
}

func (p *Role) GetClassName() string {
	return "Role"
}

type _RoleMgr struct {
}

var RoleMgr *_RoleMgr

// Get_RoleMgr returns the orm manager in case of its name starts with lower letter
func Get_RoleMgr() *_RoleMgr { return RoleMgr }

func (m *_RoleMgr) NewRole() *Role {
	rval := new(Role)
	return rval
}
