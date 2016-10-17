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

func (m *_UserMgr) ToUserId(base []*User) []int32 {
	ids := make([]int32, len(base))
	for idx, b := range base {
		ids[idx] = b.UserId
	}
	return ids
}

func (m *_UserMgr) ToUserIdQuery(base []*User) string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("`user_id` IN (")
	ids := m.ToUserIdQuery(base)
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

func (m *_UserMgr) LeftJoinBlog(userId []int32) ([]*Blog, error) {
	base := userId
	targets, err := BlogMgr.FindInUser(base)
	if err != nil {
		return nil, err
	}
	refMap := make(map[int32]*Blog, len(targets))
	for _, t := range targets {
		refMap[t.User] = t
	}

	ret := make([]*Blog, len(base))
	for idx := range base {
		ret[idx] = refMap[base[idx]]
	}
	return ret, nil
}
