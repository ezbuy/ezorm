package orm

import (
	"database/sql"
	"errors"
)

// IsErrorNotFound is a sql.ErrNoRows wrapper
func IsErrNotFound(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}
