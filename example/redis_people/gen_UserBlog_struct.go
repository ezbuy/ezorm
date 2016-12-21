package test

import "time"

var _ time.Time

type UserBlog struct {
	Value int32
}

func (p *UserBlog) GetNameSpace() string {
	return "people"
}

func (p *UserBlog) GetClassName() string {
	return "UserBlog"
}
func (p *UserBlog) GetStoreType() string {
	return "set"
}

func (p *UserBlog) GetPrimaryKey() string {
	return ""
}

func (p *UserBlog) GetIndexes() []string {
	idx := []string{}
	return idx
}

type _UserBlogMgr struct {
}

var UserBlogMgr *_UserBlogMgr

func (m *_UserBlogMgr) NewUserBlog() *UserBlog {
	rval := new(UserBlog)
	return rval
}
