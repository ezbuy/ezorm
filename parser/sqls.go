package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"regexp"
	"sort"
	"strings"
	"text/template"

	"gopkg.in/yaml.v2"
)

// sqlFile is the raw type parsed by yaml file.
type sqlFile struct {
	Namespace string                `yaml:"namespace"`
	Snippets  map[string]string     `yaml:"snippets"`
	Methods   map[string]*sqlMethod `yaml:"methods"`

	generatedStructMap map[string]struct{}
}

// sqlMethod is the raw method type parsed by yaml file.
type sqlMethod struct {
	Args []string `yaml:"args"`
	Ret  string   `yaml:"ret"`
	Sql  string   `yaml:"sql"`
}

// SqlObj uses to render "sql_method" template. It is converted
// from "sqlFile".
type SqlObj struct {
	GoPackage string

	// Namespace uses to call sql method, default is "Methods".
	// That is, user can call sql method by "<pkg>.Methods.<name>"
	// If default namespace is conflict with defined table name,
	// user can change this by "namespace" option in yaml file.
	Namespace string

	// Each method calls one sql.
	Methods []*SqlMethodObj
}

// SqlMethodObj represents one sql method. It uses to render method
// template. We will parse sql statement provided by user and generate
// this struct.
type SqlMethodObj struct {
	Name string

	// The sql stored here can be directly provided to the database
	// for execution. It comes from the sql provided by the user,
	// but the GoTemplate placeholders in it have all been processed
	// and converted into sql statements that can be finally executed.
	Sql string

	RetDef  string // Method return definition.
	RetName string // Method return type name.

	// If the length of this slice is greater than 0, it means that
	// the return structure of this method needs to be generated.
	RetFields []*Field

	ArgsDef string // Method args definition.
	ArgsUse string // Use of method args.

	DB string // Database object definition.

	ExecResult   bool // Exec sql, returns sql.Result.
	ExecLastId   bool // Exec sql, returns last inserted id.
	ExecAffected bool // Exec sql, returns rows affected.

	QueryOne     bool   // Query sql, returns one row.
	QueryMany    bool   // Query sql, returns multi rows.
	QueryPointer bool   // Query return's result is pointer or not.
	QueryType    string // Query return type name.
	QueryScan    string // Use to render rows.Scan

	sqlFlag *sqlFlag
	retFlag *sqlRetFlag
	meta    *methodMeta
}

// methodMeta stores some global meta data, comes from models
// definition.
type methodMeta struct {
	// Uses to render sql statements.
	tplData map[string]tplTableData

	// TableGoName(Hump-style) -> tableObj
	tableNameMap map[string]*Obj

	// TableName(camel-style) -> tableObj
	tableMap map[string]*Obj

	// TableName(camel-style) -> FieldName(camel-style) -> fieldObj
	fieldMap map[string]map[string]*Field

	// Some methods may use the same return structure. In order
	// to prevent them from being generated repeatedly, save the
	// struct name that has been generated here.
	genTypes map[string]struct{}
}

type sqlError struct {
	err error
	sql string
}

func (e *sqlError) Error() string {
	return fmt.Sprintf("ERROR: %v\nSQL: %s", e.err, e.sql)
}

var placeholderRe = regexp.MustCompile(`{{[^}]+}}`)

type tplTableData map[string]string

func (d tplTableData) String() string {
	return d["__name__"]
}

// use global models to build meta data.
func buildSqlMethodMeta(tables []*Obj) *methodMeta {
	meta := &methodMeta{
		tplData: make(map[string]tplTableData),

		tableNameMap: make(map[string]*Obj),

		tableMap: make(map[string]*Obj),
		fieldMap: make(map[string]map[string]*Field),

		genTypes: make(map[string]struct{}),
	}

	for _, table := range tables {
		tableData := make(tplTableData)
		tname := wrapName(table.Table)
		// use "{{ .Table }}" to represent table's name.
		tableData["__name__"] = tname

		meta.fieldMap[table.Table] = make(map[string]*Field)

		allFields := make([]string, len(table.Fields))
		for i, f := range table.Fields {
			fnameRaw := camel2name(f.Name)
			fname := wrapName(fnameRaw)
			fname = fmt.Sprintf("%s.%s", tname, fname)
			allFields[i] = fname
			tableData[f.Name] = fname
			meta.fieldMap[table.Table][fnameRaw] = f
		}

		// use "{{ .Table.allFields }}" to represent all fields list.
		tableData["allFields"] = strings.Join(allFields, ", ")

		meta.tplData[table.Name] = tableData
		meta.tableMap[table.Table] = table
		meta.tableNameMap[table.Name] = table
	}

	return meta
}

