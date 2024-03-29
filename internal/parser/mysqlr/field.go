package mysqlr

import (
	"bytes"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/ezbuy/ezorm/v2/internal/generator"
	"github.com/ezbuy/utils/container/set"
	"github.com/iancoleman/strcase"
)

var nullablePrimitiveSet = map[string]bool{
	"uint8":   true,
	"uint16":  true,
	"uint32":  true,
	"uint64":  true,
	"int8":    true,
	"int16":   true,
	"int32":   true,
	"int64":   true,
	"float32": true,
	"float64": true,
	"bool":    true,
	"string":  true,
}

type Field struct {
	Name      string
	Type      string
	sqlType   string
	sqlColumn string
	Size      int
	Flags     set.Set
	Attrs     map[string]string
	Comment   string
	Validator string
	Obj       *MetaObject
	Default   interface{}
}

func NewField() *Field {
	return &Field{
		Flags: set.NewStringSet(),
	}
}

var SupportedFieldTypes = map[string]string{
	"bool":      "bool",
	"int":       "int32",
	"int8":      "int8",
	"int16":     "int16",
	"int32":     "int32",
	"int64":     "int64",
	"uint":      "uint32",
	"uint8":     "uint8",
	"uint16":    "uint16",
	"uint32":    "uint32",
	"uint64":    "uint64",
	"float32":   "float32",
	"float64":   "float64",
	"string":    "string",
	"datetime":  "datetime",
	"timestamp": "timestamp",
	"timeint":   "timeint",
}

func (f *Field) GetName() string {
	return strcase.ToLowerCamel(f.Name)
}

func (f *Field) GetUnderlineName() string {
	return "_" + CamelName(f.Name)
}

func (f *Field) GetGoType() string {
	return f.GetType()
}

func (f *Field) SetType(t string) error {
	st, ok := SupportedFieldTypes[t]
	if !ok {
		return fmt.Errorf("%s type not support", t)
	}
	f.Type = st
	return nil
}

func (f *Field) FieldName() string {
	return f.ColumnName()
}

func (f *Field) ColumnName() string {
	if f.sqlColumn != "" {
		return f.sqlColumn
	}
	return Camel2Name(f.Name)
}

func (f *Field) IsPrimary() bool {
	return f.Flags.Contains("primary")
}

func (f *Field) IsAutoIncrement() bool {
	return f.Flags.Contains("autoinc")
}

func (f *Field) IsNullable() bool {
	return !f.IsPrimary() && f.Flags.Contains("nullable")
}

func (f *Field) IsUnique() bool {
	return f.Flags.Contains("unique")
}

func (f *Field) IsRange() bool {
	return f.Flags.Contains("range")
}

func (f *Field) IsNorange() bool {
	return f.Flags.Contains("norange")
}

func (f *Field) IsIndex() bool {
	return f.Flags.Contains("index")
}

func (f *Field) IsFullText() bool {
	return f.Flags.Contains("fulltext")
}

func (f *Field) IsEncode() bool {
	if f.IsString() {
		return f.Flags.Contains("encode") || f.Flags.Contains("base64")
	}
	return false
}

func (f *Field) IsNumber() bool {
	if transform := f.GetTransform(); transform != nil {
		if strings.HasPrefix(transform.TypeOrigin, "uint") ||
			strings.HasPrefix(transform.TypeOrigin, "int") ||
			strings.HasPrefix(transform.TypeOrigin, "bool") ||
			strings.HasPrefix(transform.TypeOrigin, "float") {
			return true
		}
	}
	if strings.HasPrefix(f.Type, "uint") ||
		strings.HasPrefix(f.Type, "int") ||
		strings.HasPrefix(f.Type, "bool") ||
		strings.HasPrefix(f.Type, "float") {
		return true
	}
	return false
}

func (f *Field) IsBool() bool {
	if transform := f.GetTransform(); transform != nil {
		return strings.HasPrefix(transform.TypeOrigin, "bool")
	}
	return strings.HasPrefix(f.Type, "bool")
}

func (f *Field) IsString() bool {
	if transform := f.GetTransform(); transform != nil {
		if strings.HasPrefix(transform.TypeOrigin, "string") {
			return true
		}
	}
	if strings.HasPrefix(f.Type, "string") {
		return true
	}
	return false
}

func (f *Field) IsTime() bool {
	switch f.Type {
	case "datetime", "timestamp", "timeint":
		return true
	}
	return false
}

func (f *Field) HasIndex() bool {
	return f.Flags.Contains("unique") ||
		f.Flags.Contains("index") ||
		f.Flags.Contains("range")
}

func (f *Field) GetType() string {
	st := f.Type
	if transform := f.GetTransform(); transform != nil {
		st = transform.TypeTarget
	}

	if f.IsNullable() {
		if st == "time.Time" {
			st = "*time.Time"
		}
	}
	return st
}

func (f *Field) GetNames() string {
	return CamelName(f.Name) + "s"
}

func (f *Field) GetUnderlineNames() string {
	return "_" + CamelName(f.Name) + "s"
}

