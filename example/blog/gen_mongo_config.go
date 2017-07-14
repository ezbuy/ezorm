package test

import (
	"errors"
	"sync"
	"sync/atomic"

	"github.com/ezbuy/ezorm/db"

	"gopkg.in/mgo.v2"
)

var (
	mgoInstances     []*mgo.Session
	mgoConfig        *db.MongoConfig
	mgoInstanceOnce  sync.Once
	mgoInstanceIndex uint32
)

var ErrOperaBeforeInit = errors.New("please set db.SetOnFinishInit when needed operating db in init()")

const mgoMaxSessions = 8

func MgoSetup(config *db.MongoConfig) {
	mgoConfig = config
	mgoInstances = make([]*mgo.Session, mgoMaxSessions)
	for i := 0; i < mgoMaxSessions; i++ {
		if mgoConfig == nil {
			panic(ErrOperaBeforeInit)
		}
		session, err := mgo.Dial(mgoConfig.MongoDB)
		if err != nil {
			panic(err)
		}
		session.SetMode(mgo.Monotonic, true)
		if err := session.Ping(); err != nil {
			panic(err)
		}
		poolLimit := config.PoolLimit
		if poolLimit <= 0 {
			poolLimit = 16
		}

		session.SetPoolLimit(poolLimit)
		mgoInstances[i] = session
	}
}

func getCol(dbName, col string) (*mgo.Session, *mgo.Collection) {
	i := atomic.AddUint32(&mgoInstanceIndex, 1)
	i = i % uint32(len(mgoInstances))
	session := mgoInstances[int(i)].Clone()
	var mgoCol *mgo.Collection
	if dbName == "" {
		mgoCol = session.DB(mgoConfig.DBName).C(col)
	} else {
		mgoCol = session.DB(dbName).C(col)
	}
	return session, mgoCol
}
