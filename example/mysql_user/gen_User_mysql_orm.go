package test

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

func (m *_UserMgr) Query(query string, args ...interface{}) (results []*User, err error) {
	return m.QueryContext(context.Background(), query, args...)
}

func (m *_UserMgr) QueryContext(ctx context.Context, query string, args ...interface{}) (results []*User, err error) {
	return m.queryLimit(ctx, query, -1, args...)
}

func (*_UserMgr) queryLimit(ctx context.Context, query string, limit int, args ...interface{}) (results []*User, err error) {
	rows, err := db.MysqlQueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("test.User query error: %v", err)
	}
	defer rows.Close()

	var Text sql.NullString

	offset := 0
	for rows.Next() {
		if limit >= 0 && offset >= limit {
			break
		}
		offset++

		var result User
		err := rows.Scan(&(result.UserId),
			&(result.Name),
			&(result.Phone),
			&(result.Age),
			&(result.Balance),
			&Text, &(result.CreateDate),
		)
		if err != nil {
			return nil, err
		}

		result.Text = Text.String

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
	query := "INSERT INTO test.user (`name`, `phone`, `age`, `balance`, `text`, `create_date`) VALUES (?, ?, ?, ?, ?, ?)"
	result, err := db.MysqlExec(query, obj.Name, obj.Phone, obj.Age, obj.Balance, obj.Text, obj.CreateDate)
	if err != nil {
		return result, err
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return result, err
	}

	obj.UserId = int64(lastInsertId)

	return result, err
}

func (m *_UserMgr) saveUpdate(obj *User) (sql.Result, error) {
	query := "UPDATE test.user SET `name`=?, `phone`=?, `age`=?, `balance`=?, `text`=?, `create_date`=? WHERE `user_id`=?"
	return db.MysqlExec(query, obj.Name, obj.Phone, obj.Age, obj.Balance, obj.Text, obj.CreateDate, obj.UserId)
}

func (m *_UserMgr) InsertBatch(objs []*User) (sql.Result, error) {
	if len(objs) == 0 {
		return nil, fmt.Errorf("Empty insert")
	}

	values := make([]string, 0, len(objs))
	params := make([]interface{}, 0, len(objs)*6)
	for _, obj := range objs {
		values = append(values, "(?, ?, ?, ?, ?, ?)")
		params = append(params, obj.Name, obj.Phone, obj.Age, obj.Balance, obj.Text, obj.CreateDate)
	}
	query := fmt.Sprintf("INSERT INTO test.user (`name`, `phone`, `age`, `balance`, `text`, `create_date`) VALUES %s", strings.Join(values, ","))
	return db.MysqlExec(query, params...)
}

func (m *_UserMgr) FindByID(id int64) (*User, error) {
	return m.FindByIDContext(context.Background(), id)
}

func (m *_UserMgr) FindByIDContext(ctx context.Context, id int64) (*User, error) {
	query := "SELECT `user_id`, `name`, `phone`, `age`, `balance`, `text`, `create_date` FROM test.user WHERE user_id=?"
	return m.queryOne(ctx, query, id)
}

func (m *_UserMgr) FindByIDs(ids []int64) ([]*User, error) {
	return m.FindByIDsContext(context.Background(), ids)
}

func (m *_UserMgr) FindByIDsContext(ctx context.Context, ids []int64) ([]*User, error) {
	idsLen := len(ids)
	placeHolders := make([]string, 0, idsLen)
	args := make([]interface{}, 0, idsLen)
	for _, id := range ids {
		placeHolders = append(placeHolders, "?")
		args = append(args, id)
	}

	query := fmt.Sprintf(
		"SELECT `user_id`, `name`, `phone`, `age`, `balance`, `text`, `create_date` FROM test.user WHERE user_id IN (%s)",
		strings.Join(placeHolders, ","))
	return m.query(ctx, query, args...)
}

func (m *_UserMgr) FindInUserId(ids []int64, sortFields ...string) ([]*User, error) {
	return m.FindInUserIdContext(context.Background(), ids, sortFields...)
}

func (m *_UserMgr) FindInUserIdContext(ctx context.Context, ids []int64, sortFields ...string) ([]*User, error) {
	return m.FindByIDsContext(ctx, ids)
}

func (m *_UserMgr) FindListUserId(UserId []int64) ([]*User, error) {
	return m.FindListUserIdContext(context.Background(), UserId)
}

func (m *_UserMgr) FindListUserIdContext(ctx context.Context, UserId []int64) ([]*User, error) {
	retmap, err := m.FindMapUserIdContext(ctx, UserId)
	if err != nil {
		return nil, err
	}
	ret := make([]*User, len(UserId))
	for idx, key := range UserId {
		ret[idx] = retmap[key]
	}
	return ret, nil
}

func (m *_UserMgr) FindMapUserId(UserId []int64, sortFields ...string) (map[int64]*User, error) {
	return m.FindMapUserIdContext(context.Background(), UserId)
}

func (m *_UserMgr) FindMapUserIdContext(ctx context.Context, UserId []int64, sortFields ...string) (map[int64]*User, error) {
	ret, err := m.FindInUserIdContext(ctx, UserId, sortFields...)
	if err != nil {
		return nil, err
	}
	retmap := make(map[int64]*User, len(ret))
	for _, n := range ret {
		retmap[n.UserId] = n
	}
	return retmap, nil
}

func (m *_UserMgr) FindInNamePhone(Name []string, Phone []string, sortFields ...string) ([]*User, error) {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("SELECT `user_id`, `name`, `phone`, `age`, `balance`, `text`, `create_date` FROM test.user WHERE ")

	buf.WriteString("`name` in ")
	stringToIds(buf, Name)
	buf.WriteString(" AND ")

	buf.WriteString("`phone` in ")
	stringToIds(buf, Phone)
	return m.query(context.Background(), buf.String()+m.GetSort(sortFields))
}

