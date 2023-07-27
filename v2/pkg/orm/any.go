package orm

import (
	"database/sql"
	"fmt"
)

type Result interface {
	sql.NullInt64 | sql.NullString
}

var ErrScanResultTypeUnsupported = fmt.Errorf("scan result type unsupported")

// ToResult is helper function which sits in the background
// when we cannot determine the underlying type of MySQL scan result(such as generic function result).

// Inspired by https://github.com/go-sql-driver/mysql/issues/86,
// if the Go SQL driver cannot determine the underlying type of scan result(`interface{}` or `any`),
// it will fallback to TEXT protocol to communicate with MySQL server,
// therefore, the result can only be `[]uint8`(`[]byte`) or raw `string`.
// But the behavior differs from drivers.
func ToResult[T Result](rawField any) (T, error) {
	var t T
	switch rawField := rawField.(type) {
	default:
		return t, fmt.Errorf("rawField type got %T", rawField)
	case int64:
		switch any(t).(type) {
		default:
			return t, ErrScanResultTypeUnsupported
		case sql.NullInt64:
			i := &sql.NullInt64{}
			if err := i.Scan(rawField); err != nil {
				return t, err
			}
			t = any(*i).(T)
		}
	case string:
		switch any(t).(type) {
		default:
			return t, ErrScanResultTypeUnsupported
		case sql.NullString:
			s := &sql.NullString{}
			if err := s.Scan(rawField); err != nil {
				return t, err
			}
			t = any(*s).(T)
		}
	case []byte:
		switch any(t).(type) {
		default:
			return t, ErrScanResultTypeUnsupported
		case sql.NullString:
			s := &sql.NullString{}
			if err := s.Scan(rawField); err != nil {
				return t, err
			}
			t = any(*s).(T)
		case sql.NullInt64:
			i := &sql.NullInt64{}
			if err := i.Scan(rawField); err != nil {
				return t, err
			}
			t = any(*i).(T)
		}
	}
	return t, nil
}
