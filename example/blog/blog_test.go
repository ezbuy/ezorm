package test

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/ezbuy/ezorm/db"
	"github.com/stretchr/testify/assert"
)

func initMongo() *db.MongoConfig {
	conf := new(db.MongoConfig)
	conf.DBName = "ezorm"
	conf.MongoDB = "mongodb://localhost"
	db.Setup(conf)
	return conf
}

func initMongoReplicaSet() *db.MongoConfig {
	conf := new(db.MongoConfig)
	conf.DBName = "ezorm"
	conf.MongoDB = "mongodb://localhost:27017,localhost:27016,localhost:27015?connect=replicaSet"
	db.Setup(conf)
	return conf
}

func logMgo() {
	mgo.SetDebug(true)
	aLogger := log.New(os.Stderr, "", log.LstdFlags)
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

// TestSessionCloneNoCloseWithRefresher 测试Refresh之后再clone session又没有关闭引起的连接池中的socket被用完的场景
//
//   每次Clone都会使用新的socket，所以会导致socket pool被用完
//   所以很有必要在Refresh之后Ping一次去重新固定socket到cluster(masterSocket/slaveSocket)里面以便之后的重用
func TestSessionRefresher_NoReplicaSet(t *testing.T) {
	conf := initMongo()
	session, col := BlogMgr.GetCol()
	defer session.Close()

	db.SetupIdleSessionRefresher(conf, []*mgo.Session{session}, 200*time.Millisecond)

	var i int
	for {
		if i == 20 {
			break
		}
		// no master or salve socket set
		// an empty socket
		cs := session.Clone()
		_, err := cs.DB("ezorm").C(col.Name).Find(bson.M{`$where`: `sleep(10000) || "true"`}).Count()
		assert.NoError(t, err)
		i++
		time.Sleep(200 * time.Millisecond)
		cs.Close()
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
