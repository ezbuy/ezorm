package db

import (
	"errors"
	"sync"

	"strings"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var config *MongoConfig
var instance *mgo.Session
var instanceOnce sync.Once

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
	return instance != nil
}

func Setup(c *MongoConfig) {
	config = c
}

func ShareSession() *mgo.Session {
	doInit := false
	instanceOnce.Do(func() {
		if instance == nil {
			if config == nil {
				panic(ErrOperaBeforeInit)
			}
			session, err := mgo.Dial(config.MongoDB)
			if err != nil {
				panic(err)
			}
			// Optional. Switch the session to a monotonic behavior.
			session.SetMode(mgo.Monotonic, true)
			instance = session
			doInit = true
		}
	})

	if doInit {
		for _, f := range afterEvents {
			f()
		}
	}
	return instance
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
	return ShareSession().Copy()
}

func NewCollection(session *mgo.Session, dbName, name string) *mgo.Collection {
	if dbName == "" {
		return session.DB(config.DBName).C(name)
	}
	return session.DB(dbName).C(name)
}

func GetCol(dbName, col string) (session *mgo.Session, collection *mgo.Collection) {
	session = NewSession()
	collection = NewCollection(session, dbName, col)
	return
}

func NewObjectId() bson.ObjectId {
	return bson.NewObjectId()
}

func IsMgoNotFound(err error) bool {
	return err == mgo.ErrNotFound
}

func IsMgoDup(err error) bool {
	return mgo.IsDup(err)
}
