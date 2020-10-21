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

func (m *_BlogMgr) queryOne(ctx context.Context, query string, args ...interface{}) (*Blog, error) {
	ret, err := m.queryLimit(ctx, query, 1, args...)
	if err != nil {
		return nil, err
	}
	if len(ret) == 0 {
		return nil, sql.ErrNoRows
	}
	return ret[0], nil
}

func (m *_BlogMgr) query(ctx context.Context, query string, args ...interface{}) (results []*Blog, err error) {
	return m.queryLimit(ctx, query, -1, args...)
}

func (m *_BlogMgr) Query(query string, args ...interface{}) (results []*Blog, err error) {
	return m.QueryContext(context.Background(), query, args...)
}

func (m *_BlogMgr) QueryContext(ctx context.Context, query string, args ...interface{}) (results []*Blog, err error) {
	return m.queryLimit(ctx, query, -1, args...)
}

func (*_BlogMgr) queryLimit(ctx context.Context, query string, limit int, args ...interface{}) (results []*Blog, err error) {
	rows, err := db.MysqlQueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("ezorm.Blog query error: %v", err)
	}
	defer rows.Close()

	var CreatedAt string
	var UpdatedAt string

	offset := 0
	for rows.Next() {
		if limit >= 0 && offset >= limit {
			break
		}
		offset++

		var result Blog
		err := rows.Scan(&(result.Id),
			&(result.UserId),
			&(result.Title),
			&(result.Content),
			&(result.Status),
			&(result.Readed),
			&CreatedAt, &UpdatedAt)
		if err != nil {
			return nil, err
		}

		result.CreatedAt = db.TimeParse(CreatedAt)

		result.UpdatedAt = db.TimeParse(UpdatedAt)

		results = append(results, &result)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ezorm.Blog fetch result error: %v", err)
	}

	return
}

func (m *_BlogMgr) Save(obj *Blog) (sql.Result, error) {
	if obj.Id == 0 {
		return m.saveInsert(obj)
	}
	return m.saveUpdate(obj)
}

func (m *_BlogMgr) saveInsert(obj *Blog) (sql.Result, error) {
	query := "INSERT INTO ezorm.blogs (`user_id`, `title`, `content`, `status`, `readed`, `created_at`, `updated_at`) VALUES (?, ?, ?, ?, ?, ?, ?)"
	result, err := db.MysqlExec(query, obj.UserId, obj.Title, obj.Content, obj.Status, obj.Readed, db.TimeFormat(obj.CreatedAt), db.TimeFormat(obj.UpdatedAt))
	if err != nil {
		return result, err
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return result, err
	}

	obj.Id = int32(lastInsertId)

	return result, err
}

func (m *_BlogMgr) saveUpdate(obj *Blog) (sql.Result, error) {
	query := "UPDATE ezorm.blogs SET `user_id`=?, `title`=?, `content`=?, `status`=?, `readed`=?, `created_at`=?, `updated_at`=? WHERE `id`=?"
	return db.MysqlExec(query, obj.UserId, obj.Title, obj.Content, obj.Status, obj.Readed, db.TimeFormat(obj.CreatedAt), db.TimeFormat(obj.UpdatedAt), obj.Id)
}

func (m *_BlogMgr) InsertBatch(objs []*Blog) (sql.Result, error) {
	if len(objs) == 0 {
		return nil, fmt.Errorf("Empty insert")
	}

	values := make([]string, 0, len(objs))
	params := make([]interface{}, 0, len(objs)*7)
	for _, obj := range objs {
		values = append(values, "(?, ?, ?, ?, ?, ?, ?)")
		params = append(params, obj.UserId, obj.Title, obj.Content, obj.Status, obj.Readed, db.TimeFormat(obj.CreatedAt), db.TimeFormat(obj.UpdatedAt))
	}
	query := fmt.Sprintf("INSERT INTO ezorm.blogs (`user_id`, `title`, `content`, `status`, `readed`, `created_at`, `updated_at`) VALUES %s", strings.Join(values, ","))
	return db.MysqlExec(query, params...)
}

