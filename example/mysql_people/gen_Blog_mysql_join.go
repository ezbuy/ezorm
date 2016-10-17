package people

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

func (m *_BlogMgr) ToUser(base []*Blog) []int32 {
	ids := make([]int32, len(base))
	for idx, b := range base {
		ids[idx] = b.User
	}
	return ids
}

func (m *_BlogMgr) ToUserQuery(base []*Blog) string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("`user` IN (")
	ids := m.ToUserQuery(base)
	for idx, n := range ids {
		buf.WriteString(fmt.Sprint(n))
		if idx == len(ids)-1 {
			break
		}
		buf.WriteString(",")
	}
	buf.WriteString(")")
	return buf.String()
}

func (m *_BlogMgr) LeftJoinUser(user []int32) ([]*User, error) {
	base := user
	targets, err := UserMgr.FindInUserId(base)
	if err != nil {
		return nil, err
	}
	refMap := make(map[int32]*User, len(targets))
	for _, t := range targets {
		refMap[t.UserId] = t
	}

	ret := make([]*User, len(base))
	for idx := range base {
		ret[idx] = refMap[base[idx]]
	}
	return ret, nil
}
