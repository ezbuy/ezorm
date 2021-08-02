package sqlm

import (
	"errors"
	"fmt"
	"strings"
)

type clause struct {
	isExec bool

	Select string
	From   string
	Joins  []string

	Full string
}

type keyword string

func (k keyword) eq(s string) bool {
	return strings.ToUpper(s) == string(k)
}

var (
	_select = keyword("SELECT")
	_from   = keyword("FROM")
	_join   = keyword("JOIN")
	_on     = keyword("ON")
	_count  = keyword("COUNT")

	// Keywords after 'FROM' clause.
	afterFroms = map[string]struct{}{
		"LEFT": {}, "RIGHT": {}, "INNER": {}, "JOIN": {},
		"WHERE": {}, "ORDER": {}, "LIMIT": {}, "GROUP": {},
	}
)

// split extracts different components in the sql, that is,
// different clauses. We mainly focus on the SELECT, FROM,
// and JOIN clauses of the query sql for subsequent analysis
// and template data generation.
func split(sql, raw string) (*clause, error) {
	tokens := strings.Fields(sql)
	if len(tokens) == 0 {
		return nil, errors.New("sql is empty")
	}
	// First part: SELECT clause.
	master := tokens[0]
	if !_select.eq(master) {
		// No query sql, no need to parse.
		return &clause{isExec: true}, nil
	}
	tokens = tokens[1:]
	var selectTokens []string
	var idx int
	for {
		if idx == len(tokens) {
			var last string
			if idx-1 >= 0 {
				last = tokens[idx-1]
			}
			// FROM not found
			return nil, &Error{
				full:  raw,
				wrong: last,
				desc:  "FROM not found in sql",
			}
		}

		t := tokens[idx]
		if _from.eq(t) {
			break
		}

		selectTokens = append(selectTokens, t)
		idx++
	}
	if len(selectTokens) == 0 {
		return nil, &Error{
			full:  raw,
			wrong: master,
			desc:  "SELECT is empty",
		}
	}
	cs := new(clause)
	cs.Full = raw
	cs.Select = strings.Join(selectTokens, " ")
	fromToken := tokens[idx]
	tokens = tokens[idx+1:]

	// Second part: FROM clause.
	var fromTokens []string
	idx = 0
	for {
		if idx == len(tokens) {
			// sql ends.
			break
		}
		t := tokens[idx]
		if _, ok := afterFroms[t]; ok {
			// meet end keyword.
			break
		}
		fromTokens = append(fromTokens, t)
		idx++
	}
	if len(fromTokens) == 0 {
		return nil, &Error{
			full:  sql,
			wrong: fromToken,
			desc:  "FROM clause is empty",
		}
	}
	cs.From = strings.Join(fromTokens, " ")
	tokens = tokens[idx:]

	// Third part: JOIN clause.
	idx = 0
	for {
		if idx == len(tokens) {
			break
		}
		t := tokens[idx]
		if !_join.eq(t) {
			idx++
			continue
		}
		var joinTokens []string
		for {
			idx++
			if idx == len(tokens) {
				return nil, &Error{
					full:  raw,
					wrong: t,
					desc:  "missing 'ON' after 'JOIN'",
				}
			}
			subt := tokens[idx]
			if _on.eq(subt) {
				break
			}
			joinTokens = append(joinTokens, subt)
		}
		join := strings.Join(joinTokens, " ")
		cs.Joins = append(cs.Joins, join)
		idx++
	}

	return cs, nil
}

type tables struct {
	names2alias map[string]string
	alias2name  map[string]string
}

// parseTables decodes all tables appeared in sql, include their
// aliases.
func (cs *clause) parseTables() (*tables, error) {
	tcap := 1 + len(cs.Joins)
	ts := &tables{
		names2alias: make(map[string]string, tcap),
		alias2name:  make(map[string]string, tcap),
	}

	tableFields := make([]*field, 1, 1+len(cs.Joins))
	fromField, err := splitField(cs.From)
	if err != nil {
		return nil, &Error{
			full:  cs.Full,
			wrong: cs.From,
			desc:  err.Error(),
		}
	}
	tableFields[0] = fromField

	for _, join := range cs.Joins {
		joinField, err := splitField(join)
		if err != nil {
			return nil, &Error{
				full:  cs.Full,
				wrong: join,
				desc:  err.Error(),
			}
		}
		tableFields = append(tableFields, joinField)
	}

	for _, f := range tableFields {
		ts.names2alias[f.Name] = f.Alias
		if f.Alias != "" {
			ts.alias2name[f.Alias] = f.Name
		}
	}
	return ts, nil
}

type queryField struct {
	TableFull  string
	TableAlias string

	Count bool

	Name  string
	Alias string
}

