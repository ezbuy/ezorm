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

// UserDetailResult is the auto-generated struct for method ListUserDetails
type UserDetailResult struct {
	UserId                 int64  `table:"user" field:"id"`
	UserName               string `table:"user" field:"name"`
	UserPhone              string `table:"user" field:"phone"`
	UserPassword           string `table:"user" field:"password"`
	UserDetailEmail        string `table:"user_detail" field:"email"`
	UserDetailIntroduction string `table:"user_detail" field:"introduction"`
	UserDetailAge          int32  `table:"user_detail" field:"age"`
	UserDetailAvatar       string `table:"user_detail" field:"avatar"`
}

const _ListUserDetailsSql = "SELECT `user`.`id`, `user`.`name`, `user`.`phone`, `user`.`password`, `user_detail`.`email`, `user_detail`.`introduction`, `user_detail`.`age`, `user_detail`.`avatar` FROM `user` JOIN `user_detail` ON `user`.`id`=`user_detail`.`user_id` LIMIT ?, ?"

func (_Methods) ListUserDetails(db sqlm.Queryable, offset int32, limit int32) (ret []*UserDetailResult, err error) {
	err = sqlm.QueryMany(db, _ListUserDetailsSql, []interface{}{offset, limit}, func(rows *sql.Rows) error {
		var v UserDetailResult
		err = rows.Scan(&v.UserId, &v.UserName, &v.UserPhone, &v.UserPassword, &v.UserDetailEmail, &v.UserDetailIntroduction, &v.UserDetailAge, &v.UserDetailAvatar)
		if err != nil {
			return err
		}
		ret = append(ret, &v)
		return nil
	})
	return
}

const _FindUsersByRoleSql = "SELECT `user`.`id`, `user`.`name`, `user`.`phone`, `user`.`password` FROM `user` JOIN `role_user` ON `role_user`.`user_id`=`user`.`id` WHERE `role_user`.`role_id`=?"

func (_Methods) FindUsersByRole(db sqlm.Queryable, roleId int64) (ret []*User, err error) {
	err = sqlm.QueryMany(db, _FindUsersByRoleSql, []interface{}{roleId}, func(rows *sql.Rows) error {
		var v User
		err = rows.Scan(&v.Id, &v.Name, &v.Phone, &v.Password)
		if err != nil {
			return err
		}
		ret = append(ret, &v)
		return nil
	})
	return
}

const _FindRolesByUserSql = "SELECT `role`.`id`, `role`.`name` FROM `role` JOIN `role_user` ON `role_user`.`role_id`=`role`.`id` WHERE `role_user`.`user_id`=?"

func (_Methods) FindRolesByUser(db sqlm.Queryable, userId int64) (ret []*Role, err error) {
	err = sqlm.QueryMany(db, _FindRolesByUserSql, []interface{}{userId}, func(rows *sql.Rows) error {
		var v Role
		err = rows.Scan(&v.Id, &v.Name)
		if err != nil {
			return err
		}
		ret = append(ret, &v)
		return nil
	})
	return
}
