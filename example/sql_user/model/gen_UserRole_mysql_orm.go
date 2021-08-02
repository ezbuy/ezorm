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

func (m *_UserRoleMgr) queryOne(ctx context.Context, query string, args ...interface{}) (*UserRole, error) {
	ret, err := m.queryLimit(ctx, query, 1, args...)
	if err != nil {
		return nil, err
	}
	if len(ret) == 0 {
		return nil, sql.ErrNoRows
	}
	return ret[0], nil
}

func (m *_UserRoleMgr) query(ctx context.Context, query string, args ...interface{}) (results []*UserRole, err error) {
	return m.queryLimit(ctx, query, -1, args...)
}

func (m *_UserRoleMgr) Query(query string, args ...interface{}) (results []*UserRole, err error) {
	return m.QueryContext(context.Background(), query, args...)
}

func (m *_UserRoleMgr) QueryContext(ctx context.Context, query string, args ...interface{}) (results []*UserRole, err error) {
	return m.queryLimit(ctx, query, -1, args...)
}

func (*_UserRoleMgr) queryLimit(ctx context.Context, query string, limit int, args ...interface{}) (results []*UserRole, err error) {
	rows, err := db.MysqlQueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("test.UserRole query error: %v", err)
	}
	defer rows.Close()

	offset := 0
	for rows.Next() {
		if limit >= 0 && offset >= limit {
			break
		}
		offset++

		var result UserRole
		err := rows.Scan(&(result.Id),
			&(result.UserId),
			&(result.RoleId),
		)
		if err != nil {
			return nil, err
		}

		results = append(results, &result)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("test.UserRole fetch result error: %v", err)
	}

	return
}
