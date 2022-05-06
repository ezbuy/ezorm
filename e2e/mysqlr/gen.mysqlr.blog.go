package mysqlr

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/ezbuy/ezorm/v2/pkg/orm"
	"gopkg.in/go-playground/validator.v9"
)

var (
	_ sql.DB
	_ time.Time
	_ fmt.Formatter
	_ strings.Reader
	_ validator.Validate
	_ context.Context
)

type Blog struct {
	Id        int64     `mysql:"id"`
	UserId    int32     `mysql:"user_id"`
	Title     string    `mysql:"title"`
	Content   string    `mysql:"content"`
	Status    int32     `mysql:"status"`
	Readed    int32     `mysql:"readed"`
	CreatedAt time.Time `mysql:"created_at"`
	UpdatedAt time.Time `mysql:"updated_at"`
}

var BlogColumns = struct {
	Id        string
	UserId    string
	Title     string
	Content   string
	Status    string
	Readed    string
	CreatedAt string
	UpdatedAt string
}{
	"id",
	"user_id",
	"title",
	"content",
	"status",
	"readed",
	"created_at",
	"updated_at",
}

type _BlogMgr struct {
}

var BlogMgr *_BlogMgr

func (m *_BlogMgr) NewBlog() *Blog {
	return &Blog{}
}

func (obj *Blog) GetNameSpace() string {
	return "mysqlr"
}

func (obj *Blog) GetClassName() string {
	return "Blog"
}

func (obj *Blog) GetTableName() string {
	return "blogs"
}

func (obj *Blog) GetColumns() []string {
	columns := []string{
		"blogs.id",
		"blogs.user_id",
		"blogs.title",
		"blogs.content",
		"blogs.status",
		"blogs.readed",
		"blogs.created_at",
		"blogs.updated_at",
	}
	return columns
}

func (obj *Blog) GetNoneIncrementColumns() []string {
	columns := []string{
		"id",
		"user_id",
		"title",
		"content",
		"status",
		"readed",
		"created_at",
		"updated_at",
	}
	return columns
}

func (obj *Blog) GetPrimaryKey() PrimaryKey {
	pk := BlogMgr.NewPrimaryKey()
	pk.Id = obj.Id
	pk.UserId = obj.UserId
	return pk
}

func (obj *Blog) Validate() error {
	validate := validator.New()
	return validate.Struct(obj)
}

type IdUserIdOfBlogPK struct {
	Id     int64
	UserId int32
}

func (m *_BlogMgr) NewPrimaryKey() *IdUserIdOfBlogPK {
	return &IdUserIdOfBlogPK{}
}

func (u *IdUserIdOfBlogPK) Key() string {
	strs := []string{
		"Id",
		fmt.Sprint(u.Id),
		"UserId",
		fmt.Sprint(u.UserId),
	}
	return strings.Join(strs, ":")
}

func (u *IdUserIdOfBlogPK) Parse(key string) error {
	arr := strings.Split(key, ":")
	if len(arr)%2 != 0 {
		return fmt.Errorf("key (%s) format error", key)
	}
	kv := map[string]string{}
	for i := 0; i < len(arr)/2; i++ {
		kv[arr[2*i]] = arr[2*i+1]
	}
	vId, ok := kv["Id"]
	if !ok {
		return fmt.Errorf("key (%s) without (Id) field", key)
	}
	if err := orm.StringScan(vId, &(u.Id)); err != nil {
		return err
	}
	vUserId, ok := kv["UserId"]
	if !ok {
		return fmt.Errorf("key (%s) without (UserId) field", key)
	}
	if err := orm.StringScan(vUserId, &(u.UserId)); err != nil {
		return err
	}
	return nil
}

func (u *IdUserIdOfBlogPK) SQLFormat() string {
	conditions := []string{
		"id = ?",
		"user_id = ?",
	}
	return orm.SQLWhere(conditions)
}

func (u *IdUserIdOfBlogPK) SQLParams() []interface{} {
	return []interface{}{
		u.Id,
		u.UserId,
	}
}

func (u *IdUserIdOfBlogPK) Columns() []string {
	return []string{
		"id",
		"user_id",
	}
}

type StatusOfBlogIDX struct {
	Status int32
	offset int
	limit  int
}

func (u *StatusOfBlogIDX) Key() string {
	strs := []string{
		"Status",
		fmt.Sprint(u.Status),
	}
	return strings.Join(strs, ":")
}

