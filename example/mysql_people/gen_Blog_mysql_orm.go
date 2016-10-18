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

func (m *_BlogMgr) queryOne(query string, args ...interface{}) (*Blog, error) {
	ret, err := m.queryLimit(query, 1, args...)
	if err != nil {
		return nil, err
	}
	if len(ret) == 0 {
		return nil, sql.ErrNoRows
	}
	return ret[0], nil
}

func (m *_BlogMgr) query(query string, args ...interface{}) (results []*Blog, err error) {
	return m.queryLimit(query, -1, args...)
}

func (m *_BlogMgr) Query(query string, args ...interface{}) (results []*Blog, err error) {
	return m.queryLimit(query, -1, args...)
}

func (*_BlogMgr) queryLimit(query string, limit int, args ...interface{}) (results []*Blog, err error) {
	rows, err := db.MysqlQuery(query, args...)
	if err != nil {
		return nil, fmt.Errorf("test.Blog query error: %v", err)
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
		return nil, fmt.Errorf("test.Blog fetch result error: %v", err)
	}

	return
}

func (m *_BlogMgr) Save(obj *Blog) (sql.Result, error) {
	if obj.BlogId == 0 {
		return m.saveInsert(obj)
	}
	return m.saveUpdate(obj)
}

func (m *_BlogMgr) saveInsert(obj *Blog) (sql.Result, error) {
	query := "INSERT INTO test.blog (`title`, `hits`, `slug`, `body`, `user`, `is_published`, `create`, `update`) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	result, err := db.MysqlExec(query, obj.Title, obj.Hits, obj.Slug, obj.Body, obj.User, obj.IsPublished, db.TimeFormat(obj.Create), db.TimeToLocalTime(obj.Update))
	if err != nil {
		return result, err
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return result, err
	}

	obj.BlogId = int32(lastInsertId)

	return result, err
}

func (m *_BlogMgr) saveUpdate(obj *Blog) (sql.Result, error) {
	query := "UPDATE test.blog SET `title`=?, `hits`=?, `slug`=?, `body`=?, `user`=?, `is_published`=?, `create`=?, `update`=? WHERE `blog_id`=?"
	return db.MysqlExec(query, obj.Title, obj.Hits, obj.Slug, obj.Body, obj.User, obj.IsPublished, db.TimeFormat(obj.Create), db.TimeToLocalTime(obj.Update), obj.BlogId)
}

func (m *_BlogMgr) InsertBatch(objs []*Blog) (sql.Result, error) {
	if len(objs) == 0 {
		return nil, fmt.Errorf("Empty insert")
	}

	values := make([]string, 0, len(objs))
	params := make([]interface{}, 0, len(objs)*8)
	for _, obj := range objs {
		values = append(values, "(?, ?, ?, ?, ?, ?, ?, ?)")
		params = append(params, obj.Title, obj.Hits, obj.Slug, obj.Body, obj.User, obj.IsPublished, db.TimeFormat(obj.Create), db.TimeToLocalTime(obj.Update))
	}
	query := fmt.Sprintf("INSERT INTO test.blog (`title`, `hits`, `slug`, `body`, `user`, `is_published`, `create`, `update`) VALUES %s", strings.Join(values, ","))
	return db.MysqlExec(query, params...)
}

func (m *_BlogMgr) FindByID(id int32) (*Blog, error) {
	query := "SELECT `blog_id`, `title`, `hits`, `slug`, `body`, `user`, `is_published`, `create`, `update` FROM test.blog WHERE blog_id=?"
	return m.queryOne(query, id)
}

func (m *_BlogMgr) FindByIDs(ids []int32) ([]*Blog, error) {
	idsLen := len(ids)
	placeHolders := make([]string, 0, idsLen)
	args := make([]interface{}, 0, idsLen)
	for _, id := range ids {
		placeHolders = append(placeHolders, "?")
		args = append(args, id)
	}

	query := fmt.Sprintf(
		"SELECT `blog_id`, `title`, `hits`, `slug`, `body`, `user`, `is_published`, `create`, `update` FROM test.blog WHERE blog_id IN (%s)",
		strings.Join(placeHolders, ","))
	return m.query(query, args...)
}

func (m *_BlogMgr) FindInBlogId(ids []int32, sortFields ...string) ([]*Blog, error) {
	return m.FindByIDs(ids)
}

func (m *_BlogMgr) FindListBlogId(BlogId []int32) ([]*Blog, error) {
	retmap, err := m.FindMapBlogId(BlogId)
	if err != nil {
		return nil, err
	}
	ret := make([]*Blog, len(BlogId))
	for idx, key := range BlogId {
		ret[idx] = retmap[key]
	}
	return ret, nil
}

