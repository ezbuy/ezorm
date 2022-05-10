package shared

import (
	"bytes"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/ezbuy/ezorm/v2/internal/generator"
	"github.com/ezbuy/utils/container/set"
)

const (
	flagNullable = "nullable"
)

var (
	nullablePrimitiveSet = map[string]bool{
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
)

var _ generator.IField = (*Field)(nil)

type Field struct {
	Attrs        map[string]string
	DefaultValue string
	Flags        set.Set
	Index        string
	Key          string
	Label        string
	PlaceHolder  string
	Name         string
	Order        string
	Tag          string
	Type         string
	Size         int    // The field size, use for DDL generation.
	Decimal      int    // The decimal size, only use for "float64" type.
	Default      string // The default value for this field, use for DDL generation.
	Widget       string
	Remark       string
	FK           *ForeignKey
	Obj          *Obj
	AsSort       bool
	Comment      string
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
	f.Flags = set.NewStringSet()
}

func (f *Field) ArgName() string {
	return strings.ToLower(f.Name[:1]) + f.Name[1:]
}

func (f *Field) IsRequired() bool {
	return false
}

func (f *Field) IsUnique() bool {
	return f.Flags.Contains("unique")
}

func (f *Field) GetThriftType() string {
	return SupportedFieldTypes[f.Type]
}

func (f *Field) getGoType(typestr string) string {
	if transform := f.GetTransformType(); transform != nil {
		return transform.TypeTarget
	}

	if typestr == "datetime" {
		return "int64"
	}

	if strings.HasPrefix(typestr, "list<") {
		innerType := typestr[5 : len(typestr)-1]
		return "[]" + f.getGoType(innerType) + ""
	}

	if strings.HasPrefix(typestr, "map[") {
		i := strings.Index(typestr, "]")
		keyType := typestr[4:i]
		valType := typestr[i+1:]
		return "map[" + f.getGoType(keyType) + "]" + f.getGoType(valType)
	}
	return typestr
}

func (f *Field) GetGoType() string {
	return f.getGoType(f.Type)
}

func (f *Field) GetNullSQLType() string {
	t := f.GetGoType()
	if t == "bool" {
		return "NullBool"
	} else if t == "string" {
		return "NullString"
	} else if strings.HasPrefix(t, "int") {
		return "NullInt64"
	} else if strings.HasPrefix(t, "float") {
		return "NullFloat64"
	}
	return t
}

func (f *Field) AttrsContains(attr string) bool {
	_, ok := f.Attrs[attr]
	return ok
}

func (f *Field) BsonTagName() string {
	if bVal, ok := f.Attrs["bsonTag"]; ok {
		return bVal
	}

	if f.Name == "ID" {
		return "_id"
	}

	return f.Name
}

func (f *Field) DbName() string {
	return camel2name(f.Name)
}

func (f *Field) GetName() string {
	return camel2name(f.Name)
}

func (f *Field) GetTag() string {
	tags := map[string]bool{}
	for _, db := range f.Obj.Dbs {
		switch db {
		case "mongo":
			tags["bson"] = true
			tags["json"] = true
		case "mysql":
			tags["db"] = false
		}
	}
	if len(tags) == 0 {
		tags["bson"] = true
		tags["json"] = true
	}

	tagstr := []string{}
	for tag, camel := range tags {
		if val, ok := f.Attrs[tag+"Tag"]; ok {
			tagstr = append(tagstr, fmt.Sprintf("%s:\"%s\"", tag, val))
		} else {
			if camel {
				tagstr = append(tagstr, fmt.Sprintf("%s:\"%s\"", tag, f.Name))
			} else {
				tagstr = append(tagstr, fmt.Sprintf("%s:\"%s\"", tag, camel2name(f.Name)))
			}
		}
	}
	sortstr := sort.StringSlice(tagstr)
	sort.Sort(sortstr)
	if len(sortstr) != 0 {
		return "`" + strings.Join(sortstr, " ") + "`"
	}
	return ""
}

func (f *Field) NullSQLTypeValue() string {
	t := f.GetGoType()
	if t == "bool" {
		return "Bool"
	} else if t == "string" {
		return "String"
	} else if strings.HasPrefix(t, "int") {
		return "Int64"
	} else if strings.HasPrefix(t, "float") {
		return "Float64"
	}
	panic("unsupported null sql type: " + t)
}

func (f *Field) NullSQLTypeNeedCast() bool {
	t := f.GetGoType()
	if strings.HasPrefix(t, "int") && t != "int64" {
		return true
	} else if strings.HasPrefix(t, "float") && t != "float64" {
		return true
	}
	return false
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

	if f.FK != nil {
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

	return f.FK.Field
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

func (f *Field) IsNullable() bool {
	return f.Flags.Contains(flagNullable)
}

func (f *Field) IsNullablePrimitive() bool {
	return f.IsNullable() && nullablePrimitiveSet[f.GetGoType()]
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
	"mysql_timestamp": { // TIMESTAMP (string, UTC)
		"string", `db.TimeParse(%v)`,
		"time.Time", `db.TimeFormat(%v)`,
	},
	"mysql_timeint": { // INT(11)
		"int64", "time.Unix(%v, 0)",
		"time.Time", "%v.Unix()",
	},
	"mysql_datetime": { // DATETIME (string, localtime)
		"string", "db.TimeParseLocalTime(%v)",
		"time.Time", "db.TimeToLocalTime(%v)",
	},
}

func (f *Field) AsArgName(prefix string) string {
	t := f.GetTransformType()
	if t == nil {
		return prefix + f.Name
	}
	return fmt.Sprintf(t.ConvertBack, prefix+f.Name)
}

func (f *Field) IsNeedTransform() bool {
	return f.GetTransformType() != nil
}

func (f *Field) GetTransformType() *Transform {
	key := fmt.Sprintf("%v_%v", f.Obj.Db, f.Type)
	t, ok := transformMap[key]
	if !ok {
		return nil
	}
	return &t
}

func (f *Field) HasIndex() bool {
	return f.Flags.Contains("index") || f.Flags.Contains("sort") || f.IsUnique()
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
					if f.Obj.DbContains("mysql") {
					} else {
						f.Type = "int64"
					}
					f.Widget = "datetime"
				}

				continue
			}
			switch key {
			case "label":
				f.Label = val
			case "fk":
				f.FK = NewForeignKey(val)
			case "widget":
				f.Widget = val
			case "remark":
				f.Remark = val
			case "comment":
				f.Comment = val
			case "default":
				f.Default = fmt.Sprintf("'%s'", val)
			case "size":
				size, err := strconv.Atoi(val)
				if err != nil {
					return fmt.Errorf("size %s is not a number", val)
				}
				f.Size = size
			default:
				return errors.New("invalid field name: " + key)
			}
		case int:
			switch key {
			case "default":
				f.Default = strconv.Itoa(val)
			case "size":
				f.Size = val
			case "decimal":
				f.Decimal = val
			default:
				f.Name = key
				f.Tag = strconv.Itoa(val)
			}
		case []interface{}:
			switch key {
			case "flags":
				for _, v := range val {
					f.Flags.Add(v.(string))
					switch v.(string) {
					case "sort":
						f.AsSort = true
					default:
					}
				}
			}
		}

		if key == "attrs" {
			attrs := make(map[string]string)
			for ki, vi := range v.(map[interface{}]interface{}) {
				attrs[ki.(string)] = vi.(string)
			}
			f.Attrs = attrs
		}
	}
	return nil
}

