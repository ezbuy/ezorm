package test

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"gopkg.in/mgo.v2"

	"github.com/ezbuy/ezorm/db"
	"github.com/stretchr/testify/assert"
)

func initMongo() {
	conf := new(db.MongoConfig)
	conf.DBName = "ezorm"
	conf.MongoDB = "mongodb://127.0.0.1"
	db.Setup(conf)
}

func logMgo() {
	mgo.SetDebug(true)
	var aLogger *log.Logger
	aLogger = log.New(os.Stderr, "", log.LstdFlags)
	mgo.SetLogger(aLogger)
}

func TestBlogSave(t *testing.T) {
	initMongo()
	p := BlogMgr.NewBlog()
	p.Title = "I like ezorm"
	p.Slug = fmt.Sprintf("ezorm_%d", time.Now().Nanosecond())

	_, err := p.Save()
	if err != nil {
		t.Fatal(err)
	}

	id := p.Id()

	_, err = BlogMgr.FindByID(id)
	assert.NoError(t, err)

	_, err = BlogMgr.RemoveAll(nil)
	assert.NoError(t, err)
}

func TestCursorWithSessionRefreshed(t *testing.T) {
	initMongo()
	session, col := BlogMgr.GetCol()
	defer session.Close()

	p := BlogMgr.NewBlog()
	p.Title = "I like ezorm"
	p.Slug = fmt.Sprintf("ezorm_%d", time.Now().Nanosecond())

	_, err := p.Save()
	assert.NoError(t, err)

	rchan := make(chan struct{})
	defer close(rchan)

	go func() {
		session.Refresh()
		rchan <- struct{}{}
	}()

	p = BlogMgr.NewBlog()
	iter := col.Find(nil).Iter()
	for iter.Next(p) {
		// always wait a refresh comming
		select {
		case <-rchan:
			assert.NotNil(t, p)
		}
	}

	err = iter.Close()
	assert.NoError(t, err)

	_, err = BlogMgr.RemoveAll(nil)
	assert.NoError(t, err)
}

func TestSessionValidationWithRefresher(t *testing.T) {
	initMongo()
	logMgo()
	session, col := BlogMgr.GetCol()
	defer session.Close()

	rchan := make(chan struct{})
	defer close(rchan)

	go func() {
		session.Refresh()
		rchan <- struct{}{}
	}()

	select {
	case <-rchan:
		var i int
		for {
			if i == 50 {
				break
			}
			cs := session.Clone()
			t.Logf("%d: acquire session: %p", i, cs)
			cs.SetSyncTimeout(10 * time.Millisecond)
			_, err := cs.DB("ezorm").C(col.Name).Count()
			t.Logf("%d: query count", i)
			assert.NoError(t, err)
			i++
		}
	}
}

func TestOperationWithSessionRefreshed(t *testing.T) {
	initMongo()
	session, col := BlogMgr.GetCol()
	defer session.Close()

	p := BlogMgr.NewBlog()
	p.Title = "I like ezorm"
	p.Slug = fmt.Sprintf("ezorm_%d", time.Now().Nanosecond())

	_, err := col.UpsertId(p.ID, p)
	assert.NoError(t, err)

	id := p.Id()

	_, err = BlogMgr.FindByID(id)
	assert.NoError(t, err)

	// session refresh at this moment
	// and the socket has been reset
	session.Refresh()

	p = BlogMgr.NewBlog()
	p.Title = "I like ezorm after session has been refreshed"
	p.Slug = fmt.Sprintf("ezorm_%d", time.Now().Nanosecond())

	_, err = col.UpsertId(p.ID, p)
	assert.NoError(t, err)

	id = p.Id()

	_, err = BlogMgr.FindByID(id)
	assert.NoError(t, err)

	_, err = BlogMgr.RemoveAll(nil)
	assert.NoError(t, err)
}

func TestBlogCount(t *testing.T) {
	initMongo()
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
	initMongo()
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
