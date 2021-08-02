package sqlm

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"regexp"
	"sort"
	"strings"

	"github.com/ezbuy/ezorm/parser"
	"gopkg.in/yaml.v2"
)

type yamlFile struct {
	Namespace string                 `yaml:"namespace"`
	SQLs      map[string]string      `yaml:"sqls"`
	Methods   map[string]*yamlMethod `yaml:"methods"`

	models []*parser.Obj
}

type yamlMethod struct {
	ArgsMaps []map[string]string `yaml:"args"`
	args     []string

	SQL string `yaml:"sql"`

	LastId   bool `yaml:"lastid"`
	Affected bool `yaml:"affected"`
}

type Obj struct {
	GoPackage string
	Namespace string
	Methods   []*Method
}

type Method struct {
	Name string

	Models  []string
	Fields  []*Field
	RetName string
	RetDef  string

	ArgsDef string
	ArgsUse string

	SQL  string
	Scan string

	ExecResult   bool
	ExecLastId   bool
	ExecAffected bool

	Query bool

	DB string
}

type Field struct {
	Name string
	Type string
	Tags string
}

func Parse(path, pkg string, models []*parser.Obj) (*Obj, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var file yamlFile
	err = yaml.Unmarshal(data, &file)
	if err != nil {
		return nil, fmt.Errorf("parse sql yaml failed: %v", err)
	}
	file.models = models

	return file.convert2obj(pkg)
}

func (file *yamlFile) convert2obj(pkg string) (*Obj, error) {
	obj := new(Obj)
	obj.Namespace = file.Namespace
	if obj.Namespace == "" {
		obj.Namespace = "Methods"
	}
	obj.GoPackage = pkg

	obj.Methods = make([]*Method, 0, len(file.Methods))
	for k, m := range file.Methods {
		mobj, err := m.convert2obj(k, file)
		if err != nil {
			return nil, &MethodError{Method: k, Err: err}
		}
		mobj.SQL = trimSpaces(mobj.SQL)
		obj.Methods = append(obj.Methods, mobj)
	}
	sort.Slice(obj.Methods, func(i, j int) bool {
		return obj.Methods[i].Name > obj.Methods[j].Name
	})

	return obj, nil
}

