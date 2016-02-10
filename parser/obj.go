package parser

import (
	"errors"
	"io"
	"strconv"
	"strings"
	"text/template"

	"github.com/ezbuy/utils/container/set"
)

var Tpl *template.Template

func init() {
	tmp, err := template.ParseGlob("src/code.1dmy.com/xyz/xuanwu/tmpl/*.tmpl")
	Tpl = tmp
	if err != nil {
		panic(err)
	}
}

type Obj struct {
	DefaultOrder string
	Extend       string
	Fields       []*Field
	Label        string
	Name         string
	GenType      string
	Package      string
	SearchIndex  string
	SearchType   string
	FilterFields []string
	TplWriter    io.Writer
}

func (o *Obj) init() {
	if strings.HasSuffix(o.Name, "Form") {
		o.GenType = "form"
	}
}

func (o *Obj) LoadTpl(tpl string) string {
	err := Tpl.ExecuteTemplate(o.TplWriter, tpl, o)
	if err != nil {
		println("LoadTpl", tpl, err.Error())
		panic(err)
	}
	return ""
}

func (o *Obj) LoadField(f *Field) string {
	err := Tpl.ExecuteTemplate(o.TplWriter, "field_"+f.Type, f)
	if err != nil {
		println("LoadField", f.Name, f.Type, err.Error())
		panic(err)
	}
	return ""
}

func (o *Obj) GetGenTypes() []string {
	switch o.GenType {
	case "form":
		return []string{"struct", "form"}
	case "enum":
		return []string{"enum"}
	default:
		return []string{"struct", "thrift_serial", "form", "orm"}
	}
}

func (o *Obj) GetFormImports() (imports []string) {
	data := set.NewStringSet()
	numberTypes := set.NewStringSet("float64", "int", "int32", "int64")
	for _, f := range o.Fields {
		if f.Type == "[]string" {
			data.Add("strings")
			continue
		}
		if numberTypes.Contains(f.Type) {
			data.Add("strconv")
			continue
		}
	}
	return data.ToArray()
}

func (o *Obj) GetOrmImports() (imports []string) {
	data := set.NewStringSet()
	for _, f := range o.Fields {
		if f.FK != "" {
			tmp := strings.SplitN(f.FK, ".", 2)
			if len(tmp) == 2 {
				packageName := tmp[0]
				data.Add(packageName)
			}
		}
	}
	return data.ToArray()
}

func (o *Obj) NeedOrm() bool {
	return true
}

func (o *Obj) NeedSearch() bool {
	return false
}

func (o *Obj) NeedIndex() bool {
	return false
}

func (o *Obj) NeedMapping() bool {
	return false
}

func (o *Obj) GetStringFilterFields() []*Field {
	return nil
}

func (o *Obj) GetListedFields() []*ListedField {
	return nil
}

func (o *Obj) GetFilterFields() []*Field {
	return nil
}

func (o *Obj) GetNonIDFields() []*Field {
	return o.Fields[1:]
}

func ToStringSlice(val []interface{}) (result []string) {
	result = make([]string, len(val))
	for i, v := range val {
		result[i] = v.(string)
	}
	return
}

func (o *Obj) Read(data map[string]interface{}) error {
	o.init()
	hasType := false
	for key, val := range data {
		switch key {
		case "type":
			o.GenType = val.(string)
			hasType = true
			break
		}
	}

	if hasType {
		delete(data, "type")
	}

	for key, val := range data {
		switch key {
		case "label":
			o.Label = val.(string)
		case "extend":
			o.Extend = val.(string)
		case "filterFields":
			o.FilterFields = ToStringSlice(val.([]interface{}))
		case "fields":
			fieldData := val.([]interface{})
			startPos := 0
			// println(o.GenType, o.Name)
			if o.GenType == "" {
				o.Fields = make([]*Field, len(fieldData)+1)
				f := new(Field)
				f.init()
				f.Obj = o
				f.Name = "ID"
				f.Tag = "1"
				f.Type = "string"
				o.Fields[0] = f
				startPos = 1
			} else {
				o.Fields = make([]*Field, len(fieldData))
			}

			for i, field := range fieldData {
				f := new(Field)
				f.Obj = o
				f.Tag = strconv.Itoa(i + 2)

				err := f.Read(field.(map[interface{}]interface{}))
				if err != nil {
					return errors.New(o.Name + " obj has " + err.Error())
				}
				o.Fields[i+startPos] = f
			}
		default:
			return errors.New(o.Name + " has invalid obj property: " + key)
		}
	}
	return nil
}
