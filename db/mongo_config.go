package db

var config *MongoConfig

type MongoConfig struct {
	MongoDB    string
	DBName     string
	PoolLimit  int
	MaxSession int
}

func Setup(c *MongoConfig) {
	config = c
}
