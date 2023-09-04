package orm

import (
	"database/sql"
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
)

// IsErrorNotFound is a sql.ErrNoRows wrapper
func IsErrNotFound(err error) bool {
	return errors.Is(err, sql.ErrNoRows) || errors.Is(err, mongo.ErrNoDocuments)
}
