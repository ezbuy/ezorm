package mysqlr

type SQL interface {
	SQLFormat(limit bool) string
	SQLParams() []interface{}
	SQLLimit() int
	Offset(n int)
	Limit(n int)
}

//! conf.orm
type PrimaryKey interface {
	Key() string
	SQLFormat() string
	SQLParams() []interface{}
	Columns() []string
	Parse(key string) error
}

type Unique interface {
	SQL
	Key() string
}
type Index interface {
	SQL
	Key() string
	PositionOffsetLimit(len int) (int, int)
}

type DBFetcher interface {
	FetchBySQL(sql string, args ...interface{}) ([]interface{}, error)
}
