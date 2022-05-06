package mysqlr

import (
	"fmt"
	"strings"
)

type PrimaryKey struct {
	Name       string
	FieldNames []string
	Fields     []*Field
	Obj        *MetaObject
}

func NewPrimaryKey(obj *MetaObject) *PrimaryKey {
	return &PrimaryKey{Obj: obj}
}

func (pk *PrimaryKey) IsSingleField() bool {
	return len(pk.Fields) == 1
}

func (pk *PrimaryKey) GetFuncParam() string {
	return Fields(pk.Fields).GetFuncParam()
}

func (pk *PrimaryKey) FirstField() *Field {
	if len(pk.Fields) > 0 {
		return pk.Fields[0]
	}
	return nil
}

func (pk *PrimaryKey) IsAutoIncrement() bool {
	if len(pk.Fields) == 1 {
		return pk.Fields[0].Flags.Contains("autoinc")
	}
	return false
}

func (pk *PrimaryKey) IsRange() bool {
	fs := make([]*Field, 0, len(pk.Fields))
	for _, f := range pk.Fields {
		if f.IsNorange() {
			continue
		}
		fs = append(fs, f)
	}
	c := len(fs)
	if c > 0 {
		return fs[c-1].IsNumber()
	}
	return false
}

func (pk *PrimaryKey) build() error {
	pk.Name = fmt.Sprintf("%sOf%sPK", strings.Join(pk.FieldNames, ""), pk.Obj.Name)
	for _, name := range pk.FieldNames {
		f := pk.Obj.FieldByName(name)
		if f == nil {
			return fmt.Errorf("%s field not exist", name)
		}
		f.Flags.Add("primary")
		pk.Fields = append(pk.Fields, f)
	}
	if len(pk.Fields) == 0 {
		return fmt.Errorf("primary key  not declare")
	}
	return nil
}

func (pk *PrimaryKey) SQLColumn() string {
	columns := make([]string, 0, len(pk.Fields))
	for _, f := range pk.Fields {
		columns = append(columns, f.SQLName())
	}
	return fmt.Sprintf("PRIMARY KEY(%s)", strings.Join(columns, ","))
}

func (pk *PrimaryKey) GetConstructor() string {
	return Fields(pk.Fields).GetConstructor()
}

func (pk *PrimaryKey) GetObjectParam() string {
	return Fields(pk.Fields).GetObjectParam()
}
