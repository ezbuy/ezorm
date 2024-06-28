package db_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/ezbuy/ezorm/v2/pkg/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/event"
)

func Test_MongoDriverConnection(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancelFn()

	db.Setup(&db.MongoConfig{
		MongoDB:   "mongodb://127.0.0.1:27017",
		DBName:    "test",
		PoolLimit: 30,
	})

	md, err := db.NewMongoDriver(ctx)
	if err != nil {
		t.Fatalf("failed to new mongo driver: %s", err)
	}
	defer md.Close()

	const (
		collectionName = "Test"
	)

	col := md.GetCol(collectionName)
	ret, err := col.InsertOne(ctx, bson.M{
		"tid":         1,
		"ezorm":       "mongo_driver_support",
		"create_date": time.Now().Unix(),
	})
	if err != nil {
		t.Fatalf("failed to insert to collection: %s", err)
	}
	t.Logf("got insert id: %s", ret.InsertedID.(primitive.ObjectID).String())
}

func Test_MongoDriverConnectionBy(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancelFn()

	db.SetupMany(&db.MongoConfig{
		MongoDB:   "mongodb://127.0.0.1:27017",
		DBName:    "test",
		PoolLimit: 30,
	}, &db.MongoConfig{
		MongoDB:   "mongodb://127.0.0.1:27017",
		DBName:    "test_2",
		PoolLimit: 30,
	})

	md, err := db.NewMongoDriverBy(ctx, db.WithDBName("test"))
	if err != nil {
		t.Fatalf("failed to new mongo driver: %s", err)
	}
	defer md.Close()

	const (
		collectionName = "Test"
	)

	col := md.GetCol(collectionName)
	ret, err := col.InsertOne(ctx, bson.M{
		"tid":         1,
		"ezorm":       "mongo_driver_support",
		"create_date": time.Now().Unix(),
	})
	if err != nil {
		t.Fatalf("failed to insert to collection: %s", err)
	}
	t.Logf("got insert id: %s", ret.InsertedID.(primitive.ObjectID).String())

	md2, err := db.NewMongoDriverBy(ctx, db.WithDBName("test_2"))
	if err != nil {
		t.Fatalf("failed to new mongo driver: %s", err)
	}
	defer md2.Close()

	const (
		collectionName2 = "Test"
	)

	col2 := md.GetCol(collectionName2)
	ret2, err := col2.InsertOne(ctx, bson.M{
		"tid":         1,
		"ezorm":       "mongo_driver_support",
		"create_date": time.Now().Unix(),
	})
	if err != nil {
		t.Fatalf("failed to insert to collection: %s", err)
	}
	t.Logf("got insert id: %s", ret2.InsertedID.(primitive.ObjectID).String())
}

func Test_MongoDriverConnPool(t *testing.T) {
	connIds := make(map[uint64]int)
	connIdsLock := new(sync.Mutex)
	monitor := &event.PoolMonitor{
		Event: func(event *event.PoolEvent) {
			connIdsLock.Lock()
			if v, e := connIds[event.ConnectionID]; !e {
				connIds[event.ConnectionID] = 1
			} else {
				v++
				connIds[event.ConnectionID] = v
			}
			connIdsLock.Unlock()
		},
	}

	db.Setup(&db.MongoConfig{
		MongoDB:   "mongodb://127.0.0.1:27017",
		DBName:    "test",
		PoolLimit: 30,
	})

	md, err := db.NewMongoDriver(context.Background(), db.WithPoolMonitor(monitor))
	if err != nil {
		t.Fatalf("failed to new mongo driver: %s", err)
	}
	defer md.Close()

	const collectionName = "Test"
	if _, err := md.GetCol(collectionName).
		InsertOne(context.Background(), bson.M{"tid": 2}); err != nil {
		t.Fatalf("failed to initial query data: %s", err)
	}

	gctx, wg := context.Background(), new(sync.WaitGroup)
	for i := 0; i < 1<<8; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			md.GetCol(collectionName).FindOne(gctx, bson.M{"tid": 2})
		}()
	}

	wg.Wait()
	t.Fatalf("got connection ids: %+v", connIds)
}
