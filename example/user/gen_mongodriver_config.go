package user

import (
	"context"
	"fmt"

	"github.com/ezbuy/ezorm/db"
	"go.mongodb.org/mongo-driver/mongo"
)

var mongoDriver *db.MongoDriver

func MgoSetup(config *db.MongoConfig) {
	db.Setup(config)

	var err error
	mongoDriver, err = db.NewMongoDriver(context.Background())
	if err != nil {
		panic(fmt.Errorf("failed to create mongodb driver: %s", err))
	}
}

func Col(col string) *mongo.Collection {
	return mongoDriver.GetCol(col)
}
