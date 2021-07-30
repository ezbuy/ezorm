package model

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

func (m *_UserMgr) queryOne(ctx context.Context, query string, args ...interface{}) (*User, error) {
	ret, err := m.queryLimit(ctx, query, 1, args...)
	if err != nil {
		return nil, err
	}
	if len(ret) == 0 {
		return nil, sql.ErrNoRows
	}
	return ret[0], nil
}

func (m *_UserMgr) query(ctx context.Context, query string, args ...interface{}) (results []*User, err error) {
	return m.queryLimit(ctx, query, -1, args...)
}

func (m *_UserMgr) Query(query string, args ...interface{}) (results []*User, err error) {
	return m.QueryContext(context.Background(), query, args...)
}

func (m *_UserMgr) QueryContext(ctx context.Context, query string, args ...interface{}) (results []*User, err error) {
	return m.queryLimit(ctx, query, -1, args...)
}

func (*_UserMgr) queryLimit(ctx context.Context, query string, limit int, args ...interface{}) (results []*User, err error) {
	rows, err := db.MysqlQueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("test.User query error: %v", err)
	}
	defer rows.Close()

	offset := 0
	for rows.Next() {
		if limit >= 0 && offset >= limit {
			break
		}
		offset++

		var result User
		err := rows.Scan(&(result.Id),
			&(result.Name),
			&(result.Phone),
			&(result.Password),
		)
		if err != nil {
			return nil, err
		}

		results = append(results, &result)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("test.User fetch result error: %v", err)
	}

	return
}
