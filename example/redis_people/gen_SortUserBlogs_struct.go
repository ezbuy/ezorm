package test

import "time"

var _ time.Time

type SortUserBlogs struct {
	UserId int32   `bson:"UserId" json:"UserId"`
	Score  float64 `bson:"Score" json:"Score"`
	BlogId int32   `bson:"BlogId" json:"BlogId"`
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
