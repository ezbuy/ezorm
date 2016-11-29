package test

import (
	"bytes"
	"fmt"
	"strconv"
)

func int32ToIds(buf *bytes.Buffer, ids []int32) {
	buf.WriteString("(")
	set := make(map[int32]struct{}, len(ids))
	for idx, id := range ids {
		if _, ok := set[id]; ok {
			continue
		}
		set[id] = struct{}{}
		if idx > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(strconv.Itoa(int(id)))
	}
	buf.WriteString(")")
}

func intToIds(buf *bytes.Buffer, ids []int) {
	buf.WriteString("(")
	set := make(map[int]struct{}, len(ids))
	for idx, id := range ids {
		if _, ok := set[id]; ok {
			continue
		}
		set[id] = struct{}{}
		if idx > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(strconv.Itoa(int(id)))
	}
	buf.WriteString(")")
}

func stringToIds(buf *bytes.Buffer, ids []string) {
	buf.WriteString("(")
	set := make(map[string]struct{}, len(ids))
	for idx, id := range ids {
		if _, ok := set[id]; ok {
			continue
		}
		set[id] = struct{}{}
		if idx > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(strconv.Quote(id))
	}
	buf.WriteString(")")
}

func boolToIds(buf *bytes.Buffer, ids []bool) {
	buf.WriteString("(")
	set := make(map[bool]struct{}, len(ids))
	for idx, id := range ids {
		if _, ok := set[id]; ok {
			continue
		}
		set[id] = struct{}{}
		if idx > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(fmt.Sprint(id))
	}
	buf.WriteString(")")
}
