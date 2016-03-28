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

	"gopkg.in/yaml.v2"

	"github.com/ezbuy/ezorm/parser"
	"github.com/spf13/cobra"
)

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate orm code",
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
		for key, obj := range objs {
			xwMetaObj := new(parser.Obj)
			xwMetaObj.Package = genPackageName
			xwMetaObj.Name = key
			err := xwMetaObj.Read(obj)
			if err != nil {
				println(err.Error())
			}

			for _, genType := range xwMetaObj.GetGenTypes() {
				file, err := os.OpenFile(output+"/gen_"+xwMetaObj.Name+"_"+genType+".go", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
				xwMetaObj.TplWriter = file
				if err != nil {
					panic(err)
				}

				err = parser.Tpl.ExecuteTemplate(file, genType, xwMetaObj)
				file.Close()
				if err != nil {
					println(err.Error())
				}
			}
		}

		oscmd := exec.Command("gofmt", "-w", output)
		oscmd.Run()
		fmt.Println("gen called")
	},
}

var input string
var output string
var genPackageName string

func init() {
	RootCmd.AddCommand(genCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	genCmd.PersistentFlags().StringVarP(&input, "input", "i", "", "input file")
	genCmd.PersistentFlags().StringVarP(&output, "output", "o", "", "output path")
	genCmd.PersistentFlags().StringVarP(&genPackageName, "package name", "p", "", "package name")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// genCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
