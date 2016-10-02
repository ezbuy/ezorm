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

func (*_BlogMgr) queryLimit(query string, limit int, args ...interface{}) (results []*Blog, err error) {
	rows, err := db.MysqlQuery(query, args...)
	if err != nil {
		return nil, fmt.Errorf("test.Blog query error: %v", err)
	}
	defer rows.Close()

	var Body sql.NullString
	var Create int64
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
		result.Create = time.Unix(Create, 0)

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
	query := "INSERT INTO test.Blog (`Title`, `Hits`, `Slug`, `Body`, `User`, `IsPublished`, `Create`, `Update`) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	result, err := db.MysqlExec(query, obj.Title, obj.Hits, obj.Slug, obj.Body, obj.User, obj.IsPublished, obj.Create.Unix(), db.TimeToLocalTime(obj.Update))
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
	query := "UPDATE test.Blog SET `Title`=?, `Hits`=?, `Slug`=?, `Body`=?, `User`=?, `IsPublished`=?, `Create`=?, `Update`=? WHERE `BlogId`=?"
	return db.MysqlExec(query, obj.Title, obj.Hits, obj.Slug, obj.Body, obj.User, obj.IsPublished, obj.Create.Unix(), db.TimeToLocalTime(obj.Update), obj.BlogId)
}

func (m *_BlogMgr) InsertBatch(objs []*Blog) (sql.Result, error) {
	if len(objs) == 0 {
		return nil, fmt.Errorf("Empty insert")
	}

	values := make([]string, 0, len(objs))
	params := make([]interface{}, 0, len(objs)*8)
	for _, obj := range objs {
		values = append(values, "(?, ?, ?, ?, ?, ?, ?, ?)")
		params = append(params, obj.Title, obj.Hits, obj.Slug, obj.Body, obj.User, obj.IsPublished, obj.Create.Unix(), db.TimeToLocalTime(obj.Update))
	}
	query := fmt.Sprintf("INSERT INTO test.Blog (Title, Hits, Slug, Body, User, IsPublished, Create, Update) VALUES %s", strings.Join(values, ","))
	return db.MysqlExec(query, params...)
}

func (m *_BlogMgr) FindByID(id int32) (*Blog, error) {
	query := "SELECT `BlogId`, `Title`, `Hits`, `Slug`, `Body`, `User`, `IsPublished`, `Create`, `Update` FROM test.Blog WHERE BlogId=?"
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
		"SELECT `BlogId`, `Title`, `Hits`, `Slug`, `Body`, `User`, `IsPublished`, `Create`, `Update` FROM test.Blog WHERE BlogId IN (%s)",
		strings.Join(placeHolders, ","))
	return m.query(query, args...)
}

func (m *_BlogMgr) FindByUserIsPublished(User int32, IsPublished bool, offset int, limit int, sortFields ...string) ([]*Blog, error) {
	orderBy := "ORDER BY %s"
	if len(sortFields) != 0 {
		orderBy = fmt.Sprintf(orderBy, strings.Join(sortFields, ","))
	} else {
		orderBy = fmt.Sprintf(orderBy, "BlogId")
	}

	query := fmt.Sprintf("SELECT `BlogId`, `Title`, `Hits`, `Slug`, `Body`, `User`, `IsPublished`, `Create`, `Update` FROM test.Blog WHERE `User`=? AND `IsPublished`=? %s LIMIT ?, ?", orderBy)

	return m.query(query, User, IsPublished, offset, limit)
}

func (m *_BlogMgr) FindOneBySlug(Slug string) (*Blog, error) {
	query := "SELECT `BlogId`, `Title`, `Hits`, `Slug`, `Body`, `User`, `IsPublished`, `Create`, `Update` FROM test.Blog WHERE Slug=?"
	return m.queryOne(query, Slug)
}

func (m *_BlogMgr) FindByIsPublished(IsPublished bool, offset int, limit int, sortFields ...string) ([]*Blog, error) {
	orderBy := "ORDER BY %s"
	if len(sortFields) != 0 {
		orderBy = fmt.Sprintf(orderBy, strings.Join(sortFields, ","))
	} else {
		orderBy = fmt.Sprintf(orderBy, "BlogId")
	}

	query := fmt.Sprintf("SELECT `BlogId`, `Title`, `Hits`, `Slug`, `Body`, `User`, `IsPublished`, `Create`, `Update` FROM test.Blog WHERE `IsPublished`=? %s LIMIT ?, ?", orderBy)

	return m.query(query, IsPublished, offset, limit)
}

func (m *_BlogMgr) FindByCreate(Create time.Time, offset int, limit int, sortFields ...string) ([]*Blog, error) {
	orderBy := "ORDER BY %s"
	if len(sortFields) != 0 {
		orderBy = fmt.Sprintf(orderBy, strings.Join(sortFields, ","))
	} else {
		orderBy = fmt.Sprintf(orderBy, "BlogId")
	}

	query := fmt.Sprintf("SELECT `BlogId`, `Title`, `Hits`, `Slug`, `Body`, `User`, `IsPublished`, `Create`, `Update` FROM test.Blog WHERE `Create`=? %s LIMIT ?, ?", orderBy)

	return m.query(query, Create.Unix(), offset, limit)
}

func (m *_BlogMgr) FindByUpdate(Update time.Time, offset int, limit int, sortFields ...string) ([]*Blog, error) {
	orderBy := "ORDER BY %s"
	if len(sortFields) != 0 {
		orderBy = fmt.Sprintf(orderBy, strings.Join(sortFields, ","))
	} else {
		orderBy = fmt.Sprintf(orderBy, "BlogId")
	}

	query := fmt.Sprintf("SELECT `BlogId`, `Title`, `Hits`, `Slug`, `Body`, `User`, `IsPublished`, `Create`, `Update` FROM test.Blog WHERE `Update`=? %s LIMIT ?, ?", orderBy)

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
	query := "SELECT `BlogId`, `Title`, `Hits`, `Slug`, `Body`, `User`, `IsPublished`, `Create`, `Update` FROM test.Blog"

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
	query := "DELETE FROM test.Blog " + where
	return db.MysqlExec(query, params...)
}

// argument example:
// set:"a=?, b=?"
// where:"c=? and d=?"
// params:[]interface{}{"a", "b", "c", "d"}...
func (m *_BlogMgr) Update(set, where string, params ...interface{}) (sql.Result, error) {
	query := fmt.Sprintf("UPDATE test.Blog SET %s", set)
	if where != "" {
		query = fmt.Sprintf("UPDATE test.Blog SET %s WHERE %s", set, where)
	}
	return db.MysqlExec(query, params...)
}

func (m *_BlogMgr) Count(where string, args ...interface{}) (int32, error) {
	query := "SELECT COUNT(*) FROM test.Blog"
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
