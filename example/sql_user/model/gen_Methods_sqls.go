package model

import (
	"database/sql"
	"github.com/ezbuy/ezorm/orm"
	"time"
)

var _ *time.Time
var _ *sql.Rows
var _ orm.Execable

type _Methods struct{}

var Methods = _Methods{}

type ListUsersResp struct {
	User
	Email string `table:"user_detail" field:"email"`
}

const _ListUsersSQL = "SELECT `user`.`id`, `user`.`name`, `user`.`phone`, `user`.`password`, `user_detail`.`email` Email FROM `user` JOIN `user_detail` ON `user`.`id`=`user_detail`.`user_id` LIMIT ?, ?"

func (_Methods) ListUsers(db orm.Queryable, offset int, limit int) (ret []*ListUsersResp, err error) {
	err = orm.ExecQuery(db, _ListUsersSQL, []interface{}{offset, limit}, func(rows *sql.Rows) error {
		var e ListUsersResp
		if err := rows.Scan(&e.User.Id, &e.User.Name, &e.User.Phone, &e.User.Password, &e.Email); err != nil {
			return err
		}
		ret = append(ret, &e)
		return nil
	})
	return
}

type ListUserSnippetResp struct {
	UserId       int64  `table:"user" field:"id"`
	UserName     string `table:"user" field:"name"`
	UserPhone    string `table:"user" field:"phone"`
	UserPassword string `table:"user" field:"password"`
}

const _ListUserSnippetSQL = "SELECT `user`.`id`, `user`.`name`, `user`.`phone`, `user`.`password` FROM `user` LIMIT ?, ?"

func (_Methods) ListUserSnippet(db orm.Queryable, offset int, limit int) (ret []*ListUserSnippetResp, err error) {
	err = orm.ExecQuery(db, _ListUserSnippetSQL, []interface{}{offset, limit}, func(rows *sql.Rows) error {
		var e ListUserSnippetResp
		if err := rows.Scan(&e.UserId, &e.UserName, &e.UserPhone, &e.UserPassword); err != nil {
			return err
		}
		ret = append(ret, &e)
		return nil
	})
	return
}

type ListUserIfNullResp struct {
	UserName     string `table:"user" field:"name"`
	UserPhone    string `table:"user" field:"phone"`
	UserPassword string `table:"user" field:"password"`
	Email        string `table:"user_detail" field:"email"`
	Text         string `table:"user_detail" field:"text"`
}

const _ListUserIfNullSQL = "SELECT `user`.`name`, `user`.`phone`, `user`.`password`, IFNULL(`user_detail`.`email`, '') Email, IFNULL(`user_detail`.`text`, '') Text FROM `user` LEFT JOIN `user_detail` ON `user`.`id`=`user_detail`.`user_id`"

func (_Methods) ListUserIfNull(db orm.Queryable) (ret []*ListUserIfNullResp, err error) {
	err = orm.ExecQuery(db, _ListUserIfNullSQL, []interface{}{}, func(rows *sql.Rows) error {
		var e ListUserIfNullResp
		if err := rows.Scan(&e.UserName, &e.UserPhone, &e.UserPassword, &e.Email, &e.Text); err != nil {
			return err
		}
		ret = append(ret, &e)
		return nil
	})
	return
}

type GetUsersByRolesResp struct {
	User
	Role
}

const _GetUsersByRolesSQL = "SELECT u.`id`, u.`name`, u.`phone`, u.`password`, r.`id`, r.`name` FROM `user` AS u JOIN `user_role` AS ur ON u.`id`=ur.`user_id` JOIN `role` AS r ON r.`id`=ur.`role_id` WHERE r.`id`=?"

func (_Methods) GetUsersByRoles(db orm.Queryable, roleId int64) (ret []*GetUsersByRolesResp, err error) {
	err = orm.ExecQuery(db, _GetUsersByRolesSQL, []interface{}{roleId}, func(rows *sql.Rows) error {
		var e GetUsersByRolesResp
		if err := rows.Scan(&e.User.Id, &e.User.Name, &e.User.Phone, &e.User.Password, &e.Role.Id, &e.Role.Name); err != nil {
			return err
		}
		ret = append(ret, &e)
		return nil
	})
	return
}

type GetUserByIdResp struct {
	UserId          int64  `table:"user" field:"id"`
	UserName        string `table:"user" field:"name"`
	UserPhone       string `table:"user" field:"phone"`
	UserPassword    string `table:"user" field:"password"`
	UserDetailEmail string `table:"user_detail" field:"email"`
	UserDetailText  string `table:"user_detail" field:"text"`
}

const _GetUserByIdSQL = "SELECT u.`id`, u.`name`, u.`phone`, u.`password`, ud.`email`, ud.`text` FROM `user` u JOIN `user_detail` ud ON u.`id`=ud.`user_id` WHERE u.`id`=?"

func (_Methods) GetUserById(db orm.Queryable, id int64) (ret []*GetUserByIdResp, err error) {
	err = orm.ExecQuery(db, _GetUserByIdSQL, []interface{}{id}, func(rows *sql.Rows) error {
		var e GetUserByIdResp
		if err := rows.Scan(&e.UserId, &e.UserName, &e.UserPhone, &e.UserPassword, &e.UserDetailEmail, &e.UserDetailText); err != nil {
			return err
		}
		ret = append(ret, &e)
		return nil
	})
	return
}

type CountUsersResp struct {
	Count0 int64 `table:"count"`
}

const _CountUsersSQL = "SELECT COUNT(1) FROM `user`"

func (_Methods) CountUsers(db orm.Queryable) (ret []*CountUsersResp, err error) {
	err = orm.ExecQuery(db, _CountUsersSQL, []interface{}{}, func(rows *sql.Rows) error {
		var e CountUsersResp
		if err := rows.Scan(&e.Count0); err != nil {
			return err
		}
		ret = append(ret, &e)
		return nil
	})
	return
}
