package people

import (
	"testing"
	"time"

	"github.com/ezbuy/ezorm/db"
)

func TestPeople(t *testing.T) {
	db.MysqlInit(&db.MysqlConfig{
		DataSource: "tcp(localhost:3306)/",
	})

	ret, err := BlogMgr.Del("1 = 1")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ret.RowsAffected(); err != nil {
		t.Fatal(err)
	}

	now := time.Now()

	blog := Blog{
		Title:       "BlogTitle1",
		Slug:        "blog-title",
		Body:        "hello! everybody!!!",
		User:        1,
		IsPublished: true,
		Create:      now,
		Update:      now,
	}

	if _, err := BlogMgr.Save(&blog); err != nil {
		t.Fatal(err)
		return
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
}
