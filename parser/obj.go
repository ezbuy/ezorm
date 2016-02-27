package parser

import (
	"errors"
	"io"
	"strconv"
	"strings"
	"text/template"

	"github.com/ezbuy/ezorm/tpl"
	"github.com/ezbuy/utils/container/set"
)

var Tpl *template.Template

func init() {
	Tpl = template.New("ezorm")
	files := []string{
		"tpl/mongo_collection.gogo",
		"tpl/mongo_foreign_key.gogo",
		"tpl/mongo_mongo.gogo",
		"tpl/mongo_orm.gogo",
		"tpl/mongo_search.gogo",
		"tpl/struct.gogo",
	}
	for _, fname := range files {
		data, err := tpl.Asset(fname)
		if err != nil {
			panic(err)
		}
		_, err = Tpl.Parse(string(data))
		if err != nil {
			panic(err)
		}
	}
}

type Obj struct {
	Extend       string
	Fields       []*Field
	Name         string
	Db           string
	Package      string
	SearchIndex  string
	SearchType   string
	FilterFields []string
	TplWriter    io.Writer
}

func (o *Obj) init() {

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

	switch o.Db {
	case "mongo":
		return []string{"struct", "mongo_orm"}
	case "enum":
		return []string{"enum"}
	default:
		return []string{"struct"}
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
	for _, f := range o.Fields {
		if f.IsUnique() {
			return true
		}
	}
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
		case "db":
			o.Db = val.(string)
			hasType = true
			break
		}
	}

	if hasType {
		delete(data, "db")
	}

	for key, val := range data {
		switch key {
		case "extend":
			o.Extend = val.(string)
		case "filterFields":
			o.FilterFields = ToStringSlice(val.([]interface{}))
		case "fields":
			fieldData := val.([]interface{})
			startPos := 0

			if o.Db == "mongo" {
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
