package test

import "time"

var _ time.Time

type UserId struct {
	Key   string `db:"key" json:"key"`
	Value int32  `db:"value" json:"value"`
	isNew bool
}

func (p *UserId) GetNameSpace() string {
	return "people"
}

func (p *UserId) GetClassName() string {
	return "UserId"
}
func (p *UserId) GetStoreType() string {
	return "list"
}

func (p *UserId) GetPrimaryKey() string {
	return ""
}

func (p *UserId) GetIndexes() []string {
	idx := []string{}
	return idx
}

type _UserIdMgr struct {
}

var UserIdMgr *_UserIdMgr

func (m *_UserIdMgr) NewUserId() *UserId {
	rval := new(UserId)
	return rval
}
