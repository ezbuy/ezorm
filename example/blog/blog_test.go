package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/ezbuy/ezorm/db"
	"github.com/stretchr/testify/assert"
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

func TestCursorWithSessionRefreshed(t *testing.T) {
	session, col := BlogMgr.GetCol()
	defer session.Close()

	rchan := make(chan struct{})

	go func() {
		for {
			session.Refresh()
			rchan <- struct{}{}
		}
	}()

	p := BlogMgr.NewBlog()
	iter := col.Find(nil).Iter()
	for iter.Next(p) {
		// always wait a refresh comming
		select {
		case <-rchan:
			assert.NotNil(t, p)
		}
	}

	err := iter.Close()
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
}
