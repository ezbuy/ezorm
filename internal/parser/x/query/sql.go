package query

import (
	"bytes"
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/ezbuy/ezorm/v2/internal/generator"
	"github.com/iancoleman/strcase"
)

type T uint8

const (
	T_PLACEHOLDER T = iota
	T_INT
	T_STRING
	T_ARRAY_STRING
	T_ARRAY_INT
	T_ANY
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
	case T_ANY:
		return "any"
	}
	panic("parser: unknown type")
}

func (t T) BaseType() string {
	switch t {
	case T_PLACEHOLDER:
		return "?"
	case T_INT, T_ARRAY_INT:
		return "int64"
	case T_STRING, T_ARRAY_STRING:
		return "string"
	case T_ANY:
		return "any"
	}
	panic("parser: unknown type")
}

type QueryField struct {
	Name   string
	Type   T
	Alias  string
	IsLike bool
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
		if qm, ok := t[table]; ok {
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

func typeMatch(fromAST T, fromSchema generator.IField) bool {
	if strings.Contains(fromAST.BaseType(), "int") && strings.Contains(fromSchema.GetGoType(), "int") {
		return true
	}
	return fromAST.BaseType() == fromSchema.GetGoType()
}

func (tm TableMetadata) Validate(tableRef map[string]map[string]generator.IField) error {
	for t, f := range tm {
		name := uglify(t.Name)
		ff, ok := tableRef[name]
		if !ok {
			return fmt.Errorf("metadata: table %s not found in tableRef(YAML)", name)
		}
		for _, p := range f.params {
			pName := uglify(p.Name)
			pName = strcase.ToLowerCamel(pName)
			col, ok := ff[pName]
			if !ok && p.Type != T_ANY && p.Name != LIMIT_COUNT && p.Name != LIMIT_OFFSET {
				return fmt.Errorf("metadata: param %s not found in table %s", pName, name)
			}
			if p.Name != LIMIT_COUNT && p.Name != LIMIT_OFFSET && !typeMatch(p.Type, col) && p.Type != T_ANY {
				return fmt.Errorf("metadata: param %s type mismatch, expect %s, got %s", pName, col.GetGoType(), p.Type.BaseType())
			}
		}
		for _, r := range f.result {
			rName := uglify(r.Name)
			rName = strcase.ToLowerCamel(rName)
			if _, ok := ff[rName]; !ok && r.Type != T_ANY {
				return fmt.Errorf("metadata: result %s not found in table %s", rName, name)
			}
		}
	}
	return nil
}

func (qm *QueryMetadata) String() string {
	var buffer bytes.Buffer
	sort.Slice(qm.params, func(i, j int) bool {
		return qm.params[i].Name < qm.params[j].Name
	})
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
	raw          *Raw
	resultFields []*QueryField
}

func (qb *QueryBuilder) IsQueryIn() bool {
	return len(qb.raw.ins) > 0
}

func (qb *QueryBuilder) rebuild() string {
	query := qb.String()

	i := strings.Index(query, "WHERE")

	rebuildQuery := bytes.NewBuffer(nil)
	rebuildQuery.WriteString(query[:i])
	rebuildQuery.WriteString("%s")
	// rebuildQuery.WriteString(query[s:])
	return rebuildQuery.String()
}

type LocationOffset struct {
	start, end int
}

type Raw struct {
	ins map[string]struct{}
}

// RawQueryParser is a parser to extract metedata from sql query
type RawQueryParser interface {
	Parse(context.Context, string) (TableMetadata, *QueryBuilder, error)
	Flush()
}