func (m *_BlogMgr) FindMapBlogId(BlogId []int32, sortFields ...string) (map[int32]*Blog, error) {
	ret, err := m.FindInBlogId(BlogId, sortFields...)
	if err != nil {
		return nil, err
	}
	retmap := make(map[int32]*Blog, len(ret))
	for _, n := range ret {
		retmap[n.BlogId] = n
	}
	return retmap, nil
}

func (m *_BlogMgr) FindAllByUserIsPublished(User int32, IsPublished bool, sortFields ...string) ([]*Blog, error) {
	return m.FindByUserIsPublished(User, IsPublished, -1, -1, sortFields...)
}

func (m *_BlogMgr) FindByUserIsPublished(User int32, IsPublished bool, offset int, limit int, sortFields ...string) ([]*Blog, error) {
	orderBy := "ORDER BY %s"
	if len(sortFields) != 0 {
		orderBy = fmt.Sprintf(orderBy, strings.Join(sortFields, ","))
	} else {
		orderBy = fmt.Sprintf(orderBy, "blog_id")
	}

	query := fmt.Sprintf("SELECT `blog_id`, `title`, `hits`, `slug`, `body`, `user`, `is_published`, `create`, `update` FROM test.blog WHERE `user`=? AND `is_published`=? %s LIMIT ?, ?", orderBy)

	return m.query(query, User, IsPublished, offset, limit)
}

func (m *_BlogMgr) FindListSlug(Slug []string) ([]*Blog, error) {
	retmap, err := m.FindMapSlug(Slug)
	if err != nil {
		return nil, err
	}
	ret := make([]*Blog, len(Slug))
	for idx, key := range Slug {
		ret[idx] = retmap[key]
	}
	return ret, nil
}

func (m *_BlogMgr) FindMapSlug(Slug []string) (map[string]*Blog, error) {
	ret, err := m.FindInSlug(Slug)
	if err != nil {
		return nil, err
	}
	retmap := make(map[string]*Blog, len(ret))
	for _, n := range ret {
		retmap[n.Slug] = n
	}
	return retmap, nil
}

func (m *_BlogMgr) FindInSlug(Slug []string, sortFields ...string) ([]*Blog, error) {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("SELECT `blog_id`, `title`, `hits`, `slug`, `body`, `user`, `is_published`, `create`, `update` FROM test.blog WHERE ")

	buf.WriteString("`slug` in ")
	stringToIds(buf, Slug)
	return m.query(buf.String())
}

func (m *_BlogMgr) FindOneBySlug(Slug string) (*Blog, error) {
	query := "SELECT `blog_id`, `title`, `hits`, `slug`, `body`, `user`, `is_published`, `create`, `update` FROM test.blog WHERE slug=?"
	return m.queryOne(query, Slug)
}

func (m *_BlogMgr) FindListUser(User []int32) ([]*Blog, error) {
	retmap, err := m.FindMapUser(User)
	if err != nil {
		return nil, err
	}
	ret := make([]*Blog, len(User))
	for idx, key := range User {
		ret[idx] = retmap[key]
	}
	return ret, nil
}

func (m *_BlogMgr) FindMapUser(User []int32) (map[int32]*Blog, error) {
	ret, err := m.FindInUser(User)
	if err != nil {
		return nil, err
	}
	retmap := make(map[int32]*Blog, len(ret))
	for _, n := range ret {
		retmap[n.User] = n
	}
	return retmap, nil
}

func (m *_BlogMgr) FindInUser(User []int32, sortFields ...string) ([]*Blog, error) {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("SELECT `blog_id`, `title`, `hits`, `slug`, `body`, `user`, `is_published`, `create`, `update` FROM test.blog WHERE ")

	buf.WriteString("`user` in ")
	int32ToIds(buf, User)
	return m.query(buf.String())
}

func (m *_BlogMgr) FindAllByUser(User int32, sortFields ...string) ([]*Blog, error) {
	return m.FindByUser(User, -1, -1, sortFields...)
}

func (m *_BlogMgr) FindByUser(User int32, offset int, limit int, sortFields ...string) ([]*Blog, error) {
	orderBy := "ORDER BY %s"
	if len(sortFields) != 0 {
		orderBy = fmt.Sprintf(orderBy, strings.Join(sortFields, ","))
	} else {
		orderBy = fmt.Sprintf(orderBy, "blog_id")
	}

	query := fmt.Sprintf("SELECT `blog_id`, `title`, `hits`, `slug`, `body`, `user`, `is_published`, `create`, `update` FROM test.blog WHERE `user`=? %s LIMIT ?, ?", orderBy)

	return m.query(query, User, offset, limit)
}

func (m *_BlogMgr) FindAllByIsPublished(IsPublished bool, sortFields ...string) ([]*Blog, error) {
	return m.FindByIsPublished(IsPublished, -1, -1, sortFields...)
}

