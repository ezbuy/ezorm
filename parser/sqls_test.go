package parser

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
	"text/template"

	"github.com/ezbuy/ezorm/tpl"
)

func TestSqlsTemplate(t *testing.T) {
	obj := &SqlObj{
		GoPackage: "test",
		Namespace: "Methods",
		Methods: []*SqlMethodObj{
			{
				Name:    "ListUsers",
				RetDef:  "[]*User",
				RetName: "User",
				RetFields: []*Field{
					{
						Name: "Name",
						Type: "string",
					},
					{
						Name: "Age",
						Type: "int32",
					},
					{
						Name: "Phone",
						Type: "string",
					},
					{
						Name: "Email",
						Type: "string",
					},
				},
			},
		},
	}

	r, err := execSqlsTpl(obj)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(r)
}

func execSqlsTpl(data interface{}) (string, error) {
	tplData, err := tpl.Asset("tpl/sql_method.gogo")
	if err != nil {
		return "", err
	}
	tpl, err := template.New("ezorm").Parse(string(tplData))
	if err != nil {
		return "", err
	}
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	err = tpl.ExecuteTemplate(w, "sql_method", data)
	if err != nil {
		return "", err
	}
	w.Flush()

	return b.String(), nil
}

func TestParseSql(t *testing.T) {
	sqls := []string{
		"SELECT `user`.`name`, IFNULL(`user`.`age`, 0) Name, COUNT(`wallet`.`left`) FROM `user` JOIN `wallet` ON `wallet`.`user_id`=`user`.`id`",
		"SELECT COUNT(1) FROM user",
		`
		SELECT user.id AS ID, user.age Age, user_detail.Text Text
		FROM user JOIN user_detail ON user.id=user_detail.user_id
		`,
	}
	for _, sql := range sqls {
		flag, err := parseSql(sql)
		if err != nil {
			fmt.Printf("parse failed: %v\n", err)
			return
		}
		fmt.Println("================")
		for _, f := range flag.queryFields {
			fmt.Printf("table=%s, name=%s, alias=%s, isCount=%v\n",
				f.table, f.name, f.alias, f.isCount)
		}
	}
}

func showObj(v interface{}) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Printf("marshal failed: %v\n", err)
		return
	}
	fmt.Println(string(data))
}

var testSqlFile = &sqlFile{
	Namespace: "Methods",
	Snippets: map[string]string{
		"userFields": "{{ .User.Name }}, {{ .User.Age }}, {{ .User.Phone }}",
		"limit":      "LIMIT {{ .args.offset }}, {{ .args.limit }}",
	},

	Methods: map[string]*sqlMethod{
		"FindUserById": {
			Args: []string{"id int64"},
			Ret:  "*UserResult",
			Sql: `
			SELECT {{ .snippets.userFields }}, {{ .UserDetail.Text }}
			FROM {{ .User }}
			JOIN {{ .UserDetail }} ON {{ .User.Id }} = {{ .UserDetail.UserId }}
			WHERE {{ .User.Id }} = {{ .args.id }}
			`,
		},
	},
}

var testTables = []*Obj{
	{
		Name:  "User",
		Table: "user",
		Fields: []*Field{
			{
				Name: "Id",
				Type: "int64",
			},
			{
				Name: "Name",
				Type: "string",
			},
			{
				Name: "Age",
				Type: "int32",
			},
			{
				Name: "Phone",
				Type: "string",
			},
			{
				Name: "Email",
				Type: "string",
			},
		},
	},
	{
		Name:  "UserDetail",
		Table: "user_detail",
		Fields: []*Field{
			{
				Name: "UserId",
				Type: "int64",
			},
			{
				Name: "Text",
				Type: "string",
			},
		},
	},
}

func init() {
	for _, t := range testTables {
		for _, f := range t.Fields {
			f.Obj = t
		}
	}
}

func TestBuildMeta(t *testing.T) {
	meta := buildSqlMethodMeta(testTables)
	fmt.Println(meta.tplData["User"]["Id"])
	fmt.Println(meta.tplData["UserDetail"])
	fmt.Println(meta.tableMap)
}

func TestParseSqlMethod(t *testing.T) {
	meta := buildSqlMethodMeta(testTables)
	m, err := testSqlFile.parseSqlMethod("FindUserById", meta)
	if err != nil {
		fmt.Printf("parse FindUserById failed: %v\n", err)
		return
	}
	fmt.Println(m)
}

func TestSqlMethodTemplate(t *testing.T) {

}
