package test

import "time"

var _ time.Time

type SortUserBlog struct {
	Key   string  `db:"Key"	json:"Key"`
	Score float64 `db:"Score"	json:"Score"`
	Value int32   `db:"Value"	json:"Value"`
	isNew bool
}

func (p *SortUserBlog) GetNameSpace() string {
	return "people"
}

func (p *SortUserBlog) GetClassName() string {
	return "SortUserBlog"
}
func (p *SortUserBlog) GetStoreType() string {
	return "zset"
}

func (p *SortUserBlog) GetPrimaryKey() string {
	return ""
}

func (p *SortUserBlog) GetIndexes() []string {
	idx := []string{}
	return idx
}

type _SortUserBlogMgr struct {
}

var SortUserBlogMgr *_SortUserBlogMgr

func (m *_SortUserBlogMgr) NewSortUserBlog() *SortUserBlog {
	rval := new(SortUserBlog)
	return rval
}