// ReadSqlFile reads the sql-yaml file, parses and converts into
// SqlObj to render "sql_method" template.
func ReadSqlFile(path, goPackage string, tables []*Obj) (*SqlObj, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var file sqlFile
	err = yaml.Unmarshal(data, &file)
	if err != nil {
		return nil, err
	}
	if len(file.Methods) == 0 {
		return nil, fmt.Errorf("method is empty")
	}
	for _, table := range tables {
		if table.Table == "" {
			table.Table = camel2name(table.Name)
		}
	}

	meta := buildSqlMethodMeta(tables)

	obj := new(SqlObj)
	obj.Namespace = file.Namespace
	if obj.Namespace == "" {
		obj.Namespace = "Methods"
	}
	obj.GoPackage = goPackage

	obj.Methods = make([]*SqlMethodObj, 0, len(file.Methods))
	for k := range file.Methods {
		mobj, err := file.parseSqlMethod(k, meta)
		if err != nil {
			return nil, fmt.Errorf("parse method %s failed: %v",
				k, err)
		}
		obj.Methods = append(obj.Methods, mobj)
	}

	sort.Slice(obj.Methods, func(i, j int) bool {
		return obj.Methods[i].Name > obj.Methods[j].Name
	})

	return obj, nil
}

func (f *sqlFile) parseSqlMethod(key string, meta *methodMeta) (
	*SqlMethodObj, error,
) {
	m := f.Methods[key]

	obj := new(SqlMethodObj)
	obj.Name = key

	m.Sql = strings.Replace(m.Sql, "\n", " ", -1)
	m.Sql = strings.Replace(m.Sql, "\t", " ", -1)

	// Because snippets may contain other placeholders, we
	// must expand snippets first to ensure that placeholders
	// inner snippets are processed.
	sql, err := f.handlePhs(m.Sql, ".snippets.", f.handleSnippets)
	if err != nil {
		return nil, &sqlError{sql: m.Sql, err: err}
	}

	// Replace args placeholders in sql into "?", check them
	// exist in the yaml args option.
	if len(m.Args) > 0 {
		obj.ArgsDef = ", " + strings.Join(m.Args, ", ")
	}
	var args []string
	sql, err = f.handlePhs(sql, ".args.", func(name string) (string, error) {
		var found bool
		for _, arg := range m.Args {
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
			return "", fmt.Errorf("can not find arg "+
				"%s in args definition", name)
		}
		args = append(args, name)
		return "?", nil
	})
	if err != nil {
		return nil, &sqlError{err: err, sql: sql}
	}
	if len(args) > 0 {
		obj.ArgsUse = fmt.Sprintf("[]interface{}{%s}",
			strings.Join(args, ", "))
	}

	// Use GoTemplate to replace all placeholders in
	// sql (mainly those fieldNames)
	tpl, err := template.New("sql-template").Parse(sql)
	if err != nil {
		return nil, &sqlError{err: err, sql: sql}
	}
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	err = tpl.Execute(w, meta.tplData)
	if err != nil {
		return nil, &sqlError{err: err, sql: sql}
	}
	w.Flush()
	sql = b.String()
	// Erase extra spaces
	sql = strings.Join(strings.Fields(sql), " ")

	flags, err := parseSql(sql)
	if err != nil {
		return nil, &sqlError{err: err, sql: sql}
	}
	obj.Sql = sql

	if flags.exec {
		obj.DB = "sqlm.Execable"
		switch m.Ret {
		case "", "result":
			obj.ExecResult = true
			obj.RetDef = "sql.Result"

		case "affected":
			obj.ExecAffected = true
			obj.RetDef = "int64"

		case "lastid":
			obj.ExecLastId = true
			obj.RetDef = "int64"

		default:
			return nil, fmt.Errorf("unknown exec ret: %s", m.Ret)
		}
		return obj, nil
	}

	retFlag, err := parseRet(m.Ret)
	if err != nil {
		return nil, err
	}
	obj.DB = "sqlm.Queryable"
	obj.sqlFlag = flags
	obj.retFlag = retFlag
	obj.meta = meta
	obj.RetDef = retFlag.def
	obj.RetName = retFlag.typeName
	obj.QueryPointer = retFlag.pointer
	obj.QueryType = retFlag.typeName

	if !retFlag.simple {
		err := f.buildSqlMethodRet(obj)
		if err != nil {
			return nil, err
		}
	} else {
		obj.QueryScan = "&v"
	}

	if !retFlag.slice {
		obj.QueryOne = true
	} else {
		obj.QueryMany = true
	}

	return obj, nil
}