func (m *_BlogMgr) FindByID(id int32) (*Blog, error) {
	return m.FindByIDContext(context.Background(), id)
}

func (m *_BlogMgr) FindByIDContext(ctx context.Context, id int32) (*Blog, error) {
	query := "SELECT `id`, `user_id`, `title`, `content`, `status`, `readed`, `created_at`, `updated_at` FROM ezorm.blogs WHERE id=?"
	return m.queryOne(ctx, query, id)
}

func (m *_BlogMgr) FindByIDs(ids []int32) ([]*Blog, error) {
	return m.FindByIDsContext(context.Background(), ids)
}

func (m *_BlogMgr) FindByIDsContext(ctx context.Context, ids []int32) ([]*Blog, error) {
	idsLen := len(ids)
	placeHolders := make([]string, 0, idsLen)
	args := make([]interface{}, 0, idsLen)
	for _, id := range ids {
		placeHolders = append(placeHolders, "?")
		args = append(args, id)
	}

	query := fmt.Sprintf(
		"SELECT `id`, `user_id`, `title`, `content`, `status`, `readed`, `created_at`, `updated_at` FROM ezorm.blogs WHERE id IN (%s)",
		strings.Join(placeHolders, ","))
	return m.query(ctx, query, args...)
}

func (m *_BlogMgr) FindInId(ids []int32, sortFields ...string) ([]*Blog, error) {
	return m.FindInIdContext(context.Background(), ids, sortFields...)
}

func (m *_BlogMgr) FindInIdContext(ctx context.Context, ids []int32, sortFields ...string) ([]*Blog, error) {
	return m.FindByIDsContext(ctx, ids)
}

func (m *_BlogMgr) FindListId(Id []int32) ([]*Blog, error) {
	return m.FindListIdContext(context.Background(), Id)
}

func (m *_BlogMgr) FindListIdContext(ctx context.Context, Id []int32) ([]*Blog, error) {
	retmap, err := m.FindMapIdContext(ctx, Id)
	if err != nil {
		return nil, err
	}
	ret := make([]*Blog, len(Id))
	for idx, key := range Id {
		ret[idx] = retmap[key]
	}
	return ret, nil
}

func (m *_BlogMgr) FindMapId(Id []int32, sortFields ...string) (map[int32]*Blog, error) {
	return m.FindMapIdContext(context.Background(), Id)
}

func (m *_BlogMgr) FindMapIdContext(ctx context.Context, Id []int32, sortFields ...string) (map[int32]*Blog, error) {
	ret, err := m.FindInIdContext(ctx, Id, sortFields...)
	if err != nil {
		return nil, err
	}
	retmap := make(map[int32]*Blog, len(ret))
	for _, n := range ret {
		retmap[n.Id] = n
	}
	return retmap, nil
}

func (m *_BlogMgr) FindListUserId(UserId []int32) ([]*Blog, error) {
	retmap, err := m.FindMapUserId(UserId)
	if err != nil {
		return nil, err
	}
	ret := make([]*Blog, len(UserId))
	for idx, key := range UserId {
		ret[idx] = retmap[key]
	}
	return ret, nil
}

func (m *_BlogMgr) FindMapUserId(UserId []int32) (map[int32]*Blog, error) {
	ret, err := m.FindInUserId(UserId)
	if err != nil {
		return nil, err
	}
	retmap := make(map[int32]*Blog, len(ret))
	for _, n := range ret {
		retmap[n.UserId] = n
	}
	return retmap, nil
}

func (m *_BlogMgr) FindInUserId(UserId []int32, sortFields ...string) ([]*Blog, error) {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("SELECT `id`, `user_id`, `title`, `content`, `status`, `readed`, `created_at`, `updated_at` FROM ezorm.blogs WHERE ")

	buf.WriteString("`user_id` in ")
	int32ToIds(buf, UserId)
	return m.query(context.Background(), buf.String()+m.GetSort(sortFields))
}

func (m *_BlogMgr) FindAllByUserId(UserId int32, sortFields ...string) ([]*Blog, error) {
	return m.FindByUserId(UserId, -1, -1, sortFields...)
}