func (u *StatusOfBlogIDX) SQLFormat(limit bool) string {
	conditions := []string{
		"status = ?",
	}
	if limit {
		return fmt.Sprintf("%s %s", orm.SQLWhere(conditions), orm.SQLOffsetLimit(u.offset, u.limit))
	}
	return orm.SQLWhere(conditions)
}

func (u *StatusOfBlogIDX) SQLParams() []interface{} {
	return []interface{}{
		u.Status,
	}
}

func (u *StatusOfBlogIDX) SQLLimit() int {
	if u.limit > 0 {
		return u.limit
	}
	return -1
}

func (u *StatusOfBlogIDX) Limit(n int) {
	u.limit = n
}

func (u *StatusOfBlogIDX) Offset(n int) {
	u.offset = n
}

func (u *StatusOfBlogIDX) PositionOffsetLimit(len int) (int, int) {
	if u.limit <= 0 {
		return 0, len
	}
	if u.offset+u.limit > len {
		return u.offset, len
	}
	return u.offset, u.limit
}

type _BlogDBMgr struct {
	db orm.DB
}

func (m *_BlogMgr) DB(db orm.DB) *_BlogDBMgr {
	return BlogDBMgr(db)
}

func BlogDBMgr(db orm.DB) *_BlogDBMgr {
	if db == nil {
		panic(fmt.Errorf("BlogDBMgr init need db"))
	}
	return &_BlogDBMgr{db: db}
}

func (m *_BlogDBMgr) Search(ctx context.Context, where string, orderby string, limit string, args ...interface{}) ([]*Blog, error) {
	obj := BlogMgr.NewBlog()

	if limit = strings.ToUpper(strings.TrimSpace(limit)); limit != "" && !strings.HasPrefix(limit, "LIMIT") {
		limit = "LIMIT " + limit
	}

	conditions := []string{where, orderby, limit}
	query := fmt.Sprintf("SELECT %s FROM blogs %s", strings.Join(obj.GetColumns(), ","), strings.Join(conditions, " "))
	return m.FetchBySQL(ctx, query, args...)
}

func (m *_BlogDBMgr) SearchConditions(ctx context.Context, conditions []string, orderby string, offset int, limit int, args ...interface{}) ([]*Blog, error) {
	obj := BlogMgr.NewBlog()
	q := fmt.Sprintf("SELECT %s FROM blogs %s %s %s",
		strings.Join(obj.GetColumns(), ","),
		orm.SQLWhere(conditions),
		orderby,
		orm.SQLOffsetLimit(offset, limit))

	return m.FetchBySQL(ctx, q, args...)
}

func (m *_BlogDBMgr) SearchCount(ctx context.Context, where string, args ...interface{}) (int64, error) {
	return m.queryCount(ctx, where, args...)
}

func (m *_BlogDBMgr) SearchConditionsCount(ctx context.Context, conditions []string, args ...interface{}) (int64, error) {
	return m.queryCount(ctx, orm.SQLWhere(conditions), args...)
}

func (m *_BlogDBMgr) FetchBySQL(ctx context.Context, q string, args ...interface{}) (results []*Blog, err error) {
	rows, err := m.db.Query(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("Blog fetch error: %v", err)
	}
	defer rows.Close()

	var CreatedAt string
	var UpdatedAt int64

	for rows.Next() {
		var result Blog
		err = rows.Scan(&(result.Id), &(result.UserId), &(result.Title), &(result.Content), &(result.Status), &(result.Readed), &CreatedAt, &UpdatedAt)
		if err != nil {
			m.db.SetError(err)
			return nil, err
		}

		result.CreatedAt = orm.TimeParse(CreatedAt)
		result.UpdatedAt = time.Unix(UpdatedAt, 0)

		results = append(results, &result)
	}
	if err = rows.Err(); err != nil {
		m.db.SetError(err)
		return nil, fmt.Errorf("Blog fetch result error: %v", err)
	}
	return
}
func (m *_BlogDBMgr) Exist(ctx context.Context, pk PrimaryKey) (bool, error) {
	c, err := m.queryCount(ctx, pk.SQLFormat(), pk.SQLParams()...)
	if err != nil {
		return false, err
	}
	return (c != 0), nil
}

