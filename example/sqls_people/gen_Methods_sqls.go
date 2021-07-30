package model

import (
	"database/sql"
	"github.com/ezbuy/ezorm/sqlm"
	"time"
)

var _ *time.Time
var _ *sql.Rows
var _ sqlm.Execable

type _Methods struct{}

var Methods = _Methods{}

const _ListUserDetailsSql = "SELECT `test_user`.`user_id` Id, `test_user`.`name`, `user_detail`.`desc`, `user_detail`.`age`, `user_detail`.`phone`, `user_detail`.`email` FROM `test_user` JOIN `user_detail` ON `test_user`.`user_id`=`user_detail`.`user_id` LIMIT ?, ?"

func (_Methods) ListUserDetails(db sqlm.Queryable, offset int32, limit int32) (ret []UserDetailResult, err error) {
	err = sqlm.QueryMany(db, _ListUserDetailsSql, []interface{}{offset, limit}, func(rows *sql.Rows) error {
		var v UserDetailResult
		err = rows.Scan(&v.Id, &v.UserName, &v.UserDetailDesc, &v.UserDetailAge, &v.UserDetailPhone, &v.UserDetailEmail)
		if err != nil {
			return err
		}
		ret = append(ret, v)
		return nil
	})
	return
}

// UserBlog is the auto-generated struct for method ListUserBlogsByName
type UserBlog struct {
	BlogId          int32  `table:"blog" field:"blog_id"`
	BlogTitle       string `table:"blog" field:"title"`
	BlogSlug        string `table:"blog" field:"slug"`
	BlogBody        string `table:"blog" field:"body"`
	BlogIsPublished bool   `table:"blog" field:"is_published"`
	UserNumber      int32  `table:"test_user" field:"user_number"`
}

const _ListUserBlogsByNameSql = "SELECT `blog`.`blog_id` BlogId, `blog`.`title`, `blog`.`slug`, `blog`.`body`, `blog`.`is_published`, `test_user`.`user_number` UserNumber FROM `test_user` JOIN `blog` ON `blog`.`user`=`test_user`.`user_id` WHERE `test_user`.`name`=? LIMIT ?, ?"

func (_Methods) ListUserBlogsByName(db sqlm.Queryable, userName string, offset int32, limit int32) (ret []*UserBlog, err error) {
	err = sqlm.QueryMany(db, _ListUserBlogsByNameSql, []interface{}{userName, offset, limit}, func(rows *sql.Rows) error {
		var v UserBlog
		err = rows.Scan(&v.BlogId, &v.BlogTitle, &v.BlogSlug, &v.BlogBody, &v.BlogIsPublished, &v.UserNumber)
		if err != nil {
			return err
		}
		ret = append(ret, &v)
		return nil
	})
	return
}

// UserDetailResult is the auto-generated struct for method GetUserDetailById
type UserDetailResult struct {
	Id              int32  `table:"test_user" field:"user_id"`
	UserName        string `table:"test_user" field:"name"`
	UserDetailDesc  string `table:"user_detail" field:"desc"`
	UserDetailAge   int32  `table:"user_detail" field:"age"`
	UserDetailPhone string `table:"user_detail" field:"phone"`
	UserDetailEmail string `table:"user_detail" field:"email"`
}

const _GetUserDetailByIdSql = "SELECT `test_user`.`user_id` Id, `test_user`.`name`, `user_detail`.`desc`, `user_detail`.`age`, `user_detail`.`phone`, `user_detail`.`email` FROM `test_user` JOIN `user_detail` ON `test_user`.`user_id`=`user_detail`.`user_id` WHERE `test_user`.`user_id`=?"

func (_Methods) GetUserDetailById(db sqlm.Queryable, id int64) (ret *UserDetailResult, err error) {
	err = sqlm.QueryOne(db, _GetUserDetailByIdSql, []interface{}{id}, func(rows *sql.Rows) error {
		var v UserDetailResult
		if err := rows.Scan(&v.Id, &v.UserName, &v.UserDetailDesc, &v.UserDetailAge, &v.UserDetailPhone, &v.UserDetailEmail); err != nil {
			return err
		}
		ret = &v
		return nil
	})
	return
}

const _GetBlogByIdSql = "SELECT `blog`.`title`, `blog`.`body` FROM `blog` WHERE `blog`.`blog_id`=?"

func (_Methods) GetBlogById(db sqlm.Queryable, bid int64) (ret *Blog, err error) {
	err = sqlm.QueryOne(db, _GetBlogByIdSql, []interface{}{bid}, func(rows *sql.Rows) error {
		var v Blog
		if err := rows.Scan(&v.Title, &v.Body); err != nil {
			return err
		}
		ret = &v
		return nil
	})
	return
}
