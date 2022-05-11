package mysql

import (
	"github.com/ezbuy/ezorm/v2/internal/generator"
	"github.com/ezbuy/ezorm/v2/internal/parser/shared"
)

var _ generator.Generator = (*MySQLGenerator)(nil)

type MySQLGenerator struct {
	*shared.Generator
}

func (mg *MySQLGenerator) DriverName() string {
	return "mysql"
}
