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
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v2"

	"github.com/ezbuy/ezorm/v2/parser"
	"github.com/spf13/cobra"
)

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate orm code from yaml file",
	RunE: func(_ *cobra.Command, _ []string) error {
		var objs map[string]map[string]interface{}
		stat, err := os.Stat(input)
		if err != nil {
			return err
		}
		if stat.IsDir() {
			if err := filepath.WalkDir(input, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if d.IsDir() {
					return nil
				}
				if strings.HasSuffix(d.Name(), ".yaml") {
					data, err := ioutil.ReadFile(path)
					if err != nil {
						return err
					}
					err = yaml.Unmarshal(data, &objs)
					if err != nil {
						return err
					}
				}
				return nil
			}); err != nil {
				return err
			}
		} else {
			data, err := ioutil.ReadFile(input)
			if err != nil {
				return err
			}
			err = yaml.Unmarshal([]byte(data), &objs)
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

		databases := make(map[string]*parser.Obj)
		dbObjs := make(map[string][]*parser.Obj)

		for key, obj := range objs {
			xwMetaObj := new(parser.Obj)
			xwMetaObj.Package = genPackageName
			xwMetaObj.GoPackage = genGoPackageName
			xwMetaObj.Name = key
			err := xwMetaObj.Read(obj)
			if err != nil {
				return err
			}

			databases[xwMetaObj.Db] = xwMetaObj
			dbObjs[xwMetaObj.Db] = append(dbObjs[xwMetaObj.Db], xwMetaObj)
			for _, genType := range xwMetaObj.GetGenTypes() {
				fileAbsPath := output + "/gen_" + xwMetaObj.Name + "_" + genType + ".go"
				executeTpl(fileAbsPath, genType, xwMetaObj)
			}
		}

		for _, obj := range databases {
			for _, t := range obj.GetConfigTemplates() {
				fileAbsPath := output + "/gen_" + t + ".go"
				executeTpl(fileAbsPath, t, obj)
			}
		}

		sqlObjs := make(map[string]*parser.Obj, len(dbObjs))
		for db, objs := range dbObjs {
			switch db {
			default:
				continue

			case "mysql":
			}

			path := fmt.Sprintf("%s/create_%s.sql", output, db)
			genType := db + "_script"
			executeTpl(path, genType, objs)

			for _, obj := range objs {
				sqlObjs[obj.Table] = obj
			}
		}

		if !disableSQLs {
			err = handleSQL(sqlObjs, genGoPackageName)
			if err != nil {
				return err
			}
		}

		oscmd := exec.Command("gofmt", "-w", output)
		oscmd.Run()
		return nil
	},
}

func handleSQL(objs map[string]*parser.Obj, pkg string) error {
	var inputDir string
	f, err := os.Stat(input)
	if err != nil {
		return err
	}
	if f.IsDir() {
		inputDir = input
	} else {
		inputDir = filepath.Dir(input)
	}
	sqlsDir := filepath.Join(inputDir, "sqls")
	stat, err := os.Stat(sqlsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if !stat.IsDir() {
		return nil
	}

	p := parser.NewSQL(objs)
	es, err := os.ReadDir(sqlsDir)
	if err != nil {
		return err
	}

	methods := make([]*parser.SQLMethod, 0, len(es))
	for _, e := range es {
		if e.IsDir() {
			continue
		}
		path := filepath.Join(sqlsDir, e.Name())
		m, err := p.Read(path)
		if err != nil {
			return err
		}
		m.FromFile = e.Name()
		methods = append(methods, m)
	}
	if len(methods) == 0 {
		return nil
	}

	file := &parser.SQLFile{
		GoPackage: pkg,
		Methods:   methods,
		Dir:       sqlsDir,
	}
	genPath := filepath.Join(output, "gen_methods.go")
	executeTpl(genPath, "sql_method", file)

	return nil
}

func executeTpl(fileAbsPath, tplName string, obj interface{}) {
	file, err := os.OpenFile(fileAbsPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	if xwMetaObj, ok := obj.(*parser.Obj); ok {
		xwMetaObj.TplWriter = file
	}
	err = parser.Tpl.ExecuteTemplate(file, tplName, obj)
	file.Close()
	if err != nil {
		panic(err)
	}
}

var input string
var output string
var genPackageName string
var genGoPackageName string
var disableSQLs bool

func init() {
	RootCmd.AddCommand(genCmd)

	genCmd.PersistentFlags().StringVarP(&input, "input", "i", "", "input file")
	genCmd.PersistentFlags().StringVarP(&output, "output", "o", "", "output path")
	genCmd.PersistentFlags().StringVarP(&genPackageName, "package name", "p", "", "package name")
	genCmd.PersistentFlags().StringVar(&genGoPackageName, "goPackage", "", "go package name")
	genCmd.PersistentFlags().BoolVarP(&disableSQLs, "disable-sql", "", false, "disable sql generate")
}
