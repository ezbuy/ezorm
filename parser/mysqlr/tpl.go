package mysqlr

import (
	"bytes"
	"embed"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"
)

//go:embed tpl/*.gogo tpl/*.sql
var content embed.FS

func init() {
	funcMap := template.FuncMap{
		"add":        Add,
		"sub":        Sub,
		"divide":     Divide,
		"multiply":   Multiply,
		"camel2name": Camel2Name,
		"camel2sep":  camel2sep,
	}

	ormTemplate = template.New("mysqlr").Funcs(funcMap)
	files := []string{
		"tpl/mysqlr_script.sql",
		"tpl/config.gogo",
		"tpl/config.gogo",
		"tpl/function.gogo",
		"tpl/index.gogo",
		"tpl/mysqlr.gogo",
		"tpl/object_query.gogo",
		"tpl/object_read.gogo",
		"tpl/object_write.gogo",
		"tpl/object.gogo",
		"tpl/primary_key.gogo",
		"tpl/unique_key.gogo",
		"tpl/orm.gogo",
	}
	_, err := ormTemplate.ParseFS(content, files...)
	if err != nil {
		panic(err)
	}
}

var ormTemplate *template.Template

func templates(obj *MetaObject) []string {
	return []string{"mysqlr"}
}

func GenerateGoTemplate(output string, obj *MetaObject) error {
	for _, tpl := range templates(obj) {
		filename := filepath.Join(output, strings.Join([]string{"gen", tpl, camel2sep(obj.Name, "."), "go"}, "."))
		fd, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		if err := ormTemplate.ExecuteTemplate(fd, tpl, obj); err != nil {
			return err
		}
		fd.Close()
		fmtCode(filename)
	}
	return nil
}

func GenerateScriptTemplate(output string, driver string, obj *MetaObject) error {
	filename := filepath.Join(output, strings.Join([]string{"gen", "script", driver, camel2sep(obj.Name, "."), "sql"}, "."))
	fd, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer fd.Close()
	if err := ormTemplate.ExecuteTemplate(fd, "mysqlr_script", obj); err != nil {
		return err
	}
	return nil
}

func GenerateConfTemplate(output string, packageName string) error {
	filename := filepath.Join(output, strings.Join([]string{"gen", "conf", "mysql", "go"}, "."))
	fd, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer fd.Close()
	if err := ormTemplate.ExecuteTemplate(fd, "mysqlr_config", map[string]interface{}{
		"GoPackage": packageName,
	}); err != nil {
		return err
	}
	filename = filepath.Join(output, strings.Join([]string{"gen", "orm", "mysql", "go"}, "."))
	fd, err = os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer fd.Close()
	if err := ormTemplate.ExecuteTemplate(fd, "mysqlr_orm", map[string]interface{}{
		"GoPackage": packageName,
	}); err != nil {
		return err
	}
	return nil
}

func camel2sep(s string, sep string) string {
	nameBuf := bytes.NewBuffer(nil)
	for i := range s {
		n := rune(s[i]) // always ASCII?
		if unicode.IsUpper(n) {
			if i > 0 {
				nameBuf.WriteString(sep)
			}
			n = unicode.ToLower(n)
		}
		nameBuf.WriteRune(n)
	}
	return nameBuf.String()
}

func fmtCode(path string) {
	oscmd := exec.Command("goimport", "-w", path)
	oscmd.Run()
}

func Add(a, b int) int {
	return a + b
}

func Sub(a, b int) int {
	return a - b
}

func Divide(a, b int) int {
	return a / b
}

func Multiply(a, b int) int {
	return a * b
}