func (m *_UserMgr) FindAllByNamePhone(Name string, Phone string, sortFields ...string) ([]*User, error) {
	return m.FindByNamePhone(Name, Phone, -1, -1, sortFields...)
}

func (m *_UserMgr) FindByNamePhone(Name string, Phone string, offset int, limit int, sortFields ...string) ([]*User, error) {
	query := fmt.Sprintf("SELECT `user_id`, `name`, `phone`, `age`, `balance`, `text`, `create_date` FROM test.user WHERE `name`=? AND `phone`=? %s%s", m.GetSort(sortFields), m.GetLimit(offset, limit))

	return m.query(context.Background(), query, Name, Phone)
}

func (m *_UserMgr) FindListCreateDate(CreateDate []int64) ([]*User, error) {
	retmap, err := m.FindMapCreateDate(CreateDate)
	if err != nil {
		return nil, err
	}
	ret := make([]*User, len(CreateDate))
	for idx, key := range CreateDate {
		ret[idx] = retmap[key]
	}
	return ret, nil
}

func (m *_UserMgr) FindMapCreateDate(CreateDate []int64) (map[int64]*User, error) {
	ret, err := m.FindInCreateDate(CreateDate)
	if err != nil {
		return nil, err
	}
	retmap := make(map[int64]*User, len(ret))
	for _, n := range ret {
		retmap[n.CreateDate] = n
	}
	return retmap, nil
}

func (m *_UserMgr) FindInCreateDate(CreateDate []int64, sortFields ...string) ([]*User, error) {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("SELECT `user_id`, `name`, `phone`, `age`, `balance`, `text`, `create_date` FROM test.user WHERE ")

	buf.WriteString("`create_date` in ")
	int64ToIds(buf, CreateDate)
	return m.query(context.Background(), buf.String()+m.GetSort(sortFields))
}

func (m *_UserMgr) FindAllByCreateDate(CreateDate int64, sortFields ...string) ([]*User, error) {
	return m.FindByCreateDate(CreateDate, -1, -1, sortFields...)
}

func (m *_UserMgr) FindByCreateDate(CreateDate int64, offset int, limit int, sortFields ...string) ([]*User, error) {
	query := fmt.Sprintf("SELECT `user_id`, `name`, `phone`, `age`, `balance`, `text`, `create_date` FROM test.user WHERE `create_date`=? %s%s", m.GetSort(sortFields), m.GetLimit(offset, limit))

	return m.query(context.Background(), query, CreateDate)
}

func (m *_UserMgr) FindOne(where string, args ...interface{}) (*User, error) {
	return m.FindOneContext(context.Background(), where, args...)
}

func (m *_UserMgr) FindOneContext(ctx context.Context, where string, args ...interface{}) (*User, error) {
	query := m.GetQuerysql(where) + m.GetLimit(0, 1)
	return m.queryOne(ctx, query, args...)
}

func (m *_UserMgr) Find(where string, args ...interface{}) ([]*User, error) {
	return m.FindContext(context.Background(), where, args...)
}

func (m *_UserMgr) FindContext(ctx context.Context, where string, args ...interface{}) ([]*User, error) {
	query := m.GetQuerysql(where)
	return m.query(ctx, query, args...)
}

func (m *_UserMgr) FindAll() (results []*User, err error) {
	return m.Find("")
}

func (m *_UserMgr) FindWithOffset(where string, offset int, limit int, args ...interface{}) ([]*User, error) {
	return m.FindWithOffsetContext(context.Background(), where, offset, limit, args...)
}

func (m *_UserMgr) FindWithOffsetContext(ctx context.Context, where string, offset int, limit int, args ...interface{}) ([]*User, error) {
	query := m.GetQuerysql(where)

	query = query + " LIMIT ?, ?"

	args = append(args, offset)
	args = append(args, limit)

	return m.query(ctx, query, args...)
}

func (m *_UserMgr) GetQuerysql(where string) string {
	query := "SELECT `user_id`, `name`, `phone`, `age`, `balance`, `text`, `create_date` FROM test.user "

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

func (m *_UserMgr) Del(where string, params ...interface{}) (sql.Result, error) {
	if where != "" {
		where = "WHERE " + where
	}
	query := "DELETE FROM test.user " + where
	return db.MysqlExec(query, params...)
}

// argument example:
// set:"a=?, b=?"
// where:"c=? and d=?"
// params:[]interface{}{"a", "b", "c", "d"}...
func (m *_UserMgr) Update(set, where string, params ...interface{}) (sql.Result, error) {
	query := fmt.Sprintf("UPDATE test.user SET %s", set)
	if where != "" {
		query = fmt.Sprintf("UPDATE test.user SET %s WHERE %s", set, where)
	}
	return db.MysqlExec(query, params...)
}

func (m *_UserMgr) Count(where string, args ...interface{}) (int32, error) {
	return m.CountContext(context.Background(), where, args...)
}

func (m *_UserMgr) CountContext(ctx context.Context, where string, args ...interface{}) (int32, error) {
	query := "SELECT COUNT(*) FROM test.user"
	if where != "" {
		query = query + " WHERE " + where
	}

	rows, err := db.MysqlQueryContext(ctx, query, args...)
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

func (m *_UserMgr) GetId2Obj(objs []*User) map[int64]*User {
	id2obj := make(map[int64]*User, len(objs))
	for _, obj := range objs {
		id2obj[obj.UserId] = obj
	}
	return id2obj
}

func (m *_UserMgr) GetIds(objs []*User) []int64 {
	ids := make([]int64, len(objs))
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
