package mysqlr

import (
	"fmt"
	"sort"
	"strings"

	"github.com/ezbuy/ezorm/v2/internal/generator"
)

var _ generator.IObject = (*MetaObject)(nil)

type MetaObject struct {
	//! package name
	// Package is deprecated, use Namespace instead.
	Package   string
	GoPackage string
	Namespace string
	//! model name
	Name string
	Tag  string
	//! dbs
	Db      string
	comment string
	//! database
	DbName  string
	DbTable string
	DbView  string
	//! fields
	fields       []*Field
	fieldNameMap map[string]*Field
	//! primary
	primary *PrimaryKey
	//! indexes
	uniques []*Index
	indexes []*Index
	ranges  []*Index
	//! importSQL
	ImportSQL string
}

func NewMetaObject(namespace string, packageName string) *MetaObject {
	return &MetaObject{
		Namespace:    namespace,
		GoPackage:    packageName,
		fieldNameMap: make(map[string]*Field),
	}
}

func (o *MetaObject) FieldByName(name string) *Field {
	if f, ok := o.fieldNameMap[name]; ok {
		return f
	}
	return nil
}

func (o *MetaObject) GetTable() string {
	return o.DbTable
}

func (o *MetaObject) FieldsMap() map[string]generator.IField {
	result := make(map[string]generator.IField, len(o.fields))
	for _, f := range o.fields {
		result[f.GetName()] = f
	}
	return result
}

func (o *MetaObject) PrimaryField() *Field {
	for _, f := range o.Fields() {
		if f.IsPrimary() {
			return f
		}
	}
	return nil
}

func (o *MetaObject) PrimaryKey() *PrimaryKey {
	return o.primary
}

func (o *MetaObject) DbSource() string {
	if o.DbTable != "" {
		return o.DbTable
	}
	return ""
}

func (o *MetaObject) FromDB() string {
	return o.DbSource()
}

func (o *MetaObject) Fields() []*Field {
	return o.fields
}

func (o *MetaObject) NoneIncrementFields() []*Field {
	fields := make([]*Field, 0, len(o.fields))
	for _, f := range o.fields {
		if !f.IsAutoIncrement() {
			fields = append(fields, f)
		}
	}
	return fields
}

func (o *MetaObject) Uniques() []*Index {
	sort.Slice(o.uniques, func(i, j int) bool {
		return strings.Compare(o.uniques[i].Name, o.uniques[j].Name) > 0
	})
	return o.uniques
}

func (o *MetaObject) NotPrimaryUniques() []*Index {
	unique := make([]*Index, 0, len(o.uniques))
	for _, u := range o.uniques {
		if u.HasPrimaryKey() {
			continue
		}
		unique = append(unique, u)
	}
	return unique
}

func (o *MetaObject) Indexes() []*Index {
	sort.Slice(o.indexes, func(i, j int) bool {
		return strings.Compare(o.indexes[i].Name, o.indexes[j].Name) > 0
	})
	return o.indexes
}

func (o *MetaObject) NotPrimaryIndexes() []*Index {
	indexes := make([]*Index, 0, len(o.indexes))
	for _, i := range o.indexes {
		if i.HasPrimaryKey() {
			continue
		}
		indexes = append(indexes, i)
	}
	return indexes
}

func (o *MetaObject) Ranges() []*Index {
	sort.Slice(o.ranges, func(i, j int) bool {
		return strings.Compare(o.ranges[i].Name, o.ranges[j].Name) > 0
	})
	return o.ranges
}

func (o *MetaObject) LastField() *Field {
	return o.fields[len(o.fields)-1]
}

func (o *MetaObject) Read(name string, data generator.Schema) error {
	o.Name = name
	for key, val := range data {
		switch key {
		case "db":
			o.Db = val.(string)
		}
	}

	for key, val := range data {
		switch key {
		case "tag":
			tag := val.(int)
			o.Tag = fmt.Sprint(tag)
		case "dbname":
			o.DbName = val.(string)
		case "dbtable":
			o.DbTable = val.(string)
		case "dbview":
			o.DbView = val.(string)
		case "comment":
			o.comment = val.(string)

		case "importSQL":
			o.ImportSQL = val.(string)
		case "fields":
			fieldData := val.([]interface{})
			o.fields = make([]*Field, len(fieldData))
			for i, field := range fieldData {
				f := NewField()
				f.Obj = o
				err := f.Read(field.(generator.Schema))
				if err != nil {
					return fmt.Errorf("object (%s) %s", o.Name, err.Error())
				}
				o.fields[i] = f
				o.fieldNameMap[f.Name] = f
			}
		case "primary":
			o.primary = NewPrimaryKey(o)
			o.primary.FieldNames = toStringSlice(val.([]interface{}))
		case "uniques":
			for _, i := range val.([]interface{}) {
				if len(i.([]interface{})) == 0 {
					continue
				}
				index := NewIndex(o)
				index.FieldNames = toStringSlice(i.([]interface{}))
				o.uniques = append(o.uniques, index)
			}
		case "indexes":
			for _, i := range val.([]interface{}) {
				if len(i.([]interface{})) == 0 {
					continue
				}
				index := NewIndex(o)
				index.FieldNames = toStringSlice(i.([]interface{}))
				o.indexes = append(o.indexes, index)
			}
		case "ranges":
			for _, i := range val.([]interface{}) {
				if len(i.([]interface{})) == 0 {
					continue
				}
				index := NewIndex(o)
				index.FieldNames = toStringSlice(i.([]interface{}))
				o.ranges = append(o.ranges, index)
			}

		}
	}

	for _, field := range o.fields {
		if field.IsPrimary() {
			if o.primary == nil {
				o.primary = NewPrimaryKey(o)
				o.primary.FieldNames = []string{}
			}
			o.primary.FieldNames = append(o.primary.FieldNames, field.Name)
		}
		if field.HasIndex() && field.IsNullable() {
			return fmt.Errorf("object (%s) field (%s) should not be nullable for indexing", o.Name, field.Name)
		}
	}

	if o.primary == nil {
		return fmt.Errorf("object (%s) needs a primary key declare", o.Name)
	} else {
		if err := o.primary.build(); err != nil {
			return fmt.Errorf("object (%s) %s", o.Name, err.Error())
		}

		if o.primary.IsRange() {
			index := NewIndex(o)
			index.FieldNames = o.primary.FieldNames
			o.ranges = append(o.ranges, index)
		}
	}

	for _, unique := range o.uniques {
		if err := unique.buildUnique(); err != nil {
			return fmt.Errorf("object (%s) %s", o.Name, err.Error())
		}
	}
	for _, index := range o.indexes {
		if err := index.buildIndex(); err != nil {
			return fmt.Errorf("object (%s) %s", o.Name, err.Error())
		}
	}
	for _, rg := range o.ranges {
		if err := rg.buildRange(); err != nil {
			return fmt.Errorf("object (%s) %s", o.Name, err.Error())
		}
	}
	return nil
}

func (m *MetaObject) Comment() string {
	if m.comment != "" {
		return m.comment
	}

	return m.DbTable
}

func toStringSlice(val []interface{}) (result []string) {
	result = make([]string, len(val))
	for i, v := range val {
		result[i] = v.(string)
	}
	return
}
