package mssql_people

import "github.com/ezbuy/ezorm/cache"

type People struct {
	Name string `db:"Name"`
	Id   int32  `db:"Id"`
	Age  int32  `db:"Age"`
}

func (p *People) GetNameSpace() string {
	return "mssql_people"
}

func (p *People) GetClassName() string {
	return "People"
}

type _PeopleMgr struct {
}

var PeopleMgr *_PeopleMgr

var PeopleCache cache.Cache

var idFieldName = "Id"

func (m *_PeopleMgr) NewPeople() *People {
	rval := new(People)
	return rval
}
