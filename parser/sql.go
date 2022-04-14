package parser

import (
	"bytes"
	"context"
	"fmt"
	"sort"
	"strings"

	_ "github.com/pingcap/parser/test_driver"
)

type T uint8

const (
	T_PLACEHOLDER T = iota
	T_INT
	T_STRING
	T_ARRAY_STRING
	T_ARRAY_INT
)

func (t T) String() string {
	switch t {
	case T_PLACEHOLDER:
		return "?"
	case T_INT:
		return "int64"
	case T_STRING:
		return "string"
	case T_ARRAY_STRING:
		return "[]string"
	case T_ARRAY_INT:
		return "[]int64"
	}
	panic("parser: unknown type")
}

type QueryField struct {
	Name string
	Type T
}

var _ fmt.Stringer = (*QueryMetadata)(nil)

type LimitOption struct {
	count, offset bool
}

type QueryMetadata struct {
	table  string
	params []*QueryField
	result []*QueryField
}

type Table struct {
	Name  string
	Alias string
}

type TableMetadata map[Table]*QueryMetadata

func (t TableMetadata) String() string {
	var buffer bytes.Buffer
	var tables []Table
	for table := range t {
		tables = append(tables, table)
	}
	sort.SliceStable(tables, func(i, j int) bool {
		return tables[i].Name < tables[j].Name
	})
	for _, table := range tables {
		if qm, ok := t[table]; !ok {
			buffer.WriteString(fmt.Sprintf("table: %s\n", table))
			buffer.WriteString(qm.String())
		}
	}
	return buffer.String()
}

func (t TableMetadata) AppendParams(table string, params ...*QueryField) {
	var key Table
	for tb := range t {
		if tb.Name == table || tb.Alias == table {
			key = tb
		}
	}
	if key.Alias == "" && key.Name == "" {
		return
	}
	if _, ok := t[key]; ok {
		t[key].params = append(t[key].params, params...)
	}
}

func (t TableMetadata) AppendResult(table string, result ...*QueryField) {
	var key Table
	for tb := range t {
		if tb.Name == table || tb.Alias == table {
			key = tb
		}
	}
	if key.Alias == "" && key.Name == "" {
		return
	}
	if _, ok := t[key]; ok {
		t[key].result = append(t[key].result, result...)
	}
}

func uglify(col string) string {
	if strings.Contains(col, ":") {
		parts := strings.Split(col, ":")
		if len(parts) != 2 {
			return col
		}
		col = parts[1]
	}
	if strings.Contains(col, ".") {
		parts := strings.Split(col, ".")
		if len(parts) != 2 {
			return col
		}
		col = parts[1]
	}
	col = strings.ReplaceAll(col, "`", "")
	return col
}

func (tm TableMetadata) Validate(tableRef map[string]map[string]*Field) error {
	for t, f := range tm {
		name := uglify(t.Name)
		ff, ok := tableRef[name]
		if !ok {
			return fmt.Errorf("metadata: table %s not found in tableRef(YAML)", name)
		}
		for _, p := range f.params {
			pName := uglify(p.Name)
			col, ok := ff[pName]
			if !ok {
				return fmt.Errorf("metadata: param %s not found in table %s", pName, name)
			}
			if col.GetGoType() != p.Type.String() {
				return fmt.Errorf("metadata: param %s type mismatch, expect %s, got %s", pName, col.GetGoType(), p.Type.String())
			}
		}
		for _, r := range f.result {
			rName := uglify(r.Name)
			if _, ok := ff[rName]; !ok {
				return fmt.Errorf("metadata: result %s not found in table %s", r.Name, name)
			}
		}
	}
	return nil
}

func (qm *QueryMetadata) String() string {
	var buffer bytes.Buffer
	for _, p := range qm.params {
		buffer.WriteString(fmt.Sprintf("param: name: %s, type: %s\n", p.Name, p.Type))
	}
	for _, r := range qm.result {
		buffer.WriteString(fmt.Sprintf("result: name: %s, type: %s\n", r.Name, r.Type))
	}
	return buffer.String()
}

type QueryBuilder struct {
	*bytes.Buffer
	raw *Raw
}

func (qb *QueryBuilder) rebuild() string {
	query := qb.String()
	los := make([]LocationOffset, len(qb.raw.lo))
	reversed := make(map[LocationOffset]string)
	var index int
	for col, lo := range qb.raw.lo {
		los[index] = lo
		index++
		reversed[lo] = col
	}
	sort.SliceStable(los, func(i, j int) bool {
		return los[i].start < los[j].start
	})

	var s int
	rebuildQuery := bytes.NewBuffer(nil)
	for _, lo := range los {
		e := lo.start
		rebuildQuery.WriteString(query[s:e])
		ins, ok := qb.raw.ins[reversed[lo]]
		switch {
		case ok:
			rebuildQuery.WriteString(ins.String())
		case reversed[lo] == "LIMIT" && qb.raw.limit.count && qb.raw.limit.offset:
			rebuildQuery.WriteString("?,?")
		default:
			rebuildQuery.WriteString("?")
		}
		s = lo.end
	}
	rebuildQuery.WriteString(query[s:])
	return rebuildQuery.String()
}

type LocationOffset struct {
	start, end int
}

type Raw struct {
	ins   map[string]*InBuilder
	lo    map[string]LocationOffset
	limit *LimitOption
}

type InBuilder struct {
	col    string
	params []any
}

func NewIn(col string, params []any) *InBuilder {
	return &InBuilder{
		col:    col,
		params: params,
	}
}

func (in *InBuilder) String() string {
	var placeholders []string
	for range in.params {
		placeholders = append(placeholders, "?")
	}
	var query string
	if len(placeholders) > 0 {
		query = strings.Join(placeholders, ",")
	}
	return query
}

// RawQueryParser is a parser to extract metedata from sql query
type RawQueryParser interface {
	Parse(context.Context, string) (TableMetadata, *QueryBuilder, error)
}