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

func (m *_BlogMgr) Query(ctx context.Context, query string, args ...interface{}) (results []*Blog, err error) {
	return m.queryLimit(ctx, query, -1, args...)
}

func (*_BlogMgr) queryLimit(ctx context.Context, query string, limit int, args ...interface{}) (results []*Blog, err error) {
	rows, err := db.MysqlQuery(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("test.Blog query error: %w", err)
	}
	defer rows.Close()

	var Body sql.NullString
	var Create string
	var Update string

	offset := 0
	for rows.Next() {
		if limit >= 0 && offset >= limit {
			break
		}
		offset++

		var result Blog
		err := rows.Scan(&(result.BlogId),
			&(result.Title),
			&(result.Hits),
			&(result.Slug),
			&Body, &(result.User),
			&(result.IsPublished),
			&(result.GroupId),
			&Create, &Update)
		if err != nil {
			return nil, err
		}

		result.Body = Body.String
		result.Create = db.TimeParse(Create)

		result.Update = db.TimeParseLocalTime(Update)

		results = append(results, &result)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("test.Blog fetch result error: %w", err)
	}

	return
}

func (m *_BlogMgr) Insert(ctx context.Context, obj *Blog) (sql.Result, error) {
	return m.saveInsert(ctx, obj)
}

func (m *_BlogMgr) UpdateObj(ctx context.Context, obj *Blog) (sql.Result, error) {
	return m.saveUpdate(ctx, obj)
}

