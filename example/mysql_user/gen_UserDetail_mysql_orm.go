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

func (m *_UserDetailMgr) queryOne(ctx context.Context, query string, args ...interface{}) (*UserDetail, error) {
	ret, err := m.queryLimit(ctx, query, 1, args...)
	if err != nil {
		return nil, err
	}
	if len(ret) == 0 {
		return nil, sql.ErrNoRows
	}
	return ret[0], nil
}

func (m *_UserDetailMgr) query(ctx context.Context, query string, args ...interface{}) (results []*UserDetail, err error) {
	return m.queryLimit(ctx, query, -1, args...)
}

func (m *_UserDetailMgr) Query(query string, args ...interface{}) (results []*UserDetail, err error) {
	return m.QueryContext(context.Background(), query, args...)
}

func (m *_UserDetailMgr) QueryContext(ctx context.Context, query string, args ...interface{}) (results []*UserDetail, err error) {
	return m.queryLimit(ctx, query, -1, args...)
}

func (*_UserDetailMgr) queryLimit(ctx context.Context, query string, limit int, args ...interface{}) (results []*UserDetail, err error) {
	rows, err := db.MysqlQueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("test.UserDetail query error: %v", err)
	}
	defer rows.Close()

	offset := 0
	for rows.Next() {
		if limit >= 0 && offset >= limit {
			break
		}
		offset++

		var result UserDetail
		err := rows.Scan(&(result.Id),
			&(result.UserId),
			&(result.Score),
			&(result.Balance),
			&(result.Text),
		)
		if err != nil {
			return nil, err
		}

		results = append(results, &result)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("test.UserDetail fetch result error: %v", err)
	}

	return
}

func (m *_UserDetailMgr) Save(obj *UserDetail) (sql.Result, error) {
	if obj.Id == 0 {
		return m.saveInsert(obj)
	}
	return m.saveUpdate(obj)
}

func (m *_UserDetailMgr) saveInsert(obj *UserDetail) (sql.Result, error) {
	query := "INSERT INTO test.user_detail (`user_id`, `score`, `balance`, `text`) VALUES (?, ?, ?, ?)"
	result, err := db.MysqlExec(query, obj.UserId, obj.Score, obj.Balance, obj.Text)
	if err != nil {
		return result, err
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return result, err
	}

	obj.Id = int64(lastInsertId)

	return result, err
}

func (m *_UserDetailMgr) saveUpdate(obj *UserDetail) (sql.Result, error) {
	query := "UPDATE test.user_detail SET `user_id`=?, `score`=?, `balance`=?, `text`=? WHERE `id`=?"
	return db.MysqlExec(query, obj.UserId, obj.Score, obj.Balance, obj.Text, obj.Id)
}

func (m *_UserDetailMgr) InsertBatch(objs []*UserDetail) (sql.Result, error) {
	if len(objs) == 0 {
		return nil, fmt.Errorf("Empty insert")
	}

	values := make([]string, 0, len(objs))
	params := make([]interface{}, 0, len(objs)*4)
	for _, obj := range objs {
		values = append(values, "(?, ?, ?, ?)")
		params = append(params, obj.UserId, obj.Score, obj.Balance, obj.Text)
	}
	query := fmt.Sprintf("INSERT INTO test.user_detail (`user_id`, `score`, `balance`, `text`) VALUES %s", strings.Join(values, ","))
	return db.MysqlExec(query, params...)
}

func (m *_UserDetailMgr) FindByID(id int64) (*UserDetail, error) {
	return m.FindByIDContext(context.Background(), id)
}

func (m *_UserDetailMgr) FindByIDContext(ctx context.Context, id int64) (*UserDetail, error) {
	query := "SELECT `id`, `user_id`, `score`, `balance`, `text` FROM test.user_detail WHERE id=?"
	return m.queryOne(ctx, query, id)
}