type phHandler func(name string) (string, error)

func (f *sqlFile) buildSqlMethodRet(obj *SqlMethodObj) error {
	// The struct of tables has been generated by the orm component.
	// In order to prevent the generation of duplicate structs, they
	// need to be recorded here, and these structs need to be skipped
	// later.
	retTable := obj.meta.tableNameMap[obj.retFlag.typeName]
	if retTable != nil {
		return f.buildSqlMethodRetByTable(obj, retTable)
	}

	// The return struct of this method may need to be generated.
	// The following needs to ensure that no method has generated
	// this struct before assigning values to obj.RetFields to
	// generate the struct.
	scans := make([]string, len(obj.sqlFlag.queryFields))
	_, exists := obj.meta.genTypes[obj.retFlag.typeName]
	if !exists {
		obj.meta.genTypes[obj.retFlag.typeName] = struct{}{}
		obj.RetFields = make([]*Field, len(obj.sqlFlag.queryFields))
	}
	for i, f := range obj.sqlFlag.queryFields {
		table := obj.meta.fieldMap[f.table]
		tobj := obj.meta.tableMap[f.table]
		if table == nil || tobj == nil {
			return fmt.Errorf("can not find "+
				"table `%s`, please check your "+
				"SELECT clause", f.table)
		}
		fobj := table[f.name]
		if fobj == nil {
			return fmt.Errorf("can not find "+
				"field `%s` in table `%s`, please "+
				"check your SELECT clause", f.name, f.table)
		}
		genf := new(Field)
		*genf = *fobj
		genf.Name = f.alias
		if genf.Name == "" {
			genf.Name = tobj.Name + fobj.Name
		}
		scan := fmt.Sprintf("&v.%s", genf.Name)
		genf.Tag = fmt.Sprintf("`table:%q field:%q`",
			f.table, f.name)
		scans[i] = scan
		if !exists {
			obj.RetFields[i] = genf
		}
	}
	obj.QueryScan = strings.Join(scans, ", ")
	return nil
}

func (f *sqlFile) buildSqlMethodRetByTable(obj *SqlMethodObj, table *Obj) error {
	tableFieldMap := make(map[string]string, len(table.Fields))
	for _, f := range table.Fields {
		tableFieldMap[camel2name(f.Name)] = f.Name
	}
	scans := make([]string, len(obj.sqlFlag.queryFields))
	for i, f := range obj.sqlFlag.queryFields {
		if f.table != table.Table {
			return fmt.Errorf("can not find "+
				"`%s`.`%s` in %s, please check your"+
				" SELECT clause", f.table, f.name, obj.Name)
		}
		name := tableFieldMap[f.name]
		if name == "" {
			return fmt.Errorf("can not find field "+
				"`%s` in table `%s`", f.name, f.table)
		}
		scan := fmt.Sprintf("&v.%s", name)
		scans[i] = scan
	}
	obj.QueryScan = strings.Join(scans, ", ")
	return nil
}

