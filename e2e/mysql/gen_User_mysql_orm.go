package mysql

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/ezbuy/ezorm/v2/pkg/db"
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

func (m *_UserMgr) Query(ctx context.Context, query string, args ...interface{}) (results []*User, err error) {
	return m.queryLimit(ctx, query, -1, args...)
}

func (*_UserMgr) queryLimit(ctx context.Context, query string, limit int, args ...interface{}) (results []*User, err error) {
	rows, err := db.MysqlQuery(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("test.User query error: %w", err)
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
		return nil, fmt.Errorf("test.User fetch result error: %w", err)
	}

	return
}

func (m *_UserMgr) Save(ctx context.Context, obj *User) (sql.Result, error) {
	if obj.UserId == 0 {
		return m.saveInsert(ctx, obj)
	}
	return m.saveUpdate(ctx, obj)
}

func (m *_UserMgr) saveInsert(ctx context.Context, obj *User) (sql.Result, error) {
	query := "INSERT INTO test.test_user (`user_number`, `name`) VALUES (?, ?)"
	result, err := db.MysqlExec(ctx, query, obj.UserNumber, obj.Name)
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

func (m *_UserMgr) saveUpdate(ctx context.Context, obj *User) (sql.Result, error) {
	query := "UPDATE test.test_user SET `user_number`=?, `name`=? WHERE `user_id`=?"
	return db.MysqlExec(ctx, query, obj.UserNumber, obj.Name, obj.UserId)
}

func (m *_UserMgr) InsertBatch(ctx context.Context, objs []*User) (sql.Result, error) {
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
	return db.MysqlExec(ctx, query, params...)
}

func (m *_UserMgr) FindByID(ctx context.Context, id int32) (*User, error) {
	query := "SELECT `user_id`, `user_number`, `name` FROM test.test_user WHERE user_id=?"
	return m.queryOne(ctx, query, id)
}

func (m *_UserMgr) FindByIDs(ctx context.Context, ids []int32) ([]*User, error) {
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
	return m.query(ctx, query, args...)
}

func (m *_UserMgr) FindInUserId(ctx context.Context, ids []int32, sortFields ...string) ([]*User, error) {
	return m.FindByIDs(ctx, ids)
}

func (m *_UserMgr) FindListUserId(ctx context.Context, UserId []int32) ([]*User, error) {
	retmap, err := m.FindMapUserId(ctx, UserId)
	if err != nil {
		return nil, err
	}
	ret := make([]*User, len(UserId))
	for idx, key := range UserId {
		ret[idx] = retmap[key]
	}
	return ret, nil
}

func (m *_UserMgr) FindMapUserId(ctx context.Context, UserId []int32, sortFields ...string) (map[int32]*User, error) {
	ret, err := m.FindInUserId(ctx, UserId, sortFields...)
	if err != nil {
		return nil, err
	}
	retmap := make(map[int32]*User, len(ret))
	for _, n := range ret {
		retmap[n.UserId] = n
	}
	return retmap, nil
}

func (m *_UserMgr) FindListUserNumber(ctx context.Context, UserNumber []int32) ([]*User, error) {
	retmap, err := m.FindMapUserNumber(ctx, UserNumber)
	if err != nil {
		return nil, err
	}
	ret := make([]*User, len(UserNumber))
	for idx, key := range UserNumber {
		ret[idx] = retmap[key]
	}
	return ret, nil
}

func (m *_UserMgr) FindMapUserNumber(ctx context.Context, UserNumber []int32) (map[int32]*User, error) {
	ret, err := m.FindInUserNumber(ctx, UserNumber)
	if err != nil {
		return nil, err
	}
	retmap := make(map[int32]*User, len(ret))
	for _, n := range ret {
		retmap[n.UserNumber] = n
	}
	return retmap, nil
}

func (m *_UserMgr) FindInUserNumber(ctx context.Context, UserNumber []int32, sortFields ...string) ([]*User, error) {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("SELECT `user_id`, `user_number`, `name` FROM test.test_user WHERE ")

	buf.WriteString("`user_number` in ")
	int32ToIds(buf, UserNumber)
	return m.query(ctx, buf.String()+m.GetSort(sortFields))
}

func (m *_UserMgr) FindOneByUserNumber(ctx context.Context, UserNumber int32) (*User, error) {
	query := "SELECT `user_id`, `user_number`, `name` FROM test.test_user WHERE user_number=?"
	return m.queryOne(ctx, query, UserNumber)
}

func (m *_UserMgr) FindListName(ctx context.Context, Name []string) ([]*User, error) {
	retmap, err := m.FindMapName(ctx, Name)
	if err != nil {
		return nil, err
	}
	ret := make([]*User, len(Name))
	for idx, key := range Name {
		ret[idx] = retmap[key]
	}
	return ret, nil
}

func (m *_UserMgr) FindMapName(ctx context.Context, Name []string) (map[string]*User, error) {
	ret, err := m.FindInName(ctx, Name)
	if err != nil {
		return nil, err
	}
	retmap := make(map[string]*User, len(ret))
	for _, n := range ret {
		retmap[n.Name] = n
	}
	return retmap, nil
}

func (m *_UserMgr) FindInName(ctx context.Context, Name []string, sortFields ...string) ([]*User, error) {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("SELECT `user_id`, `user_number`, `name` FROM test.test_user WHERE ")

	buf.WriteString("`name` in ")
	stringToIds(buf, Name)
	return m.query(ctx, buf.String()+m.GetSort(sortFields))
}

func (m *_UserMgr) FindAllByName(ctx context.Context, Name string, sortFields ...string) ([]*User, error) {
	return m.FindByName(ctx, Name, -1, -1, sortFields...)
}

func (m *_UserMgr) FindByName(ctx context.Context, Name string, offset int, limit int, sortFields ...string) ([]*User, error) {
	query := fmt.Sprintf("SELECT `user_id`, `user_number`, `name` FROM test.test_user WHERE `name`=? %s%s", m.GetSort(sortFields), m.GetLimit(offset, limit))

	return m.query(ctx, query, Name)
}

func (m *_UserMgr) FindOne(ctx context.Context, where string, args ...interface{}) (*User, error) {
	query := m.GetQuerysql(where) + m.GetLimit(0, 1)
	return m.queryOne(ctx, query, args...)
}

func (m *_UserMgr) Find(ctx context.Context, where string, args ...interface{}) ([]*User, error) {
	query := m.GetQuerysql(where)
	return m.query(ctx, query, args...)
}

func (m *_UserMgr) FindAll(ctx context.Context) (results []*User, err error) {
	return m.Find(ctx, "")
}

func (m *_UserMgr) FindWithOffset(ctx context.Context, where string, offset int, limit int, args ...interface{}) ([]*User, error) {
	query := m.GetQuerysql(where)

	query = query + " LIMIT ?, ?"

	args = append(args, offset)
	args = append(args, limit)

	return m.query(ctx, query, args...)
}

func (m *_UserMgr) GetQuerysql(where string) string {
	query := "SELECT `user_id`, `user_number`, `name` FROM test.test_user "

	where = strings.TrimSpace(where)
	if where != "" {
		upwhere := strings.ToUpper(where)

		if !strings.HasPrefix(upwhere, "WHERE") && !strings.HasPrefix(upwhere, "ORDER BY") {
			where = " WHERE " + where
		}

		query = query + where
	}

	return query
}

func (m *_UserMgr) Del(ctx context.Context, where string, params ...interface{}) (sql.Result, error) {
	if where != "" {
		where = "WHERE " + where
	}
	query := "DELETE FROM test.test_user " + where
	return db.MysqlExec(ctx, query, params...)
}

// argument example:
// set:"a=?, b=?"
// where:"c=? and d=?"
// params:[]interface{}{"a", "b", "c", "d"}...
func (m *_UserMgr) Update(ctx context.Context, set, where string, params ...interface{}) (sql.Result, error) {
	query := fmt.Sprintf("UPDATE test.test_user SET %s", set)
	if where != "" {
		query = fmt.Sprintf("UPDATE test.test_user SET %s WHERE %s", set, where)
	}
	return db.MysqlExec(ctx, query, params...)
}

func (m *_UserMgr) Count(ctx context.Context, where string, args ...interface{}) (int32, error) {
	query := "SELECT COUNT(*) FROM test.test_user"
	if where != "" {
		query = query + " WHERE " + where
	}

	rows, err := db.MysqlQuery(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var count int32
	if rows.Next() {
		err = rows.Scan(&count)
	}
	if err := rows.Err(); err != nil {
		return 0, err
	}

	return count, nil
}

func (m *_UserMgr) GetSort(sorts []string) string {
	if len(sorts) == 0 {
		return ""
	}
	buf := bytes.NewBuffer(nil)
	buf.WriteString(" ORDER BY ")
	for idx, s := range sorts {
		if len(s) == 0 {
			continue
		}
		if s[0] == '-' {
			buf.WriteString(s[1:] + " DESC")
		} else {
			buf.WriteString(s)
		}
		if idx == len(sorts)-1 {
			break
		}
		buf.WriteString(",")
	}
	return buf.String()
}

func (m *_UserMgr) GetId2Obj(objs []*User) map[int32]*User {
	id2obj := make(map[int32]*User, len(objs))
	for _, obj := range objs {
		id2obj[obj.UserId] = obj
	}
	return id2obj
}

func (m *_UserMgr) GetIds(objs []*User) []int32 {
	ids := make([]int32, len(objs))
	for i, obj := range objs {
		ids[i] = obj.UserId
	}
	return ids
}

func (m *_UserMgr) GetLimit(offset, limit int) string {
	if limit <= 0 {
		return ""
	}
	if offset <= 0 {
		return fmt.Sprintf(" LIMIT %d", limit)
	}
	return fmt.Sprintf(" LIMIT %d, %d", offset, limit)
}
