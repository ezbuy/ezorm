package people

import (
	"database/sql"
	"fmt"
	"github.com/ezbuy/ezorm/db"
	"strings"
	"time"
)

var (
	_ time.Time
)

// -----------------------------------------------------------------------------

func (m *_UserMgr) queryOne(query string, args ...interface{}) (*User, error) {
	ret, err := m.queryLimit(query, 1, args...)
	if err != nil {
		return nil, err
	}
	if len(ret) == 0 {
		return nil, sql.ErrNoRows
	}
	return ret[0], nil
}

func (m *_UserMgr) query(query string, args ...interface{}) (results []*User, err error) {
	return m.queryLimit(query, -1, args...)
}

func (*_UserMgr) queryLimit(query string, limit int, args ...interface{}) (results []*User, err error) {
	rows, err := db.MysqlQuery(query, args...)
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
		err := rows.Scan(&(result.UserId),
			&(result.Name),
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

func (m *_UserMgr) Save(obj *User) (sql.Result, error) {
	if obj.UserId == 0 {
		return m.saveInsert(obj)
	}
	return m.saveUpdate(obj)
}

func (m *_UserMgr) saveInsert(obj *User) (sql.Result, error) {
	query := "INSERT INTO test.test_blog (`Name`) VALUES (?)"
	result, err := db.MysqlExec(query, obj.Name)
	if err != nil {
		return result, err
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return result, err
	}

	obj.UserId = int32(lastInsertId)

	return result, err
}

func (m *_UserMgr) saveUpdate(obj *User) (sql.Result, error) {
	query := "UPDATE test.test_blog SET `Name`=? WHERE `UserId`=?"
	return db.MysqlExec(query, obj.Name, obj.UserId)
}

func (m *_UserMgr) InsertBatch(objs []*User) (sql.Result, error) {
	if len(objs) == 0 {
		return nil, fmt.Errorf("Empty insert")
	}

	values := make([]string, 0, len(objs))
	params := make([]interface{}, 0, len(objs)*1)
	for _, obj := range objs {
		values = append(values, "(?)")
		params = append(params, obj.Name)
	}
	query := fmt.Sprintf("INSERT INTO test.test_blog (Name) VALUES %s", strings.Join(values, ","))
	return db.MysqlExec(query, params...)
}

func (m *_UserMgr) FindByID(id int32) (*User, error) {
	query := "SELECT `UserId`, `Name` FROM test.test_blog WHERE UserId=?"
	return m.queryOne(query, id)
}

func (m *_UserMgr) FindByIDs(ids []int32) ([]*User, error) {
	idsLen := len(ids)
	placeHolders := make([]string, 0, idsLen)
	args := make([]interface{}, 0, idsLen)
	for _, id := range ids {
		placeHolders = append(placeHolders, "?")
		args = append(args, id)
	}

	query := fmt.Sprintf(
		"SELECT `UserId`, `Name` FROM test.test_blog WHERE UserId IN (%s)",
		strings.Join(placeHolders, ","))
	return m.query(query, args...)
}

func (m *_UserMgr) FindByName(Name string, offset int, limit int, sortFields ...string) ([]*User, error) {
	orderBy := "ORDER BY %s"
	if len(sortFields) != 0 {
		orderBy = fmt.Sprintf(orderBy, strings.Join(sortFields, ","))
	} else {
		orderBy = fmt.Sprintf(orderBy, "UserId")
	}

	query := fmt.Sprintf("SELECT `UserId`, `Name` FROM test.test_blog WHERE `Name`=? %s LIMIT ?, ?", orderBy)

	return m.query(query, Name, offset, limit)
}

func (m *_UserMgr) FindOne(where string, args ...interface{}) (*User, error) {
	query := m.getQuerysql(true, where)
	return m.queryOne(query, args...)
}

func (m *_UserMgr) Find(where string, args ...interface{}) ([]*User, error) {
	query := m.getQuerysql(false, where)
	return m.query(query, args...)
}

func (m *_UserMgr) FindAll() (results []*User, err error) {
	return m.Find("")
}

func (m *_UserMgr) FindWithOffset(where string, offset int, limit int, args ...interface{}) ([]*User, error) {
	query := m.getQuerysql(false, where)

	query = query + " LIMIT ?, ?"

	args = append(args, offset)
	args = append(args, limit)

	return m.query(query, args...)
}

func (m *_UserMgr) getQuerysql(topOne bool, where string) string {
	query := "SELECT `UserId`, `Name` FROM test.test_blog"

	where = strings.TrimSpace(where)
	if where != "" {
		upwhere := strings.ToUpper(where)

		if !strings.HasPrefix(upwhere, "WHERE") && !strings.HasPrefix(upwhere, "ORDER BY") {
			where = " WHERE " + where
		}

		query = query + where
	}

	if topOne {
		query += " LIMIT 1"
	}
	return query
}

func (m *_UserMgr) Del(where string, params ...interface{}) (sql.Result, error) {
	if where != "" {
		where = "WHERE " + where
	}
	query := "DELETE FROM test.test_blog " + where
	return db.MysqlExec(query, params...)
}

// argument example:
// set:"a=?, b=?"
// where:"c=? and d=?"
// params:[]interface{}{"a", "b", "c", "d"}...
func (m *_UserMgr) Update(set, where string, params ...interface{}) (sql.Result, error) {
	query := fmt.Sprintf("UPDATE test.test_blog SET %s", set)
	if where != "" {
		query = fmt.Sprintf("UPDATE test.test_blog SET %s WHERE %s", set, where)
	}
	return db.MysqlExec(query, params...)
}

func (m *_UserMgr) Count(where string, args ...interface{}) (int32, error) {
	query := "SELECT COUNT(*) FROM test.test_blog"
	if where != "" {
		query = query + " WHERE " + where
	}

	rows, err := db.MysqlQuery(query, args...)
	if err != nil {
		return 0, err
	}

	var count int32
	if rows.Next() {
		err = rows.Scan(&count)
	}

	return count, err
}
