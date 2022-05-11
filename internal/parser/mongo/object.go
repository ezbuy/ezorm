package mongo

import (
	"github.com/ezbuy/ezorm/v2/internal/generator"
	"github.com/ezbuy/ezorm/v2/internal/parser/shared"
)

var _ generator.IObject = (*MongoObject)(nil)

type MongoObject struct {
	*shared.Obj
}

func NewMongoObject(pkg string, templateName string) *MongoObject {
	return &MongoObject{
		Obj: &shared.Obj{
			Package:   pkg,
			GoPackage: pkg,
			Name:      templateName,
		},
	}
}
