package parser

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
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
	for table, qm := range t {
		buffer.WriteString(fmt.Sprintf("table: %s\n", table))
		buffer.WriteString(qm.String())
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

type SelectStmt struct {
	Fields []*SelectField
	Tables []*SelectTable
}

type SelectField struct {
	Table string
	Name  string
	Alias string

	IsAggregate bool
}

func (f *SelectField) String() string {
	return fmt.Sprintf("`%s`.`%s`", f.Table, f.Name)
}

type SelectTable struct {
	Name  string
	Alias string
}

// the SELECT clause vistor impl ast.Vistor.
type selectAstVistor struct {
	parsingField *SelectField
	parsingTable *SelectTable

	Fields []*SelectField
	Tables []*SelectTable
}

func (v *selectAstVistor) Enter(n ast.Node) (ast.Node, bool) {
	switch n := n.(type) {
	case *ast.SelectField:
		v.parsingField = &SelectField{Alias: n.AsName.String()}

	case *ast.ColumnName:
		if v.parsingField == nil {
			return n, true
		}
		v.parsingField.Name = n.Name.String()
		v.parsingField.Table = n.Table.String()
		v.Fields = append(v.Fields, v.parsingField)
		// reset the parsing field to let vistor parse next field.
		v.parsingField = nil

	case *ast.AggregateFuncExpr:
		v.parsingField.IsAggregate = true
		v.parsingField.Name = strings.ToLower(n.F)
		v.Fields = append(v.Fields, v.parsingField)
		v.parsingField = nil

	case *ast.TableSource:
		v.parsingTable = &SelectTable{Alias: n.AsName.String()}

	case *ast.TableName:
		if v.parsingTable == nil {
			return n, true
		}
		v.parsingTable.Name = n.Name.String()
		v.Tables = append(v.Tables, v.parsingTable)
		v.parsingTable = nil
	}
	return n, false
}

func (v *selectAstVistor) Leave(n ast.Node) (ast.Node, bool) {
	return n, true
}

func ParseSelect(sql string) (*SelectStmt, error) {
	// quick check if this sql is a query.
	if sql == "" {
		return nil, errors.New("parser: received empty sql")
	}
	lead := strings.Fields(sql)[0]
	if strings.ToUpper(lead) != "SELECT" {
		return nil, errors.New("parser: the sql is not a query")
	}

	p := parser.New()
	nodes, warns, err := p.Parse(sql, "", "")
	if err != nil {
		return nil, err
	}
	if len(warns) > 0 {
		for _, warn := range warns {
			fmt.Fprintf(os.Stdout, "parser: WARN %q\n", warn)
		}
	}
	if len(nodes) == 0 {
		return nil, errors.New("parser: internal: empty nodes")
	}
	rootNode := nodes[0]

	v := new(selectAstVistor)
	_, ok := rootNode.Accept(v)
	if !ok {
		return nil, errors.New("parser: internal: failed to vist the ast node")
	}
	if len(v.Fields) == 0 {
		return nil, errors.New("parser: empty select field")
	}
	if len(v.Tables) == 0 {
		return nil, errors.New("parser: empty select table")
	}

	// We should convert the TableAlias in fields to TableName.
	tableMap := make(map[string]struct{}, len(v.Tables))
	aliasMap := make(map[string]string, len(v.Tables))
	for _, t := range v.Tables {
		tableMap[t.Name] = struct{}{}
		aliasMap[t.Alias] = t.Name
	}

	for _, f := range v.Fields {
		if f.IsAggregate {
			continue
		}
		_, ok := tableMap[f.Table]
		if ok {
			continue
		}
		name := aliasMap[f.Table]
		if name == "" {
			return nil, fmt.Errorf("parser: parse field %v failed: "+
				"cannot find table named %q", f, f.Table)
		}
		f.Table = name
	}

	return &SelectStmt{
		Fields: v.Fields,
		Tables: v.Tables,
	}, nil
}