func (f *Field) IsNullablePrimitive() bool {
	return f.IsNullable() && nullablePrimitiveSet[f.GetType()]
}

func (f *Field) GetNullSQLType() string {
	origin_type := f.Type
	if transform := f.GetTransform(); transform != nil {
		origin_type = transform.TypeOrigin
	}

	if f.IsNullable() {
		if origin_type == "bool" {
			return "NullBool"
		} else if origin_type == "string" {
			return "NullString"
		} else if strings.HasPrefix(origin_type, "int") {
			return "NullInt64"
		} else if strings.HasPrefix(origin_type, "float") {
			return "NullFloat64"
		}
	}
	return origin_type
}

func (f *Field) NullSQLTypeValue() string {
	origin_type := f.Type
	if transform := f.GetTransform(); transform != nil {
		origin_type = transform.TypeOrigin
	}
	if origin_type == "bool" {
		return "Bool"
	} else if origin_type == "string" {
		return "String"
	} else if strings.HasPrefix(origin_type, "int") {
		return "Int64"
	} else if strings.HasPrefix(origin_type, "float") {
		return "Float64"
	}
	panic("unsupported null sql type: " + origin_type)
}

func (f *Field) NullSQLTypeNeedCast() bool {
	t := f.GetType()
	if strings.HasPrefix(t, "int") && t != "int64" {
		return true
	} else if strings.HasPrefix(t, "float") && t != "float64" {
		return true
	}
	return false
}

type Transform struct {
	TypeOrigin  string
	ConvertTo   string
	TypeTarget  string
	ConvertBack string
}

// convert `TypeOrigin` in datebase to `TypeTarget` when quering
// convert `TypeTarget` back to `TypeOrigin` when updating/inserting
var transformMap = map[string]Transform{
	"mysqlr_timestamp": { // TIMESTAMP (string, UTC)
		"string", `orm.TimeParse(%v)`,
		"time.Time", `orm.TimeFormat(%v)`,
	},
	"mysqlr_timeint": { // INT(11)
		"int64", "time.Unix(%v, 0)",
		"time.Time", "%v.Unix()",
	},
	"mysqlr_datetime": { // DATETIME (string, localtime)
		"string", "orm.TimeParseLocalTime(%v)",
		"time.Time", "orm.TimeToLocalTime(%v)",
	},
}

func (f *Field) IsNeedTransform() bool {
	return f.GetTransform() != nil
}

func (f *Field) GetTransform() *Transform {
	key := fmt.Sprintf("%v_%v", f.Obj.Db, f.Type)
	t, ok := transformMap[key]
	if !ok {
		return nil
	}
	return &t
}

func (f *Field) GetTransformValue(prefix string) string {
	t := f.GetTransform()
	if t == nil {
		return prefix + f.Name
	}
	return fmt.Sprintf(t.ConvertBack, prefix+f.Name)
}

func (f *Field) GetTag() string {
	tags := map[string]bool{
		"mysql": false,
	}

	tagstr := []string{}
	for tag, camel := range tags {
		if val, ok := f.Attrs[tag+"Tag"]; ok {
			tagstr = append(tagstr, fmt.Sprintf("%s:\"%s\"", tag, val))
			continue
		}
		switch {
		// use `sqlcolumn` option first
		case tag == "db":
			tagstr = append(tagstr, fmt.Sprintf("%s:\"%s\"", tag, f.ColumnName()))
		case camel:
			tagstr = append(tagstr, fmt.Sprintf("%s:\"%s\"", tag, f.Name))
		default:
			tagstr = append(tagstr, fmt.Sprintf("%s:\"%s\"", tag, Camel2Name(f.Name)))
		}
	}
	if f.Validator != "" {
		tagstr = append(tagstr, fmt.Sprintf("validate:\"%s\"", f.Validator))
	}
	sortstr := sort.StringSlice(tagstr)
	sort.Sort(sortstr)
	if len(sortstr) != 0 {
		return "`" + strings.Join(sortstr, " ") + "`"
	}
	return ""
}

func (f *Field) Read(data generator.Schema) error {
	foundName := false

	for k, v := range data {
		key := string(k)

		if isUpperCase(key[0:1]) {
			if foundName {
				return errors.New("invalid field name: " + key)
			}
			f.Name = key
			if err := f.SetType(v.(string)); err != nil {
				return err
			}

			continue
		}

		switch key {
		case "size":
			f.Size = v.(int)
		case "sqltype":
			f.sqlType = v.(string)
		case "sqlcolumn":
			f.sqlColumn = v.(string)
		case "comment":
			f.Comment = v.(string)
		case "validator":
			f.Validator = strings.ToLower(v.(string))
		case "attrs":
			attrs := make(map[string]string)
			for ki, vi := range v.(map[interface{}]interface{}) {
				attrs[ki.(string)] = vi.(string)
			}
			f.Attrs = attrs
		case "flags":
			for _, flag := range v.([]interface{}) {
				f.Flags.Add(flag.(string))
			}

		case "default":
			f.Default = v
		default:
			return errors.New("invalid field name: " + key)
		}
	}

	//! single field primary adjust for redis ops
	if f.IsUnique() {
		index := NewIndex(f.Obj)
		index.FieldNames = []string{f.Name}
		f.Obj.uniques = append(f.Obj.uniques, index)
	}
	if f.IsIndex() {
		index := NewIndex(f.Obj)
		index.FieldNames = []string{f.Name}
		f.Obj.indexes = append(f.Obj.indexes, index)
	}
	if f.IsRange() {
		index := NewIndex(f.Obj)
		index.FieldNames = []string{f.Name}
		f.Obj.ranges = append(f.Obj.ranges, index)
	}
	return nil
}

