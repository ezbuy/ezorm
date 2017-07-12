// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
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
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/ezbuy/ezorm/parser"
	"github.com/spf13/cobra"
)

// mssqlconfigCmd represents the genmssqlconfig command
var mongoconfigCmd = &cobra.Command{
	Use:   "mongoconfig",
	Short: "Generate mongo config code",
	Long:  `Generate mongo config code`,
	Run: func(cmd *cobra.Command, args []string) {
		genPackageName = strings.TrimSpace(genPackageName)

		if genPackageName == "" {
			log.Fatalln("package name must be provided")
		}

		regx := regexp.MustCompile(`^[a-zA-Z]+$`)
		if !regx.MatchString(genPackageName) {
			log.Fatalln("package name invalid")
		}

		if output == "" {
			log.Fatalln("output folder must be provided")
		}

		if err := os.MkdirAll(output, os.ModeDir|os.ModePerm); err != nil {
			log.Fatalf("create folder %q error: %s", output, err)
		}

		metaObj := new(parser.Obj)
		metaObj.GoPackage = genPackageName

		fileAbsPath := output + "/gen_mongo_config.go"
		executeTpl(fileAbsPath, "mongo_config", metaObj)
	},
}

func init() {
	mongoconfigCmd.PersistentFlags().StringVarP(&output, "output", "o", "", "output path")
	mongoconfigCmd.PersistentFlags().StringVarP(&genPackageName, "package", "p", "", "package name")

	RootCmd.AddCommand(mongoconfigCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mssqlconfigCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// mssqlconfigCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
