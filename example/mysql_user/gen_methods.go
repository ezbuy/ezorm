package test

import (
	"context"
	"time"

	"github.com/ezbuy/ezorm/db"
)

var (
	_ time.Time
	_ context.Context
)

type sqlMethods struct{}

var SQL = &sqlMethods{}

type CountUsersByNameResp struct {
	Count0 int64
}

const _CountUsersByNameSQL = "SELECT /* count_users_by_name */ COUNT(1) FROM user WHERE name LIKE ?"

func (*sqlMethods) CountUsersByName(ctx context.Context, args ...interface{}) ([]*CountUsersByNameResp, error) {
	rows, err := db.MysqlQuery(_CountUsersByNameSQL, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*CountUsersByNameResp
	for rows.Next() {
		var o CountUsersByNameResp
		err = rows.Scan(&o.Count0)
		if err != nil {
			return nil, err
		}
		results = append(results, &o)
	}
	return results, nil
}

type CountUsersByPhoneResp struct {
	UserCount int64
}

const _CountUsersByPhoneSQL = "SELECT /* count_users_by_phone */ COUNT(0) user_count FROM user WHERE phone LIKE ?"

func (*sqlMethods) CountUsersByPhone(ctx context.Context, args ...interface{}) ([]*CountUsersByPhoneResp, error) {
	rows, err := db.MysqlQuery(_CountUsersByPhoneSQL, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*CountUsersByPhoneResp
	for rows.Next() {
		var o CountUsersByPhoneResp
		err = rows.Scan(&o.UserCount)
		if err != nil {
			return nil, err
		}
		results = append(results, &o)
	}
	return results, nil
}

type FindUsersByNameResp struct {
	Id                int64
	UserName          string
	UserPhone         string
	UserAge           int32
	UserBalance       float64
	UserText          string
	UserCreateDate    int64
	UserDetailScore   int32
	UserDetailBalance int32
	DetailText        string
}

const _FindUsersByNameSQL = "SELECT /* find_users_by_name */ u.user_id id, u.name, u.phone, u.age, u.balance, u.text, u.create_date, IFNULL(ud.score, 0), IFNULL(ud.balance, 0), IFNULL(ud.text, '') detail_text FROM user u JOIN user_detail ud ON u.user_id=ud.user_id WHERE u.name LIKE ? LIMIT ?, ?"

func (*sqlMethods) FindUsersByName(ctx context.Context, args ...interface{}) ([]*FindUsersByNameResp, error) {
	rows, err := db.MysqlQuery(_FindUsersByNameSQL, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*FindUsersByNameResp
	for rows.Next() {
		var o FindUsersByNameResp
		err = rows.Scan(&o.Id, &o.UserName, &o.UserPhone, &o.UserAge, &o.UserBalance, &o.UserText, &o.UserCreateDate, &o.UserDetailScore, &o.UserDetailBalance, &o.DetailText)
		if err != nil {
			return nil, err
		}
		results = append(results, &o)
	}
	return results, nil
}

type GetUserByPhoneResp struct {
	Id                int64
	UserName          string
	UserPhone         string
	UserAge           int32
	UserBalance       float64
	UserText          string
	UserCreateDate    int64
	UserDetailScore   int32
	UserDetailBalance int32
	DetailText        string
}

const _GetUserByPhoneSQL = "SELECT /* get_user_by_phone */ u.user_id id, u.name, u.phone, u.age, u.balance, u.text, u.create_date, IFNULL(ud.score, 0), IFNULL(ud.balance, 0), IFNULL(ud.text, '') detail_text FROM user u JOIN user_detail ud ON u.user_id=ud.user_id WHERE u.phone LIKE ?"

func (*sqlMethods) GetUserByPhone(ctx context.Context, args ...interface{}) ([]*GetUserByPhoneResp, error) {
	rows, err := db.MysqlQuery(_GetUserByPhoneSQL, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*GetUserByPhoneResp
	for rows.Next() {
		var o GetUserByPhoneResp
		err = rows.Scan(&o.Id, &o.UserName, &o.UserPhone, &o.UserAge, &o.UserBalance, &o.UserText, &o.UserCreateDate, &o.UserDetailScore, &o.UserDetailBalance, &o.DetailText)
		if err != nil {
			return nil, err
		}
		results = append(results, &o)
	}
	return results, nil
}