func (m *yamlMethod) convert2obj(name string, file *yamlFile) (*Method, error) {
	obj := new(Method)
	obj.Name = name
	if err := checkPhs(m.SQL); err != nil {
		return nil, err
	}

	sql := strings.Replace(m.SQL, "\n", " ", -1)
	sql = strings.Replace(sql, "\t", " ", -1)

	// Because snippets may contain other placeholders, we
	// must expand snippets first to ensure that placeholders
	// inner snippets are processed.
	sql, err := handlePhs(sql, ".sqls.", func(name string) (string, error) {
		if file.SQLs != nil {
			snippet, ok := file.SQLs[name]
			if ok {
				return snippet, nil
			}
		}
		return "", &Error{
			full:  m.SQL,
			wrong: name,
			desc:  fmt.Sprintf("cannot find sql named %q", name),
		}
	})
	if err != nil {
		return nil, err
	}

	m.args = make([]string, 0, len(m.ArgsMaps))
	for _, argMap := range m.ArgsMaps {
		for k, v := range argMap {
			m.args = append(m.args,
				fmt.Sprintf("%s %s", k, v))
		}
	}

	// Replace args placeholders in sql into "?", check them
	// exist in the yaml args option.
	if len(m.args) > 0 {
		obj.ArgsDef = ", " + strings.Join(m.args, ", ")
	}
	var args []string
	sql, err = handlePhs(sql, ".args.", func(name string) (string, error) {
		var found bool
		for _, arg := range m.args {
			// The prefix instead of equality is used here
			// because the arg defined in yaml may be a
			// struct or map, and when used in sql, it may
			// be a specific element of a structure or map.
			// For example, the arg defined by yaml is
			// 		"m map[string]string"
			// "{{.args.m["name"]}}" is used in sql
			if strings.HasPrefix(arg, name) {
				found = true
				break
			}
		}
		if !found {
			return "", &Error{
				full:  m.SQL,
				wrong: name,
				desc:  fmt.Sprintf("cannot find arg named %q", name),
			}
		}
		args = append(args, name)
		return "?", nil
	})
	if err != nil {
		return nil, err
	}
	if len(args) > 0 {
		obj.ArgsUse = fmt.Sprintf("[]interface{}{%s}",
			strings.Join(args, ", "))
	} else {
		obj.ArgsUse = "[]interface{}{}"
	}

	// Parse SQL.
	cs, err := split(sql, m.SQL)
	if err != nil {
		return nil, err
	}
	if cs.isExec {
		md := createModelsData(file.models, nil)
		sql, err := execTpl(sql, md.tplData)
		if err != nil {
			return nil, err
		}
		obj.SQL = sql
		obj.DB = "orm.Execable"
		obj.ExecAffected = m.Affected
		obj.ExecLastId = m.LastId
		obj.RetDef = "int64"
		if !obj.ExecAffected && !obj.ExecLastId {
			obj.ExecResult = true
			obj.RetDef = "sql.Result"
		}

		return obj, nil
	}

	obj.DB = "orm.Queryable"
	obj.RetName = fmt.Sprintf("%sResp", obj.Name)
	obj.RetDef = fmt.Sprintf("[]*%s", obj.RetName)
	obj.Query = true

	tables, err := cs.parseTables()
	if err != nil {
		return nil, err
	}
	md := createModelsData(file.models, tables)
	sql, err = execTpl(sql, md.tplData)
	if err != nil {
		return nil, err
	}
	obj.SQL = sql

	fields, err := cs.parseQuery(tables)
	if err != nil {
		return nil, err
	}
	var scans []string
	obj.Fields = make([]*Field, 0, len(fields))
	for _, f := range fields {
		if f.Count {
			of := &Field{
				Name: f.Name,
				Type: "int64",
				Tags: fmt.Sprintf("`table:%q`", "count"),
			}
			obj.Fields = append(obj.Fields, of)
			scan := fmt.Sprintf("&e.%s", f.Name)
			scans = append(scans, scan)
			continue
		}
		model := md.mmap[f.TableFull]
		fmap := md.fmap[f.TableFull]
		if model == nil || fmap == nil {
			return nil, &Error{
				full:  m.SQL,
				wrong: f.TableFull,
				desc: fmt.Sprintf("cannot find "+
					"model %q in models file", f.TableFull),
			}
		}
		if f.Name == "allFields" {
			obj.Models = append(obj.Models, model.Name)
			for _, mf := range model.Fields {
				scan := fmt.Sprintf("&e.%s.%s", model.Name, mf.Name)
				scans = append(scans, scan)
			}
			continue
		}

		mf := fmap[f.Name]
		if mf == nil {
			return nil, &Error{
				full:  m.SQL,
				wrong: f.Name,
				desc: fmt.Sprintf("cannot find field %q in model %q",
					f.Name, model.Name),
			}
		}
		of := new(Field)
		of.Name = f.Alias
		if of.Name == "" {
			of.Name = model.Name + mf.Name
		}
		of.Type = mf.Type
		of.Tags = fmt.Sprintf("`table:%q field:%q`",
			model.Table, parser.Camel2name(mf.Name))
		obj.Fields = append(obj.Fields, of)
		scan := fmt.Sprintf("&e.%s", of.Name)
		scans = append(scans, scan)
	}
	obj.Scan = strings.Join(scans, ", ")

	return obj, nil
}

type modelTplData map[string]string

func (d modelTplData) String() string {
	return d["__name__"]
}

type modelsData struct {
	tplData map[string]modelTplData

	fmap map[string]map[string]*parser.Field
	mmap map[string]*parser.Obj
}

