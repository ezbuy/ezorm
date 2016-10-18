package people

import (
	"bytes"
	"fmt"
	"github.com/ezbuy/ezorm/db"
)

var (
	_ db.M
	_ = fmt.Println
	_ bytes.Buffer
)

type BlogGroupUnAssigned struct {
	COUNT       []int
	Hits        []int32
	IsPublished []bool
}

func (m *_BlogMgr) GroupByUnAssigned(user []int32, offset, limit int, sorts ...string) (results *BlogGroupUnAssigned, err error) {
	from := user
	buf := bytes.NewBuffer(nil)
	buf.WriteString("SELECT COUNT(*),hits,is_published")
	buf.WriteString(" FROM test.blog")
	if from != nil {
		buf.WriteString(" WHERE `user` in ")
		int32ToIds(buf, from)
	}
	buf.WriteString(" GROUP BY hits,is_published ")
	buf.WriteString(m.getLimitQuery(offset, limit, sorts))

	rows, err := db.MysqlQuery(buf.String())
	if err != nil {
		return nil, fmt.Errorf("test.Blog query error: %v", err)
	}
	defer rows.Close()

	results = new(BlogGroupUnAssigned)
	for rows.Next() {
		var (
			COUNT       int
			Hits        int32
			IsPublished bool
		)
		err := rows.Scan(&(COUNT),
			&(Hits),
			&(IsPublished),
		)
		if err != nil {
			return nil, err
		}

		results.COUNT = append(results.COUNT, COUNT)
		results.Hits = append(results.Hits, Hits)
		results.IsPublished = append(results.IsPublished, IsPublished)

	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("test.Blog fetch result error: %v", err)
	}

	return

}
