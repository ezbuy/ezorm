package db

import "sync"

var config *MongoConfig

type MongoConfig struct {
	MongoDB    string
	DBName     string
	PoolLimit  int
	MaxSession int
}

// Setup setup one mongo config
// For multiple configs, use `SetupMany` instead
func Setup(c *MongoConfig) {
	config = c
}

var multiConfigs sync.Map

// SetupMany setup many mongo configs
// For singleton , use `Setup` instead
func SetupMany(
	cs ...*MongoConfig,
) {
	if config == nil && len(cs) == 1 {
		Setup(cs[0])
		multiConfigs.Store(cs[0].DBName, cs[0])
		return
	}
	for _, c := range cs {
		multiConfigs.Store(c.DBName, c)
	}
}

func GetConfigByName(name string) *MongoConfig {
	if v, ok := multiConfigs.Load(name); ok {
		return v.(*MongoConfig)
	}
	return nil
}

func GetConfig(name string) *MongoConfig {
	if config != nil {
		return config
	}
	return GetConfigByName(name)
}
