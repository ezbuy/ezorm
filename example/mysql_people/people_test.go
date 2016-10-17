package people

import (
	"reflect"
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
		blog2 := &Blog{
			Title:       "BlogTitile2",
			Slug:        "blog-title-2",
			User:        2,
			IsPublished: false,
		}
		if _, err := BlogMgr.Save(blog2); err != nil {
			t.Fatal(err)
			return
		}
	}
	{
		blog3 := &Blog{
			Title:       "BlogTitle3",
			Slug:        "blog-title-3",
			User:        1,
			IsPublished: true,
		}
		if _, err := BlogMgr.Save(blog3); err != nil {
			t.Fatal(err)
			return
		}
	}

	grouped, err := BlogMgr.GroupByUnAssigned([]int32{1, 2}, -1, -1)
	if err != nil {
		t.Fatal(err)
		return
	}
	expectedGroup := &BlogGroupUnAssigned{
		COUNT:       []int{1, 2},
		Hits:        []int32{0, 0},
		IsPublished: []bool{false, true},
	}
	if !reflect.DeepEqual(expectedGroup, grouped) {
		t.Fatal("result not expect")
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
	testJoin(t)
}

func testJoin(t *testing.T) {
	blogs, err := UserMgr.LeftJoinBlog([]int32{1, 2})
	if err != nil {
		t.Fatal(err)
	}
	if len(blogs) != 2 {
		t.Fatal("result not expected")
	}
}
