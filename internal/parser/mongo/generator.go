package mongo

import (
	"github.com/ezbuy/ezorm/v2/internal/generator"
	"github.com/ezbuy/ezorm/v2/internal/parser/shared"
)

var _ generator.Generator = (*MongoGenerator)(nil)

type MongoGenerator struct {
	*shared.Generator
}

func (mg *MongoGenerator) DriverName() string {
	return "mongo"
}
