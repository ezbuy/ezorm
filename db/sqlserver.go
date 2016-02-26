package db

import (
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
)

func GetSqlServer() *sql.DB {
	sqlServer, err := sql.Open("mssql", "")
	if err != nil {
		fmt.Println("DB error", err)
	}
	return sqlServer
}
