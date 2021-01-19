package db

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func NewMongoDriver(ctx context.Context) (*MongoDriver, error) {
	if config == nil {
		return nil, errors.New("db: initialize config before new mongo driver")
	}
	cli, err := mongo.NewClient(options.Client().ApplyURI(config.MongoDB).
		SetMaxPoolSize(uint64(config.PoolLimit)))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongodb server: %w", err)
	}

	if err = cli.Connect(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to mongodb server: %w", err)
	}

	if err = cli.Ping(ctx, readpref.PrimaryPreferred()); err != nil {
		return nil, fmt.Errorf("failed to ping remote mongodb server: %w", err)
	}

	return &MongoDriver{cli: cli}, nil
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
