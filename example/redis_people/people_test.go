package test

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/bmizerany/assert"

	redis "gopkg.in/redis.v5"
)

func TestPeopleObject(t *testing.T) {

	var cmd redis.Cmdable

	cmd = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := cmd.Ping().Result()
	fmt.Println(pong, err)

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

	blog4 := Blog{
		BlogId:      4,
		Title:       "BlogTitile2",
		Slug:        "blog-title-2",
		User:        1,
		IsPublished: false,
		Create:      now,
		Update:      now,
	}
	assert.Equal(t, BlogMgr.DelBlog(&blog4), nil)
	assert.Equal(t, BlogMgr.SetBlog(&blog4), nil)

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
	log.Println("get blogs =>", blogs)
	assert.Equal(t, err, nil)
	assert.Equal(t, len(blogs), 3)

	blogs2, err := BlogMgr.GetBlogsByIndexes(map[string]interface{}{
		"User":        1,
		"IsPublished": false,
	})
	assert.Equal(t, err, nil)
	log.Println("get blogs =>", blogs2)
	assert.Equal(t, len(blogs2), 2)

	allblogs, err := BlogMgr.ListRange(0, -1)
	assert.Equal(t, err, nil)
	log.Println("all blogs =>", allblogs)
	count, err := BlogMgr.ListCount()
	assert.Equal(t, err, nil)
	assert.Equal(t, int64(4), count)

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

	//! set test
	for i := 1; i < 5; i++ {
		ub := UserBlogMgr.NewUserBlog()
		ub.Value = int32(i)
		err := UserBlogMgr.SetAdd(fmt.Sprint(1), ub)
		assert.Equal(t, err, nil)
	}

	relations, err := UserBlogMgr.SetGet(fmt.Sprint(1))
	assert.Equal(t, err, nil)
	assert.Equal(t, len(relations), 4)

	relatedblogs, err := UserBlogMgr.RelatedBlogs(fmt.Sprint(1))
	assert.Equal(t, err, nil)
	assert.Equal(t, len(relatedblogs), 4)
	log.Println("relatedblogs =>", relatedblogs)

	urm := UserBlogMgr.NewUserBlog()
	urm.Value = int32(1)
	assert.Equal(t, UserBlogMgr.SetRem(fmt.Sprint(1), urm), nil)

	relations, err = UserBlogMgr.SetGet(fmt.Sprint(1))
	assert.Equal(t, err, nil)
	assert.Equal(t, len(relations), 3)

	assert.Equal(t, UserBlogMgr.SetDel(fmt.Sprint(1)), nil)
	relations, err = UserBlogMgr.SetGet(fmt.Sprint(1))
	assert.Equal(t, err, nil)
	assert.Equal(t, len(relations), 0)

	//! zset test
	for i := 1; i < 5; i++ {
		ub := SortUserBlogMgr.NewSortUserBlog()
		ub.Value = int32(i)
		switch i {
		case 1:
			ub.Score = 0.1
			err := SortUserBlogMgr.ZAdd(fmt.Sprint(1), ub)
			assert.Equal(t, err, nil)
		case 2:
			ub.Score = 1.1
			err := SortUserBlogMgr.ZAdd(fmt.Sprint(1), ub)
			assert.Equal(t, err, nil)
		case 3:
			ub.Score = 2.1
			err := SortUserBlogMgr.ZAdd(fmt.Sprint(1), ub)
			assert.Equal(t, err, nil)
		case 4:
			ub.Score = 3.1
			err := SortUserBlogMgr.ZAdd(fmt.Sprint(1), ub)
			assert.Equal(t, err, nil)
		}
	}

	zrelations, err := SortUserBlogMgr.ZRangeByScore(fmt.Sprint(1), 2, 5)
	assert.Equal(t, err, nil)
	assert.Equal(t, len(zrelations), 2)

	//! list test
	UserIdMgr.ListDel("key")
	for i := 1; i < 5; i++ {
		uid := UserIdMgr.NewUserId()
		uid.Value = int32(i)
		length, err := UserIdMgr.ListLPush("key", uid)
		assert.Equal(t, err, nil)
		assert.Equal(t, length, int64(i))
	}

	uids, err := UserIdMgr.ListRange("key", 0, 3)
	assert.Equal(t, err, nil)
	assert.Equal(t, 4, len(uids))

	//! geo test

}
