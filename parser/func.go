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
	"float":     "Float64",
	"float32":   "Float64",
	"float64":   "Float64",
}

func getNullType(gotype string) string {
	return NullTypes[gotype]
}
