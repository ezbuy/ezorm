// Copyright © 2021 NAME HERE ezbuy TEAM
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
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/ezbuy/ezorm/v2/internal/template"
	"github.com/spf13/cobra"

	"gopkg.in/yaml.v2"
)

var driver string
var outputDirpath string
var inputDirpath string

var RootCmd = &cobra.Command{
	Use:   "ezorm",
	Short: "ezorm",
	Long: `ezorm is an code-generation based ORM lib for golang, supporting mongodb/sql server/mysql.
data model is defined with YAML file`,
	Run: func(_ *cobra.Command, _ []string) {
		if err := filepath.WalkDir(inputDirpath, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			ext := filepath.Ext(path)
			if ext != ".yaml" && ext != ".yml" {
				return nil
			}
			data, err := os.ReadFile(d.Name())
			if err != nil {
				return err
			}
			value := make(map[string]interface{})
			if err := yaml.Unmarshal(data, value); err != nil {
				return err
			}
			g := template.NewGenerator(outputDirpath)
			dname, ok := value["driver"]
			if !ok {
				return errors.New("parse yaml: driver field not set")
			}
			dr, err := template.GetDriver(dname.(string))
			if err != nil {
				return err
			}
			ll := filepath.SplitList(outputDirpath)
			value["GoPackage"] = ll[len(ll)-1]
			if err := g.Generate(context.TODO(), dr, value); err != nil {
				return err
			}
			return nil
		}); err != nil {
			log.Fatal(err)
		}
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&driver, "driver", "D", "", "The driver of ezorm generate template")
	RootCmd.PersistentFlags().StringVarP(&outputDirpath, "output", "O", "", "output directory path")
	RootCmd.PersistentFlags().StringVarP(&inputDirpath, "input", "I", "", "input directory path")
}
