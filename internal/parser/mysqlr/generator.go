package mysqlr

import (
	"fmt"

	"github.com/ezbuy/ezorm/v2/internal/generator"
	"github.com/ezbuy/ezorm/v2/internal/parser"
)

var _ generator.Generator = (*MySQLRGenerator)(nil)

type MySQLRGenerator struct{}

func (g *MySQLRGenerator) Generate(meta generator.TMetadata) error {
	var hasDriver bool
	if err := meta.Meta.Each(func(tn generator.TemplateName, om generator.Schema) error {
		d, err := om.GetDriver()
		if err != nil {
			return err
		}
		if d != g.DriverName() {
			return nil
		}
		hasDriver = true
		m := NewMetaObject(meta.Pkg)
		if err := m.Read(string(tn), om); err != nil {
			return fmt.Errorf("%s: %w", d, err)
		}
		if err := GenerateGoTemplate(meta.Output, m); err != nil {
			return fmt.Errorf("%s: %w", d, err)
		}
		if err := GenerateScriptTemplate(meta.Output, meta.Pkg, m); err != nil {
			return fmt.Errorf("%s: %w", d, err)
		}
		return nil
	}); err != nil {
		return err
	}
	if hasDriver {
		if err := GenerateConfTemplate(meta.Output, meta.Pkg); err != nil {
			return fmt.Errorf("mysqlr: %w", err)
		}
	}
	return nil
}

func (g *MySQLRGenerator) DriverName() string {
	return parser.MySQLRGeneratorName
}
