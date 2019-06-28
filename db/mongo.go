package db

import (
	"errors"
	"log"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var config *MongoConfig
var instances []*mgo.Session
var instanceOnce sync.Once
var instancesIndex uint32

const mgoMaxSessions = 8

type M bson.M

func (m M) Update(qs ...M) M {
	for _, q := range qs {
		for k, v := range q {
			m[k] = v
		}
	}
	return m
}

var ErrInitResource = errors.New("ezorm/db: failed to initial mongo resource")
var ErrOperaBeforeInit = errors.New("please set db.SetOnFinishInit when needed operating db in init()")

// non-multhreads
var afterEvents []func()
var indexEvents []func()

func EnsureAllIndex() {
	for _, f := range indexEvents {
		f()
	}
}

func SetOnEnsureIndex(f func()) {
	indexEvents = append(indexEvents, f)
}

func SetOnFinishInit(f func()) {
	if IsFinishInit() {
		f()
		return
	}
	afterEvents = append(afterEvents, f)
}

func IsFinishInit() bool {
	return instances != nil
}

func Setup(c *MongoConfig) {
	config = c
}

func ShareSession() *mgo.Session {
	doInit := false
	instanceOnce.Do(func() {
		instances = MustNewMgoSessions(config)
	})

	if doInit {
		for _, f := range afterEvents {
			f()
		}
	}

	i := atomic.AddUint32(&instancesIndex, 1)
	i = i % uint32(len(instances))
	return instances[int(i)].Clone()
}

func MustNewMgoSessions(config *MongoConfig) []*mgo.Session {
	sessions := make([]*mgo.Session, mgoMaxSessions)
	for i := 0; i < mgoMaxSessions; i++ {
		if config == nil {
			log.Fatal(ErrOperaBeforeInit)
		}
		session, err := mgo.Dial(config.MongoDB)
		if err != nil {
			log.Println("failed to dial mongo:", err)
			panic(ErrInitResource)
		}
		// Optional. Switch the session to a monotonic behavior.
		session.SetMode(mgo.Monotonic, true)
		if err := session.Ping(); err != nil {
			log.Println("failed to ping mongo:", err)
			panic(ErrInitResource)
		}

		poolLimit := config.PoolLimit
		if poolLimit <= 0 {
			poolLimit = 16
		}

		session.SetPoolLimit(poolLimit)
		sessions[i] = session

		// Refresh session in case of network error.
		go func(s *mgo.Session) {
			for {
				if err := s.Ping(); err != nil {
					s.Refresh()
					s.Ping()
				}
				time.Sleep(time.Second)
			}
		}(session)
	}

	return sessions
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
	return ShareSession()
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
