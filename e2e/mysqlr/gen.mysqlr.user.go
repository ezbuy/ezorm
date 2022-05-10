package mysqlr


import (
	"fmt"
	"time"
	"strings"
	"database/sql"
	"context"

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

type User struct {
	Id  int64 `mysql:"id"`
	UserId  int32 `mysql:"user_id"`
	Name  string `mysql:"name"`
	CreatedAt  int64 `mysql:"created_at"`
	UpdatedAt  int64 `mysql:"updated_at"`
}

var UserColumns = struct{
	Id  string
	UserId  string
	Name  string
	CreatedAt  string
	UpdatedAt  string
}{
	"id",
	"user_id",
	"name",
	"created_at",
	"updated_at",
}

type _UserMgr struct {
}
var UserMgr *_UserMgr

func (m *_UserMgr) NewUser() *User {
	return &User{}
}





func (obj *User) GetNameSpace() string {
	return "mysqlr"
}

func (obj *User) GetClassName() string {
	return "User"
}

func (obj *User) GetTableName() string {
	return "users"
}

func (obj *User) GetColumns() []string {
	columns := []string{
	"users.id",
	"users.user_id",
	"users.name",
	"users.created_at",
	"users.updated_at",
	}
	return columns
}

func (obj *User) GetNoneIncrementColumns() []string {
	columns := []string{
	"id",
	"user_id",
	"name",
	"created_at",
	"updated_at",
	}
	return columns
}

func (obj *User) GetPrimaryKey() PrimaryKey {
	pk := UserMgr.NewPrimaryKey()
	pk.Id = obj.Id
	pk.UserId = obj.UserId
	return pk
}

func (obj *User) Validate() error {
	validate := validator.New()
	return validate.Struct(obj)
}







type IdUserIdOfUserPK struct{
	Id int64
	UserId int32
}

func (m *_UserMgr) NewPrimaryKey() *IdUserIdOfUserPK {
		return &IdUserIdOfUserPK{}
}

func (u *IdUserIdOfUserPK) Key() string {
	strs := []string{
		"Id",
			fmt.Sprint(u.Id),
		"UserId",
			fmt.Sprint(u.UserId),
	}
	return  strings.Join(strs, ":")
}

func (u *IdUserIdOfUserPK) Parse(key string) error {
	arr := strings.Split(key, ":")
	if len(arr) % 2 != 0 {
		return fmt.Errorf("key (%s) format error", key)
	}
	kv := map[string]string{}
	for i := 0; i < len(arr) / 2; i++ {
		kv[arr[2*i]] = arr[2*i + 1]
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

func (u *IdUserIdOfUserPK) SQLFormat() string {
	conditions := []string{
		"id = ?",
		"user_id = ?",
	}
	return orm.SQLWhere(conditions)
}

func (u *IdUserIdOfUserPK) SQLParams() []interface{} {
	return []interface{}{
		u.Id,
		u.UserId,
	}
}

func (u *IdUserIdOfUserPK) Columns() []string {
	return []string{
		"id",
		"user_id",
	}
}










type _UserDBMgr struct {
	db orm.DB
}

func (m *_UserMgr) DB(db orm.DB) *_UserDBMgr {
	return UserDBMgr(db)
}

func UserDBMgr(db orm.DB) *_UserDBMgr {
	if db == nil {
		panic(fmt.Errorf("UserDBMgr init need db"))
	}
	return &_UserDBMgr{db: db}
}

func (m *_UserDBMgr) Search (ctx context.Context, where string, orderby string, limit string, args ...interface{}) ([]*User, error) {
	obj := UserMgr.NewUser()

    if limit = strings.ToUpper(strings.TrimSpace(limit)); limit !="" && !strings.HasPrefix(limit, "LIMIT") {
	    limit = "LIMIT " + limit
	}

	conditions := []string{where, orderby, limit}
	query := fmt.Sprintf("SELECT %s FROM users %s", strings.Join(obj.GetColumns(), ","), strings.Join(conditions, " "))
	return m.FetchBySQL(ctx, query, args...)
}

func (m *_UserDBMgr) SearchConditions(ctx context.Context,conditions []string, orderby string, offset int, limit int, args ...interface{}) ([]*User, error) {
	obj := UserMgr.NewUser()
	q := fmt.Sprintf("SELECT %s FROM users %s %s %s",
			strings.Join(obj.GetColumns(), ","),
			orm.SQLWhere(conditions),
			orderby,
			orm.SQLOffsetLimit(offset, limit))

	return m.FetchBySQL(ctx,q, args...)
}


func (m *_UserDBMgr) SearchCount(ctx context.Context,where string, args ...interface{}) (int64, error){
	return m.queryCount(ctx,where, args...)
}

func (m *_UserDBMgr) SearchConditionsCount(ctx context.Context,conditions []string, args ...interface{}) (int64, error){
	return m.queryCount(ctx,orm.SQLWhere(conditions), args...)
}

func (m *_UserDBMgr) FetchBySQL(ctx context.Context,q string, args ... interface{}) (results []*User, err error) {
	rows, err := m.db.Query(ctx,q, args...)
	if err != nil {
		return nil, fmt.Errorf("User fetch error: %v", err)
	}
	defer rows.Close()

	

	for rows.Next() {
		var result User
		err = rows.Scan(&(result.Id),&(result.UserId),&(result.Name),&(result.CreatedAt),&(result.UpdatedAt),)
		if err != nil {
			m.db.SetError(err)
			return nil, err
		}

		
		
		
		
		
		
		results = append(results, &result)
	}
	if err = rows.Err() ;err != nil {
		m.db.SetError(err)
		return nil, fmt.Errorf("User fetch result error: %v", err)
	}
	return
}
func (m *_UserDBMgr) Exist(ctx context.Context, pk PrimaryKey) (bool, error) {
	c, err := m.queryCount(ctx, pk.SQLFormat(), pk.SQLParams()...)
	if err != nil {
		return false, err
	}
	return (c != 0), nil
}

// FetchByPrimaryKey fetches a single User by its primary key
// it returns the specific error type(sql.ErrNoRows) when no rows found
func (m *_UserDBMgr) FetchByPrimaryKey(ctx context.Context,id int64, userId int32) (*User, error) {
	obj := UserMgr.NewUser()
	pk := &IdUserIdOfUserPK{
	Id : id,
UserId : userId,

	}

	query := fmt.Sprintf("SELECT %s FROM users %s", strings.Join(obj.GetColumns(), ","), pk.SQLFormat())
	objs, err := m.FetchBySQL(ctx,query, pk.SQLParams()...)
	if err != nil {
		return nil, err
	}
	if len(objs) > 0 {
		return objs[0], nil
	}
	return nil, sql.ErrNoRows
}


// err not found check
func (m *_UserDBMgr) IsErrNotFound(err error) bool {
	return strings.Contains(err.Error(), "not found") || err == sql.ErrNoRows
}

// indexes

// uniques

func (m *_UserDBMgr) FindOne(ctx context.Context,unique Unique) (PrimaryKey, error) {
	objs, err := m.queryLimit(ctx,unique.SQLFormat(true), unique.SQLLimit(), unique.SQLParams()...)
	if err != nil {
		return nil, err
	}
	if len(objs) > 0 {
		return objs[0], nil
	}
	return nil, fmt.Errorf("User find record not found")
}

func (m *_UserDBMgr) FindFetch(ctx context.Context,index Index) (int64, []*User, error) {
	total, err := m.queryCount(ctx,index.SQLFormat(false), index.SQLParams()...)
	if err != nil {
		return total, nil, err
	}

	obj := UserMgr.NewUser()
	query := fmt.Sprintf("SELECT %s FROM users %s", strings.Join(obj.GetColumns(), ","), index.SQLFormat(true))
	results, err := m.FetchBySQL(ctx, query, index.SQLParams()...)
	if err != nil {
		return total, nil, err
	}
	return total, results, nil
}


func (m *_UserDBMgr) queryLimit(ctx context.Context,where string, limit int, args ...interface{}) (results []PrimaryKey, err error){
	pk := UserMgr.NewPrimaryKey()
	query := fmt.Sprintf("SELECT %s FROM users %s", strings.Join(pk.Columns(), ","), where)
	rows, err := m.db.Query(ctx,query, args...)
	if err != nil {
		return nil, fmt.Errorf("User query limit error: %v", err)
	}
	defer rows.Close()

	offset :=0

	for rows.Next() {
		if limit >= 0 && offset >= limit {
			break
		}
		offset++

		result := UserMgr.NewPrimaryKey()
		err = rows.Scan(&(result.Id),&(result.UserId),)
		if err != nil {
			m.db.SetError(err)
			return nil, err
		}

		
		
		
		results = append(results, result)
	}
	if err := rows.Err() ;err != nil {
		m.db.SetError(err)
		return nil, fmt.Errorf("User query limit result error: %v", err)
	}
	return
}


func (m *_UserDBMgr) queryCount(ctx context.Context,where string, args ...interface{}) (int64, error){
	query := fmt.Sprintf("SELECT count( id ) FROM users %s", where)
	rows, err := m.db.Query(ctx,query, args...)
	if err != nil {
		return 0, fmt.Errorf("User query count error: %v", err)
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







func (m *_UserDBMgr) BatchCreate(ctx context.Context,objs []*User) (int64, error) {
	if len(objs) == 0 {
		return 0, nil
	}

	params := make([]string, 0, len(objs))
	values := make([]interface{}, 0, len(objs)*5)
	for _, obj := range objs {
		params = append(params, fmt.Sprintf("(%s)", strings.Join(orm.NewStringSlice(5, "?"), ",")))
					values = append(values, obj.Id)
					values = append(values, obj.UserId)
					values = append(values, obj.Name)
					values = append(values, obj.CreatedAt)
					values = append(values, obj.UpdatedAt)
	}
	query := fmt.Sprintf("INSERT INTO users(%s) VALUES %s", strings.Join(objs[0].GetNoneIncrementColumns(), ","), strings.Join(params, ","))
	result, err := m.db.Exec(ctx,query, values...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// argument example:
// set:"a=?, b=?"
// where:"c=? and d=?"
// params:[]interface{}{"a", "b", "c", "d"}...
func (m *_UserDBMgr) UpdateBySQL(ctx context.Context,set, where string, args ...interface{}) (int64, error) {
	query := fmt.Sprintf("UPDATE users SET %s", set)
	if where != "" {
		query = fmt.Sprintf("UPDATE users SET %s WHERE %s", set, where)
	}
	result, err := m.db.Exec(ctx,query, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (m *_UserDBMgr) Create(ctx context.Context,obj *User) (int64, error) {
	params := orm.NewStringSlice(5, "?")
	q := fmt.Sprintf("INSERT INTO users(%s) VALUES(%s)",
		strings.Join(obj.GetNoneIncrementColumns(), ","),
		strings.Join(params, ","))

	values := make([]interface{}, 0, 5)
				values = append(values, obj.Id)
				values = append(values, obj.UserId)
				values = append(values, obj.Name)
				values = append(values, obj.CreatedAt)
				values = append(values, obj.UpdatedAt)
	result, err := m.db.Exec(ctx,q, values...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (m *_UserDBMgr) Update(ctx context.Context,obj *User) (int64, error) {
	columns := []string{
		"name = ?",
		"created_at = ?",
		"updated_at = ?",
	}

	pk := obj.GetPrimaryKey()
	q := fmt.Sprintf("UPDATE users SET %s %s", strings.Join(columns, ","), pk.SQLFormat())
	values := make([]interface{}, 0, 5 - 2)
					values = append(values, obj.Name)
					values = append(values, obj.CreatedAt)
					values = append(values, obj.UpdatedAt)
	values = append(values, pk.SQLParams()...)

	result, err := m.db.Exec(ctx,q, values...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (m *_UserDBMgr) Save(ctx context.Context,obj *User) (int64, error) {
	affected, err := m.Update(ctx,obj)
	if err != nil {
		return affected, err
	}
	if affected == 0 {
		return m.Create(ctx,obj)
	}
	return affected, err
}

func (m *_UserDBMgr) Delete(ctx context.Context,obj *User) (int64, error) {
	return m.DeleteByPrimaryKey(ctx,obj.Id, obj.UserId)
}

func (m *_UserDBMgr) DeleteByPrimaryKey(ctx context.Context,id int64, userId int32) (int64, error) {
	pk:= &IdUserIdOfUserPK{
	Id : id,
UserId : userId,

	}
	q := fmt.Sprintf("DELETE FROM users %s", pk.SQLFormat())
	result, err := m.db.Exec(ctx,q , pk.SQLParams()...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (m *_UserDBMgr) DeleteBySQL(ctx context.Context,where string, args ...interface{}) (int64, error) {
	query := "DELETE FROM users"
	if where != "" {
		query = fmt.Sprintf("DELETE FROM users WHERE %s", where)
	}
	result, err := m.db.Exec(ctx,query, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}


