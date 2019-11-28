package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/ezbuy/ezorm/db"
)

func init() {
	conf := new(db.MongoConfig)
	conf.DBName = "ezorm"
	conf.MongoDB = "mongodb://127.0.0.1"
	db.Setup(conf)
}

func TestBlogSave(t *testing.T) {
	p := BlogMgr.NewBlog()
	p.Title = "I like ezorm"
	p.Slug = fmt.Sprintf("ezorm_%d", time.Now().Nanosecond())

	_, err := p.Save()
	if err != nil {
		t.Fatal(err)
	}

	id := p.Id()

	b, err := BlogMgr.FindByID(id)
	if err != nil {
		t.Fatal(err.Error())
	}

	fmt.Printf("get blog ok: %#v", b)
}

func TestBlogCount(t *testing.T) {
	now := time.Now() // get current time.Time as `Save` and `Count` condition
	p := Get_BlogMgr().NewBlog()
	p.Title = "ezorm counter"
	p.Slug = fmt.Sprintf("%d", now.Nanosecond())

	if _, err := p.Save(); err != nil {
		t.Fatalf("failed to save blog: %v", err)
	}

	result := Get_BlogMgr().Count(db.M{
		BlogMgoFieldTitle: "ezorm counter",
		BlogMgoFieldSlug:  fmt.Sprintf("%d", now.Nanosecond()),
	})
	if result != 1 {
		t.Fatalf("return value of Count should equal 1, got: %d", result)
	}

	Get_BlogMgr().RemoveBySlug(p.Slug) // cleanup the test environment
}

func TestBlogCountE(t *testing.T) {
	now := time.Now() // get current time.Time as `Save` and `CountE` condition
	p := Get_BlogMgr().NewBlog()
	p.Title = "ezorm counter"
	p.Slug = fmt.Sprintf("%d", now.Nanosecond())

	if _, err := p.Save(); err != nil {
		t.Fatalf("failed to save blog: %v", err)
	}

	result, err := Get_BlogMgr().CountE(db.M{
		BlogMgoFieldTitle: "ezorm counter",
		BlogMgoFieldSlug:  fmt.Sprintf("%d", now.Nanosecond()),
	})
	if err != nil {
		t.Fatalf("failed to run CountE: %v", err)
	}
	if result != 1 {
		t.Fatalf("return value of CountE should equal 1, got: %d", result)
	}

	Get_BlogMgr().RemoveBySlug(p.Slug) // cleanup the test environment
}
