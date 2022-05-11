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
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v2"

	"github.com/ezbuy/ezorm/v2/internal/generator"
	"github.com/ezbuy/ezorm/v2/internal/parser/mongo"
	"github.com/ezbuy/ezorm/v2/internal/parser/mysql"
	"github.com/ezbuy/ezorm/v2/internal/parser/mysqlr"
	"github.com/ezbuy/ezorm/v2/internal/parser/x/query"

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
				err = yaml.Unmarshal(data, &objs)
				if err != nil {
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
			err = yaml.Unmarshal(data, &objs)
			if err != nil {
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
			generators = append(generators, generator.NewPluginGenerator(plugin))
		}

		if err := generator.Render(generator.TMetadata{
			Meta:   objs,
			Pkg:    genGoPackageName,
			Output: output,
			Input:  input,
		}, generators...); err != nil {
			return err
		}

		oscmd := exec.Command("goimport", "-w", output)
		oscmd.Run()
		return nil
	},
}

var input string
var output string
var genPackageName string
var genGoPackageName string
var disableSQLs bool
var plugin string

func init() {
	RootCmd.AddCommand(genCmd)

	genCmd.PersistentFlags().StringVarP(&input, "input", "i", "", "input file")
	genCmd.PersistentFlags().StringVarP(&output, "output", "o", "", "output path")
	genCmd.PersistentFlags().StringVarP(&genPackageName, "package name", "p", "", "package name")
	genCmd.PersistentFlags().StringVar(&genGoPackageName, "goPackage", "", "go package name")
	genCmd.PersistentFlags().BoolVarP(&disableSQLs, "disable-sql", "", false, "disable sql generate")
	genCmd.PersistentFlags().StringVar(&plugin, "plugin", "", "The external generation plugin")
}
