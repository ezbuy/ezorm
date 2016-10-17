package parser

import (
	"fmt"
	"strings"
)

type GroupByItem struct {
	Fields    []*GroupByKey
	FromField *Field
}

func (g *GroupByItem) Add(name string, f *Field) {
	g.Fields = append(g.Fields, NewGroupByKey(name, f))
}

func (g *GroupByItem) HasAggGroupFields() bool {
	for _, f := range g.Fields {
		if f.IsFunc {
			return true
		}
	}
	return false
}

func (g *GroupByItem) GetGroupFieldKey() []string {
	ret := make([]string, 0, len(g.Fields))
	for _, f := range g.Fields {
		if f.IsFunc {
			continue
		}
		ret = append(ret, f.Projection)
	}
	return ret

}

func (g *GroupByItem) GetGroupFieldProjection() []string {
	ret := make([]string, 0, len(g.Fields))
	for _, f := range g.Fields {
		ret = append(ret, f.Projection)
	}
	return ret
}

type GroupByKey struct {
	Field      string
	FieldType  string
	Projection string
	IsFunc     bool
}

func NewGroupByKey(name string, f *Field) *GroupByKey {
	if f != nil {
		return &GroupByKey{
			Field:      f.Name,
			FieldType:  f.Type,
			Projection: camel2name(f.Name),
			IsFunc:     false,
		}
	}
	idx := strings.Index(name, "(")
	if idx <= 0 {
		panic(fmt.Sprintf("field %v is not found, are you putting `fields` behind `groupby` ?", name))
		return nil
	}

	return &GroupByKey{
		Field:      strings.ToUpper(name[:idx]),
		FieldType:  "int",
		Projection: name,
		IsFunc:     true,
	}
}
