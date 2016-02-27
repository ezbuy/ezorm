package example

import (
	"fmt"
	"testing"

	"github.com/ezbuy/ezorm/db"
)

func TestGetRefIf(t *testing.T) {
	conf := new(db.MongoConfig)
	conf.DBName = "ezorm"
	conf.MongoDB = "mongodb://127.0.0.1"
	db.Setup(conf)
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
