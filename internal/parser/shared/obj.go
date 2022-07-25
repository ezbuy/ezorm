package shared

import (
	"embed"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"text/template"

	"github.com/ezbuy/ezorm/v2/internal/generator"
	"github.com/ezbuy/utils/container/set"
)

var Tpl *template.Template

//go:embed tpl/*.gogo tpl/*.sql
var content embed.FS

func init() {
	funcMap := template.FuncMap{
		"minus":         minus,
		"getNullType":   getNullType,
		"join":          strings.Join,
		"preSuffixJoin": preSuffixJoin,
		"repeatJoin":    repeatJoin,
		"camel2list":    camel2list,
		"camel2name":    camel2name,
		"strDefault":    strDefault,
		"strif":         strif,
		"toids":         toIds,
		"add":           add,
	}

	files := []string{
		"tpl/struct.gogo",
		"tpl/mysql_config.gogo",
		"tpl/mysql_orm.gogo",
		"tpl/mysql_fk.gogo",
		"tpl/mongo_config.gogo",
		"tpl/mysql_script.sql",
		"tpl/sql_method.gogo",
		"tpl/mongo_config.gogo",
		"tpl/mongo_orm.gogo",
	}

	Tpl = template.New("ezorm").Funcs(funcMap)

	if _, err := Tpl.ParseFS(content, files...); err != nil {
		panic(err)
	}

}

func (f *Field) BJTag() string {
	var bjTag string
	if bVal, ok := f.Attrs["bsonTag"]; ok {
		if bVal != "" {
			bjTag = fmt.Sprintf("`bson:\"%s\"", bVal)
		}
	} else {
		bjTag = fmt.Sprintf("`bson:\"%s\"", f.Name)
	}

	if jVal, ok := f.Attrs["jsonTag"]; ok {
		if jVal != "" {
			bjTag += fmt.Sprintf(" json:\"%s\"`", jVal)
		}
	} else {
		bjTag += fmt.Sprintf(" json:\"%s\"`", f.Name)
	}
	return bjTag
}

var _ generator.IObject = (*Obj)(nil)

type Obj struct {
	Db           string
	Dbs          []string
	Extend       string
	Fields       []*Field
	FieldNameMap map[string]*Field
	FilterFields []string
	Indexes      []*Index
	Name         string
	Namespace    string
	// Package is deprecated , use Namespace instead
	Package     string
	GoPackage   string
	SearchIndex string
	SearchType  string
	Table       string
	TplWriter   io.Writer
	DbName      string
	StoreType   string
	ValueType   string
	ValueField  *Field
	ModelType   string
	ImportSQL   string
	Comment     string

	// IsEmbed describes whether the object should embedded in another object.
	// The embedded object will not generator the orm method but only the db struct.
	IsEmbed bool
}

func (o *Obj) init() {
	if o.Namespace == "" {
		o.Namespace = o.GoPackage
	}
	o.FieldNameMap = make(map[string]*Field)
}

func (o *Obj) GetTable() string {
	return o.Table
}

func (o *Obj) GetFieldNameWithDB(name string) string {
	if o.DbName != "" {
		dbname := o.DbName
		if o.DbContains("mysql") {
			dbname = camel2name(o.DbName)
		}
		return fmt.Sprintf("%s.%s", dbname, name)
	}
	return name
}

func (o *Obj) FieldsMap() map[string]generator.IField {
	m := make(map[string]generator.IField)
	for _, f := range o.Fields {
		m[f.GetName()] = f
	}
	return m
}

func (o *Obj) GetPrimaryKeyName() string {
	k := o.GetPrimaryKey()
	if k != nil {
		return k.Name
	}
	return ""
}

func (o *Obj) GetPrimaryKey() *Field {
	for _, f := range o.Fields {
		if f.Name == o.Name+"Id" {
			return f
		}
	}
	for _, f := range o.Fields {
		if !strings.HasPrefix(f.Type, "int") {
			continue
		}
		if f.Flags.Contains("unique") || f.Flags.Contains("primary") {
			return f
		}
	}
	return nil
}

