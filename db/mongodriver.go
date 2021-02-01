package db

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoDriver(ctx context.Context, opts ...MongoDriverOption) (*MongoDriver, error) {
	if config == nil {
		return nil, errors.New("db: initialize config before new mongo driver")
	}

	uri := config.MongoDB
	if !strings.HasPrefix(uri, "mongodb") {
		uri = "mongodb://" + config.MongoDB
	}
	fmt.Println("uri:", uri)

	cliOpts := options.Client().ApplyURI(uri).SetMaxPoolSize(uint64(config.PoolLimit))
	for _, opt := range opts {
		opt(cliOpts)
	}

	cli, err := mongo.NewClient(cliOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to create mongodb client: %w", err)
	}

	if err = cli.Connect(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to mongodb server: %w", err)
	}

	if err = cli.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping remote mongodb server: %w", err)
	}

	return &MongoDriver{cli: cli}, nil
}

type MongoDriverOption func(*options.ClientOptions)

func WithPoolMonitor(m *event.PoolMonitor) MongoDriverOption {
	return func(opt *options.ClientOptions) {
		opt.SetPoolMonitor(m)
	}
}

type MongoDriver struct {
	cli *mongo.Client
}

func (md *MongoDriver) GetCol(cname string) *mongo.Collection {
	return md.cli.Database(config.DBName).Collection(cname)
}

func (md *MongoDriver) Close() error {
	if err := md.cli.Disconnect(context.Background()); err != nil {
		return fmt.Errorf("failed to disconncet to mongodb server: %w", err)
	}
	return nil
}
