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

func (i *Index) GetFieldList() string {
	return strings.Join(i.FieldNames, `","`)
}

func (i *Index) GetFuncParam() string {
	var params []string
	for _, f := range i.Fields {
		params = append(params, f.Name+" "+f.GetGoType())
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
