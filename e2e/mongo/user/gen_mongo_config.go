package user

import (
	"context"
	"fmt"

	"github.com/ezbuy/ezorm/v2/db"
	"github.com/ezbuy/ezorm/v2/pkg/orm"
	"github.com/ezbuy/wrapper/database"

	"go.mongodb.org/mongo-driver/mongo"
)

type SetupOption struct {
	monitor   database.Monitor
	postHooks []func()
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

var mongoDriver *db.MongoDriver

func MgoSetup(config *db.MongoConfig, opts ...SetupOptionFn) {
	sopt := &SetupOption{}
	for _, opt := range opts {
		opt(sopt)
	}
	// setup the indexes
	sopt.postHooks = append(sopt.postHooks, orm.PostSetupHooks...)
	var dopt []db.MongoDriverOption
	if sopt.monitor != nil {
		dopt = append(dopt, db.WithPoolMonitor(database.NewMongoDriverMonitor(sopt.monitor)))
	}
	db.Setup(config)

	var err error
	mongoDriver, err = db.NewMongoDriver(
		context.Background(),
		dopt...,
	)
	if err != nil {
		panic(fmt.Errorf("failed to create mongodb driver: %s", err))
	}
	for _, hook := range sopt.postHooks {
		hook()
	}
}

func Col(col string) *mongo.Collection {
	return mongoDriver.GetCol(col)
}