// parseQuery decodes all query fields in SELECT clause.
// This will parse their table alias name into full table name.
func (cs *clause) parseQuery(ts *tables) ([]*queryField, error) {
	rawFields := splitWithFunc(cs.Select, ',')
	fs := make([]*queryField, len(rawFields))
	cntIdx := 0
	for i, rawField := range rawFields {
		rawField = strings.TrimSpace(rawField)
		f, err := splitField(rawField)
		if err != nil {
			return nil, &Error{
				full:  cs.Full,
				wrong: rawField,
				desc:  err.Error(),
			}
		}
		qf := new(queryField)
		qf.Count = _count.eq(f.Func)
		if qf.Count {
			qf.Alias = f.Alias
			qf.Name = fmt.Sprintf("Count%d", cntIdx)
			cntIdx++
			fs[i] = qf
			continue
		}

		tmp := strings.Split(f.Name, ".")
		if len(tmp) != 2 {
			return nil, &Error{
				full:  cs.Full,
				wrong: rawField,
				desc:  "field bad format, expect <table/alias>.<field>",
			}
		}
		tableRaw := tmp[0]
		// We are not sure whether the user write the full
		// table name or alias, here need to confirm.
		_, ok := ts.names2alias[tableRaw]
		if !ok {
			full, ok := ts.alias2name[tableRaw]
			if !ok {
				// Can not find user write table name.
				return nil, &Error{
					full:  cs.Full,
					wrong: rawField,
					desc:  fmt.Sprintf("cannot find table %q in your sql", tableRaw),
				}
			}
			qf.TableFull = full
			qf.TableAlias = tableRaw
		} else {
			qf.TableFull = tableRaw
		}

		// Here we will not check whether this field appears
		// in the table, this step is left to the outside.
		qf.Name = tmp[1]
		qf.Alias = f.Alias

		fs[i] = qf
	}

	return fs, nil
}

type field struct {
	Func  string
	Name  string
	Alias string
}

// splitField splits field string into field struct.
// The format is:
//   [<Func>][(] {{ .<Name> }} [, ...] [)] [AS] [<Alias>]
//
// The "{{ .<Name> }}" indicates that GoTemplate must be used
// to represent the name. Except for "COUNT(1)" and "COUNT(*)"
// will be handled specially.
func splitField(s string) (*field, error) {
	if s == "" {
		return nil, errors.New("field is empty")
	}
	tmp := splitWithFunc(s, ' ')

	f := new(field)
	main := tmp[0]
	bkIdx := strings.Index(main, "(")
	var nameph string
	if bkIdx >= 0 {
		// Format: FUNC_NAME(FIELD) [AS] ALIAS
		funcName := main[:bkIdx]
		f.Func = strings.ToUpper(funcName)

		bkEndIdx := strings.Index(main, ")")
		if bkEndIdx <= bkIdx+1 {
			return nil, errors.New("Func bad format")
		}
		funcBody := main[bkIdx+1 : bkEndIdx]
		if funcBody == "" {
			return nil, errors.New("Func body is empty")
		}
		nameph = strings.Split(funcBody, ",")[0]
	} else {
		nameph = main
	}

	name, err := splitPlaceholder(nameph)
	if err != nil {
		return nil, err
	}
	if name == "" {
		return nil, errors.New("Field name is empty")
	}
	if name == "1" || name == "*" {
		// Only allow "COUNT(1)" or "COUNT(*)"
		if f.Func != "COUNT" {
			return nil, errors.New("SELECT * is not allowed")
		}
	}
	f.Name = name

	if len(tmp) > 1 {
		f.Alias = trimName(tmp[len(tmp)-1])
	}

	return f, nil
}

// Unlike ordinary split, splitWithFunc treats those splitters
// that appear in functions or placeholders as a whole.
func splitWithFunc(query string, s rune) []string {
	allRunes := []rune(query)
	var bucket []rune
	var inFunc bool
	var inPh bool

	var fs []string
	for _, r := range allRunes {
		if r == '(' {
			inFunc = true
			bucket = append(bucket, r)
			continue
		}
		if r == ')' {
			inFunc = false
			bucket = append(bucket, r)
			continue
		}
		if r == '{' {
			inPh = true
			bucket = append(bucket, r)
			continue
		}
		if r == '}' {
			inPh = false
			bucket = append(bucket, r)
			continue
		}
		if r == s && !inFunc && !inPh {
			fs = append(fs, string(bucket))
			bucket = nil
			continue
		}
		bucket = append(bucket, r)
	}
	if len(bucket) > 0 {
		fs = append(fs, string(bucket))
	}
	return fs
}

var errNotPlaceholder = errors.New("field or table must be GoTemplate placeholder")

// Extract GoTemplate placeholder's name.
//    {{ .<Name> }}
func splitPlaceholder(ph string) (string, error) {
	if ph == "1" || ph == "*" {
		return ph, nil
	}
	ph = strings.TrimSpace(ph)
	if !strings.HasPrefix(ph, "{{") {
		return "", errNotPlaceholder
	}
	if !strings.HasSuffix(ph, "}}") {
		return "", errNotPlaceholder
	}
	ph = strings.TrimPrefix(ph, "{{")
	ph = strings.TrimSuffix(ph, "}}")
	ph = strings.TrimSpace(ph)
	if !strings.HasPrefix(ph, ".") {
		return "", errors.New("missing '.' in placeholder's head")
	}
	ph = strings.TrimPrefix(ph, ".")
	return ph, nil
}

func trimName(s string) string {
	return strings.Trim(s, "`")
}
