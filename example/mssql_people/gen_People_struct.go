package people

type People struct {
	NonIndexA string `db:"NonIndexA"`

	NonIndexB string `db:"NonIndexB"`

	PeopleId int32 `db:"PeopleId"`

	Age int32 `db:"Age"`

	Name string `db:"Name"`

	IndexAPart1 int32 `db:"IndexAPart1"`

	IndexAPart2 int32 `db:"IndexAPart2"`

	IndexAPart3 int32 `db:"IndexAPart3"`
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

func (m *_PeopleMgr) NewPeople() *People {
	rval := new(People)
	return rval
}
