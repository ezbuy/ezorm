package generator

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

var _ Generator = (*PluginGenerator)(nil)

type PluginGenerator struct {
	bin  string
	args map[string]string
}

type PluginGeneratorOption func(*PluginGenerator)

func WithGeneratorArgs(args map[string]string) PluginGeneratorOption {
	return func(g *PluginGenerator) {
		g.args = args
	}
}

func NewPluginGenerator(bin string, opts ...PluginGeneratorOption) *PluginGenerator {
	g := &PluginGenerator{
		bin: bin,
	}
	for _, opt := range opts {
		opt(g)
	}
	return g
}

func (g *PluginGenerator) Generate(t TMetadata) error {
	m, err := t.Encode(g.args)
	if err != nil {
		return err
	}
	path, err := exec.LookPath(fmt.Sprintf("ezorm-gen-%s", g.bin))
	if err != nil {
		return err
	}
	cmd := exec.Command(path)
	cmd.Stdin = bytes.NewBuffer(m)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (g *PluginGenerator) DriverName() string {
	return filepath.Base(g.bin)
}