// FetchByPrimaryKey fetches a single Blog by its primary key
// it returns the specific error type(sql.ErrNoRows) when no rows found
func (m *_BlogDBMgr) FetchByPrimaryKey(ctx context.Context, id int64, userId int32) (*Blog, error) {
	obj := BlogMgr.NewBlog()
	pk := &IdUserIdOfBlogPK{
		Id:     id,
		UserId: userId,
	}

	query := fmt.Sprintf("SELECT %s FROM blogs %s", strings.Join(obj.GetColumns(), ","), pk.SQLFormat())
	objs, err := m.FetchBySQL(ctx, query, pk.SQLParams()...)
	if err != nil {
		return nil, err
	}
	if len(objs) > 0 {
		return objs[0], nil
	}
	return nil, sql.ErrNoRows
}

// err not found check
func (m *_BlogDBMgr) IsErrNotFound(err error) bool {
	return strings.Contains(err.Error(), "not found") || err == sql.ErrNoRows
}

// indexes

func (m *_BlogDBMgr) FindByStatus(ctx context.Context, status int32, limit int, offset int) ([]*Blog, error) {
	obj := BlogMgr.NewBlog()
	idx_ := &StatusOfBlogIDX{
		Status: status,
		limit:  limit,
		offset: offset,
	}
	query := fmt.Sprintf("SELECT %s FROM blogs %s", strings.Join(obj.GetColumns(), ","), idx_.SQLFormat(true))
	return m.FetchBySQL(ctx, query, idx_.SQLParams()...)
}

func (m *_BlogDBMgr) FindAllByStatus(ctx context.Context, status int32) ([]*Blog, error) {
	obj := BlogMgr.NewBlog()
	idx_ := &StatusOfBlogIDX{
		Status: status,
	}

	query := fmt.Sprintf("SELECT %s FROM blogs %s", strings.Join(obj.GetColumns(), ","), idx_.SQLFormat(true))
	return m.FetchBySQL(ctx, query, idx_.SQLParams()...)
}

func (m *_BlogDBMgr) FindByStatusGroup(ctx context.Context, items []int32) ([]*Blog, error) {
	obj := BlogMgr.NewBlog()
	if len(items) == 0 {
		return nil, nil
	}
	params := make([]interface{}, 0, len(items))
	for _, item := range items {
		params = append(params, item)
	}
	query := fmt.Sprintf("SELECT %s FROM blogs where status in (?", strings.Join(obj.GetColumns(), ",")) +
		strings.Repeat(",?", len(items)-1) + ")"
	return m.FetchBySQL(ctx, query, params...)
}

// uniques

func (m *_BlogDBMgr) FindOne(ctx context.Context, unique Unique) (PrimaryKey, error) {
	objs, err := m.queryLimit(ctx, unique.SQLFormat(true), unique.SQLLimit(), unique.SQLParams()...)
	if err != nil {
		return nil, err
	}
	if len(objs) > 0 {
		return objs[0], nil
	}
	return nil, fmt.Errorf("Blog find record not found")
}

func (m *_BlogDBMgr) FindFetch(ctx context.Context, index Index) (int64, []*Blog, error) {
	total, err := m.queryCount(ctx, index.SQLFormat(false), index.SQLParams()...)
	if err != nil {
		return total, nil, err
	}

	obj := BlogMgr.NewBlog()
	query := fmt.Sprintf("SELECT %s FROM blogs %s", strings.Join(obj.GetColumns(), ","), index.SQLFormat(true))
	results, err := m.FetchBySQL(ctx, query, index.SQLParams()...)
	if err != nil {
		return total, nil, err
	}
	return total, results, nil
}

func (m *_BlogDBMgr) queryLimit(ctx context.Context, where string, limit int, args ...interface{}) (results []PrimaryKey, err error) {
	pk := BlogMgr.NewPrimaryKey()
	query := fmt.Sprintf("SELECT %s FROM blogs %s", strings.Join(pk.Columns(), ","), where)
	rows, err := m.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("Blog query limit error: %v", err)
	}
	defer rows.Close()

	offset := 0

	for rows.Next() {
		if limit >= 0 && offset >= limit {
			break
		}
		offset++

		result := BlogMgr.NewPrimaryKey()
		err = rows.Scan(&(result.Id), &(result.UserId))
		if err != nil {
			m.db.SetError(err)
			return nil, err
		}

		results = append(results, result)
	}
	if err := rows.Err(); err != nil {
		m.db.SetError(err)
		return nil, fmt.Errorf("Blog query limit result error: %v", err)
	}
	return
}

func (m *_BlogDBMgr) queryCount(ctx context.Context, where string, args ...interface{}) (int64, error) {
	query := fmt.Sprintf("SELECT count( id ) FROM blogs %s", where)
	rows, err := m.db.Query(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("Blog query count error: %v", err)
	}
	defer rows.Close()

	var count int64
	if rows.Next() {
		if err = rows.Scan(&count); err != nil {
			m.db.SetError(err)
			return 0, err
		}
	}
	return count, nil
}