func (o *Obj) GetByFields(f []*Field) []string {
	newFields := make([]string, 0, len(f))
	for _, ff := range f {
		if ff == nil {
			continue
		}
		newFields = append(newFields, ff.Name)
	}
	return newFields
}

func (o *Obj) GetForeignKeys() []*Field {
	newFields := make([]*Field, 0, len(o.Fields))
	for _, f := range o.Fields {
		if f.HasForeign() {
			newFields = append(newFields, f)
		}
	}
	if len(newFields) == 0 {
		return nil
	}
	return newFields
}

func (o *Obj) GetFieldNames() []string {
	fieldNames := make([]string, 0, len(o.Fields))
	for _, f := range o.Fields {
		fieldNames = append(fieldNames, f.Name)
	}

	return fieldNames
}

func (o *Obj) GetAllNamesAsArgs(prefix string) []string {
	fieldNames := make([]string, 0, len(o.Fields))
	for _, f := range o.Fields {
		fieldNames = append(fieldNames, f.AsArgName(prefix))
	}

	return fieldNames
}

func (o *Obj) GetFieldNamesAsArgs(prefix string) []string {
	fieldNames := make([]string, 0, len(o.Fields))
	pf := o.GetPrimaryKey()
	for _, f := range o.Fields {
		if pf == nil || f.Name != pf.Name {
			fieldNames = append(fieldNames, f.AsArgName(prefix))
		}
	}
	return fieldNames
}

func (o *Obj) GetNonIdFieldNames() []string {
	pf := o.GetPrimaryKey()
	fieldNames := make([]string, 0, len(o.Fields))
	for _, f := range o.Fields {
		if pf == nil || f.Name != pf.Name {
			fieldNames = append(fieldNames, f.Name)
		}
	}

	return fieldNames
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
	var gens []string
	for _, db := range o.Dbs {
		switch db {
		case "mongo":
			gens = append(gens, "struct")
			if !o.IsEmbed {
				gens = append(gens, "mongo_orm")
			}
		case "enum":
			gens = append(gens, "enum")
		case "mysql":
			gens = append(gens, "struct", "mysql_orm", "mysql_fk")
		}
	}
	if len(gens) == 0 {
		gens = append(gens, "struct")
	}

	return gens
}

func (o *Obj) GetConfigTemplates() []string {
	tpls := []string{}
	for _, db := range o.Dbs {
		switch db {
		case "mysql":
			tpls = append(tpls, "mysql_config")

		case "mongo":
			tpls = append(tpls, "mongo_config")

		}
	}
	return tpls
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
	// gen import for across packages foreign keys
	return nil
}

func (o *Obj) NeedOrm() bool {
	return true
}

func (o *Obj) NeedSearch() bool {
	return false
}