func (m *_BlogMgr) Save(ctx context.Context, obj *Blog) (sql.Result, error) {
	// upsert
	result, err := m.saveUpdate(ctx, obj)
	if err != nil {
		return nil, err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if n != 0 {
		return result, nil
	}
	return m.saveInsert(ctx, obj)

}

func (m *_BlogMgr) saveInsert(ctx context.Context, obj *Blog) (sql.Result, error) {
	if obj.BlogId == 0 {
		return nil, fmt.Errorf("missing Id: BlogId")
	}

	query := "INSERT INTO test.blog (`blog_id`, `title`, `hits`, `slug`, `body`, `user`, `is_published`, `group_id`, `create`, `update`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	result, err := db.MysqlExec(ctx, query, obj.BlogId, obj.Title, obj.Hits, obj.Slug, obj.Body, obj.User, obj.IsPublished, obj.GroupId, db.TimeFormat(obj.Create), db.TimeToLocalTime(obj.Update))
	if err != nil {
		return result, err
	}

	return result, err
}

func (m *_BlogMgr) saveUpdate(ctx context.Context, obj *Blog) (sql.Result, error) {
	query := "UPDATE test.blog SET `title`=?, `hits`=?, `slug`=?, `body`=?, `user`=?, `is_published`=?, `group_id`=?, `create`=?, `update`=? WHERE `blog_id`=?"
	return db.MysqlExec(ctx, query, obj.Title, obj.Hits, obj.Slug, obj.Body, obj.User, obj.IsPublished, obj.GroupId, db.TimeFormat(obj.Create), db.TimeToLocalTime(obj.Update), obj.BlogId)
}

func (m *_BlogMgr) InsertBatch(ctx context.Context, objs []*Blog) (sql.Result, error) {
	if len(objs) == 0 {
		return nil, fmt.Errorf("Empty insert")
	}

	values := make([]string, 0, len(objs))
	params := make([]interface{}, 0, len(objs)*9)
	for _, obj := range objs {
		values = append(values, "(?, ?, ?, ?, ?, ?, ?, ?, ?)")
		params = append(params, obj.Title, obj.Hits, obj.Slug, obj.Body, obj.User, obj.IsPublished, obj.GroupId, db.TimeFormat(obj.Create), db.TimeToLocalTime(obj.Update))
	}
	query := fmt.Sprintf("INSERT INTO test.blog (`title`, `hits`, `slug`, `body`, `user`, `is_published`, `group_id`, `create`, `update`) VALUES %s", strings.Join(values, ","))
	return db.MysqlExec(ctx, query, params...)
}

func (m *_BlogMgr) FindByID(ctx context.Context, id int32) (*Blog, error) {
	query := "SELECT `blog_id`, `title`, `hits`, `slug`, `body`, `user`, `is_published`, `group_id`, `create`, `update` FROM test.blog WHERE blog_id=?"
	return m.queryOne(ctx, query, id)
}

func (m *_BlogMgr) FindByIDs(ctx context.Context, ids []int32) ([]*Blog, error) {
	idsLen := len(ids)
	placeHolders := make([]string, 0, idsLen)
	args := make([]interface{}, 0, idsLen)
	for _, id := range ids {
		placeHolders = append(placeHolders, "?")
		args = append(args, id)
	}

	query := fmt.Sprintf(
		"SELECT `blog_id`, `title`, `hits`, `slug`, `body`, `user`, `is_published`, `group_id`, `create`, `update` FROM test.blog WHERE blog_id IN (%s)",
		strings.Join(placeHolders, ","))
	return m.query(ctx, query, args...)
}

func (m *_BlogMgr) FindInBlogId(ctx context.Context, ids []int32, sortFields ...string) ([]*Blog, error) {
	return m.FindByIDs(ctx, ids)
}

func (m *_BlogMgr) FindListBlogId(ctx context.Context, BlogId []int32) ([]*Blog, error) {
	retmap, err := m.FindMapBlogId(ctx, BlogId)
	if err != nil {
		return nil, err
	}
	ret := make([]*Blog, len(BlogId))
	for idx, key := range BlogId {
		ret[idx] = retmap[key]
	}
	return ret, nil
}

func (m *_BlogMgr) FindMapBlogId(ctx context.Context, BlogId []int32, sortFields ...string) (map[int32]*Blog, error) {
	ret, err := m.FindInBlogId(ctx, BlogId, sortFields...)
	if err != nil {
		return nil, err
	}
	retmap := make(map[int32]*Blog, len(ret))
	for _, n := range ret {
		retmap[n.BlogId] = n
	}
	return retmap, nil
}

func (m *_BlogMgr) FindAllByUserIsPublished(ctx context.Context, User int32, IsPublished bool, sortFields ...string) ([]*Blog, error) {
	return m.FindByUserIsPublished(ctx, User, IsPublished, -1, -1, sortFields...)
}

func (m *_BlogMgr) FindByUserIsPublished(ctx context.Context, User int32, IsPublished bool, offset int, limit int, sortFields ...string) ([]*Blog, error) {
	query := fmt.Sprintf("SELECT `blog_id`, `title`, `hits`, `slug`, `body`, `user`, `is_published`, `group_id`, `create`, `update` FROM test.blog WHERE `user`=? AND `is_published`=? %s%s", m.GetSort(sortFields), m.GetLimit(offset, limit))

	return m.query(ctx, query, User, IsPublished)
}

func (m *_BlogMgr) FindListSlug(ctx context.Context, Slug []string) ([]*Blog, error) {
	retmap, err := m.FindMapSlug(ctx, Slug)
	if err != nil {
		return nil, err
	}
	ret := make([]*Blog, len(Slug))
	for idx, key := range Slug {
		ret[idx] = retmap[key]
	}
	return ret, nil
}

func (m *_BlogMgr) FindMapSlug(ctx context.Context, Slug []string) (map[string]*Blog, error) {
	ret, err := m.FindInSlug(ctx, Slug)
	if err != nil {
		return nil, err
	}
	retmap := make(map[string]*Blog, len(ret))
	for _, n := range ret {
		retmap[n.Slug] = n
	}
	return retmap, nil
}

func (m *_BlogMgr) FindInSlug(ctx context.Context, Slug []string, sortFields ...string) ([]*Blog, error) {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("SELECT `blog_id`, `title`, `hits`, `slug`, `body`, `user`, `is_published`, `group_id`, `create`, `update` FROM test.blog WHERE ")

	buf.WriteString("`slug` in ")
	stringToIds(buf, Slug)
	return m.query(ctx, buf.String()+m.GetSort(sortFields))
}

func (m *_BlogMgr) FindOneBySlug(ctx context.Context, Slug string) (*Blog, error) {
	query := "SELECT `blog_id`, `title`, `hits`, `slug`, `body`, `user`, `is_published`, `group_id`, `create`, `update` FROM test.blog WHERE slug=?"
	return m.queryOne(ctx, query, Slug)
}

func (m *_BlogMgr) FindListUser(ctx context.Context, User []int32) ([]*Blog, error) {
	retmap, err := m.FindMapUser(ctx, User)
	if err != nil {
		return nil, err
	}
	ret := make([]*Blog, len(User))
	for idx, key := range User {
		ret[idx] = retmap[key]
	}
	return ret, nil
}

func (m *_BlogMgr) FindMapUser(ctx context.Context, User []int32) (map[int32]*Blog, error) {
	ret, err := m.FindInUser(ctx, User)
	if err != nil {
		return nil, err
	}
	retmap := make(map[int32]*Blog, len(ret))
	for _, n := range ret {
		retmap[n.User] = n
	}
	return retmap, nil
}

func (m *_BlogMgr) FindInUser(ctx context.Context, User []int32, sortFields ...string) ([]*Blog, error) {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("SELECT `blog_id`, `title`, `hits`, `slug`, `body`, `user`, `is_published`, `group_id`, `create`, `update` FROM test.blog WHERE ")

	buf.WriteString("`user` in ")
	int32ToIds(buf, User)
	return m.query(ctx, buf.String()+m.GetSort(sortFields))
}

func (m *_BlogMgr) FindAllByUser(ctx context.Context, User int32, sortFields ...string) ([]*Blog, error) {
	return m.FindByUser(ctx, User, -1, -1, sortFields...)
}

func (m *_BlogMgr) FindByUser(ctx context.Context, User int32, offset int, limit int, sortFields ...string) ([]*Blog, error) {
	query := fmt.Sprintf("SELECT `blog_id`, `title`, `hits`, `slug`, `body`, `user`, `is_published`, `group_id`, `create`, `update` FROM test.blog WHERE `user`=? %s%s", m.GetSort(sortFields), m.GetLimit(offset, limit))

	return m.query(ctx, query, User)
}

func (m *_BlogMgr) FindAllByIsPublished(ctx context.Context, IsPublished bool, sortFields ...string) ([]*Blog, error) {
	return m.FindByIsPublished(ctx, IsPublished, -1, -1, sortFields...)
}

func (m *_BlogMgr) FindByIsPublished(ctx context.Context, IsPublished bool, offset int, limit int, sortFields ...string) ([]*Blog, error) {
	query := fmt.Sprintf("SELECT `blog_id`, `title`, `hits`, `slug`, `body`, `user`, `is_published`, `group_id`, `create`, `update` FROM test.blog WHERE `is_published`=? %s%s", m.GetSort(sortFields), m.GetLimit(offset, limit))

	return m.query(ctx, query, IsPublished)
}

func (m *_BlogMgr) FindListGroupId(ctx context.Context, GroupId []int64) ([]*Blog, error) {
	retmap, err := m.FindMapGroupId(ctx, GroupId)
	if err != nil {
		return nil, err
	}
	ret := make([]*Blog, len(GroupId))
	for idx, key := range GroupId {
		ret[idx] = retmap[key]
	}
	return ret, nil
}

func (m *_BlogMgr) FindMapGroupId(ctx context.Context, GroupId []int64) (map[int64]*Blog, error) {
	ret, err := m.FindInGroupId(ctx, GroupId)
	if err != nil {
		return nil, err
	}
	retmap := make(map[int64]*Blog, len(ret))
	for _, n := range ret {
		retmap[n.GroupId] = n
	}
	return retmap, nil
}

func (m *_BlogMgr) FindInGroupId(ctx context.Context, GroupId []int64, sortFields ...string) ([]*Blog, error) {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("SELECT `blog_id`, `title`, `hits`, `slug`, `body`, `user`, `is_published`, `group_id`, `create`, `update` FROM test.blog WHERE ")

	buf.WriteString("`group_id` in ")
	int64ToIds(buf, GroupId)
	return m.query(ctx, buf.String()+m.GetSort(sortFields))
}

func (m *_BlogMgr) FindAllByGroupId(ctx context.Context, GroupId int64, sortFields ...string) ([]*Blog, error) {
	return m.FindByGroupId(ctx, GroupId, -1, -1, sortFields...)
}

func (m *_BlogMgr) FindByGroupId(ctx context.Context, GroupId int64, offset int, limit int, sortFields ...string) ([]*Blog, error) {
	query := fmt.Sprintf("SELECT `blog_id`, `title`, `hits`, `slug`, `body`, `user`, `is_published`, `group_id`, `create`, `update` FROM test.blog WHERE `group_id`=? %s%s", m.GetSort(sortFields), m.GetLimit(offset, limit))

	return m.query(ctx, query, GroupId)
}

func (m *_BlogMgr) FindAllByCreate(ctx context.Context, Create time.Time, sortFields ...string) ([]*Blog, error) {
	return m.FindByCreate(ctx, Create, -1, -1, sortFields...)
}

func (m *_BlogMgr) FindByCreate(ctx context.Context, Create time.Time, offset int, limit int, sortFields ...string) ([]*Blog, error) {
	query := fmt.Sprintf("SELECT `blog_id`, `title`, `hits`, `slug`, `body`, `user`, `is_published`, `group_id`, `create`, `update` FROM test.blog WHERE `create`=? %s%s", m.GetSort(sortFields), m.GetLimit(offset, limit))

	return m.query(ctx, query, db.TimeFormat(Create))
}

func (m *_BlogMgr) FindAllByUpdate(ctx context.Context, Update time.Time, sortFields ...string) ([]*Blog, error) {
	return m.FindByUpdate(ctx, Update, -1, -1, sortFields...)
}

func (m *_BlogMgr) FindByUpdate(ctx context.Context, Update time.Time, offset int, limit int, sortFields ...string) ([]*Blog, error) {
	query := fmt.Sprintf("SELECT `blog_id`, `title`, `hits`, `slug`, `body`, `user`, `is_published`, `group_id`, `create`, `update` FROM test.blog WHERE `update`=? %s%s", m.GetSort(sortFields), m.GetLimit(offset, limit))

	return m.query(ctx, query, db.TimeToLocalTime(Update))
}

func (m *_BlogMgr) FindOne(ctx context.Context, where string, args ...interface{}) (*Blog, error) {
	query := m.GetQuerysql(where) + m.GetLimit(0, 1)
	return m.queryOne(ctx, query, args...)
}

func (m *_BlogMgr) Find(ctx context.Context, where string, args ...interface{}) ([]*Blog, error) {
	query := m.GetQuerysql(where)
	return m.query(ctx, query, args...)
}

func (m *_BlogMgr) FindAll(ctx context.Context) (results []*Blog, err error) {
	return m.Find(ctx, "")
}

func (m *_BlogMgr) FindWithOffset(ctx context.Context, where string, offset int, limit int, args ...interface{}) ([]*Blog, error) {
	query := m.GetQuerysql(where)

	query = query + " LIMIT ?, ?"

	args = append(args, offset)
	args = append(args, limit)

	return m.query(ctx, query, args...)
}

func (m *_BlogMgr) GetQuerysql(where string) string {
	query := "SELECT `blog_id`, `title`, `hits`, `slug`, `body`, `user`, `is_published`, `group_id`, `create`, `update` FROM test.blog "

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

func (m *_BlogMgr) Del(ctx context.Context, where string, params ...interface{}) (sql.Result, error) {
	if where != "" {
		where = "WHERE " + where
	}
	query := "DELETE FROM test.blog " + where
	return db.MysqlExec(ctx, query, params...)
}

// argument example:
// set:"a=?, b=?"
// where:"c=? and d=?"
// params:[]interface{}{"a", "b", "c", "d"}...
func (m *_BlogMgr) Update(ctx context.Context, set, where string, params ...interface{}) (sql.Result, error) {
	query := fmt.Sprintf("UPDATE test.blog SET %s", set)
	if where != "" {
		query = fmt.Sprintf("UPDATE test.blog SET %s WHERE %s", set, where)
	}
	return db.MysqlExec(ctx, query, params...)
}

func (m *_BlogMgr) Count(ctx context.Context, where string, args ...interface{}) (int32, error) {
	query := "SELECT COUNT(*) FROM test.blog"
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
		id2obj[obj.BlogId] = obj
	}
	return id2obj
}

func (m *_BlogMgr) GetIds(objs []*Blog) []int32 {
	ids := make([]int32, len(objs))
	for i, obj := range objs {
		ids[i] = obj.BlogId
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
