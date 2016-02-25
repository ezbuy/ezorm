package parser

func minus(a, b int) int {
	return a - b
}

var NullTypes = map[string]string{
	"string":    "String",
	"bool":      "Bool",
	"int":       "Int64",
	"int32":     "Int64",
	"int64":     "Int64",
	"bit":       "Bool",
	"time.Time": "String",
}

func getNullType(gotype string) string {
	return NullTypes[gotype]
}