func (o *Obj) NeedIndex() bool {
	return len(o.Indexes) > 0
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

func (o *Obj) HasTimeFields() bool {
	for _, f := range o.Fields {
		if f.GetGoType() == "*time.Time" {
			return true
		}
	}
	return false
}

func (o *Obj) GetTimeFields() []*Field {
	timeFields := make([]*Field, 0)
	for _, f := range o.Fields {
		if f.GetGoType() == "*time.Time" {
			timeFields = append(timeFields, f)
		}
	}
	return timeFields
}

func ToStringSlice(val []interface{}) (result []string) {
	result = make([]string, len(val))
	for i, v := range val {
		result[i] = v.(string)
	}
	return
}

func (o *Obj) setIndexes() {
	for _, i := range o.Indexes {
		for _, name := range i.FieldNames {
			i.Fields = append(i.Fields, o.FieldNameMap[name])
		}
	}

	for _, f := range o.Fields {
		if f.HasIndex() {
			index := new(Index)
			index.FieldNames = []string{f.Name}
			index.Fields = []*Field{f}
			index.IsUnique = f.IsUnique()
			index.IsSparse = !f.Flags.Contains("sort")
			index.Name = f.Name
			o.Indexes = append(o.Indexes, index)
		}
	}
}

func (o *Obj) DbContains(db string) bool {
	for _, v := range o.Dbs {
		if strings.EqualFold(v, db) {
			return true
		}
	}
	return false
}

//! for the multiple dbs support struct template switch
func (o *Obj) DbSwitch(db string) bool {
	for _, v := range o.Dbs {
		if strings.EqualFold(v, db) {
			o.Db = db
			return true
		}
	}
	return false
}

func (o *Obj) Read(name string, data generator.Schema) error {
	o.init()
	for key, val := range data {
		switch key {
		case "db":
			o.Db = val.(string)
			dbs := []string{}
			dbs = append(dbs, o.Db)
			dbs = append(dbs, o.Dbs...)
			o.Dbs = dbs
		case "dbs":
			o.Dbs = ToStringSlice(val.([]interface{}))
			if len(o.Dbs) != 0 {
				o.Db = o.Dbs[0]
			}
		}
	}

	for key, val := range data {
		switch key {
		case "embed":
			v, ok := val.(bool)
			if !ok {
				return fmt.Errorf("embed option must be true or false")
			}
			o.IsEmbed = v
		case "indexes":
			for _, i := range val.([]interface{}) {
				index := new(Index)
				index.FieldNames = ToStringSlice(i.([]interface{}))
				index.Name = strings.Join(index.FieldNames, "")
				index.IsSparse = true
				o.Indexes = append(o.Indexes, index)
			}
		case "uniques":
			for _, i := range val.([]interface{}) {
				index := new(Index)
				index.FieldNames = ToStringSlice(i.([]interface{}))
				index.Name = strings.Join(index.FieldNames, "")
				index.IsSparse = true
				index.IsUnique = true
				o.Indexes = append(o.Indexes, index)
			}
		case "extend":
			o.Extend = val.(string)
		case "table":
			o.Table = val.(string)
		case "dbname":
			o.DbName = val.(string)
		case "storetype":
			o.StoreType = val.(string)
		case "valuetype":
			o.ValueType = val.(string)
		case "modeltype":
			o.ModelType = val.(string)
		case "importSQL":
			o.ImportSQL = val.(string)
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
				f.AsSort = true
				o.Fields[0] = f
				startPos = 1
			} else {
				o.Fields = make([]*Field, len(fieldData))
			}

			for i, field := range fieldData {
				f := new(Field)
				f.Obj = o
				f.Tag = strconv.Itoa(i + 2)

				err := f.Read(field.(generator.Schema))
				if err != nil {
					return errors.New(o.Name + " obj has " + err.Error())
				}
				o.Fields[i+startPos] = f
				o.FieldNameMap[f.Name] = f
			}
		case "comment":
			o.Comment = val.(string)
		case "db", "dbs":
		default:
			return errors.New(o.Name + " has invalid obj property: " + string(key))
		}
	}
	if o.ValueType != "" {
		switch strings.ToLower(o.StoreType) {
		case "set", "list":
			o.Fields = make([]*Field, 2)
			f1 := new(Field)
			f1.init()
			f1.Obj = o
			f1.Name = "Key"
			f1.Tag = "1"
			f1.Type = "string"
			o.Fields[0] = f1

			f2 := new(Field)
			f2.init()
			f2.Obj = o
			f2.Name = "Value"
			f2.Tag = "2"
			f2.Type = o.ValueType
			o.Fields[1] = f2
			o.ValueField = f2
		default:
			return errors.New("please specify `storetype` to " + o.Name)
		}
	}
	if o.DbContains("mysql") && o.Table == "" {
		o.Table = camel2name(o.Name)
	}
	// all mysql dbs share the same connection pool
	if o.DbContains("mysql") && o.DbName == "" {
		return errors.New("please specify `dbname` to " + o.Name)
	}
	o.setIndexes()
	return nil
}
