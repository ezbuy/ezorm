// Package template is the metadata for orm code generation
package template

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/valyala/fasttemplate"
)

var driverSet []Driver

func RegisterDriver(d Driver) {
	driverSet = append(driverSet, d)
}

type Template struct {
	Data     []byte
	Filename string
}

type Driver interface {
	// Name should be same as the driver defines in yaml
	Name() string
	// Template is the pre-defined template for generation.
	Templates() ([]Template, error)
}

func NewGenerator(output string) Generator {
	return Generator{
		path: output,
	}
}

type Generator struct {
	path string
}

func (g Generator) Generate(ctx context.Context, d Driver, value map[string]interface{}) error {
	tpls, err := d.Templates()
	if err != nil {
		return fmt.Errorf("generate: %w", err)
	}
	for _, tpl := range tpls {
		if err := g.generate(ctx, tpl, value); err != nil {
			return fmt.Errorf("generate: %w", err)
		}
	}
	return nil
}

func (g Generator) generate(ctx context.Context, t Template, v map[string]interface{}) error {
	b := bytes.NewBuffer(t.Data)
	if _, err := fasttemplate.New("ezorm.v2", "{{", "}}").Execute(b, v); err != nil {
		return fmt.Errorf("%w", err)
	}
	return os.WriteFile(filepath.Join(g.path, t.Filename), b.Bytes(), 0644)
}

func GetDriver(name string) (Driver, error) {
	for _, d := range driverSet {
		if d.Name() == name {
			return d, nil
		}
	}
	return nil, fmt.Errorf("template: not found with name: %s", name)
}
