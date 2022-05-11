package shared

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ezbuy/ezorm/v2/internal/generator"
	"github.com/ezbuy/ezorm/v2/internal/parser"
)

var _ generator.Generator = (*Generator)(nil)

// Generator is a generator that supports mysql and mongo drivers
// which share the same metadata format.
type Generator struct{}

func (g *Generator) Generate(meta generator.TMetadata) error {
	drivers := make(map[string]*Obj)
	var dbs []*Obj
	if err := meta.Meta.Each(func(tn generator.TemplateName, s generator.Schema) error {
		d, err := s.GetDriver()
		if err != nil {
			return fmt.Errorf("template: %s: %q", string(tn), err)
		}
		var o *Obj
		switch d {
		case parser.MongoGeneratorName:
			o = &Obj{
				Package:   meta.Pkg,
				GoPackage: meta.Pkg,
				Name:      string(tn),
			}
			if err := o.Read(string(tn), s); err != nil {
				return err
			}
		case parser.MySQLGeneratorName:
			o = &Obj{
				Package:   meta.Pkg,
				GoPackage: meta.Pkg,
				Name:      string(tn),
			}
			if err := o.Read(string(tn), s); err != nil {
				return err
			}
			dbs = append(dbs, o)
		default:
			return nil
		}
		for _, gt := range o.GetGenTypes() {
			fileAbsPath := filepath.Join(meta.Output, fmt.Sprintf("gen_%s_%s.go", string(tn), gt))
			if err := render(fileAbsPath, gt, o); err != nil {
				return err
			}
		}
		drivers[d] = o
		return nil
	}); err != nil {
		return err
	}
	if len(dbs) > 0 {
		if err := render(
			filepath.Join(meta.Output, "create_mysql.sql"),
			"mysql_script",
			dbs); err != nil {
			return err
		}
	}
	for _, d := range drivers {
		for _, t := range d.GetConfigTemplates() {
			fileAbsPath := filepath.Join(meta.Output, fmt.Sprintf("gen_%s.go", t))
			if err := render(fileAbsPath, t, d); err != nil {
				return err
			}
		}
	}
	return nil
}

func render(path string, name string, obj any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return Tpl.ExecuteTemplate(bytes.NewBuffer(data), name, obj)
}

func (g *Generator) DriverName() string {
	return "shared_generator"
}