func (m *_UserDetailMgr) FindByIDs(ids []int64) ([]*UserDetail, error) {
	return m.FindByIDsContext(context.Background(), ids)
}

func (m *_UserDetailMgr) FindByIDsContext(ctx context.Context, ids []int64) ([]*UserDetail, error) {
	idsLen := len(ids)
	placeHolders := make([]string, 0, idsLen)
	args := make([]interface{}, 0, idsLen)
	for _, id := range ids {
		placeHolders = append(placeHolders, "?")
		args = append(args, id)
	}

	query := fmt.Sprintf(
		"SELECT `id`, `user_id`, `score`, `balance`, `text` FROM test.user_detail WHERE id IN (%s)",
		strings.Join(placeHolders, ","))
	return m.query(ctx, query, args...)
}

func (m *_UserDetailMgr) FindInId(ids []int64, sortFields ...string) ([]*UserDetail, error) {
	return m.FindInIdContext(context.Background(), ids, sortFields...)
}

func (m *_UserDetailMgr) FindInIdContext(ctx context.Context, ids []int64, sortFields ...string) ([]*UserDetail, error) {
	return m.FindByIDsContext(ctx, ids)
}

func (m *_UserDetailMgr) FindListId(Id []int64) ([]*UserDetail, error) {
	return m.FindListIdContext(context.Background(), Id)
}

func (m *_UserDetailMgr) FindListIdContext(ctx context.Context, Id []int64) ([]*UserDetail, error) {
	retmap, err := m.FindMapIdContext(ctx, Id)
	if err != nil {
		return nil, err
	}
	ret := make([]*UserDetail, len(Id))
	for idx, key := range Id {
		ret[idx] = retmap[key]
	}
	return ret, nil
}

func (m *_UserDetailMgr) FindMapId(Id []int64, sortFields ...string) (map[int64]*UserDetail, error) {
	return m.FindMapIdContext(context.Background(), Id)
}

func (m *_UserDetailMgr) FindMapIdContext(ctx context.Context, Id []int64, sortFields ...string) (map[int64]*UserDetail, error) {
	ret, err := m.FindInIdContext(ctx, Id, sortFields...)
	if err != nil {
		return nil, err
	}
	retmap := make(map[int64]*UserDetail, len(ret))
	for _, n := range ret {
		retmap[n.Id] = n
	}
	return retmap, nil
}

func (m *_UserDetailMgr) FindListUserId(UserId []int64) ([]*UserDetail, error) {
	retmap, err := m.FindMapUserId(UserId)
	if err != nil {
		return nil, err
	}
	ret := make([]*UserDetail, len(UserId))
	for idx, key := range UserId {
		ret[idx] = retmap[key]
	}
	return ret, nil
}

func (m *_UserDetailMgr) FindMapUserId(UserId []int64) (map[int64]*UserDetail, error) {
	ret, err := m.FindInUserId(UserId)
	if err != nil {
		return nil, err
	}
	retmap := make(map[int64]*UserDetail, len(ret))
	for _, n := range ret {
		retmap[n.UserId] = n
	}
	return retmap, nil
}

func (m *_UserDetailMgr) FindInUserId(UserId []int64, sortFields ...string) ([]*UserDetail, error) {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("SELECT `id`, `user_id`, `score`, `balance`, `text` FROM test.user_detail WHERE ")

	buf.WriteString("`user_id` in ")
	int64ToIds(buf, UserId)
	return m.query(context.Background(), buf.String()+m.GetSort(sortFields))
}

func (m *_UserDetailMgr) FindAllByUserId(UserId int64, sortFields ...string) ([]*UserDetail, error) {
	return m.FindByUserId(UserId, -1, -1, sortFields...)
}