func createModelsData(models []*parser.Obj, ts *tables) *modelsData {
	mds := &modelsData{
		tplData: make(map[string]modelTplData),
		fmap:    make(map[string]map[string]*parser.Field),
		mmap:    make(map[string]*parser.Obj),
	}

	for _, model := range models {
		if model.Table == "" {
			model.Table = parser.Camel2name(model.Name)
		}
		tableName := wrapName(model.Table)
		md := make(modelTplData)
		md["__name__"] = tableName

		var aliasMD modelTplData
		var alias string
		if ts != nil {
			alias = ts.names2alias[model.Name]
			if alias != "" {
				aliasMD = make(modelTplData)
				aliasMD["__name__"] = alias
			}
		}

		allFields := make([]string, len(model.Fields))
		mds.fmap[model.Name] = make(map[string]*parser.Field)
		for i, f := range model.Fields {
			name := parser.Camel2name(f.Name)
			name = wrapName(name)
			allFields[i] = name

			full := fmt.Sprintf("%s.%s", tableName, name)
			md[f.Name] = full

			if alias != "" {
				aliasMD[f.Name] = fmt.Sprintf("%s.%s", alias, name)
			}
			mds.fmap[model.Name][f.Name] = f
		}

		joinFields := func(table string) string {
			fs := make([]string, len(allFields))
			for i, f := range allFields {
				fs[i] = fmt.Sprintf("%s.%s", table, f)
			}
			return strings.Join(fs, ", ")
		}

		md["allFields"] = joinFields(tableName)
		mds.tplData[model.Name] = md
		if alias != "" {
			aliasMD["allFields"] = joinFields(alias)
			mds.tplData[alias] = aliasMD
		}
		mds.mmap[model.Name] = model
	}

	return mds
}

func wrapName(s string) string {
	return fmt.Sprintf("`%s`", s)
}

var placeholderRe = regexp.MustCompile(`{{[^}]+}}`)

type phHandler func(name string) (string, error)

func handlePhs(sql, prefix string, fn phHandler) (r string, err error) {
	r = placeholderRe.ReplaceAllStringFunc(sql, func(ph string) string {
		if err != nil {
			// Last placeholder's parse exceeded errors, we
			// just skip this one.
			return ph
		}
		val := trimPlaceholder(ph)
		if !strings.HasPrefix(val, prefix) {
			// This placeholder is not the target, let's keey its
			// prototype.
			return ph
		}
		// Extract the suffix(name) of this placeholder
		name := strings.TrimPrefix(val, prefix)
		val, err = fn(name)
		if err != nil {
			return ph
		}
		return val
	})
	return
}

func checkPhs(sql string) error {
	phs := placeholderRe.FindAllString(sql, -1)
	for _, ph := range phs {
		val := trimPlaceholder(ph)
		if !strings.HasPrefix(val, ".") {
			return &Error{
				full:  sql,
				wrong: ph,
				desc: fmt.Sprintf("missing '.', "+
					"do you mean %q?", "."+val),
			}
		}
	}
	return nil
}

func trimPlaceholder(s string) string {
	s = strings.TrimLeft(s, "{{")
	s = strings.TrimRight(s, "}}")
	return strings.TrimSpace(s)
}

func execTpl(sql string, data map[string]modelTplData) (string, error) {
	// Use GoTemplate to replace all placeholders in
	// sql (mainly those fieldNames)
	tpl, err := template.New("sql-template").Parse(sql)
	if err != nil {
		return "", fmt.Errorf("parse template for sql failed: %v", err)
	}
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	err = tpl.Execute(w, data)
	if err != nil {
		return "", fmt.Errorf("exec template for sql failed: %v", err)
	}
	w.Flush()
	return b.String(), nil
}

func trimSpaces(sql string) string {
	return strings.Join(strings.Fields(sql), " ")
}
