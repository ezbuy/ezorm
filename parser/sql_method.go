package parser

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"
)

type SQL struct {
	fieldMap       map[string]map[string]*Field
	RawQueryParser RawQueryParser
}

type SQLFile struct {
	GoPackage string
	Methods   []*SQLMethod
}

type SQLMethod struct {
	Name   string
	Fields []*SQLMethodField
	Result []*SQLMethodField
	SQL    string

	Assign string
}

type SQLMethodField struct {
	Name string
	Raw  string
	Type string
}

func NewSQL(objs map[string]*Obj) *SQL {
	fieldMap := make(map[string]map[string]*Field, len(objs))
	for table, obj := range objs {
		name := camel2name(table)
		fieldMap[name] = make(map[string]*Field, len(obj.Fields))
		for _, f := range obj.Fields {
			fname := camel2name(f.Name)
			fieldMap[name][fname] = f
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
	}

	result.SQL = builder.rebuild()
	return result, nil
}
