package example

import (
	"fmt"
	"testing"

	"github.com/ezbuy/ezorm/db"
)

func init() {
	conf := new(db.MongoConfig)
	conf.DBName = "ezorm"
	conf.MongoDB = "mongodb://127.0.0.1"
	db.Setup(conf)
}

func TestBlog(t *testing.T) {
	b, _ := BlogMgr.FindOneBySlug("bingo")
	if b != nil {
		BlogMgr.RemoveByID(b.Id())
	}

	b = BlogMgr.NewBlog()
	b.Slug = "bingo"
	_, err := b.Save()
	if err != nil {
		t.Error(err)
	}

	b = BlogMgr.NewBlog()
	b.Slug = "bingo"
	_, err = b.Save()
	if err == nil {
		t.Error(err)
	}

	b, err = BlogMgr.FindOneBySlug("bingo")
	if err != nil {
		t.Error(err)
	}
	BlogMgr.RemoveByID(b.Id())
}

func TestPage(t *testing.T) {
	p := PageMgr.NewPage()
	p.Hits = 19
	p.Title = "bingo"
	p.Sections = make([]Section, 1)
	section := Section{}
	section.Key = "key1"
	section.Val = 2
	section.Data = make(map[string]string)
	section.Data["foo"] = "bar"
	p.Sections[0] = section
	p.Slug = "ezorm"
	p.Save()

	p, err := PageMgr.FindOneBySlug("ezorm")
	if err != nil {
		t.Error("find fail")
	}
	fmt.Println("%v", p)
	PageMgr.RemoveByID(p.Id())

	_, err = PageMgr.FindOneBySlug("ezorm")
	if err == nil {
		t.Error("delete fail")
	}
}
