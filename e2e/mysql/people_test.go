package mysql

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/ezbuy/ezorm/v2/db"
	"github.com/stretchr/testify/assert"
)

func mysqlConfigFromEnv() *db.MysqlFieldConfig {
	return &db.MysqlFieldConfig{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("MYSQL_HOST"), os.Getenv("MYSQL_PORT")),
		UserName: "ezbuy",
		Password: "ezbuyisthebest",
		Database: "test",

		Options: map[string]string{
			"multiStatements": "true",
		},
	}
}

func TestMain(m *testing.M) {
	db.MysqlInitByField(mysqlConfigFromEnv())

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
		GroupId:     1,
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
			GroupId:     1,
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
			GroupId:     2,
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
	{
		blogs, err := BlogMgr.FindInGroupId([]int64{1, 2})
		if err != nil {
			t.Fatal(err)
		}
		if len(blogs) != 3 {
			t.Fatal("not expected")
		}

		for _, blog := range blogs {
			switch blog.BlogId {
			case 1:
				if blog.Slug != "blog-title" {
					t.Fatal("not expected id 1")
				}

			case 2:
				if blog.Slug != "blog-title-2" {
					t.Fatal("not expected id 2")
				}

			case 3:
				if blog.Slug != "blog-title-3" {
					t.Fatal("not expected id 3")
				}

			default:
				t.Fatalf("not expected id %d", blog.BlogId)
			}
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
			t.Fatalf("result not expected: %d", idx)
		}
	}
}

func TestRawQuery(t *testing.T) {
	ctx := context.Background()
	if _, err := UserMgr.Del("1=1"); err != nil {
		t.Fatal(err)
	}

	user1 := &User{
		UserNumber: 1,
		Name:       "user1",
	}
	if _, err := UserMgr.Save(user1); err != nil {
		t.Fatal(err)
	}
	resp, err := SQL.GetUser(ctx, &GetUserReq{
		Name: "user1",
	})
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, len(resp), 1)
	assert.Equal(t, resp[0].Name, "user1")
}
