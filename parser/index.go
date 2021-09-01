package parser

import (
	"bytes"
	"fmt"
	"strings"
)

type Index struct {
	Name       string
	Fields     []*Field
	FieldNames []string
	IsUnique   bool
	IsSparse   bool
}

func (i *Index) GetFindInIds(idx int, bufName, name string) string {
	return toIds(bufName, i.Fields[idx].Type, name)
}

func (i *Index) GetFirstField() *Field {
	return i.Fields[0]
}

func (i *Index) IsFindInType(field *Field) bool {
	switch field.Type {
	case "int", "int32", "string", "int64":
		return true
	default:
		return false
	}
}

func (i *Index) CanUseFindList() bool {
	return len(i.Fields) == 1 && i.CanUseFindIn()
}

func (i *Index) CanUseFindIn() bool {
	for _, field := range i.Fields {
		if !i.IsFindInType(field) {
			return false
		}
	}
	return true
}

func (i *Index) GetFieldList() string {
	return strings.Join(i.FieldNames, `","`)
}

func (i *Index) GetFuncParamIn() string {
	var params []string
	for _, f := range i.Fields {
		params = append(params, f.Name+" []"+f.GetGoType())
	}
	return strings.Join(params, ", ")
}

func (i *Index) GetFuncParam() string {
	var params []string
	for _, f := range i.Fields {
		params = append(params, f.Name+" "+f.GetGoType())
	}
	return strings.Join(params, ", ")
}

func (i *Index) GetFuncParamOriNames() string {
	var params []string
	for _, f := range i.Fields {
		params = append(params, f.Name)
	}
	return strings.Join(params, ", ")
}

func (i *Index) GetFuncParamNames(prefixs ...string) string {
	buf := bytes.NewBuffer(nil)
	prefix := ""
	if len(prefixs) > 0 {
		prefix = prefixs[0]
	}
	for _, f := range i.Fields {
		tf := f.GetTransformType()
		if tf == nil {
			buf.WriteString(prefix + f.Name)
		} else {
			buf.WriteString(fmt.Sprintf(tf.ConvertBack, prefix+f.Name))
		}
		buf.WriteString(",")
	}
	length := buf.Len()
	if length == 0 {
		return ""
	}
	return string(buf.Bytes()[:length-1])
}
