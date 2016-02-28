package mysql

import (
	"errors"
	"strings"
	"database/sql"
	_ "gopkg.in/go-sql-driver/mysql.v1"

	"gopkg.in/mgo.v2/bson"
)

var config *MySQLConfig
var Db *sql.DB

type M bson.M

func (m M) Update(qs ...M) M {
	for _, q := range qs {
		for k, v := range q {
			m[k] = v
		}
	}
	return m
}

var ErrOperaBeforeInit = errors.New("please set db.SetOnFinishInit when needed operating db in init()")

// non-multhreads
var afterEvents []func()

func SetOnFinishInit(f func()) {
	if IsFinishInit() {
		f()
		return
	}
	afterEvents = append(afterEvents, f)
}

func IsFinishInit() bool {
	return Db != nil
}

func Setup(c *MySQLConfig) {
	config = c
	db, err := sql.Open("mysql", config.MySQLDB)

    if err != nil {
        panic(err.Error())
    }
    // defer db.Close()

	Db = db
	for _, f := range afterEvents {
		f()
	}
}

func InID(ids []string) (ret M) {
	return M{"_id": M{"$in": ObjectIds(ids)}}
}

func In(ids []string) M {
	return M{"$in": ObjectIds(ids)}
}

func ObjectIds(ids []string) (ret []bson.ObjectId) {
	ret = make([]bson.ObjectId, 0, len(ids))
	for _, id := range ids {
		id = strings.TrimSpace(id)
		if bson.IsObjectIdHex(id) {
			ret = append(ret, bson.ObjectIdHex(id))
		}
	}
	return
}

func GetCol() (db *sql.DB) {
	return Db
}

func NewObjectId() bson.ObjectId {
	return bson.NewObjectId()
}
