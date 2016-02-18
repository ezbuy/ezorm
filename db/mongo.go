package db

import (
	"errors"

	"strings"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var config *MongoConfig
var ShareSession *mgo.Session

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
	return ShareSession != nil
}

func Setup(c *MongoConfig) {
	config = c
	session, err := mgo.Dial(config.MongoDB)
	if err != nil {
		panic(err)
	}
	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	ShareSession = session
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

func NewSession() (session *mgo.Session) {
	if ShareSession == nil {
		panic(ErrOperaBeforeInit)
	}
	return ShareSession.Copy()
}

func NewCollection(session *mgo.Session, name string) *mgo.Collection {
	return session.DB(config.DBName).C(name)
}

func GetCol(col string) (session *mgo.Session, collection *mgo.Collection) {
	session = NewSession()
	collection = NewCollection(session, col)
	return
}

func NewObjectId() bson.ObjectId {
	return bson.NewObjectId()
}
