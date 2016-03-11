package people

import "github.com/ezbuy/ezorm/cache"

type People struct {
	Age int32 `db:"Age"`

	Name string `db:"Name"`

	PeopleId int32 `db:"PeopleId"`
}

func (p *People) GetNameSpace() string {
	return "people"
}

func (p *People) GetClassName() string {
	return "People"
}

type _PeopleMgr struct {
}

var PeopleMgr *_PeopleMgr

var PeopleCache cache.Cache

func (m *_PeopleMgr) NewPeople() *People {
	rval := new(People)
	return rval
}
