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
	bin string
}

func NewPluginGenerator(bin string) *PluginGenerator {
	return &PluginGenerator{bin: bin}
}

func (g *PluginGenerator) Generate(t TMetadata) error {
	m, err := t.Encode()
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
