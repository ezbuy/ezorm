package parser

import (
	"errors"
	"strconv"
	"strings"

	"github.com/ezbuy/utils/container/set"
)

type Field struct {
	DefaultValue string
	Attrs        set.Set
	Index        string
	Key          string
	Label        string
	PlaceHolder  string
	Name         string
	Order        string
	Tag          string
	Type         string
	Widget       string
	Remark       string
	FK           string
	Obj          *Obj
}

type ListedField struct {
	Key      string
	ObjName  string
	ObjField string
}

var SupportedFieldTypes = map[string]string{
	"[]string": "list<string>",
	"bool":     "bool",
	"datetime": "i64",
	"float64":  "double",
	"int":      "i32",
	"int32":    "i32",
	"int64":    "i64",
	"string":   "string",
}

func isUpperCase(c string) bool {
	return c == strings.ToUpper(c)
}

func (f *Field) init() {
	f.Attrs = set.NewStringSet()
}

func (f *Field) IsRequired() bool {
	return false
}

func (f *Field) GetThriftType() string {
	return SupportedFieldTypes[f.Type]
}

func GetGoType(typestr string) string {
	if typestr == "datetime" {
		return "int64"
	}

	if strings.HasPrefix(typestr, "list<") {
		innerType := typestr[5 : len(typestr)-1]
		return "[]" + GetGoType(innerType) + ""
	}

	if strings.HasPrefix(typestr, "map[") {
		i := strings.Index(typestr, "]")
		keyType := typestr[4:i]
		valType := typestr[i+1:]
		return "map[" + GetGoType(keyType) + "]" + GetGoType(valType)
	}
	return typestr
}

func (f *Field) GetGoType() string {
	return GetGoType(f.Type)
}

func (f *Field) HasDefaultValue() bool {
	return f.DefaultValue != "" && f.DefaultValue != "currentUser"
}

func (f *Field) HasRule() bool {
	return false
}

func (f *Field) HasStringList() bool {
	return false
}

func (f *Field) HasForeign() bool {
	if f.Name == "ID" {
		return false
	}
	if strings.HasSuffix(f.Name, "ID") {
		return true
	}

	if f.FK != "" {
		return true
	}
	return false
}

func (f *Field) Foreign() string {
	if strings.HasSuffix(f.Name, "ID") {
		return f.Name[:len(f.Name)-2]
	}

	if strings.HasSuffix(f.Name, "Id") {
		return f.Name[:len(f.Name)-2]
	}

	return f.Name
}

func (f *Field) ForeignType() string {
	if strings.HasSuffix(f.Name, "ID") {
		return f.Name[:len(f.Name)-2]
	}

	tmp := strings.Split(f.FK, "/")

	return tmp[len(tmp)-1]
}

func (f *Field) HasBindData() bool {
	return false
}

func (f *Field) HasDisable() bool {
	return false
}

func (f *Field) HasHidden() bool {
	return false
}

func (f *Field) HasReadOnly() bool {
	return false
}

func (f *Field) HasMeta() bool {
	return false
}

func (f *Field) HasEnums() bool {
	return false
}

func (f *Field) HasIndex() bool {
	return f.Index != ""
}

func (f *Field) Read(data map[interface{}]interface{}) error {
	f.init()
	foundName := false
	for k, v := range data {
		key := k.(string)

		switch val := v.(type) {
		case string:
			if isUpperCase(key[0:1]) {
				if foundName {
					return errors.New("invalid field name: " + key)
				}
				foundName = true
				f.Name = key

				// if _, ok := SupportedFieldTypes[val]; !ok {
				// 	return errors.New(key + " has invalid type: " + val)
				// }

				f.Type = val
				if f.Type == "int" {
					f.Type = "int32"
				} else if f.Type == "datetime" {
					f.Type = "int64"
					f.Widget = "datetime"
				}

				continue
			}
			switch key {
			case "label":
				f.Label = val
			case "fk":
				f.FK = val
			case "widget":
				f.Widget = val
			case "remark":
				f.Remark = val
			default:
				return errors.New("invalid field name: " + key)
			}
		case int:
			f.Name = key
			f.Tag = strconv.Itoa(val)
		case []interface{}:
			switch key {
			case "attrs":
				for _, v := range val {
					f.Attrs.Add(v.(string))
				}
			}
		}

	}
	return nil
}

func DbToGoType(colType string) string {
	var typeStr string
	switch colType {
	case "nvarchar", "timestamp", "text", "cursor", "uniqueidentifier", "sysname", "real",
		"binary", "varbinary", "nchar", "char":
		typeStr = "string"
	case "datetime", "smalldatetime":
		typeStr = "string"
	case "decimal", "numeric", "float":
		typeStr = "float64"
	case "smallint", "tinyint":
		typeStr = "int8"
	case "int":
		typeStr = "int32"
	case "bigint":
		typeStr = "int64"
	case "money", "smallmoney":
		typeStr = "float32"
	case "bit":
		typeStr = "bool"
	case "image":
		typeStr = "[]byte"
	}
	return typeStr
}
