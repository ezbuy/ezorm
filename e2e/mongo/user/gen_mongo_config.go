package user

import (
	"context"
	"fmt"
	"sync"

	"github.com/ezbuy/ezorm/v2/pkg/db"
	"github.com/ezbuy/wrapper/database"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

type SetupOption struct {
	monitor   database.Monitor
	postHooks []func()
	// mockStub is for mocking purpose
	mockStub *mtest.T
}

type SetupOptionFn func(opts *SetupOption)

func WithStatsDMonitor(app string) SetupOptionFn {
	return func(opts *SetupOption) {
		opts.monitor = database.NewStatsDPoolMonitor(app)
	}
}

func WithPrometheusMonitor(app, gatewayAddress string) SetupOptionFn {
	return func(opts *SetupOption) {
		opts.monitor = database.NewPrometheusPoolMonitor(app, gatewayAddress)
	}
}

func WithPostHooks(fn ...func()) SetupOptionFn {
	return func(opts *SetupOption) {
		opts.postHooks = append(opts.postHooks, fn...)
	}
}

func WithMockStub(t *mtest.T) SetupOptionFn {
	return func(opts *SetupOption) {
		opts.mockStub = t
	}
}

var mongoDriver *db.MongoDriver
var mongoDriverOnce sync.Once

func MgoSetup(config *db.MongoConfig, opts ...SetupOptionFn) {
	sopt := &SetupOption{}
	for _, opt := range opts {
		opt(sopt)
	}
	if sopt.mockStub != nil {
		// reset the mongo driver if it was already initialized
		mongoDriverOnce = sync.Once{}
		mongoDriver = nil
		mongoDriver = db.NewMockMongoDriver(sopt.mockStub)
		for _, hook := range sopt.postHooks {
			hook()
		}
		return
	}
	// setup the indexes
	sopt.postHooks = append(sopt.postHooks,
		UserIndexesFunc,
		UserBlogIndexesFunc,
		UserNullableIndexesFunc,
	)
	var dopt []db.MongoDriverConnOptionFn
	if sopt.monitor != nil {
		clientOpt := db.WithClientOption(db.WithPoolMonitor(database.NewMongoDriverMonitor(sopt.monitor)))
		dopt = append(dopt, clientOpt)
	}
	if config.DBName == "" {
		panic("db name is required")
	}
	db.SetupMany(config)
	dopt = append(dopt, db.WithDBName(config.DBName))

	mongoDriverOnce.Do(func() {
		var err error
		mongoDriver, err = db.NewMongoDriverBy(
			context.Background(),
			dopt...,
		)
		if err != nil {
			panic(fmt.Errorf("failed to create mongodb driver: %s", err))
		}
		for _, hook := range sopt.postHooks {
			hook()
		}
	})
}

func Col(col string) *mongo.Collection {
	return mongoDriver.GetCol(col)
}
