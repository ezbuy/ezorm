package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/bmizerany/assert"
	"github.com/ezbuy/ezorm/db"
)

func TestPeopleObject(t *testing.T) {
	dsn := fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s?charset=utf8&autocommit=true&parseTime=True",
		"root",
		"147011",
		"ezorm")

	db.MysqlInit(&db.MysqlConfig{
		DataSource: dsn,
	})
	now := time.Now()
	user1 := UserMgr.NewUser()
	user1.Name = "user01"
	user1.Mailbox = "user01@sss.fff"
	user1.HeadUrl = "aaaa.png"
	user1.Password = "123456"
	user1.CreatedAt = now
	user1.UpdatedAt = now
	user1.Longitude = 103.754
	user1.Latitude = 1.3282

	UserMgr.Del("id > 0")
	BlogMgr.Del("id > 0")
	_, err := UserMgr.Save(user1)
	assert.Equal(t, nil, err)

	blog11 := Blog{
		UserId:    user1.Id,
		Title:     "BlogTitle1",
		Content:   "hello! everybody!!!",
		Status:    1,
		Readed:    10,
		CreatedAt: now,
		UpdatedAt: now,
	}
	BlogMgr.Save(&blog11)

	blog12 := Blog{
		UserId:    user1.Id,
		Title:     "BlogTitle1222",
		Content:   "hello! everybody!!!",
		Status:    1,
		Readed:    10,
		CreatedAt: now,
		UpdatedAt: now,
	}
	BlogMgr.Save(&blog12)

	user2 := UserMgr.NewUser()
	user2.Name = "user02"
	user2.Mailbox = "user201@sss.fff"
	user2.HeadUrl = "aaaa.png"
	user2.Password = "123456"
	user2.CreatedAt = now
	user2.UpdatedAt = now
	user2.Longitude = 103.754
	user2.Latitude = 1.3282

	_, err = UserMgr.Save(user2)
	assert.Equal(t, nil, err)

	blog21 := Blog{
		UserId:    user2.Id,
		Title:     "BlogTitle1",
		Content:   "hello! everybody!!!",
		Status:    1,
		Readed:    10,
		CreatedAt: now,
		UpdatedAt: now,
	}
	BlogMgr.Save(&blog21)

	blog22 := Blog{
		UserId:    user2.Id,
		Title:     "BlogTitle1222",
		Content:   "hello! everybody!!!",
		Status:    1,
		Readed:    12,
		CreatedAt: now,
		UpdatedAt: now,
	}
	BlogMgr.Save(&blog22)

	blog23 := Blog{
		UserId:    user2.Id,
		Title:     "BlogTitle1222",
		Content:   "hello! everybody!!!",
		Status:    1,
		Readed:    18,
		CreatedAt: now,
		UpdatedAt: now,
	}
	BlogMgr.Save(&blog23)

	user3 := UserMgr.NewUser()
	user3.Name = "user03"
	user3.Mailbox = "use301@sss.fff"
	user3.HeadUrl = "aaaa.png"
	user3.Password = "123456"
	user3.CreatedAt = now
	user3.UpdatedAt = now
	user3.Longitude = 103.754
	user3.Latitude = 1.3282

	_, err = UserMgr.Save(user3)
	assert.Equal(t, nil, err)

	RedisSetUp(&RedisConfig{
		Host: "127.0.0.1",
		Port: 6379,
	})

	//! redis 数据准备
	BlogMgr.Clear()
	assert.Equal(t, BlogMgr.Import(), err)

	UserMgr.Clear()
	assert.Equal(t, UserMgr.Import(), err)

	UserBlogMgr.Clear()
	assert.Equal(t, UserBlogMgr.Import(), err)

	SortUserBlogMgr.Clear()
	assert.Equal(t, SortUserBlogMgr.Import(), err)

	UserLocationMgr.Clear()
	assert.Equal(t, UserLocationMgr.Import(), err)

	c1, err := BlogMgr.ListCount()
	fmt.Println("BlogMgr.ListCount =>", c1, err)

	c2, err := UserMgr.ListCount()
	fmt.Println("UserMgr.ListCount =>", c2, err)

	//! redis 数据查询
	blogs1, err := BlogMgr.GetByUserId(user2.Id)
	assert.Equal(t, 3, len(blogs1))

	blogs2, err := UserBlogMgr.RelatedBlogs(user1.Id)
	assert.Equal(t, 2, len(blogs2))

	//! redis 数据增量同步
	blog24 := Blog{
		UserId:    user2.Id,
		Title:     "BlogTitle1222444",
		Content:   "hello! everybody!!!",
		Status:    1,
		Readed:    19,
		CreatedAt: now,
		UpdatedAt: now,
	}
	r, err := BlogMgr.Save(&blog24)
	assert.Equal(t, err, nil)
	fmt.Println("save =>", r, blog24, user2.Id)

	BlogMgr.AddBySQL("SELECT `id`,`user_id`,`title`,`content`,`status`,`readed`, `created_at`, `updated_at` FROM blogs WHERE id = ?", blog24.Id)
	UserBlogMgr.AddBySQL("SELECT user_id, id FROM blogs WHERE id = ?", blog24.Id)

	ids, err := UserBlogMgr.SetGet(fmt.Sprint(user2.Id))
	assert.Equal(t, 4, len(ids))

	blogs3, err := UserBlogMgr.RelatedBlogs(user2.Id)
	assert.Equal(t, 4, len(blogs3))
}
