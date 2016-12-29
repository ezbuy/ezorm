package test

import "time"

var _ time.Time

type Blog struct {
	Id        int32     `db:"id" json:"id"`
	UserId    int32     `db:"user_id" json:"user_id"`
	Title     string    `db:"title" json:"title"`
	Content   string    `db:"content" json:"content"`
	Status    int32     `db:"status" json:"status"`
	Readed    int32     `db:"readed" json:"readed"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	isNew     bool
}

func (p *Blog) GetNameSpace() string {
	return "people"
}

func (p *Blog) GetClassName() string {
	return "Blog"
}
func (p *Blog) GetStoreType() string {
	return "hash"
}

func (p *Blog) GetPrimaryKey() string {
	return "Id"
}

func (p *Blog) GetIndexes() []string {
	idx := []string{
		"UserId",
	}
	return idx
}

type _BlogMgr struct {
}

var BlogMgr *_BlogMgr

func (m *_BlogMgr) NewBlog() *Blog {
	rval := new(Blog)
	return rval
}
