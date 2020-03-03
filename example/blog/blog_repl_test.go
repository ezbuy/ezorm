package test

import (
	"testing"
	"time"

	"github.com/ezbuy/ezorm/db"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func TestSessionRefresher_ReplicaSet(t *testing.T) {
	t.Skip()
	conf := initMongoReplicaSet()
	session, col := BlogMgr.GetCol()
	defer session.Close()

	db.SetupIdleSessionRefresher(conf, []*mgo.Session{session}, 200*time.Millisecond)

	var i int
	for {
		if i == 20 {
			break
		}
		// no master or salve socket set
		// an empty socket
		cs := session.Clone()
		_, err := cs.DB("ezorm").C(col.Name).Find(bson.M{`$where`: `sleep(10000) || "true"`}).Count()
		assert.NoError(t, err)
		i++
		time.Sleep(200 * time.Millisecond)
		cs.Close()
	}
}
