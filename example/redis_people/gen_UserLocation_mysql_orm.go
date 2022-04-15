package test

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/ezbuy/ezorm/v2/db"
)

var (
	_ time.Time
	_ bytes.Buffer
	_ = strings.Index
)

// -----------------------------------------------------------------------------

func (m *_UserLocationMgr) queryOne(ctx context.Context, query string, args ...interface{}) (*UserLocation, error) {
	ret, err := m.queryLimit(ctx, query, 1, args...)
	if err != nil {
		return nil, err
	}
	if len(ret) == 0 {
		return nil, sql.ErrNoRows
	}
	return ret[0], nil
}

func (m *_UserLocationMgr) query(ctx context.Context, query string, args ...interface{}) (results []*UserLocation, err error) {
	return m.queryLimit(ctx, query, -1, args...)
}

func (m *_UserLocationMgr) Query(query string, args ...interface{}) (results []*UserLocation, err error) {
	return m.QueryContext(context.Background(), query, args...)
}

func (m *_UserLocationMgr) QueryContext(ctx context.Context, query string, args ...interface{}) (results []*UserLocation, err error) {
	return m.queryLimit(ctx, query, -1, args...)
}

func (*_UserLocationMgr) queryLimit(ctx context.Context, query string, limit int, args ...interface{}) (results []*UserLocation, err error) {
	rows, err := db.MysqlQueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("ezorm.UserLocation query error: %v", err)
	}
	defer rows.Close()

	offset := 0
	for rows.Next() {
		if limit >= 0 && offset >= limit {
			break
		}
		offset++

		var result UserLocation
		err := rows.Scan(&(result.Key),
			&(result.Longitude),
			&(result.Latitude),
			&(result.Value),
		)
		if err != nil {
			return nil, err
		}

		results = append(results, &result)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ezorm.UserLocation fetch result error: %v", err)
	}

	return
}
