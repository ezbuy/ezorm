package test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/ezbuy/ezorm/db"
)

const tablePrefix = `
USE test;

DROP TABLE IF EXISTS user;
`

func TestMain(m *testing.M) {
	db.MysqlInitByField(&db.MysqlFieldConfig{
		Addr:     "localhost:3306",
		UserName: "root",
		Password: "",
		Database: "",

		Options: map[string]string{
			"multiStatements": "true",
		},
	})

	table, err := ioutil.ReadFile("create_mysql.sql")
	if err != nil {
		panic(fmt.Errorf("failed to read create_mysql.sql script: %v", err))
	}
	if _, err := db.MysqlExec("CREATE DATABASE IF NOT EXISTS test"); err != nil {
		panic(fmt.Errorf("failed to create database: %s", err))
	}

	create := string(table)
	create = tablePrefix + create

	if _, err := db.MysqlExec(create); err != nil {
		panic(fmt.Errorf("failed to create table: %s", err))
	}

	os.Exit(m.Run())
}

func TestUser(t *testing.T) {
	user := &User{
		Name:       "tang",
		Phone:      "110",
		Age:        24,
		Balance:    10.23,
		CreateDate: time.Now().Unix(),
	}
	if _, err := UserMgr.Save(user); err != nil {
		t.Fatal(err)
		return
	}

	users, err := UserMgr.FindByNamePhone("tang", "110", 0, 1)
	if err != nil {
		t.Fatal(err)
		return
	}
	if len(users) != 1 {
		t.Fatalf("user length not expected")
	}

	user = users[0]

	if user.Age != 24 || user.Balance != 10.23 || user.Text != "" {
		t.Fatalf("user info not expected")
	}
}
