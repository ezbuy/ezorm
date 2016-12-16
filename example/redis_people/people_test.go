package test

import (
	"log"
	"testing"
	"time"

	"github.com/bmizerany/assert"
)

func TestPeopleObject(t *testing.T) {
	RedisSetUp(&RedisConfig{
		Host: "127.0.0.1",
		Port: 6379,
		DB:   1,
	})

	now := time.Now()

	blog := Blog{
		BlogId:      1,
		Title:       "BlogTitle1",
		Slug:        "blog-title",
		Body:        "hello! everybody!!!",
		User:        1,
		IsPublished: false,
		Create:      now,
		Update:      now,
	}
	assert.Equal(t, BlogMgr.DelBlog(&blog), nil)
	assert.Equal(t, BlogMgr.SetBlog(&blog), nil)

	blog2 := Blog{
		BlogId:      2,
		Title:       "BlogTitile2",
		Slug:        "blog-title-2",
		User:        2,
		IsPublished: false,
		Create:      now,
		Update:      now,
	}
	assert.Equal(t, BlogMgr.DelBlog(&blog2), nil)
	assert.Equal(t, BlogMgr.SetBlog(&blog2), nil)

	blog3 := Blog{
		BlogId:      3,
		Title:       "BlogTitle3",
		Slug:        "blog-title-3",
		User:        1,
		IsPublished: true,
		Create:      now,
		Update:      now,
	}
	assert.Equal(t, BlogMgr.DelBlog(&blog3), nil)
	assert.Equal(t, BlogMgr.SetBlog(&blog3), nil)

	b := BlogMgr.NewBlog()
	b.BlogId = 2
	assert.Equal(t, BlogMgr.GetBlog(b), nil)
	assert.Equal(t, blog2.Title, b.Title)
	assert.Equal(t, blog2.Slug, b.Slug)
	assert.Equal(t, blog2.User, b.User)
	assert.Equal(t, blog2.IsPublished, b.IsPublished)
	assert.Equal(t, blog2.Create.Unix(), b.Create.Unix())
	assert.Equal(t, blog2.Update.Unix(), b.Update.Unix())

	log.Println("get blog =>", b)

	blogs, err := BlogMgr.GetBlogsByUser(1)
	assert.Equal(t, err, nil)
	assert.Equal(t, len(blogs), 2)

	blogs2, err := BlogMgr.GetBlogsByIndexes(map[string]interface{}{
		"User":        1,
		"IsPublished": false,
	})
	assert.Equal(t, err, nil)
	assert.Equal(t, len(blogs2), 1)

	user := UserMgr.NewUser()
	user.UserId = 101
	user.Name = "username"
	user.UserNumber = 9527
	user.Create = now
	user.Update = now
	assert.Equal(t, UserMgr.DelUser(user), nil)
	assert.Equal(t, UserMgr.SetUser(user), nil)

	u2 := UserMgr.NewUser()
	u2.UserId = 101

	assert.Equal(t, UserMgr.GetUser(u2), nil)
	assert.Equal(t, user.UserId, u2.UserId)
	assert.Equal(t, user.Name, u2.Name)
	assert.Equal(t, user.UserNumber, u2.UserNumber)
	assert.Equal(t, user.Create.Unix(), u2.Create.Unix())
	assert.Equal(t, user.Update.Unix(), u2.Update.Unix())

	r := UserBlogsMgr.NewUserBlogs()
	r.UserId = 101
	for i := 1; i <= 3; i++ {
		r.BlogId = int32(i)
		assert.Equal(t, UserBlogsMgr.SetUserBlogs(r), nil)
	}

	sr := SortUserBlogsMgr.NewSortUserBlogs()
	sr.UserId = 101
	for i := 1; i <= 3; i++ {
		sr.BlogId = int32(i)
		assert.Equal(t, SortUserBlogsMgr.SetSortUserBlogs(sr), nil)
	}

	pos := UserLocationMgr.NewUserLocation()
	pos.RegionId = 100000000
	pos.Longitude = 103.0232
	pos.Latitude = 30.0343
	pos.UserId = 100
	for i := 1; i <= 3; i++ {
		pos.UserId++
		assert.Equal(t, UserLocationMgr.SetUserLocation(pos), nil)
	}
}
