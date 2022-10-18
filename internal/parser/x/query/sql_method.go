package query

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/ezbuy/ezorm/v2/internal/generator"
	"github.com/iancoleman/strcase"
)

type SQL struct {
	fieldMap       map[string]map[string]generator.IField
	RawQueryParser RawQueryParser
}

type SQLFile struct {
	GoPackage string
	Methods   []*SQLMethod
	Dir       string
}

type SQLMethod struct {
	Name   string
	Fields []*SQLMethodField
	Result []*SQLMethodField
	SQL    string

	Assign   string
	FromFile string
	QueryIn  bool
}

type SQLMethodField struct {
	Name string
	Raw  string
	Type string
}

func NewSQL(objs map[string]generator.IObject) *SQL {
	fieldMap := make(map[string]map[string]generator.IField, len(objs))
	for table, obj := range objs {
		name := camel2name(table)
		fieldMap[name] = make(map[string]generator.IField, len(obj.FieldsMap()))
		for _, f := range obj.FieldsMap() {
			fieldMap[name][f.GetName()] = f
		}
	}
	return &SQL{
		fieldMap:       fieldMap,
		RawQueryParser: NewTiDBParser(),
	}
}

func (p *SQL) retypeResult(table string, col string) (string, error) {
	t, ok := p.fieldMap[table]
	if !ok {
		return "", fmt.Errorf("res: retype: table: %s not found", table)
	}
	f, ok := t[col]
	if !ok {
		return "", fmt.Errorf("res: retype: field: %s not found", col)
	}
	return f.GetGoType(), nil
}

func (p *SQL) Read(path string) (*SQLMethod, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	sql := string(data)
	ctx := context.Background()
	meta, builder, err := p.RawQueryParser.Parse(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer p.RawQueryParser.Flush()

	// Insert the name to the raw sql as an internal comment.
	name := filepath.Base(path)
	name = strings.TrimSuffix(name, ".sql")
	if name == "" {
		return nil, errors.New("parse sql: the filename is empty")
	}

	name = strcase.ToCamel(name)

	result := &SQLMethod{
		Name: name,
		SQL:  sql,
	}
	// validation: validate if every table and column are defined in yaml file(TableRef).
	if err := meta.Validate(p.fieldMap); err != nil {
		return nil, err
	}
	for t, f := range meta {
		for _, c := range f.params {
			name := uglify(c.Name)
			result.Fields = append(result.Fields, &SQLMethodField{
				Name: strcase.ToCamel(name),
				Raw:  name,
				Type: c.Type.String(),
			})
		}
		for _, c := range f.result {
			name := uglify(c.Name)
			if c.Type == T_ANY {
				result.Result = append(result.Result, &SQLMethodField{
					Name: strcase.ToCamel(name),
					Type: c.Type.String(),
				})
				continue
			}
			tp, err := p.retypeResult(t.Name, name)
			if err != nil {
				return nil, err
			}
			result.Result = append(result.Result, &SQLMethodField{
				Name: strcase.ToCamel(name),
				Raw:  name,
				Type: tp,
			})
		}
		result.QueryIn = builder.IsQueryIn()
	}

	var scan bytes.Buffer
	for _, r := range result.Result {
		scan.WriteString(fmt.Sprintf("&o.%s, ", r.Name))
	}

	result.Assign = scan.String()
	result.SQL = builder.rebuild()
	return result, nil
}

func camel2name(s string) string {
	nameBuf := bytes.NewBuffer(nil)
	afterSpace := false
	for i, c := range s {
		if unicode.IsUpper(c) && unicode.IsLetter(c) {
			if i > 0 && !afterSpace {
				nameBuf.WriteRune('_')
			}
			c = unicode.ToLower(c)
		}
		nameBuf.WriteRune(c)
		afterSpace = unicode.IsSpace(c)
	}
	return nameBuf.String()
}
