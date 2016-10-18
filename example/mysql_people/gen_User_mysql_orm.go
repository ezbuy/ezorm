package people

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

func (m *_UserMgr) Query(query string, args ...interface{}) (results []*User, err error) {
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
			&(result.UserNumber),
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
	query := "INSERT INTO test.test_user (`user_number`, `name`) VALUES (?, ?)"
	result, err := db.MysqlExec(query, obj.UserNumber, obj.Name)
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
	query := "UPDATE test.test_user SET `user_number`=?, `name`=? WHERE `user_id`=?"
	return db.MysqlExec(query, obj.UserNumber, obj.Name, obj.UserId)
}

func (m *_UserMgr) InsertBatch(objs []*User) (sql.Result, error) {
	if len(objs) == 0 {
		return nil, fmt.Errorf("Empty insert")
	}

	values := make([]string, 0, len(objs))
	params := make([]interface{}, 0, len(objs)*2)
	for _, obj := range objs {
		values = append(values, "(?, ?)")
		params = append(params, obj.UserNumber, obj.Name)
	}
	query := fmt.Sprintf("INSERT INTO test.test_user (`user_number`, `name`) VALUES %s", strings.Join(values, ","))
	return db.MysqlExec(query, params...)
}

func (m *_UserMgr) FindByID(id int32) (*User, error) {
	query := "SELECT `user_id`, `user_number`, `name` FROM test.test_user WHERE user_id=?"
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
		"SELECT `user_id`, `user_number`, `name` FROM test.test_user WHERE user_id IN (%s)",
		strings.Join(placeHolders, ","))
	return m.query(query, args...)
}

func (m *_UserMgr) FindInUserId(ids []int32, sortFields ...string) ([]*User, error) {
	return m.FindByIDs(ids)
}

func (m *_UserMgr) FindListUserId(UserId []int32) ([]*User, error) {
	retmap, err := m.FindMapUserId(UserId)
	if err != nil {
		return nil, err
	}
	ret := make([]*User, len(UserId))
	for idx, key := range UserId {
		ret[idx] = retmap[key]
	}
	return ret, nil
}

func (m *_UserMgr) FindMapUserId(UserId []int32, sortFields ...string) (map[int32]*User, error) {
	ret, err := m.FindInUserId(UserId, sortFields...)
	if err != nil {
		return nil, err
	}
	retmap := make(map[int32]*User, len(ret))
	for _, n := range ret {
		retmap[n.UserId] = n
	}
	return retmap, nil
}

func (m *_UserMgr) FindListUserNumber(UserNumber []int32) ([]*User, error) {
	retmap, err := m.FindMapUserNumber(UserNumber)
	if err != nil {
		return nil, err
	}
	ret := make([]*User, len(UserNumber))
	for idx, key := range UserNumber {
		ret[idx] = retmap[key]
	}
	return ret, nil
}

func (m *_UserMgr) FindMapUserNumber(UserNumber []int32) (map[int32]*User, error) {
	ret, err := m.FindInUserNumber(UserNumber)
	if err != nil {
		return nil, err
	}
	retmap := make(map[int32]*User, len(ret))
	for _, n := range ret {
		retmap[n.UserNumber] = n
	}
	return retmap, nil
}

func (m *_UserMgr) FindInUserNumber(UserNumber []int32, sortFields ...string) ([]*User, error) {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("SELECT `user_id`, `user_number`, `name` FROM test.test_user WHERE ")

	buf.WriteString("`user_number` in ")
	int32ToIds(buf, UserNumber)
	return m.query(buf.String())
}

func (m *_UserMgr) FindOneByUserNumber(UserNumber int32) (*User, error) {
	query := "SELECT `user_id`, `user_number`, `name` FROM test.test_user WHERE user_number=?"
	return m.queryOne(query, UserNumber)
}

func (m *_UserMgr) FindListName(Name []string) ([]*User, error) {
	retmap, err := m.FindMapName(Name)
	if err != nil {
		return nil, err
	}
	ret := make([]*User, len(Name))
	for idx, key := range Name {
		ret[idx] = retmap[key]
	}
	return ret, nil
}

func (m *_UserMgr) FindMapName(Name []string) (map[string]*User, error) {
	ret, err := m.FindInName(Name)
	if err != nil {
		return nil, err
	}
	retmap := make(map[string]*User, len(ret))
	for _, n := range ret {
		retmap[n.Name] = n
	}
	return retmap, nil
}

func (m *_UserMgr) FindInName(Name []string, sortFields ...string) ([]*User, error) {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("SELECT `user_id`, `user_number`, `name` FROM test.test_user WHERE ")

	buf.WriteString("`name` in ")
	stringToIds(buf, Name)
	return m.query(buf.String())
}

func (m *_UserMgr) FindAllByName(Name string, sortFields ...string) ([]*User, error) {
	return m.FindByName(Name, -1, -1, sortFields...)
}

func (m *_UserMgr) FindByName(Name string, offset int, limit int, sortFields ...string) ([]*User, error) {
	orderBy := "ORDER BY %s"
	if len(sortFields) != 0 {
		orderBy = fmt.Sprintf(orderBy, strings.Join(sortFields, ","))
	} else {
		orderBy = fmt.Sprintf(orderBy, "user_id")
	}

	query := fmt.Sprintf("SELECT `user_id`, `user_number`, `name` FROM test.test_user WHERE `name`=? %s LIMIT ?, ?", orderBy)

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
	query := "SELECT `user_id`, `user_number`, `name` FROM test.test_user"

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
	query := "DELETE FROM test.test_user " + where
	return db.MysqlExec(query, params...)
}

// argument example:
// set:"a=?, b=?"
// where:"c=? and d=?"
// params:[]interface{}{"a", "b", "c", "d"}...
func (m *_UserMgr) Update(set, where string, params ...interface{}) (sql.Result, error) {
	query := fmt.Sprintf("UPDATE test.test_user SET %s", set)
	if where != "" {
		query = fmt.Sprintf("UPDATE test.test_user SET %s WHERE %s", set, where)
	}
	return db.MysqlExec(query, params...)
}

func (m *_UserMgr) Count(where string, args ...interface{}) (int32, error) {
	query := "SELECT COUNT(*) FROM test.test_user"
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

func (m *_UserMgr) getLimitQuery(offset, limit int, sorts []string) string {
	orderBy := ""
	if len(sorts) != 0 {
		orderBy = fmt.Sprintf("ORDER BY %s", strings.Join(sorts, ","))
	}
	if limit > 0 {
		return orderBy + fmt.Sprintf(" LIMIT %d, %d", offset, limit)
	}
	return orderBy
}
