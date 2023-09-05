// Copyright © 2022 ezbuy & LITB team
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v3"

	"github.com/ezbuy/ezorm/v2/internal/generator"
	"github.com/ezbuy/ezorm/v2/internal/parser/mongo"
	"github.com/ezbuy/ezorm/v2/internal/parser/mysql"
	"github.com/ezbuy/ezorm/v2/internal/parser/mysqlr"
	"github.com/ezbuy/ezorm/v2/internal/parser/x/query"
	"github.com/ezbuy/ezorm/v2/pkg/tools"

	"github.com/spf13/cobra"
)

var generators = []generator.Generator{
	&mysql.MySQLGenerator{},
	&mongo.MongoGenerator{},
	&mysqlr.MySQLRGenerator{},
	&query.RawQueryGenerator{},
}

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate orm code from yaml file",
	RunE: func(_ *cobra.Command, _ []string) error {
		var objs generator.Metadata
		stat, err := os.Stat(input)
		if err != nil {
			return err
		}
		switch {
		case stat.IsDir():
			if err := filepath.WalkDir(input, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if d.IsDir() {
					return nil
				}
				if filepath.Ext(path) != ".yaml" {
					return nil
				}
				data, err := ioutil.ReadFile(path)
				if err != nil {
					return err
				}
				if err := decode(data, &objs); err != nil {
					return err
				}
				return nil
			}); err != nil {
				return err
			}
		default:
			data, err := ioutil.ReadFile(input)
			if err != nil {
				return err
			}
			if err := decode(data, &objs); err != nil {
				return err
			}
		}
		if genPackageName == "" {
			genPackageName = strings.Split(stat.Name(), ".")[0]
		}
		if genGoPackageName == "" {
			genGoPackageName = genPackageName
		}
		if plugin != "" {
			if pluginOnly {
				generators = []generator.Generator{}
			}
			generators = append(generators, generator.NewPluginGenerator(plugin))
		}

		if err := generator.Render(generator.TMetadata{
			Meta:      objs,
			Pkg:       genGoPackageName,
			Output:    output,
			Input:     input,
			Namespace: namespace,
		}, generators...); err != nil {
			return err
		}
		errBuffer := bytes.NewBuffer(nil)
		tools.FmtCode(output, errBuffer)
		if errBuffer.Len() > 0 {
			fmt.Fprintf(os.Stderr, "run fmt tool: goimports: %s", errBuffer.String())
			// fallthrough ,do not return here
		}

		return nil
	},
}

func decode(data []byte, meta *generator.Metadata) error {
	dc := yaml.NewDecoder(bytes.NewReader(data))
	for {
		if err := dc.Decode(meta); err != nil {
			if err == io.EOF {
				return nil
			} else {
				return err
			}
		}
	}
}

var (
	input            string
	output           string
	genPackageName   string
	genGoPackageName string
	disableSQLs      bool
	plugin           string
	namespace        string
	pluginOnly       bool
)

func init() {
	RootCmd.AddCommand(genCmd)

	genCmd.PersistentFlags().StringVarP(&input, "input", "i", "", "input file")
	genCmd.PersistentFlags().StringVarP(&output, "output", "o", "", "output path")
	genCmd.PersistentFlags().StringVarP(&genPackageName, "package name", "p", "", "package name")
	genCmd.PersistentFlags().StringVar(&genGoPackageName, "goPackage", "", "go package name")
	genCmd.PersistentFlags().BoolVarP(&disableSQLs, "disable-sql", "", false, "disable sql generate")
	genCmd.PersistentFlags().StringVar(&plugin, "plugin", "", "The external generation plugin")
	genCmd.PersistentFlags().BoolVar(&pluginOnly, "plugin-only", false, "to generate plugin only")
	genCmd.PersistentFlags().MarkDeprecated("package name", "package flag is deprecated , should use namespace instead")
	genCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "",
		"the namespace for the generated table (or collection) , users can use `GetNamespace()` and `GetClassName()`to build their own table(collection) name if `table` not set")
}
