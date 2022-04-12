package test

import (
	"context"
	"time"

	"github.com/ezbuy/ezorm/v2/db"
)

var (
	_ time.Time
	_ context.Context
)

type sqlMethods struct{}

var SQL = &sqlMethods{}

type CountBlogsResp struct {
	Count0 int64
}

const _CountBlogsSQL = "SELECT /* count_blogs */ COUNT(1) FROM test_user u JOIN blog b ON u.user_id=b.blog_id WHERE u.name = ?"

func (*sqlMethods) CountBlogs(ctx context.Context, args ...interface{}) ([]*CountBlogsResp, error) {
	rows, err := db.MysqlQuery(_CountBlogsSQL, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*CountBlogsResp
	for rows.Next() {
		var o CountBlogsResp
		err = rows.Scan(&o.Count0)
		if err != nil {
			return nil, err
		}
		results = append(results, &o)
	}
	return results, nil
}

type FindBlogsResp struct {
	Id                 int32
	BlogTitle          string
	BlogHits           int32
	BlogSlug           string
	BlogBody           string
	Published          bool
	BlogGroupId        int64
	BlogCreate         time.Time
	BlogUpdate         time.Time
	TestUserUserId     int32
	TestUserUserNumber int32
	TestUserName       string
}

const _FindBlogsSQL = "SELECT /* find_blogs */ b.blog_id ID, b.title, b.hits, b.slug, IFNULL(b.body, ''), IFNULL(b.is_published, 0) published, b.group_id, b.create, b.update, u.user_id, u.user_number, u.name FROM test_user u JOIN blog b ON u.user_id=b.blog_id WHERE u.name = ? LIMIT ?, ?"

func (*sqlMethods) FindBlogs(ctx context.Context, args ...interface{}) ([]*FindBlogsResp, error) {
	rows, err := db.MysqlQuery(_FindBlogsSQL, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*FindBlogsResp
	for rows.Next() {
		var o FindBlogsResp
		err = rows.Scan(&o.Id, &o.BlogTitle, &o.BlogHits, &o.BlogSlug, &o.BlogBody, &o.Published, &o.BlogGroupId, &o.BlogCreate, &o.BlogUpdate, &o.TestUserUserId, &o.TestUserUserNumber, &o.TestUserName)
		if err != nil {
			return nil, err
		}
		results = append(results, &o)
	}
	return results, nil
}
