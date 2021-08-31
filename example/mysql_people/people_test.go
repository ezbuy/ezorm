package test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/ezbuy/ezorm/db"
)

func TestMain(m *testing.M) {
	// db.MysqlInit(&db.MysqlConfig{
	// DataSource: "root@tcp(localhost:3306)/?multiStatements=true",
	// })
	db.MysqlInitByField(&db.MysqlFieldConfig{
		Addr:     "localhost:3306",
		UserName: "root",
		Password: "",
		Database: "",

		Options: map[string]string{
			"multiStatements": "true",
		},
	})

	// initialize mysql database environment for running test below
	table, err := ioutil.ReadFile("people.sql")
	if err != nil {
		panic(fmt.Errorf("failed to read people table script: %s", err))
	}
	if _, err := db.MysqlExec("CREATE DATABASE IF NOT EXISTS test"); err != nil {
		panic(fmt.Errorf("failed to create database: %s", err))
	}
	if _, err := db.MysqlExec(string(table)); err != nil {
		panic(fmt.Errorf("failed to create table: %s", err))
	}

	os.Exit(m.Run())
}

func TestPeople(t *testing.T) {

	ret, err := BlogMgr.Del("1 = 1")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ret.RowsAffected(); err != nil {
		t.Fatal(err)
	}

	now := time.Now()

	blog := Blog{
		BlogId:      1,
		Title:       "BlogTitle1",
		Slug:        "blog-title",
		Body:        "hello! everybody!!!",
		User:        1,
		IsPublished: true,
		Create:      now,
		Update:      now,
	}

	if _, err := BlogMgr.Save(&blog); err != nil {
		t.Fatal(err)
		return
	}

	{
		blog2 := &Blog{
			BlogId:      2,
			Title:       "BlogTitile2",
			Slug:        "blog-title-2",
			User:        2,
			IsPublished: false,
			Create:      now,
			Update:      now,
		}
		if _, err := BlogMgr.Save(blog2); err != nil {
			t.Fatal(err)
			return
		}
	}
	{
		blog3 := &Blog{
			BlogId:      3,
			Title:       "BlogTitle3",
			Slug:        "blog-title-3",
			User:        1,
			IsPublished: true,
			Create:      now,
			Update:      now,
		}
		if _, err := BlogMgr.Save(blog3); err != nil {
			t.Fatal(err)
			return
		}
	}

	{
		b, err := BlogMgr.FindByUser(1, 0, 1, "-slug")
		if err != nil {
			t.Fatal(err)
		}
		if b[0].Title != "BlogTitle3" {
			t.Fatal("not expected")
		}
	}

	{
		blog, err := BlogMgr.FindOneBySlug("blog-title")
		if err != nil {
			t.Fatal(err)
		}
		if blog.Slug != "blog-title" {
			t.Fatal("not expected")
		}

		if blog.Create.Unix() != now.Unix() {
			t.Fatal("not expected createtime")
		}
		if blog.Update.Unix() != now.Unix() {
			t.Fatal("not expected updatetime")
		}
	}
	testForeignKey(t)
}

func testForeignKey(t *testing.T) {
	if _, err := UserMgr.Del("1=1"); err != nil {
		t.Fatal(err)
	}

	user1 := &User{
		UserNumber: 1,
		Name:       "user1",
	}
	user2 := &User{
		UserNumber: 2,
		Name:       "user2",
	}
	if _, err := UserMgr.Save(user1); err != nil {
		t.Fatal(err)
	}
	if _, err := UserMgr.Save(user2); err != nil {
		t.Fatal(err)
	}

	blogs, err := BlogMgr.FindAll()
	if err != nil {
		t.Fatal(err)
	}
	userNumbers := BlogMgr.ToFieldUser(blogs)
	users, err := UserMgr.FindListUserNumber(userNumbers)
	if err != nil {
		t.Fatal(err)
	}
	for idx, b := range blogs {
		if b.User != users[idx].UserNumber {
			t.Fatal(fmt.Sprintf("result not expected: %v", idx))
		}
	}
}
