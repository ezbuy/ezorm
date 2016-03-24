package parser

import "strings"

type Index struct {
	Name       string
	Fields     []*Field
	FieldNames []string
	IsUnique   bool
	IsSparse   bool
}

func (i *Index) GetFieldList() string {
	return strings.Join(i.FieldNames, `","`)
}

func (i *Index) GetFuncParam() string {
	var params []string
	for _, f := range i.Fields {
		params = append(params, f.Name+" "+GetGoType(f.Type))
	}
	return strings.Join(params, ", ")
}

func (i *Index) GetFuncParamNames() string {
	return strings.Join(i.FieldNames, ", ")
}
