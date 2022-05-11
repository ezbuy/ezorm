package shared

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"
)

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

func camel2list(s []string) []string {
	s2 := make([]string, len(s))
	for idx := range s {
		s2[idx] = camel2name(s[idx])
	}
	return s2
}

func toIds(bufName, typeName, name string) string {
	switch typeName {
	case "int":
		return "intToIds(" + bufName + "," + name + ")"
	case "int32":
		return "int32ToIds(" + bufName + "," + name + ")"
	case "bool":
		return "boolToIds(" + bufName + "," + name + ")"
	case "string":
		return "stringToIds(" + bufName + "," + name + ")"
	case "int64":
		return "int64ToIds(" + bufName + "," + name + ")"
	}
	return name
}

func strif(a bool, b, c string) string {
	if a {
		return b
	}
	return c
}

func strDefault(a, b string) string {
	if a == "" {
		return b
	}
	return a
}

func camel2name(s string) string {
	nameBuf := bytes.NewBuffer(nil)
	afterSpace := false
	for i, c := range s {
		if unicode.IsUpper(c) && unicode.IsLetter(c) {
			if i > 0 && !afterSpace {
				nameBuf.WriteRune('_')
			}
			c = unicode.ToLower(c)
		}
		nameBuf.WriteRune(c)
		afterSpace = unicode.IsSpace(c)
	}
	return nameBuf.String()
}

func getNullType(gotype string) string {
	return NullTypes[gotype]
}

func preSuffixJoin(s []string, prefix, suffix, sep string) string {
	sNew := make([]string, 0, len(s))
	for _, each := range s {
		sNew = append(sNew, fmt.Sprintf("%s%s%s", prefix, each, suffix))
	}
	return strings.Join(sNew, sep)
}

func repeatJoin(n int, repeatStr, sep string) string {
	a := make([]string, 0, n)
	for i := 0; i < n; i++ {
		a = append(a, repeatStr)
	}
	return strings.Join(a, sep)
}