func (m *_BlogMgr) FindByIsPublished(IsPublished bool, offset int, limit int, sortFields ...string) ([]*Blog, error) {
	orderBy := "ORDER BY %s"
	if len(sortFields) != 0 {
		orderBy = fmt.Sprintf(orderBy, strings.Join(sortFields, ","))
	} else {
		orderBy = fmt.Sprintf(orderBy, "blog_id")
	}

	query := fmt.Sprintf("SELECT `blog_id`, `title`, `hits`, `slug`, `body`, `user`, `is_published`, `create`, `update` FROM test.blog WHERE `is_published`=? %s LIMIT ?, ?", orderBy)

	return m.query(query, IsPublished, offset, limit)
}

func (m *_BlogMgr) FindAllByCreate(Create time.Time, sortFields ...string) ([]*Blog, error) {
	return m.FindByCreate(Create, -1, -1, sortFields...)
}

func (m *_BlogMgr) FindByCreate(Create time.Time, offset int, limit int, sortFields ...string) ([]*Blog, error) {
	orderBy := "ORDER BY %s"
	if len(sortFields) != 0 {
		orderBy = fmt.Sprintf(orderBy, strings.Join(sortFields, ","))
	} else {
		orderBy = fmt.Sprintf(orderBy, "blog_id")
	}

	query := fmt.Sprintf("SELECT `blog_id`, `title`, `hits`, `slug`, `body`, `user`, `is_published`, `create`, `update` FROM test.blog WHERE `create`=? %s LIMIT ?, ?", orderBy)

	return m.query(query, db.TimeFormat(Create), offset, limit)
}

func (m *_BlogMgr) FindAllByUpdate(Update time.Time, sortFields ...string) ([]*Blog, error) {
	return m.FindByUpdate(Update, -1, -1, sortFields...)
}

func (m *_BlogMgr) FindByUpdate(Update time.Time, offset int, limit int, sortFields ...string) ([]*Blog, error) {
	orderBy := "ORDER BY %s"
	if len(sortFields) != 0 {
		orderBy = fmt.Sprintf(orderBy, strings.Join(sortFields, ","))
	} else {
		orderBy = fmt.Sprintf(orderBy, "blog_id")
	}

	query := fmt.Sprintf("SELECT `blog_id`, `title`, `hits`, `slug`, `body`, `user`, `is_published`, `create`, `update` FROM test.blog WHERE `update`=? %s LIMIT ?, ?", orderBy)

	return m.query(query, db.TimeToLocalTime(Update), offset, limit)
}

func (m *_BlogMgr) FindOne(where string, args ...interface{}) (*Blog, error) {
	query := m.getQuerysql(true, where)
	return m.queryOne(query, args...)
}

func (m *_BlogMgr) Find(where string, args ...interface{}) ([]*Blog, error) {
	query := m.getQuerysql(false, where)
	return m.query(query, args...)
}

func (m *_BlogMgr) FindAll() (results []*Blog, err error) {
	return m.Find("")
}

func (m *_BlogMgr) FindWithOffset(where string, offset int, limit int, args ...interface{}) ([]*Blog, error) {
	query := m.getQuerysql(false, where)

	query = query + " LIMIT ?, ?"

	args = append(args, offset)
	args = append(args, limit)

	return m.query(query, args...)
}

func (m *_BlogMgr) getQuerysql(topOne bool, where string) string {
	query := "SELECT `blog_id`, `title`, `hits`, `slug`, `body`, `user`, `is_published`, `create`, `update` FROM test.blog"

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

func (m *_BlogMgr) Del(where string, params ...interface{}) (sql.Result, error) {
	if where != "" {
		where = "WHERE " + where
	}
	query := "DELETE FROM test.blog " + where
	return db.MysqlExec(query, params...)
}

// argument example:
// set:"a=?, b=?"
// where:"c=? and d=?"
// params:[]interface{}{"a", "b", "c", "d"}...
func (m *_BlogMgr) Update(set, where string, params ...interface{}) (sql.Result, error) {
	query := fmt.Sprintf("UPDATE test.blog SET %s", set)
	if where != "" {
		query = fmt.Sprintf("UPDATE test.blog SET %s WHERE %s", set, where)
	}
	return db.MysqlExec(query, params...)
}

func (m *_BlogMgr) Count(where string, args ...interface{}) (int32, error) {
	query := "SELECT COUNT(*) FROM test.blog"
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

func (m *_BlogMgr) getLimitQuery(offset, limit int, sorts []string) string {
	orderBy := ""
	if len(sorts) != 0 {
		orderBy = fmt.Sprintf("ORDER BY %s", strings.Join(sorts, ","))
	}
	if limit > 0 {
		return orderBy + fmt.Sprintf(" LIMIT %d, %d", offset, limit)
	}
	return orderBy
}
