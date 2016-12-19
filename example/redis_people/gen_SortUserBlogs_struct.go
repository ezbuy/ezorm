package test

import "time"

var _ time.Time

type SortUserBlogs struct {
	UserId int32   `json:"user_id"`
	Score  float64 `json:"score"`
	BlogId int32   `json:"blog_id"`
	isNew  bool
}

func (p *SortUserBlogs) GetNameSpace() string {
	return "people"
}

func (p *SortUserBlogs) GetClassName() string {
	return "SortUserBlogs"
}
func (p *SortUserBlogs) GetStoreType() string {
	return "zset"
}

func (p *SortUserBlogs) GetPrimaryKey() string {
	return "UserId"
}

func (p *SortUserBlogs) GetIndexes() []string {
	idx := []string{}
	return idx
}

type _SortUserBlogsMgr struct {
}

var SortUserBlogsMgr *_SortUserBlogsMgr

func (m *_SortUserBlogsMgr) NewSortUserBlogs() *SortUserBlogs {
	rval := new(SortUserBlogs)
	return rval
}
