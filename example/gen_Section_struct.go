package example

import "time"

var _ time.Time

type Section struct {
	Key string `bson:"Key" json:"Key"`

	Val int32 `bson:"Val" json:"Val"`

	Data map[string]string `bson:"Data" json:"Data"`
}

func (p *Section) GetNameSpace() string {
	return "example"
}

func (p *Section) GetClassName() string {
	return "Section"
}

type _SectionMgr struct {
}

var SectionMgr *_SectionMgr

func (m *_SectionMgr) NewSection() *Section {
	rval := new(Section)
	return rval
}
