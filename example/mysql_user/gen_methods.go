package test

import (
	"time"

	"github.com/ezbuy/ezorm/db"
)

var (
	_ time.Time
)

type sqlMethods struct{}

var SQL = &sqlMethods{}

type FindUserByPhoneResp struct {
	UserUserId     int64
	UserName       string
	UserPhone      string
	UserAge        int32
	UserBalance    float64
	UserText       string
	UserCreateDate int64
}

const _FindUserByPhoneSQL = "SELECT /* find_user_by_phone */ user_id, name, phone, age, balance, text, create_date FROM user WHERE phone=?"

func (*sqlMethods) FindUserByPhone(args ...interface{}) ([]*FindUserByPhoneResp, error) {
	rows, err := db.MysqlQuery(_FindUserByPhoneSQL, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*FindUserByPhoneResp
	for rows.Next() {
		var o FindUserByPhoneResp
		err = rows.Scan(&o.UserUserId, &o.UserName, &o.UserPhone, &o.UserAge, &o.UserBalance, &o.UserText, &o.UserCreateDate)
		if err != nil {
			return nil, err
		}
		results = append(results, &o)
	}
	return results, nil
}

type SearchUserResp struct {
	UserUserId int64
	UserName   string
	UserPhone  string
	UserAge    int32
}

const _SearchUserSQL = "SELECT /* search_user */ user_id, name, phone, age FROM user WHERE name LIKE ? LIMIT ?, ?"

func (*sqlMethods) SearchUser(args ...interface{}) ([]*SearchUserResp, error) {
	rows, err := db.MysqlQuery(_SearchUserSQL, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*SearchUserResp
	for rows.Next() {
		var o SearchUserResp
		err = rows.Scan(&o.UserUserId, &o.UserName, &o.UserPhone, &o.UserAge)
		if err != nil {
			return nil, err
		}
		results = append(results, &o)
	}
	return results, nil
}