// ! field SQL script functions
func (f *Field) SQLColumn() string {
	columns := make([]string, 0, 6)
	columns = append(columns, f.SQLName())
	columns = append(columns, f.SQLType())
	columns = append(columns, f.SQLNull())
	if f.IsAutoIncrement() {
		columns = append(columns, "AUTO_INCREMENT")
	} else {
		columns = append(columns, f.SQLDefault())
	}
	if f.Comment != "" {
		columns = append(columns, "COMMENT", "'"+f.Comment+"'")
	}
	return strings.Join(columns, " ")
}

func (f *Field) SQLName() string {
	return "`" + f.ColumnName() + "`"
}

func (f *Field) SQLType() string {
	if f.sqlType != "" {
		return strings.ToUpper(f.sqlType)
	}
	if f.IsNumber() {
		switch f.GetType() {
		case "bool":
			return "TINYINT(1) UNSIGNED"
		case "uint8":
			return "SMALLINT UNSIGNED"
		case "uint16":
			return "MEDIUMINT UNSIGNED"
		case "uint32":
			return "INT(11) UNSIGNED"
		case "uint64":
			return "BIGINT UNSIGNED"
		case "int8":
			return "SMALLINT"
		case "int16":
			return "MEDIUMINT"
		case "int32", "int":
			return "INT(11)"
		case "int64":
			return "BIGINT(20)"
		case "float32", "float64":
			return "FLOAT"
		case "time.Time", "*time.Time":
			return "BIGINT(20)"
		}
	}
	if f.IsString() {
		switch f.Type {
		case "datetime":
			return "DATETIME"
		case "timestamp", "timeint":
			return "TIMESTAMP"
		}
		if f.Size == 0 {
			return "VARCHAR(100)"
		}
		return fmt.Sprintf("VARCHAR(%d)", f.Size)
	}
	return f.GetType()
}

func (f *Field) SQLNull() string {
	if f.IsNullable() {
		return "NULL"
	}
	return "NOT NULL"
}

func (f *Field) SQLDefault() string {
	if f.IsNullable() {
		return ""
	}
	if f.IsTime() {
		if f.IsString() {
			return "DEFAULT CURRENT_TIMESTAMP"
		}
		if f.IsNumber() {
			return "DEFAULT '0'"
		}
	}

	if f.IsBool() {
		switch v, _ := f.Default.(bool); v {
		case true:
			return "DEFAULT '1'"
		default:
			return "DEFAULT '0'"
		}
	}

	if f.IsNumber() {
		return "DEFAULT '0'"
	}
	if f.IsString() {
		return "DEFAULT ''"
	}
	return ""
}

type Fields []*Field

func (fs Fields) GetFuncParam() string {
	var params []string
	for _, f := range fs {
		params = append(params, "_"+CamelName(f.Name)+" "+f.GetType())
	}
	return strings.Join(params, ", ")
}

func (fs Fields) GetObjectParam() string {
	var params []string
	for _, f := range fs {
		params = append(params, "obj."+f.Name)
	}
	return strings.Join(params, ", ")
}

func (fs Fields) GetConstructor() string {
	params := make([]string, 0, len(fs)+1)
	for _, f := range fs {
		params = append(params, f.Name+" : "+"_"+CamelName(f.Name))
	}
	params = append(params, "")
	return strings.Join(params, ",\n")
}

func (fs Fields) GetFieldNames() string {
	var names []string
	for _, f := range fs {
		names = append(names, strconv.Quote(f.FieldName()))
	}
	return fmt.Sprintf("[]string{%s}", strings.Join(names, ", "))
}

func Camel2Name(s string) string {
	nameBuf := bytes.NewBuffer(nil)
	before := false
	for i := range s {
		n := rune(s[i]) // always ASCII?
		if unicode.IsUpper(n) {
			if !before && i > 0 {
				nameBuf.WriteRune('_')
			}
			n = unicode.ToLower(n)
			before = true
		} else {
			before = false
		}
		nameBuf.WriteRune(n)
	}
	return nameBuf.String()
}

func CamelName(argName string) string {
	size := len(argName)
	if size <= 0 {
		return "nilArgs"
	}
	fl := argName[0]
	if fl >= 65 && fl <= 90 {
		return string([]byte{byte(fl + byte(32))}) + string(argName[1:])
	}
	return argName
}

func isUpperCase(c string) bool {
	return c == strings.ToUpper(c)
}
