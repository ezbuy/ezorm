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
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"text/template"

	yaml "gopkg.in/yaml.v2"

	"github.com/ezbuy/ezorm/parser"
	"github.com/ezbuy/ezorm/parser/sqlm"
	"github.com/ezbuy/ezorm/tpl"
	"github.com/spf13/cobra"
)

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate orm code from yaml file",
	Run: func(cmd *cobra.Command, args []string) {
		var objs map[string]map[string]interface{}
		data, _ := ioutil.ReadFile(input)
		stat, err := os.Stat(input)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = yaml.Unmarshal([]byte(data), &objs)

		if err != nil {
			println(err.Error())
			return
		}

		if genPackageName == "" {
			genPackageName = strings.Split(stat.Name(), ".")[0]
		}
		if genGoPackageName == "" {
			genGoPackageName = genPackageName
		}

		databases := make(map[string]*parser.Obj)

		allObjs := make([]*parser.Obj, 0, len(objs))
		for key, obj := range objs {
			xwMetaObj := new(parser.Obj)
			xwMetaObj.Package = genPackageName
			xwMetaObj.GoPackage = genGoPackageName
			xwMetaObj.Name = key
			err := xwMetaObj.Read(obj)
			if err != nil {
				println(err.Error())
				return
			}

			databases[xwMetaObj.Db] = xwMetaObj
			for _, genType := range xwMetaObj.GetGenTypes() {
				fileAbsPath := output + "/gen_" + xwMetaObj.Name + "_" + genType + ".go"
				executeTpl(fileAbsPath, genType, xwMetaObj)
			}
			allObjs = append(allObjs, xwMetaObj)
		}

		for _, obj := range databases {
			for _, t := range obj.GetConfigTemplates() {
				fileAbsPath := output + "/gen_" + t + ".go"
				executeTpl(fileAbsPath, t, obj)
			}
		}

		if inputSql != "" {
			obj, err := sqlm.Parse(inputSql, genGoPackageName, allObjs)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			path := output + "/gen_" + obj.Namespace + "_sqls.go"
			executeSQLTpl(path, obj)
		}

		oscmd := exec.Command("gofmt", "-w", output)
		oscmd.Run()

	},
}

func executeTpl(fileAbsPath, tplName string, xwMetaObj *parser.Obj) {
	file, err := os.OpenFile(fileAbsPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	xwMetaObj.TplWriter = file
	err = parser.Tpl.ExecuteTemplate(file, tplName, xwMetaObj)
	file.Close()
	if err != nil {
		panic(err)
	}
}

func executeSQLTpl(path string, obj *sqlm.Obj) {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	tplText := tpl.MustAsset("tpl/sql_method.gogo")
	tpl, err := template.New("sql-tpl").Parse(string(tplText))
	if err != nil {
		fmt.Printf("parse template failed: %v\n", err)
		os.Exit(1)
	}
	err = tpl.Execute(file, obj)
	if err != nil {
		fmt.Printf("execute template failed: %v\n", err)
		os.Exit(1)
	}
}

var input string
var output string
var genPackageName string
var genGoPackageName string

var inputSql string

func init() {
	RootCmd.AddCommand(genCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	genCmd.PersistentFlags().StringVarP(&input, "input", "i", "", "input model file")
	genCmd.PersistentFlags().StringVarP(&inputSql, "input-sql", "s", "", "input sql file")
	genCmd.PersistentFlags().StringVarP(&output, "output", "o", "", "output path")
	genCmd.PersistentFlags().StringVarP(&genPackageName, "package name", "p", "", "package name")
	genCmd.PersistentFlags().StringVar(&genGoPackageName, "goPackage", "", "go package name")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// genCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
