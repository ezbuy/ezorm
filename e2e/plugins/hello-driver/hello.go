package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ezbuy/ezorm/v2/pkg/plugin"
)

func main() {
	r := bufio.NewScanner(os.Stdin)
	if r.Scan() {
		req, err := plugin.Decode([]byte(r.Text()))
		if err != nil {
			fmt.Fprintf(os.Stdout, "decode meta error: %v\n", err)
			return
		}
		f, err := os.Create(filepath.Join(req.GetOutputPath(), "metadata.json"))
		if err != nil {
			fmt.Fprintf(os.Stdout, "create file error: %v\n", err)
			return
		}
		defer f.Close()
		if _, err := f.WriteString(r.Text()); err != nil {
			fmt.Fprintf(os.Stdout, "write file error: %v\n", err)
			return
		}
		if err := req.Each(func(_ plugin.TemplateName, s plugin.Schema) error {
			d, err := s.GetDriver()
			if err != nil {
				return err
			}
			fmt.Fprintf(os.Stdout, "%s\n", d)
			return nil
		}); err != nil {
			fmt.Fprintf(os.Stdout, "each error: %v\n", err)
			return
		}
	}

	if err := r.Err(); err != nil {
		fmt.Fprintf(os.Stdout, "scan error: %v\n", err)
	}
}
