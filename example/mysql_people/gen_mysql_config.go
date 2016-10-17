package people

import (
	"bytes"
	"strconv"
)

func int32ToIds(buf *bytes.Buffer, ids []int32) {
	buf.WriteString("(")
	for idx, id := range ids {
		buf.WriteString(strconv.Itoa(int(id)))
		if idx == len(ids)-1 {
			continue
		}
		buf.WriteString(",")
	}
	buf.WriteString(")")
}

func intToIds(buf *bytes.Buffer, ids []int) {
	buf.WriteString("(")
	for idx, id := range ids {
		buf.WriteString(strconv.Itoa(int(id)))
		if idx == len(ids)-1 {
			continue
		}
		buf.WriteString(",")
	}
	buf.WriteString(")")
}

func stringToIds(buf *bytes.Buffer, ids []string) {
	buf.WriteString("(")
	for idx, id := range ids {
		buf.WriteString(strconv.Quote(id))
		if idx == len(ids)-1 {
			continue
		}
		buf.WriteString(",")
	}
	buf.WriteString(")")
}