func (f *Field) IsAutoInc() bool {
	return f.Flags.Contains("autoinc")
}

func (f *Field) DisableAutoInc() bool {
	return f.Flags.Contains("noinc")
}

func (f *Field) MysqlCreation() string {
	var buffer bytes.Buffer
	name := fmt.Sprintf("`%s` ", camel2name(f.Name))
	buffer.WriteString(name)
	buffer.WriteString(f.mysqlDbType())

	if !f.IsNullable() {
		buffer.WriteString(" NOT NULL ")
		if f.Default == "" && !f.IsAutoInc() {
			f.Default = f.mysqlDefaultValue()
		}
	}
	if f.Default != "" {
		buffer.WriteString(" DEFAULT ")
		buffer.WriteString(f.Default)
		buffer.WriteByte(' ')
	}
	if f.IsAutoInc() {
		buffer.WriteString(" AUTO_INCREMENT ")
	}
	if f.Comment != "" {
		comment := fmt.Sprintf(" COMMENT '%s' ", f.Comment)
		buffer.WriteString(comment)
	}
	line := buffer.String()
	line = strings.TrimSpace(line)

	return strings.Join(strings.Fields(line), " ")
}

func (f *Field) mysqlDbType() string {
	var basic string
	switch f.Type {
	case "int8":
		basic = "smallint"
	case "int32":
		basic = "int"
	case "int64":
		basic = "bigint"
	case "float64", "float32":
		decimal := f.Decimal
		if decimal <= 0 {
			decimal = 4
		}
		size := f.Size
		if size <= 0 {
			size = 11
		}
		if size < decimal {
			// For mysql, Size cannot smaller than Decimal
			size = decimal + 4
		}
		return fmt.Sprintf("DECIMAL(%d, %d)", size, decimal)

	case "string":
		basic = "varchar"

	case "[]byte":
		basic = "binary"

	case "bool":
		basic = "tinyint"

	case "time.Time", "datetime":
		return "DATETIME"

	case "timestamp":
		return "TIMESTAMP"

	default:
		return strings.ToUpper(f.Type)
	}

	if f.Size > 0 {
		basic = fmt.Sprintf("%s(%d)", basic, f.Size)
	} else {
		if f.Type == "string" || f.Type == "[]byte" {
			fmt.Printf("WARNING: [mysql-script] Use default size 200 for "+
				"field %s.%s, please consider add size for it.\n",
				f.Obj.Name, f.Name)
			basic = fmt.Sprintf("%s(200)", basic)
		}
	}

	return strings.ToUpper(basic)
}

func (f *Field) mysqlDefaultValue() string {
	switch f.Type {
	case "int8", "int32", "int64", "bool":
		return "0"
	case "string", "[]byte":
		return "''"
	case "float32", "float64":
		return "'0.00'"
	case "timestamp", "datetime":
		return "CURRENT_TIMESTAMP"
	}
	return "''"
}

func DbToGoType(colType string) string {
	var typeStr string
	switch colType {
	case "nvarchar", "timestamp", "text", "cursor", "uniqueidentifier", "sysname", "real",
		"binary", "varbinary", "nchar", "char", "varchar":
		typeStr = "string"
	case "datetime", "smalldatetime":
		// Use pointer type to avoid null value panic
		typeStr = "*time.Time"
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

// -----------------------------------------------------------------------------

type ForeignKey struct {
	Tbl   string
	Field string
}

func NewForeignKey(name string) *ForeignKey {
	sp := strings.Split(name, ".")
	return &ForeignKey{
		Tbl:   sp[0],
		Field: sp[1],
	}
}
