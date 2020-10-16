package test

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/ezbuy/ezorm/db"
)

var (
	_ time.Time
	_ bytes.Buffer
	_ = strings.Index
)

// -----------------------------------------------------------------------------

func (m *_SortUserBlogMgr) queryOne(ctx context.Context, query string, args ...interface{}) (*SortUserBlog, error) {
	ret, err := m.queryLimit(ctx, query, 1, args...)
	if err != nil {
		return nil, err
	}
	if len(ret) == 0 {
		return nil, sql.ErrNoRows
	}
	return ret[0], nil
}

func (m *_SortUserBlogMgr) query(ctx context.Context, query string, args ...interface{}) (results []*SortUserBlog, err error) {
	return m.queryLimit(ctx, query, -1, args...)
}

func (m *_SortUserBlogMgr) Query(query string, args ...interface{}) (results []*SortUserBlog, err error) {
	return m.queryLimit(context.Background(), query, -1, args...)
}

func (*_SortUserBlogMgr) queryLimit(ctx context.Context, query string, limit int, args ...interface{}) (results []*SortUserBlog, err error) {
	rows, err := db.MysqlQueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("ezorm.SortUserBlog query error: %v", err)
	}
	defer rows.Close()

	offset := 0
	for rows.Next() {
		if limit >= 0 && offset >= limit {
			break
		}
		offset++

		var result SortUserBlog
		err := rows.Scan(&(result.Key),
			&(result.Score),
			&(result.Value),
		)
		if err != nil {
			return nil, err
		}

		results = append(results, &result)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ezorm.SortUserBlog fetch result error: %v", err)
	}

	return
}
