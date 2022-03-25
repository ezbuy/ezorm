package parser

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"
)

type SQL struct {
	fieldMap map[string]map[string]*Field
}

type SQLFile struct {
	GoPackage string
	Methods   []*SQLMethod
}

type SQLMethod struct {
	Name   string
	Fields []*SQLMethodField
	SQL    string

	Assign string
}

type SQLMethodField struct {
	Name string
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
	return &SQL{fieldMap: fieldMap}
}

func (p *SQL) Read(path string) (*SQLMethod, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	// Trim the comment lines.
	raws := strings.Split(string(data), "\n")
	lines := make([]string, 0, len(raws))
	for _, line := range raws {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "--") {
			continue
		}
		lines = append(lines, line)
	}
	sql := strings.Join(lines, " ")

	stmt, err := ParseSelect(sql)
	if err != nil {
		return nil, fmt.Errorf("parse sql %s: %v", path, err)
	}

	// Insert the name to the raw sql as an internal comment.
	name := filepath.Base(path)
	name = strings.TrimSuffix(name, ".sql")
	if name == "" {
		return nil, errors.New("parse sql: the filename is empty")
	}
	sql = insertCommentToSQL(sql, name)
	name = strcase.ToCamel(name)

	result := &SQLMethod{
		Name:   name,
		SQL:    sql,
		Fields: make([]*SQLMethodField, len(stmt.Fields)),
	}
	var aggregateIdx int
	for i, rawField := range stmt.Fields {
		if rawField.IsAggregate {
			name := strcase.ToCamel(rawField.Alias)
			if name == "" {
				name = fmt.Sprintf("%s%d", strcase.ToCamel(rawField.Name),
					aggregateIdx)
				aggregateIdx++
			}
			result.Fields[i] = &SQLMethodField{
				Name: name,
				Type: "int64",
			}
			continue
		}
		m := p.fieldMap[rawField.Table]
		if m == nil {
			return nil, fmt.Errorf("parse %s: cannot find table `%s` in your yaml file",
				path, rawField.Table)
		}
		objField := m[rawField.Name]
		if objField == nil {
			return nil, fmt.Errorf("parse %s: cannot find field `%s` in table `%s`",
				path, rawField.Name, rawField.Table)
		}
		var fullname = strcase.ToCamel(rawField.Alias)
		if fullname == "" {
			fullname = strcase.ToCamel(rawField.Table) + strcase.ToCamel(rawField.Name)
		}
		result.Fields[i] = &SQLMethodField{
			Name: fullname,
			Type: objField.GetGoType(),
		}
	}

	assigns := make([]string, len(result.Fields))
	for i, f := range result.Fields {
		assigns[i] = fmt.Sprintf("&o.%s", f.Name)
	}
	result.Assign = strings.Join(assigns, ", ")

	return result, nil
}

func insertCommentToSQL(sql, comment string) string {
	rawFields := strings.Fields(sql)
	switch len(rawFields) {
	case 0:
		return sql

	case 1:
		return fmt.Sprintf("%s /* %s */", sql, comment)
	}
	return fmt.Sprintf("%s /* %s */ %s", rawFields[0],
		comment, strings.Join(rawFields[1:], " "))
}
