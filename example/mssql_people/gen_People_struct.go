package test

import "time"

var _ time.Time

type People struct {
	NonIndexA   string     `db:"non_index_a"`
	NonIndexB   string     `db:"non_index_b"`
	PeopleId    int32      `db:"people_id"`
	Age         int32      `db:"age"`
	Name        string     `db:"name"`
	IndexAPart1 int64      `db:"index_a_part1"`
	IndexAPart2 int32      `db:"index_a_part2"`
	IndexAPart3 int32      `db:"index_a_part3"`
	UniquePart1 int32      `db:"unique_part1"`
	UniquePart2 int32      `db:"unique_part2"`
	CreateDate  *time.Time `db:"create_date"`
	UpdateDate  *time.Time `db:"update_date"`
}

func (p *People) GetNameSpace() string {
	return "people"
}

func (p *People) GetClassName() string {
	return "People"
}

func (p *People) GetIndexes() []string {
	idx := []string{
		"Age",
		"Name",
	}
	return idx
}

type _PeopleMgr struct {
}

var PeopleMgr *_PeopleMgr

func (m *_PeopleMgr) NewPeople() *People {
	rval := new(People)
	return rval
}
