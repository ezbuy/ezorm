package people

type People struct {
	NonIndexA string `db:"NonIndexA"`

	NonIndexB string `db:"NonIndexB"`

	PeopleId int32 `db:"PeopleId"`

	Age int32 `db:"Age"`

	Name string `db:"Name"`

	IndexAPart1 int64 `db:"IndexAPart1"`

	IndexAPart2 int32 `db:"IndexAPart2"`

	IndexAPart3 int32 `db:"IndexAPart3"`

	UniquePart1 int32 `db:"UniquePart1"`

	UniquePart2 int32 `db:"UniquePart2"`
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
