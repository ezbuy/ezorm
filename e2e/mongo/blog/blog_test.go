package blog

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/ezbuy/ezorm/v2/db"
	"github.com/stretchr/testify/assert"
)

func getConfigFromEnv() *db.MongoConfig {
	return &db.MongoConfig{
		DBName: "ezorm",
		MongoDB: fmt.Sprintf(
			"mongodb://%s:%s@%s:%s",
			os.Getenv("MONGO_USER"),
			os.Getenv("MONGO_PASSWORD"),
			os.Getenv("MONGO_HOST"),
			os.Getenv("MONGO_PORT"),
		),
	}
}

func init() {
	db.Setup(getConfigFromEnv())
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

	_, err = BlogMgr.FindByID(id)
	assert.NoError(t, err)

	_, err = BlogMgr.RemoveAll(nil)
	assert.NoError(t, err)
}

func TestCursorWithSessionRefreshed(t *testing.T) {
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

func TestOperationWithSessionRefreshed(t *testing.T) {
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
