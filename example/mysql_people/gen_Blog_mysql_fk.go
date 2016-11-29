package test

import (
	"bytes"
	"fmt"
	"strings"
)

var (
	_ = fmt.Println
	_ = strings.Index
	_ bytes.Buffer
)

func (m *_BlogMgr) ToFieldUser(base []*Blog) []int32 {
	ids := make([]int32, len(base))
	for idx, b := range base {
		ids[idx] = b.User
	}
	return ids
}
