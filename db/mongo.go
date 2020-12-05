package db

import (
	"errors"
	"log"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ezbuy/statsd"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var config *MongoConfig
var instances []*mgo.Session
var instanceOnce sync.Once
var instancesIndex uint32

const mgoMaxSessions = 8

var monitorInterval = 10 * time.Second

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
		SetupIdleSessionRefresher(config, instances, 3*time.Minute)
	})

	if doInit {
		for _, f := range afterEvents {
			f()
		}
	}

	i := atomic.AddUint32(&instancesIndex, 1)
	i = i % uint32(len(instances))
	// use Clone here to avoid socket refresh
	return instances[int(i)].Clone()
}

// SetupIdleSessionRefresher will always refresh idle sessions if it is possible
func SetupIdleSessionRefresher(config *MongoConfig, instances []*mgo.Session, every time.Duration) {
	go func() {
		var cursor uint64
		// NOTE: `instances` itself is read-only
		instanceLen := uint64(len(instances))
		for {
			time.Sleep(every)
			// Acquire one idle session ,and do refresh(re-dail)
			// move the session instance cursor to next
			idleSession := instances[int(cursor%instanceLen)]
			// reset session socket , and the socket will be re-allocated in next mongo operation
			idleSession.Refresh()
			// Ping the server after refresh sockets, it will pre-allocate the socket
			// if ping fails , the socket will not be pre-allocated
			if err := idleSession.Ping(); err != nil {
				log.Printf("sessionRefreher: %q", err)
			}
			// move to next instance
			cursor++
		}
	}()
}

func MustNewMgoSessions(config *MongoConfig) []*mgo.Session {
	maxSession := config.MaxSession
	if maxSession == 0 {
		maxSession = mgoMaxSessions
	}
	sessions := make([]*mgo.Session, maxSession)
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

// MustSetupMgoMonitor set up mongo connection pool monitor metrics
func MustSetupMgoMonitor(srv string) {
	if srv == "" {
		panic("ezorm: srv name must set at app layer")
	}
	mgo.SetStats(true)
	go func() {
		for {
			monitorMongoStats(srv)
			time.Sleep(monitorInterval)
		}
	}()
}

func monitorMongoStats(srv string) {
	st := mgo.GetStats()
	statsd.Gauge("infra.db.mongo."+srv+".clusterNode", int64(st.Clusters))
	statsd.Gauge("infra.db.mongo."+srv+"masterConn", int64(st.MasterConns))
	statsd.Gauge("infra.db.mongo."+srv+"slaveConn", int64(st.SlaveConns))
	statsd.Gauge("infra.db.mongo."+srv+"socketRefs", int64(st.SocketRefs))
	statsd.Gauge("infra.db.mongo."+srv+"socketAlive", int64(st.SocketsAlive))
	statsd.Gauge("infra.db.mongo."+srv+"socketInUse", int64(st.SocketsInUse))
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
