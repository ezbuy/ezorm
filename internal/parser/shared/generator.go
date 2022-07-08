package shared

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/ezbuy/ezorm/v2/internal/generator"
	"github.com/ezbuy/ezorm/v2/internal/parser"
)

var _ generator.Generator = (*Generator)(nil)

// Generator is a generator that supports mysql and mongo drivers
// which share the same metadata format.
type Generator struct{}

func (g *Generator) Generate(meta generator.TMetadata) error {
	drivers := make(map[string][]*Obj)
	var dbs []*Obj
	if err := meta.Meta.Each(func(tn generator.TemplateName, s generator.Schema) error {
		d, err := s.GetDriver()
		if err != nil {
			return fmt.Errorf("template: %s: %q", string(tn), err)
		}
		var o *Obj
		ns := meta.Namespace
		if meta.Namespace == "" {
			ns = meta.Pkg
		}
		switch d {
		case parser.MongoGeneratorName:
			o = &Obj{
				Namespace: ns,
				GoPackage: meta.Pkg,
				Name:      string(tn),
			}
			if err := o.Read(string(tn), s); err != nil {
				return err
			}
		case parser.MySQLGeneratorName:
			o = &Obj{
				Namespace: ns,
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
		drivers[d] = append(drivers[d], o)
		return nil
	}); err != nil {
		return err
	}
	sort.SliceStable(dbs, func(i, j int) bool {
		return dbs[i].Name < dbs[j].Name
	})
	if len(dbs) > 0 {
		if err := render(
			filepath.Join(meta.Output, "create_mysql.sql"),
			"mysql_script",
			dbs); err != nil {
			return err
		}
	}
	for _, d := range drivers {
		if len(d) <= 0 {
			continue
		}
		for _, t := range d[0].GetConfigTemplates() {
			fileAbsPath := filepath.Join(meta.Output, fmt.Sprintf("gen_%s.go", t))
			if err := render(fileAbsPath, t, d); err != nil {
				return err
			}
		}
	}
	return nil
}

func render(path string, name string, obj any) error {
	fd, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	return Tpl.ExecuteTemplate(fd, name, obj)
}

func (g *Generator) DriverName() string {
	return "shared_generator"
}
