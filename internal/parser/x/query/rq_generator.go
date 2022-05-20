package query

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/ezbuy/ezorm/v2/internal/generator"
	"github.com/ezbuy/ezorm/v2/internal/parser"
	"github.com/ezbuy/ezorm/v2/internal/parser/mysql"
	"github.com/ezbuy/ezorm/v2/internal/parser/mysqlr"
	"github.com/ezbuy/ezorm/v2/internal/parser/shared"
)

var _ generator.Generator = (*RawQueryGenerator)(nil)

// RawQueryGenerator is a generator for raw query.
// It should be available for all drivers.
type RawQueryGenerator struct{}

func (rg *RawQueryGenerator) Generate(meta generator.TMetadata) error {
	tableSchema := make(map[string]generator.IObject)
	if err := meta.Meta.Each(func(tn generator.TemplateName, om generator.Schema) error {
		dr, err := om.GetDriver()
		if err != nil {
			return err
		}
		table, err := om.GetTable(dr)
		if err != nil {
			log.Printf("rawquery: warning: %q\n", err)
			return nil
		}

		switch dr {
		case parser.MySQLGeneratorName:
			s := mysql.NewMySQLObject(meta.Pkg, string(tn))
			if err := s.Read(string(tn), om); err != nil {
				return err
			}
			tableSchema[table] = s
		case parser.MySQLRGeneratorName:
			s := mysqlr.NewMetaObject(meta.Pkg)
			if err := s.Read(string(tn), om); err != nil {
				return err
			}
			tableSchema[table] = s
		}

		return nil
	}); err != nil {
		return err
	}
	if len(tableSchema) == 0 {
		return nil
	}
	// validate input
	input := meta.Input
	var inputDir string
	f, err := os.Stat(input)
	if err != nil {
		return err
	}
	if f.IsDir() {
		inputDir = input
	} else {
		inputDir = filepath.Dir(input)
	}
	sqlsDir := filepath.Join(inputDir, "sqls")
	stat, err := os.Stat(sqlsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if !stat.IsDir() {
		return nil
	}
	p := NewSQL(tableSchema)
	var methods []*SQLMethod
	if err := filepath.WalkDir(sqlsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		m, err := p.Read(path)
		if err != nil {
			return err
		}
		m.FromFile = d.Name()
		methods = append(methods, m)
		return nil
	}); err != nil {
		return err
	}
	file := &SQLFile{
		GoPackage: meta.Pkg,
		Methods:   methods,
		Dir:       sqlsDir,
	}
	goFile := filepath.Join(meta.Output, "gen_methods.go")
	fd, err := os.OpenFile(goFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	return shared.Tpl.ExecuteTemplate(fd, "sql_method", file)
}

func (rg *RawQueryGenerator) DriverName() string {
	return "common_raw_query"
}
