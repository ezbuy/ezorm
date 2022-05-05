package mysqlr

import (
	"fmt"
	"strings"
)

type IndexArray []*Index

func (a IndexArray) Len() int      { return len(a) }
func (a IndexArray) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a IndexArray) Less(i, j int) bool {
	if strings.Compare(a[i].Name, a[j].Name) > 0 {
		return true
	}
	return false
}

type Index struct {
	Name       string
	Fields     []*Field
	FieldNames []string
	Obj        *MetaObject
}

func NewIndex(obj *MetaObject) *Index {
	return &Index{Obj: obj}
}

func (idx *Index) IsSingleField() bool {
	if len(idx.Fields) == 1 {
		return true
	}
	return false
}

func (idx *Index) HasPrimaryKey() bool {
	for _, f := range idx.Fields {
		if f.IsPrimary() {
			return true
		}
	}
	return false
}

func (idx *Index) GetFuncParam() string {
	return Fields(idx.Fields).GetFuncParam()
}

func (idx *Index) GetFuncName() string {
	params := make([]string, len(idx.Fields))
	for i, f := range idx.Fields {
		params[i] = f.Name
	}
	return strings.Join(params, "")
}

func (idx *Index) FirstField() *Field {
	return idx.Fields[0]
}

func (idx *Index) LastField() *Field {
	return idx.Fields[len(idx.Fields)-1]
}

func (idx *Index) buildUnique() error {
	return idx.build("UK")
}
func (idx *Index) buildIndex() error {
	return idx.build("IDX")
}
func (idx *Index) buildRange() error {
	err := idx.build("RNG")
	if err != nil {
		return err
	}
	if !idx.LastField().IsNumber() {
		return fmt.Errorf("range <%s> field <%s> is not number type", idx.Name, idx.LastField().Name)
	}
	return nil
}
func (idx *Index) build(suffix string) error {
	idx.Name = fmt.Sprintf("%sOf%s%s", strings.Join(idx.FieldNames, ""), idx.Obj.Name, suffix)
	for _, name := range idx.FieldNames {
		f := idx.Obj.FieldByName(name)
		if f == nil {
			return fmt.Errorf("%s field not exist", name)
		}
		idx.Fields = append(idx.Fields, f)
	}

	return nil
}

func (idx *Index) GetConstructor() string {
	return Fields(idx.Fields).GetConstructor()
}
