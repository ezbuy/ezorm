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
	"github.com/spf13/cobra"
)

// genmsyamlCmd represents the genmsyaml command
var genmsyamlCmd = &cobra.Command{
	Use:   "genmsyaml",
	Short: "Generate sql server yaml file",
	Long:  "dbConfig eg: -d=\"server=...;user id=...;password=...;DATABASE=...\"",
	Run: func(cmd *cobra.Command, args []string) {
		sqlServer := db.GetSqlServer(dbConfig)
		if table != "all" {
			create(table, sqlServer)
		} else {
			tables := getAllTables(sqlServer)
			for _, t := range tables {
				create(t, sqlServer)
			}
		}

		fmt.Println("genmsyaml called")
	},
}

func create(table string, sqlServer *db.SqlServer) {
	columnsinfo := getColumnInfo(table, sqlServer)
	createYamlFile(table, columnsinfo)
	//generate(table)
}

func init() {
	RootCmd.AddCommand(genmsyamlCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// genmsyamlCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// genmsyamlCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	genmsyamlCmd.PersistentFlags().StringVarP(&table, "table", "t", "all", "table name, 'all' meaning all tables")
	genmsyamlCmd.PersistentFlags().StringVarP(&outputYaml, "output", "o", "", "output path")
	genmsyamlCmd.PersistentFlags().StringVarP(&dbConfig, "db config", "d", "", "database configuration")
	genmsyamlCmd.PersistentFlags().StringVarP(&packageName, "package name", "p", "", "package name")
}
