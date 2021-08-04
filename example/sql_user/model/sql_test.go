package model

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/ezbuy/ezorm/db"
)

func TestMain(m *testing.M) {
	db.MysqlInit(&db.MysqlConfig{
		DataSource: "root:19971008@tcp(localhost:3306)/?multiStatements=true",
	})

	table, err := ioutil.ReadFile("user.sql")
	if err != nil {
		panic(fmt.Errorf("failed to read user table script: %s", err))
	}
	if _, err := db.MysqlExec("CREATE DATABASE IF NOT EXISTS `test_sql`"); err != nil {
		panic(fmt.Errorf("failed to create database: %s", err))
	}
	if _, err := db.MysqlExec(string(table)); err != nil {
		panic(fmt.Errorf("failed to create table: %s", err))
	}

	os.Exit(m.Run())
}

func TestInsert(t *testing.T) {
	ctx := context.Background()
	db := db.GetMysql()

	id, err := Methods.InsertUser(ctx, db, &User{
		Name:     "JieGe",
		Phone:    "110",
		Password: "TestPassword",
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("inserted user_id = %d\n", id)

	affected, err := Methods.InsertUserDetail(ctx, db, &UserDetail{
		UserId: id,
		Email:  "jiege@qq.com",
		Text:   "My house is very big.",
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("insert affected: %d\n", affected)

	r, err := Methods.InsertRole(ctx, db, "TestRole")
	if err != nil {
		t.Fatal(err)
	}
	rid, err := r.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}

	_, err = Methods.InsertUserRole(ctx, db, id, rid)
	if err != nil {
		t.Fatal(err)
	}
}

func TestQuery(t *testing.T) {
	ctx := context.Background()
	db := db.GetMysql()

	uds, err := Methods.ListUsers(ctx, db, 0, 10)
	if err != nil {
		t.Fatal(err)
	}

	for i, ud := range uds {
		fmt.Printf("query%d: %+v\n", i, ud)
	}

	cntResps, err := Methods.CountUsers(ctx, db)
	if err != nil {
		t.Fatal(err)
	}
	if len(cntResps) == 0 {
		t.Fatalf("not returns from count")
	}
	fmt.Printf("user count = %d\n", cntResps[0].Count0)
}
