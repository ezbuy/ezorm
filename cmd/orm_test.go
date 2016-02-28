package cmd

import (
	"fmt"
	"testing"
	"github.com/ezbuy/ezorm/db"
	"github.com/ezbuy/ezorm/page"
)



	conf := new(db.MongoConfig)
	conf.DBName = "ezorm"
	conf.MongoDB = "mongodb://127.0.0.1"
	db.Setup(conf)
	p := page.PageMgr.NewPage()
	p.Hits = 19
	p.Title = "bingo"
	p.Sections = make([]page.Section, 1)
	section := page.Section{}
	section.Key = "key1"
	section.Val = 2
	section.Data = make(map[string]string)
	section.Data["foo"] = "bar"
	p.Sections[0] = section
	p.Slug = "ezorm"
	p.Save()

	p, err := page.PageMgr.FindBySlug("ezorm")
	if err != nil {
		t.Error("find fail")
	}
	fmt.Println("%v", p)
	// page.PageMgr.RemoveByID(p.Id())

	// _, err = page.PageMgr.FindBySlug("ezorm")
	// if err == nil {
	// 	t.Error("delete fail")
	// }
}
