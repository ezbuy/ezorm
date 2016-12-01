package db

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"
)

func requestTimeLogger(queryer Queryer, query string, args ...interface{}) Queryer {
	return func(query string, args ...interface{}) (interface{}, error) {
		start := time.Now()
		defer func() {
			fmt.Printf("query time: %.6f seconds\n", time.Now().Sub(start).Seconds())
		}()

		time.Sleep(time.Millisecond * time.Duration(rand.Int31n(100)))
		return queryer(query, args...)
	}
}

var _testSqlSvr *SqlServer
var (
	host     = os.Getenv("host")
	userId   = os.Getenv("userId")
	password = os.Getenv("password")
	database = os.Getenv("database")
)

func init() {
	dsn := fmt.Sprintf("server=%s;user id=%s;password=%s;DATABASE=%s",
		host, userId, password, database)

	_testSqlSvr = GetSqlServer(dsn)
	_testSqlSvr.AddQueryWrapper(requestTimeLogger)
}

func TestQuery(t *testing.T) {
	var result int32
	randNumber := rand.Int31()
	err := _testSqlSvr.Query(&result, "SELECT ?", randNumber)
	if err != nil {
		t.Error(err)
	}

	if result != randNumber {
		t.Errorf("invalid result: %d != %d", result, randNumber)
	}
}

func TestExec(t *testing.T) {
	result, err := _testSqlSvr.Exec(`
	BEGIN
	DECLARE @VarId INT
	SET @VarId = ?
	END`, rand.Int31())
	if err != nil {
		t.Error(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Error(err)
	}

	if rowsAffected != 1 {
		t.Error("rowsAffected!=1")
	}
}