func (m *_UserDetailMgr) FindByUserId(UserId int64, offset int, limit int, sortFields ...string) ([]*UserDetail, error) {
	query := fmt.Sprintf("SELECT `id`, `user_id`, `score`, `balance`, `text` FROM test.user_detail WHERE `user_id`=? %s%s", m.GetSort(sortFields), m.GetLimit(offset, limit))

	return m.query(context.Background(), query, UserId)
}

func (m *_UserDetailMgr) FindOne(where string, args ...interface{}) (*UserDetail, error) {
	return m.FindOneContext(context.Background(), where, args...)
}

func (m *_UserDetailMgr) FindOneContext(ctx context.Context, where string, args ...interface{}) (*UserDetail, error) {
	query := m.GetQuerysql(where) + m.GetLimit(0, 1)
	return m.queryOne(ctx, query, args...)
}

func (m *_UserDetailMgr) Find(where string, args ...interface{}) ([]*UserDetail, error) {
	return m.FindContext(context.Background(), where, args...)
}

func (m *_UserDetailMgr) FindContext(ctx context.Context, where string, args ...interface{}) ([]*UserDetail, error) {
	query := m.GetQuerysql(where)
	return m.query(ctx, query, args...)
}

func (m *_UserDetailMgr) FindAll() (results []*UserDetail, err error) {
	return m.Find("")
}

func (m *_UserDetailMgr) FindWithOffset(where string, offset int, limit int, args ...interface{}) ([]*UserDetail, error) {
	return m.FindWithOffsetContext(context.Background(), where, offset, limit, args...)
}

func (m *_UserDetailMgr) FindWithOffsetContext(ctx context.Context, where string, offset int, limit int, args ...interface{}) ([]*UserDetail, error) {
	query := m.GetQuerysql(where)

	query = query + " LIMIT ?, ?"

	args = append(args, offset)
	args = append(args, limit)

	return m.query(ctx, query, args...)
}

func (m *_UserDetailMgr) GetQuerysql(where string) string {
	query := "SELECT `id`, `user_id`, `score`, `balance`, `text` FROM test.user_detail "

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

func (m *_UserDetailMgr) Del(where string, params ...interface{}) (sql.Result, error) {
	if where != "" {
		where = "WHERE " + where
	}
	query := "DELETE FROM test.user_detail " + where
	return db.MysqlExec(query, params...)
}

// argument example:
// set:"a=?, b=?"
// where:"c=? and d=?"
// params:[]interface{}{"a", "b", "c", "d"}...
func (m *_UserDetailMgr) Update(set, where string, params ...interface{}) (sql.Result, error) {
	query := fmt.Sprintf("UPDATE test.user_detail SET %s", set)
	if where != "" {
		query = fmt.Sprintf("UPDATE test.user_detail SET %s WHERE %s", set, where)
	}
	return db.MysqlExec(query, params...)
}

func (m *_UserDetailMgr) Count(where string, args ...interface{}) (int32, error) {
	return m.CountContext(context.Background(), where, args...)
}

func (m *_UserDetailMgr) CountContext(ctx context.Context, where string, args ...interface{}) (int32, error) {
	query := "SELECT COUNT(*) FROM test.user_detail"
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

func (m *_UserDetailMgr) GetSort(sorts []string) string {
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

func (m *_UserDetailMgr) GetId2Obj(objs []*UserDetail) map[int64]*UserDetail {
	id2obj := make(map[int64]*UserDetail, len(objs))
	for _, obj := range objs {
		id2obj[obj.Id] = obj
	}
	return id2obj
}

func (m *_UserDetailMgr) GetIds(objs []*UserDetail) []int64 {
	ids := make([]int64, len(objs))
	for i, obj := range objs {
		ids[i] = obj.Id
	}
	return ids
}

func (m *_UserDetailMgr) GetLimit(offset, limit int) string {
	if limit <= 0 {
		return ""
	}
	if offset <= 0 {
		return fmt.Sprintf(" LIMIT %d", limit)
	}
	return fmt.Sprintf(" LIMIT %d, %d", offset, limit)
}
