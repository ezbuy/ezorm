package mysql

import (
	"github.com/ezbuy/ezorm/v2/internal/generator"
	"github.com/ezbuy/ezorm/v2/internal/parser/shared"
)

var _ generator.IObject = (*MySQLObject)(nil)

type MySQLObject struct {
	*shared.Obj
}

func NewMySQLObject(pkg string, templateName string) *MySQLObject {
	return &MySQLObject{
		Obj: &shared.Obj{
			Namespace: pkg,
			GoPackage: pkg,
			Name:      templateName,
		},
	}
}
