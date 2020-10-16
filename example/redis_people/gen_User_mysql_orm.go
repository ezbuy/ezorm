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
	return m.queryLimit(context.Background(), query, -1, args...)
}

func (*_UserMgr) queryLimit(ctx context.Context, query string, limit int, args ...interface{}) (results []*User, err error) {
	rows, err := db.MysqlQueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("ezorm.User query error: %v", err)
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

		var result User
		err := rows.Scan(&(result.Id),
			&(result.Name),
			&(result.Mailbox),
			&(result.Sex),
			&(result.Longitude),
			&(result.Latitude),
			&(result.Description),
			&(result.Password),
			&(result.HeadUrl),
			&(result.Status),
			&CreatedAt, &UpdatedAt)
		if err != nil {
			return nil, err
		}

		result.CreatedAt = db.TimeParse(CreatedAt)

		result.UpdatedAt = db.TimeParse(UpdatedAt)

		results = append(results, &result)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ezorm.User fetch result error: %v", err)
	}

	return
}

func (m *_UserMgr) Save(obj *User) (sql.Result, error) {
	if obj.Id == 0 {
		return m.saveInsert(obj)
	}
	return m.saveUpdate(obj)
}

func (m *_UserMgr) saveInsert(obj *User) (sql.Result, error) {
	query := "INSERT INTO ezorm.users (`name`, `mailbox`, `sex`, `longitude`, `latitude`, `description`, `password`, `head_url`, `status`, `created_at`, `updated_at`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	result, err := db.MysqlExec(query, obj.Name, obj.Mailbox, obj.Sex, obj.Longitude, obj.Latitude, obj.Description, obj.Password, obj.HeadUrl, obj.Status, db.TimeFormat(obj.CreatedAt), db.TimeFormat(obj.UpdatedAt))
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

func (m *_UserMgr) saveUpdate(obj *User) (sql.Result, error) {
	query := "UPDATE ezorm.users SET `name`=?, `mailbox`=?, `sex`=?, `longitude`=?, `latitude`=?, `description`=?, `password`=?, `head_url`=?, `status`=?, `created_at`=?, `updated_at`=? WHERE `id`=?"
	return db.MysqlExec(query, obj.Name, obj.Mailbox, obj.Sex, obj.Longitude, obj.Latitude, obj.Description, obj.Password, obj.HeadUrl, obj.Status, db.TimeFormat(obj.CreatedAt), db.TimeFormat(obj.UpdatedAt), obj.Id)
}

func (m *_UserMgr) InsertBatch(objs []*User) (sql.Result, error) {
	if len(objs) == 0 {
		return nil, fmt.Errorf("Empty insert")
	}

	values := make([]string, 0, len(objs))
	params := make([]interface{}, 0, len(objs)*11)
	for _, obj := range objs {
		values = append(values, "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		params = append(params, obj.Name, obj.Mailbox, obj.Sex, obj.Longitude, obj.Latitude, obj.Description, obj.Password, obj.HeadUrl, obj.Status, db.TimeFormat(obj.CreatedAt), db.TimeFormat(obj.UpdatedAt))
	}
	query := fmt.Sprintf("INSERT INTO ezorm.users (`name`, `mailbox`, `sex`, `longitude`, `latitude`, `description`, `password`, `head_url`, `status`, `created_at`, `updated_at`) VALUES %s", strings.Join(values, ","))
	return db.MysqlExec(query, params...)
}

func (m *_UserMgr) FindByID(id int32) (*User, error) {
	return m.FindByIDContext(context.Background(), id)
}

func (m *_UserMgr) FindByIDContext(ctx context.Context, id int32) (*User, error) {
	query := "SELECT `id`, `name`, `mailbox`, `sex`, `longitude`, `latitude`, `description`, `password`, `head_url`, `status`, `created_at`, `updated_at` FROM ezorm.users WHERE id=?"
	return m.queryOne(ctx, query, id)
}

func (m *_UserMgr) FindByIDs(ids []int32) ([]*User, error) {
	return m.FindByIDsContext(context.Background(), ids)
}

func (m *_UserMgr) FindByIDsContext(ctx context.Context, ids []int32) ([]*User, error) {
	idsLen := len(ids)
	placeHolders := make([]string, 0, idsLen)
	args := make([]interface{}, 0, idsLen)
	for _, id := range ids {
		placeHolders = append(placeHolders, "?")
		args = append(args, id)
	}

	query := fmt.Sprintf(
		"SELECT `id`, `name`, `mailbox`, `sex`, `longitude`, `latitude`, `description`, `password`, `head_url`, `status`, `created_at`, `updated_at` FROM ezorm.users WHERE id IN (%s)",
		strings.Join(placeHolders, ","))
	return m.query(ctx, query, args...)
}

func (m *_UserMgr) FindInId(ids []int32, sortFields ...string) ([]*User, error) {
	return m.FindInIdContext(context.Background(), ids, sortFields...)
}

func (m *_UserMgr) FindInIdContext(ctx context.Context, ids []int32, sortFields ...string) ([]*User, error) {
	return m.FindByIDsContext(ctx, ids)
}

func (m *_UserMgr) FindListId(Id []int32) ([]*User, error) {
	return m.FindListIdContext(context.Background(), Id)
}

func (m *_UserMgr) FindListIdContext(ctx context.Context, Id []int32) ([]*User, error) {
	retmap, err := m.FindMapIdContext(ctx, Id)
	if err != nil {
		return nil, err
	}
	ret := make([]*User, len(Id))
	for idx, key := range Id {
		ret[idx] = retmap[key]
	}
	return ret, nil
}

func (m *_UserMgr) FindMapId(Id []int32, sortFields ...string) (map[int32]*User, error) {
	return m.FindMapIdContext(context.Background(), Id)
}

func (m *_UserMgr) FindMapIdContext(ctx context.Context, Id []int32, sortFields ...string) (map[int32]*User, error) {
	ret, err := m.FindInIdContext(ctx, Id, sortFields...)
	if err != nil {
		return nil, err
	}
	retmap := make(map[int32]*User, len(ret))
	for _, n := range ret {
		retmap[n.Id] = n
	}
	return retmap, nil
}

func (m *_UserMgr) FindInMailboxPassword(Mailbox []string, Password []string, sortFields ...string) ([]*User, error) {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("SELECT `id`, `name`, `mailbox`, `sex`, `longitude`, `latitude`, `description`, `password`, `head_url`, `status`, `created_at`, `updated_at` FROM ezorm.users WHERE ")

	buf.WriteString("`mailbox` in ")
	stringToIds(buf, Mailbox)
	buf.WriteString(" AND ")

	buf.WriteString("`password` in ")
	stringToIds(buf, Password)
	return m.query(context.Background(), buf.String()+m.GetSort(sortFields))
}

func (m *_UserMgr) FindAllByMailboxPassword(Mailbox string, Password string, sortFields ...string) ([]*User, error) {
	return m.FindByMailboxPassword(Mailbox, Password, -1, -1, sortFields...)
}

func (m *_UserMgr) FindByMailboxPassword(Mailbox string, Password string, offset int, limit int, sortFields ...string) ([]*User, error) {
	query := fmt.Sprintf("SELECT `id`, `name`, `mailbox`, `sex`, `longitude`, `latitude`, `description`, `password`, `head_url`, `status`, `created_at`, `updated_at` FROM ezorm.users WHERE `mailbox`=? AND `password`=? %s%s", m.GetSort(sortFields), m.GetLimit(offset, limit))

	return m.query(context.Background(), query, Mailbox, Password)
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
	buf.WriteString("SELECT `id`, `name`, `mailbox`, `sex`, `longitude`, `latitude`, `description`, `password`, `head_url`, `status`, `created_at`, `updated_at` FROM ezorm.users WHERE ")

	buf.WriteString("`name` in ")
	stringToIds(buf, Name)
	return m.query(context.Background(), buf.String()+m.GetSort(sortFields))
}

func (m *_UserMgr) FindAllByName(Name string, sortFields ...string) ([]*User, error) {
	return m.FindByName(Name, -1, -1, sortFields...)
}

func (m *_UserMgr) FindByName(Name string, offset int, limit int, sortFields ...string) ([]*User, error) {
	query := fmt.Sprintf("SELECT `id`, `name`, `mailbox`, `sex`, `longitude`, `latitude`, `description`, `password`, `head_url`, `status`, `created_at`, `updated_at` FROM ezorm.users WHERE `name`=? %s%s", m.GetSort(sortFields), m.GetLimit(offset, limit))

	return m.query(context.Background(), query, Name)
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
	query := "SELECT `id`, `name`, `mailbox`, `sex`, `longitude`, `latitude`, `description`, `password`, `head_url`, `status`, `created_at`, `updated_at` FROM ezorm.users "

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
	query := "DELETE FROM ezorm.users " + where
	return db.MysqlExec(query, params...)
}

// argument example:
// set:"a=?, b=?"
// where:"c=? and d=?"
// params:[]interface{}{"a", "b", "c", "d"}...
func (m *_UserMgr) Update(set, where string, params ...interface{}) (sql.Result, error) {
	query := fmt.Sprintf("UPDATE ezorm.users SET %s", set)
	if where != "" {
		query = fmt.Sprintf("UPDATE ezorm.users SET %s WHERE %s", set, where)
	}
	return db.MysqlExec(query, params...)
}

func (m *_UserMgr) Count(where string, args ...interface{}) (int32, error) {
	return m.CountContext(context.Background(), where, args...)
}

func (m *_UserMgr) CountContext(ctx context.Context, where string, args ...interface{}) (int32, error) {
	query := "SELECT COUNT(*) FROM ezorm.users"
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

func (m *_UserMgr) GetId2Obj(objs []*User) map[int32]*User {
	id2obj := make(map[int32]*User, len(objs))
	for _, obj := range objs {
		id2obj[obj.Id] = obj
	}
	return id2obj
}

func (m *_UserMgr) GetIds(objs []*User) []int32 {
	ids := make([]int32, len(objs))
	for i, obj := range objs {
		ids[i] = obj.Id
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
