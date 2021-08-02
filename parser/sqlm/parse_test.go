package sqlm

import (
	"fmt"
	"strings"
	"testing"
)

func TestSplit(t *testing.T) {
	sqls := []string{
		`
		SELECT {{ .User.Name }}, {{ .User.Age }}, {{ .ud.Phone }}
		FROM {{ .User }}
		JOIN {{ .UserDetail }} ud ON {{ .User.Id }}={{ ud.UserId }}
		WHERE {{ .User.Id }}={{ .args.id }}
		`,
		`
		SELECT COUNT(1) FROM {{ .User }} u
		WHERE {{ u.Age }}>{{ .args.age }}
		`,
		`
		SELECT
			{{ .u.Name }},
			{{ .u.Phone }},
			{{ .u.Password }},
			{{ .ud.Email }},
			IFNULL({{ ud.Text }}, '') Text
		FROM {{ .User }} AS u
		JOIN {{ .UserDetail }} AS ud ON {{ u.Id }}={{ ud.UserId }}
		LIMIT {{ .args.offset }}, {{ .args.limit }}
		`,
	}

	for _, sql := range sqls {
		sql = strings.Replace(sql, "\n", " ", -1)
		cs, err := split(sql)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Printf("SELECT: %s\nFROM: %s\nJOINs: %v\n=====\n",
			cs.Select, cs.From, cs.Joins)
	}
}

func TestSplitField(t *testing.T) {
	sqls := []string{
		"{{ .User.Name }}",
		"{{ .ud.Phone }}",
		"{{.User.Phone}}",
		"{{.u.Phone}}",
		"{{ .User.Name }} AS name",
		"{{ .User.Name }} Name",
		"COUNT(1)",
		"COUNT({{.User.Name}}) AS UserCount",
		"COUNT({{.User.Name}}) UserCount",
		"IFNULL({{.ud.Text}}, '') AS Text",
		"IFNULL({{.ud.Text}}, '')",

		"{{ .User }}",
		"{{.UserDetail}}",
		"{{ .User }} AS u",
		"{{ .UserDetail }} ud",
	}

	for _, sql := range sqls {
		f, err := splitField(sql)
		if err != nil {
			fmt.Printf("parse %s failed: %v\n", sql, err)
			continue
		}
		fmt.Println("===================")
		fmt.Println(sql)
		fmt.Printf("Func: %s\nName: %s\nAlias: %s\n",
			f.Func, f.Name, f.Alias)
		fmt.Println("===================")
	}
}

func TestParseSQL(t *testing.T) {
	sqls := []string{
		`
		SELECT
			{{ .u.Id }},
			{{ .u.Name }},
			{{ .u.Phone }},
			{{ .u.Password }},
			{{ .r.Id }},
			{{ .r.Name }},
			{{ .r.Act }}

		FROM {{ .User }} u
		JOIN {{ .RoleUser }} ru ON {{ .ru.UserId }}={{ .u.Id }}
		JOIN {{ .Role }} r ON {{ .r.Id }}={{ .ru.RoleId }}
		WHERE {{ .u.Id }}={{ .args.uid }}
		`,

		`
		SELECT
			{{ .User.Id }} UID,
			IFNULL({{ .ud.Text }}, '') AS Text,
			IFNULL({{ .ud.Desc }}, ''),
			{{ .User.Name }} AS UName,

		FROM {{ .User }}
		JOIN {{ .UserDetail }} ud ON {{ .ud.UserId }}={{ .User.Id }}
		LIMIT {{ .args.offset }}, {{ .args.limit }}
		`,

		`
		SELECT COUNT(1)
		FROM {{ .User }}
		WHERE {{ .User.Age }} > 18
		`,
	}

	for _, sql := range sqls {
		cs, err := split(sql)
		if err != nil {
			fmt.Printf("split failed: %v\n", err)
			return
		}

		ts, err := cs.parseTables()
		if err != nil {
			fmt.Printf("parse tables failed: %v\n", err)
			return
		}

		fs, err := cs.parseQuery(ts)
		if err != nil {
			fmt.Printf("parse query failed: %v\n", err)
			return
		}

		for _, f := range fs {
			fmt.Printf("tableFull: %s, tableAlias: %s, "+
				"Count: %v, Name: %s, Alias: %s\n",
				f.TableFull, f.TableAlias, f.Count, f.Name, f.Alias)
		}
		fmt.Println("=============================")
	}

}