func (f *sqlFile) handlePhs(sql, prefix string, fn phHandler) (r string, err error) {
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

func (f *sqlFile) handleSnippets(name string) (string, error) {
	sql, ok := f.Snippets[name]
	if !ok {
		return "", fmt.Errorf("can not "+
			"find snippet %s, please check your sql", name)
	}
	return sql, nil
}

type sqlFlag struct {
	exec bool

	queryFields []*sqlQueryField
	queryFrom   string
	queryJoins  []string
}

type sqlQueryField struct {
	table string
	name  string
	alias string

	isCount bool
}

type sqlToken string

func (t sqlToken) Eq(k string) bool {
	up := strings.ToUpper(string(t))
	return k == up
}

func parseSql(sql string) (*sqlFlag, error) {
	fs := strings.Fields(sql)
	if len(fs) == 0 {
		return nil, fmt.Errorf("sql is empty")
	}
	tokens := make([]sqlToken, len(fs))
	for i, f := range fs {
		tokens[i] = sqlToken(f)
	}
	if !tokens[0].Eq("SELECT") {
		// no-query sql, no need to parse
		return &sqlFlag{exec: true}, nil
	}
	flag := new(sqlFlag)
	tokens = tokens[1:]
	var queryFieldStrs []string
	var selectScanned bool
	for idx, token := range tokens {
		if token.Eq("FROM") || token.Eq("JOIN") {
			idx++
			if idx >= len(tokens) {
				return nil, fmt.Errorf("expect " +
					"<table-name> after 'FROM' or 'JOIN', found EOF")
			}
			next := trimName(string(tokens[idx]))
			if token.Eq("FROM") {
				flag.queryFrom = next
			}
			if token.Eq("JOIN") {
				flag.queryJoins = append(flag.queryJoins, next)
			}
			selectScanned = true
		}
		if !selectScanned {
			queryFieldStrs = append(queryFieldStrs, string(token))
		}
	}

	queryFieldStr := strings.Join(queryFieldStrs, " ")
	queryFieldStrs = splitSqlFields(queryFieldStr, ',')

	flag.queryFields = make([]*sqlQueryField, len(queryFieldStrs))
	for i, fstr := range queryFieldStrs {
		fstr = strings.TrimSpace(fstr)
		f := new(sqlQueryField)
		tmp := splitSqlFields(fstr, ' ')
		if len(tmp) == 0 {
			return nil, fmt.Errorf("field at %d is empty", i)
		}
		def := tmp[0]
		if len(tmp) > 1 {
			f.alias = trimName(tmp[len(tmp)-1])
		}
		fn, name, err := extractSqlFieldDef(def)
		if err != nil {
			return nil, err
		}
		fnToken := sqlToken(fn)
		if fnToken.Eq("COUNT") {
			f.isCount = true
		}
		if name == "*" && !f.isCount {
			return nil, fmt.Errorf(`"SELECT *" is not allowed!`)
		}

		tmp = strings.Split(name, ".")
		switch len(tmp) {
		case 1:
			f.name = trimName(name)

		case 2:
			f.table = trimName(tmp[0])
			f.name = trimName(tmp[1])

		default:
			return nil, fmt.Errorf("field %s is bad formatted", fstr)
		}
		if f.table == "" || f.name == "" {
			return nil, fmt.Errorf("field %s: missing table or fieldname", fstr)
		}
		flag.queryFields[i] = f
	}

	return flag, nil
}

func splitSqlFields(s string, t rune) []string {
	var bucket []rune
	var fs []string
	rs := []rune(s)
	inFunc := false
	for _, r := range rs {
		if r == t && !inFunc {
			fs = append(fs, string(bucket))
			bucket = nil
			continue
		}
		bucket = append(bucket, r)
		if r == '(' {
			inFunc = true
		}
		if r == ')' {
			inFunc = false
		}
	}
	if len(bucket) > 0 {
		fs = append(fs, string(bucket))
	}
	return fs
}

func extractSqlFieldDef(f string) (string, string, error) {
	startIdx := strings.Index(f, "(")
	if startIdx < 0 {
		return "", f, nil
	}
	fnName := f[:startIdx]
	endIdx := strings.Index(f, ")")
	if endIdx < 0 {
		return "", "", fmt.Errorf("field %s is bad formatted", f)
	}
	if startIdx+1 >= len(f) {
		return "", "", fmt.Errorf("found EOF after '('")
	}
	name := f[startIdx+1 : endIdx]
	tmp := strings.Split(name, ",")
	name = trimName(tmp[0])
	return fnName, name, nil
}

func trimPlaceholder(s string) string {
	s = strings.TrimLeft(s, "{{")
	s = strings.TrimRight(s, "}}")
	return strings.TrimSpace(s)
}

func wrapName(s string) string {
	return "`" + s + "`"
}

func trimName(s string) string {
	return strings.Trim(s, "`")
}

type sqlRetFlag struct {
	typeName string
	def      string

	simple  bool
	slice   bool
	pointer bool
}

func parseRet(ret string) (*sqlRetFlag, error) {
	f := new(sqlRetFlag)
	if strings.HasPrefix(ret, "list<") {
		ret = strings.TrimPrefix(ret, "list<")
		ret = strings.TrimSuffix(ret, ">")
		f.slice = true
	}

	if strings.HasPrefix(ret, "*") {
		ret = strings.TrimPrefix(ret, "*")
		f.pointer = true
	}
	if !f.slice && !f.pointer {
		f.pointer = true
	}
	f.typeName = ret
	f.simple = isGoTypeSimple(ret)
	f.def = ret
	if f.pointer {
		f.def = "*" + f.def
	}
	if f.slice {
		f.def = "[]" + f.def
	}
	return f, nil
}

var simpleGoTypePrefix = []string{
	"int", "float", "complex", "string", "bool",
	"sql.NullString", "sql.NullInt", "sql.NullString",
	"sql.NullFloat", "sql.NullBool",
	"[]byte", "byte", "rune", "[]rune",
	"time.Time", "sql.NullTime",
}

func isGoTypeSimple(t string) bool {
	for _, prefix := range simpleGoTypePrefix {
		if strings.HasPrefix(t, prefix) {
			return true
		}
	}
	return false
}
