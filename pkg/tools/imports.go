package tools

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

func FmtCode(path string, errOut io.Writer) {
	oscmd := exec.Command("goimports", "-format-only", "-w", path)
	oscmd.Stderr = errOut
	if err := oscmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "run fmt tool: goimports: %w", err)
		return
	}
}
