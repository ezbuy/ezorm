package parser

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
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

type QueryMetadata struct {
	params []*QueryField
	result []*QueryField
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
	raw string
}

func (q *QueryBuilder) rebuild() (string, error) {
	return "", nil
}

// RawQueryParser is a parser to extract metedata from sql query
type RawQueryParser interface {
	Parse(context.Context, string) (*QueryMetadata, error)
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
