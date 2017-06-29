package test

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/ezbuy/ezorm/db"
	"strings"
	"time"
)

var (
	_ time.Time
	_ bytes.Buffer
	_ = strings.Index
)

// -----------------------------------------------------------------------------

func (m *_BlogTempMgr) queryOne(query string, args ...interface{}) (*BlogTemp, error) {
	ret, err := m.queryLimit(query, 1, args...)
	if err != nil {
		return nil, err
	}
	if len(ret) == 0 {
		return nil, sql.ErrNoRows
	}
	return ret[0], nil
}

func (m *_BlogTempMgr) query(query string, args ...interface{}) (results []*BlogTemp, err error) {
	return m.queryLimit(query, -1, args...)
}

func (m *_BlogTempMgr) Query(query string, args ...interface{}) (results []*BlogTemp, err error) {
	return m.queryLimit(query, -1, args...)
}

func (*_BlogTempMgr) queryLimit(query string, limit int, args ...interface{}) (results []*BlogTemp, err error) {
	rows, err := db.MysqlQuery(query, args...)
	if err != nil {
		return nil, fmt.Errorf("test.BlogTemp query error: %v", err)
	}
	defer rows.Close()

	var Body sql.NullString

	offset := 0
	for rows.Next() {
		if limit >= 0 && offset >= limit {
			break
		}
		offset++

		var result BlogTemp
		err := rows.Scan(&(result.BlogId),
			&(result.Title),
			&(result.Hits),
			&(result.Slug),
			&Body, &(result.User),
			&(result.GhostNumber),
		)
		if err != nil {
			return nil, err
		}

		result.Body = Body.String

		results = append(results, &result)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("test.BlogTemp fetch result error: %v", err)
	}

	return
}
