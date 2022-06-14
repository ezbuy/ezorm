package sql

import (
	"fmt"
	"strings"
)

type InBuilder struct {
	params int
}

func NewIn(params int) *InBuilder {
	return &InBuilder{
		params: params,
	}
}

func (in *InBuilder) String() string {
	var placeholders []string
	for i := 0; i < in.params; i++ {
		placeholders = append(placeholders, "?")
	}
	var query string
	if len(placeholders) > 0 {
		query = strings.Join(placeholders, ",")
	}
	return fmt.Sprintf("(%s)", query)
}