func (m *_BlogMgr) FindByUserId(UserId int32, offset int, limit int, sortFields ...string) ([]*Blog, error) {
	query := fmt.Sprintf("SELECT `id`, `user_id`, `title`, `content`, `status`, `readed`, `created_at`, `updated_at` FROM ezorm.blogs WHERE `user_id`=? %s%s", m.GetSort(sortFields), m.GetLimit(offset, limit))

	return m.query(context.Background(), query, UserId)
}

func (m *_BlogMgr) FindOne(where string, args ...interface{}) (*Blog, error) {
	return m.FindOneContext(context.Background(), where, args...)
}

func (m *_BlogMgr) FindOneContext(ctx context.Context, where string, args ...interface{}) (*Blog, error) {
	query := m.GetQuerysql(where) + m.GetLimit(0, 1)
	return m.queryOne(ctx, query, args...)
}

func (m *_BlogMgr) Find(where string, args ...interface{}) ([]*Blog, error) {
	return m.FindContext(context.Background(), where, args...)
}

func (m *_BlogMgr) FindContext(ctx context.Context, where string, args ...interface{}) ([]*Blog, error) {
	query := m.GetQuerysql(where)
	return m.query(ctx, query, args...)
}

func (m *_BlogMgr) FindAll() (results []*Blog, err error) {
	return m.Find("")
}

func (m *_BlogMgr) FindWithOffset(where string, offset int, limit int, args ...interface{}) ([]*Blog, error) {
	return m.FindWithOffsetContext(context.Background(), where, offset, limit, args...)
}

func (m *_BlogMgr) FindWithOffsetContext(ctx context.Context, where string, offset int, limit int, args ...interface{}) ([]*Blog, error) {
	query := m.GetQuerysql(where)

	query = query + " LIMIT ?, ?"

	args = append(args, offset)
	args = append(args, limit)

	return m.query(ctx, query, args...)
}

func (m *_BlogMgr) GetQuerysql(where string) string {
	query := "SELECT `id`, `user_id`, `title`, `content`, `status`, `readed`, `created_at`, `updated_at` FROM ezorm.blogs "

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

func (m *_BlogMgr) Del(where string, params ...interface{}) (sql.Result, error) {
	if where != "" {
		where = "WHERE " + where
	}
	query := "DELETE FROM ezorm.blogs " + where
	return db.MysqlExec(query, params...)
}

// argument example:
// set:"a=?, b=?"
// where:"c=? and d=?"
// params:[]interface{}{"a", "b", "c", "d"}...
func (m *_BlogMgr) Update(set, where string, params ...interface{}) (sql.Result, error) {
	query := fmt.Sprintf("UPDATE ezorm.blogs SET %s", set)
	if where != "" {
		query = fmt.Sprintf("UPDATE ezorm.blogs SET %s WHERE %s", set, where)
	}
	return db.MysqlExec(query, params...)
}

func (m *_BlogMgr) Count(where string, args ...interface{}) (int32, error) {
	return m.CountContext(context.Background(), where, args...)
}

func (m *_BlogMgr) CountContext(ctx context.Context, where string, args ...interface{}) (int32, error) {
	query := "SELECT COUNT(*) FROM ezorm.blogs"
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

func (m *_BlogMgr) GetSort(sorts []string) string {
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

func (m *_BlogMgr) GetId2Obj(objs []*Blog) map[int32]*Blog {
	id2obj := make(map[int32]*Blog, len(objs))
	for _, obj := range objs {
		id2obj[obj.Id] = obj
	}
	return id2obj
}

func (m *_BlogMgr) GetIds(objs []*Blog) []int32 {
	ids := make([]int32, len(objs))
	for i, obj := range objs {
		ids[i] = obj.Id
	}
	return ids
}

func (m *_BlogMgr) GetLimit(offset, limit int) string {
	if limit <= 0 {
		return ""
	}
	if offset <= 0 {
		return fmt.Sprintf(" LIMIT %d", limit)
	}
	return fmt.Sprintf(" LIMIT %d, %d", offset, limit)
}