func (m *_BlogDBMgr) BatchCreate(ctx context.Context, objs []*Blog) (int64, error) {
	if len(objs) == 0 {
		return 0, nil
	}

	params := make([]string, 0, len(objs))
	values := make([]interface{}, 0, len(objs)*8)
	for _, obj := range objs {
		params = append(params, fmt.Sprintf("(%s)", strings.Join(orm.NewStringSlice(8, "?"), ",")))
		values = append(values, obj.Id)
		values = append(values, obj.UserId)
		values = append(values, obj.Title)
		values = append(values, obj.Content)
		values = append(values, obj.Status)
		values = append(values, obj.Readed)
		values = append(values, orm.TimeFormat(obj.CreatedAt))
		values = append(values, obj.UpdatedAt.Unix())
	}
	query := fmt.Sprintf("INSERT INTO blogs(%s) VALUES %s", strings.Join(objs[0].GetNoneIncrementColumns(), ","), strings.Join(params, ","))
	result, err := m.db.Exec(ctx, query, values...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// argument example:
// set:"a=?, b=?"
// where:"c=? and d=?"
// params:[]interface{}{"a", "b", "c", "d"}...
func (m *_BlogDBMgr) UpdateBySQL(ctx context.Context, set, where string, args ...interface{}) (int64, error) {
	query := fmt.Sprintf("UPDATE blogs SET %s", set)
	if where != "" {
		query = fmt.Sprintf("UPDATE blogs SET %s WHERE %s", set, where)
	}
	result, err := m.db.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (m *_BlogDBMgr) Create(ctx context.Context, obj *Blog) (int64, error) {
	params := orm.NewStringSlice(8, "?")
	q := fmt.Sprintf("INSERT INTO blogs(%s) VALUES(%s)",
		strings.Join(obj.GetNoneIncrementColumns(), ","),
		strings.Join(params, ","))

	values := make([]interface{}, 0, 8)
	values = append(values, obj.Id)
	values = append(values, obj.UserId)
	values = append(values, obj.Title)
	values = append(values, obj.Content)
	values = append(values, obj.Status)
	values = append(values, obj.Readed)
	values = append(values, orm.TimeFormat(obj.CreatedAt))
	values = append(values, obj.UpdatedAt.Unix())
	result, err := m.db.Exec(ctx, q, values...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (m *_BlogDBMgr) Update(ctx context.Context, obj *Blog) (int64, error) {
	columns := []string{
		"title = ?",
		"content = ?",
		"status = ?",
		"readed = ?",
		"created_at = ?",
		"updated_at = ?",
	}

	pk := obj.GetPrimaryKey()
	q := fmt.Sprintf("UPDATE blogs SET %s %s", strings.Join(columns, ","), pk.SQLFormat())
	values := make([]interface{}, 0, 8-2)
	values = append(values, obj.Title)
	values = append(values, obj.Content)
	values = append(values, obj.Status)
	values = append(values, obj.Readed)
	values = append(values, orm.TimeFormat(obj.CreatedAt))
	values = append(values, obj.UpdatedAt.Unix())
	values = append(values, pk.SQLParams()...)

	result, err := m.db.Exec(ctx, q, values...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (m *_BlogDBMgr) Save(ctx context.Context, obj *Blog) (int64, error) {
	affected, err := m.Update(ctx, obj)
	if err != nil {
		return affected, err
	}
	if affected == 0 {
		return m.Create(ctx, obj)
	}
	return affected, err
}

func (m *_BlogDBMgr) Delete(ctx context.Context, obj *Blog) (int64, error) {
	return m.DeleteByPrimaryKey(ctx, obj.Id, obj.UserId)
}

func (m *_BlogDBMgr) DeleteByPrimaryKey(ctx context.Context, id int64, userId int32) (int64, error) {
	pk := &IdUserIdOfBlogPK{
		Id:     id,
		UserId: userId,
	}
	q := fmt.Sprintf("DELETE FROM blogs %s", pk.SQLFormat())
	result, err := m.db.Exec(ctx, q, pk.SQLParams()...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (m *_BlogDBMgr) DeleteBySQL(ctx context.Context, where string, args ...interface{}) (int64, error) {
	query := "DELETE FROM blogs"
	if where != "" {
		query = fmt.Sprintf("DELETE FROM blogs WHERE %s", where)
	}
	result, err := m.db.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
