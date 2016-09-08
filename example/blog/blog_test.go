package blog

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
