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
	"github.com/ezbuy/ezorm/db"
	"github.com/ezbuy/ezorm/parser"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// genmsormCmd represents the genmsorm command
var genmsormCmd = &cobra.Command{
	Use:   "genmsorm",
	Short: "Generate sql server orm code",
	Run: func(cmd *cobra.Command, args []string) {

		if table != "all" {
			columnsinfo := getColumnInfo(table)
			createYamlFile(table, columnsinfo)
			generate(table, columnsinfo)
		} else {

			generate("table", nil)
		}

		fmt.Println("genmsorm called")
	},
}

var table string
var outputYaml string

type ColumnInfo struct {
	ColumnName   string
	DataType     string
	MaxLength    int
	Nullable     bool
	IsPrimaryKey bool
}

func createYamlFile(table string, columns []ColumnInfo) {
	objs := mapper(table, columns)
	bs, err := yaml.Marshal(objs)
	if err != nil {
		println(err)
		return
	}
	fileName := "example/" + table + "_mssql.yaml"
	ioutil.WriteFile(fileName, bs, 0644)
}

func mapper(table string, columns []ColumnInfo) map[string]map[string]interface{} {
	objs := make(map[string] map[string]interface{})
	db := make(map[string] interface{})
	db["db"]="mssql"
 	objs[table]= db 
 	fields:= make(map[string] interface{})
	for i,v := range columns{
		 fields[v.ColumnName] =  parser.DbToGoType(v.DataType)
	}
	db["fields"] = fields
}

func getColumnInfo(table string) []ColumnInfo {
	query := `SELECT DISTINCT c.name AS ColumnName, t.Name AS DataType, c.max_length AS MaxLength,
    c.is_nullable AS Nullable, ISNULL(i.is_primary_key, 0) AS IsPrimaryKey
	FROM    
    sys.columns c
	INNER JOIN 
    sys.types t ON c.user_type_id = t.user_type_id
	LEFT OUTER JOIN 
    sys.index_columns ic ON ic.object_id = c.object_id AND ic.column_id = c.column_id
	LEFT OUTER JOIN 
    sys.indexes i ON ic.object_id = i.object_id AND ic.index_id = i.index_id
	WHERE
    c.object_id = OBJECT_ID(?)`

	server := db.GetSqlServer()
	rows, err := server.Query(query, table)
	server.Close()
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	defer rows.Close()
	var columninfos []ColumnInfo
	for rows.Next() {
		curent := ColumnInfo{}
		rows.Scan(&curent.ColumnName, &curent.DataType, &curent.MaxLength, &curent.Nullable, &curent.IsPrimaryKey)
		columninfos = append(columninfos, curent)
	}
	return columninfos
}

func generate(table string, columnsInfo []ColumnInfo) {
	var objs map[string]map[string]interface{}
	data, _ := ioutil.ReadFile(outputYaml)
	stat, err := os.Stat(outputYaml)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = yaml.Unmarshal([]byte(data), &objs)

	for key, obj := range objs {
			xwMetaObj := new(parser.Obj)
			xwMetaObj.Package = packageName
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

func init() {
	RootCmd.AddCommand(genmsormCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// genmsormCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// genmsormCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	genmsormCmd.PersistentFlags().StringVarP(&table, "table", "t", "all", "table name, 'all' meaning all tables")
	genmsormCmd.PersistentFlags().StringVarP(&output, "output", "o", "", "output path")
	genmsormCmd.PersistentFlags().StringVarP(&outputYaml, "output yaml", "oy", "", "output *.yaml path")
}
