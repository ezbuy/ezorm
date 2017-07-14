package test

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

func (m *_PeopleMgr) query(query string, args ...interface{}) ([]*People, error) {
	rows, err := mssqlQuery(query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var Age sql.NullInt64
	var IndexAPart2 sql.NullInt64

	var results []*People
	for rows.Next() {
		var result People
		err := rows.Scan(&(result.NonIndexA), &(result.NonIndexB), &(result.PeopleId), &Age, &(result.Name), &(result.IndexAPart1), &IndexAPart2, &(result.IndexAPart3), &(result.UniquePart1), &(result.UniquePart2), &(result.CreateDate), &(result.UpdateDate))
		if err != nil {
			return nil, err
		}

		result.Age = int32(Age.Int64)
		result.IndexAPart2 = int32(IndexAPart2.Int64)

		results = append(results, &result)
	}

	// 目前sql server保存的都是local time
	for _, r := range results {
		r.CreateDate = m.timeConvToLocal(r.CreateDate)
		r.UpdateDate = m.timeConvToLocal(r.UpdateDate)
	}

	return results, nil
}

func (m *_PeopleMgr) queryOne(query string, args ...interface{}) (*People, error) {
	rows, err := m.query(query, args...)
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, sql.ErrNoRows
	}

	return rows[0], err
}

func (m *_PeopleMgr) Save(obj *People) (sql.Result, error) {
	if obj.PeopleId == 0 {
		return m.saveInsert(obj)
	}
	return m.saveUpdate(obj)
}

func (m *_PeopleMgr) saveInsert(obj *People) (sql.Result, error) {
	query := "INSERT INTO [dbo].[People] (NonIndexA, NonIndexB, Age, Name, IndexAPart1, IndexAPart2, IndexAPart3, UniquePart1, UniquePart2, CreateDate, UpdateDate) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	result, err := mssqlExec(query, obj.NonIndexA, obj.NonIndexB, obj.Age, obj.Name, obj.IndexAPart1, obj.IndexAPart2, obj.IndexAPart3, obj.UniquePart1, obj.UniquePart2, obj.CreateDate, obj.UpdateDate)
	if err != nil {
		return result, err
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return result, err
	}

	obj.PeopleId = int32(lastInsertId)

	return result, err
}

func (m *_PeopleMgr) saveUpdate(obj *People) (sql.Result, error) {
	query := "UPDATE [dbo].[People] SET NonIndexA=?, NonIndexB=?, Age=?, Name=?, IndexAPart1=?, IndexAPart2=?, IndexAPart3=?, UniquePart1=?, UniquePart2=?, CreateDate=?, UpdateDate=? WHERE PeopleId=?"
	return mssqlExec(query, obj.NonIndexA, obj.NonIndexB, obj.Age, obj.Name, obj.IndexAPart1, obj.IndexAPart2, obj.IndexAPart3, obj.UniquePart1, obj.UniquePart2, obj.CreateDate, obj.UpdateDate, obj.PeopleId)
}

func (m *_PeopleMgr) InsertBatch(objs []*People) (sql.Result, error) {
	if len(objs) == 0 {
		return nil, errors.New("Empty insert")
	}

	values := make([]string, 0, len(objs))
	params := make([]interface{}, 0, len(objs)*11)
	for _, obj := range objs {
		values = append(values, "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		params = append(params, obj.NonIndexA, obj.NonIndexB, obj.Age, obj.Name, obj.IndexAPart1, obj.IndexAPart2, obj.IndexAPart3, obj.UniquePart1, obj.UniquePart2, obj.CreateDate, obj.UpdateDate)
	}
	query := fmt.Sprintf("INSERT INTO [dbo].[People] (NonIndexA, NonIndexB, Age, Name, IndexAPart1, IndexAPart2, IndexAPart3, UniquePart1, UniquePart2, CreateDate, UpdateDate) VALUES %s", strings.Join(values, ","))
	return mssqlExec(query, params...)
}

func (m *_PeopleMgr) GetId2Obj(objs []*People) map[int32]*People {
	id2obj := make(map[int32]*People, len(objs))
	for _, obj := range objs {
		id2obj[obj.PeopleId] = obj
	}
	return id2obj
}

func (m *_PeopleMgr) GetIds(objs []*People) []int32 {
	ids := make([]int32, len(objs))
	for i, obj := range objs {
		ids[i] = obj.PeopleId
	}
	return ids
}

func (m *_PeopleMgr) FindByID(id int32) (*People, error) {
	query := "SELECT NonIndexA, NonIndexB, PeopleId, Age, Name, IndexAPart1, IndexAPart2, IndexAPart3, UniquePart1, UniquePart2, CreateDate, UpdateDate FROM [dbo].[People] WHERE PeopleId=?"
	return m.queryOne(query, id)
}

func (m *_PeopleMgr) FindByIDs(ids []int32) ([]*People, error) {
	idsArray := m.getSplitIds(ids)

	var vals []*People
	for _, idsBySep := range idsArray {
		placeHolders, args := m.getPlaceHolderAndParameter(idsBySep)

		query := fmt.Sprintf("SELECT NonIndexA, NonIndexB, PeopleId, Age, Name, IndexAPart1, IndexAPart2, IndexAPart3, UniquePart1, UniquePart2, CreateDate, UpdateDate FROM [dbo].[People] WHERE PeopleId IN (%s)", placeHolders)
		val, err := m.query(query, args...)
		if err != nil {
			return nil, err
		}

		vals = append(vals, val...)
	}

	return vals, nil
}

func (m *_PeopleMgr) getPlaceHolderAndParameter(idsBySep []int32) (string, []interface{}) {
	params := strings.Repeat("?,", len(idsBySep))
	if len(params) > 0 {
		params = params[:len(params)-1]
	}

	val := make([]interface{}, 0, len(idsBySep))

	for _, id := range idsBySep {
		val = append(val, id)
	}

	return params, val
}

func (m *_PeopleMgr) getSplitIds(ids []int32) [][]int32 {
	re := [][]int32{}
	if len(ids) <= 0 {
		return re
	}

	maxLimit := 2000
	idsBySep := []int32{}
	for i, id := range ids {
		idsBySep = append(idsBySep, id)
		if (i+1)%maxLimit == 0 {
			ns := make([]int32, len(idsBySep))
			copy(ns, idsBySep)
			idsBySep = idsBySep[:0]
			re = append(re, ns)
		}
	}
	if len(idsBySep) > 0 {
		ns := make([]int32, len(idsBySep))
		copy(ns, idsBySep)
		re = append(re, ns)
	}
	return re
}

func (m *_PeopleMgr) FindByIndexAPart1IndexAPart2IndexAPart3(IndexAPart1 int64, IndexAPart2 int32, IndexAPart3 int32, offset int, limit int, sortFields ...string) ([]*People, error) {
	orderBy := "ORDER BY %s"
	if len(sortFields) != 0 {
		orderBy = fmt.Sprintf(orderBy, strings.Join(sortFields, ","))
	} else {
		orderBy = fmt.Sprintf(orderBy, "PeopleId")
	}

	query := fmt.Sprintf("SELECT NonIndexA, NonIndexB, PeopleId, Age, Name, IndexAPart1, IndexAPart2, IndexAPart3, UniquePart1, UniquePart2, CreateDate, UpdateDate FROM [dbo].[People] WHERE IndexAPart1=? AND IndexAPart2=? AND IndexAPart3=? %s  OFFSET ? Rows FETCH NEXT ? Rows ONLY", orderBy)

	return m.query(query, IndexAPart1, IndexAPart2, IndexAPart3, offset, limit)
}

func (m *_PeopleMgr) FindOneByUniquePart1UniquePart2(UniquePart1 int32, UniquePart2 int32) (*People, error) {
	query := "SELECT NonIndexA, NonIndexB, PeopleId, Age, Name, IndexAPart1, IndexAPart2, IndexAPart3, UniquePart1, UniquePart2, CreateDate, UpdateDate FROM [dbo].[People] WHERE UniquePart1=? AND UniquePart2=?"
	return m.queryOne(query, UniquePart1, UniquePart2)
}

func (m *_PeopleMgr) FindByAge(Age int32, offset int, limit int, sortFields ...string) ([]*People, error) {
	orderBy := "ORDER BY %s"
	if len(sortFields) != 0 {
		orderBy = fmt.Sprintf(orderBy, strings.Join(sortFields, ","))
	} else {
		orderBy = fmt.Sprintf(orderBy, "PeopleId")
	}

	query := fmt.Sprintf("SELECT NonIndexA, NonIndexB, PeopleId, Age, Name, IndexAPart1, IndexAPart2, IndexAPart3, UniquePart1, UniquePart2, CreateDate, UpdateDate FROM [dbo].[People] WHERE Age=? %s  OFFSET ? Rows FETCH NEXT ? Rows ONLY", orderBy)

	return m.query(query, Age, offset, limit)
}

func (m *_PeopleMgr) FindOneByName(Name string) (*People, error) {
	query := "SELECT NonIndexA, NonIndexB, PeopleId, Age, Name, IndexAPart1, IndexAPart2, IndexAPart3, UniquePart1, UniquePart2, CreateDate, UpdateDate FROM [dbo].[People] WHERE Name=?"
	return m.queryOne(query, Name)
}

func (m *_PeopleMgr) FindOne(where string, args ...interface{}) (*People, error) {
	query := m.getQuerysql(true, where)
	return m.queryOne(query, args...)
}

func (m *_PeopleMgr) Find(where string, args ...interface{}) ([]*People, error) {
	query := m.getQuerysql(false, where)
	return m.query(query, args...)
}

func (m *_PeopleMgr) FindAll() (results []*People, err error) {
	return m.Find("")
}

func (m *_PeopleMgr) FindWithOffset(where string, offset int, limit int, args ...interface{}) ([]*People, error) {
	query := m.getQuerysql(false, where)

	query = query + " OFFSET ? Rows FETCH NEXT ? Rows ONLY"

	args = append(args, offset)
	args = append(args, limit)

	return m.query(query, args...)
}

func (m *_PeopleMgr) getQuerysql(topOne bool, where string) string {
	query := `SELECT `
	if topOne {
		query = query + ` TOP 1 `
	}
	query = query + ` NonIndexA, NonIndexB, PeopleId, Age, Name, IndexAPart1, IndexAPart2, IndexAPart3, UniquePart1, UniquePart2, CreateDate, UpdateDate FROM [dbo].[People] WITH(NOLOCK) `

	where = strings.Trim(where, " ")
	if where != "" {
		upwhere := strings.ToUpper(where)

		if !strings.HasPrefix(upwhere, "WHERE") && !strings.HasPrefix(upwhere, "ORDER BY") {
			where = " WHERE " + where
		}

		query = query + where
	}
	return query
}

func (m *_PeopleMgr) Del(where string, params ...interface{}) (sql.Result, error) {
	query := "DELETE FROM [dbo].[People]"
	if where != "" {
		query = fmt.Sprintf("DELETE FROM People WHERE " + where)
	}
	return mssqlExec(query, params...)
}

// argument example:
// set:"a=?, b=?"
// where:"c=? and d=?"
// params:[]interface{}{"a", "b", "c", "d"}...
func (m *_PeopleMgr) Update(set, where string, params ...interface{}) (sql.Result, error) {
	query := fmt.Sprintf("UPDATE [dbo].[People] SET %s", set)
	if where != "" {
		query = fmt.Sprintf("UPDATE [dbo].[People] SET %s WHERE %s", set, where)
	}
	return mssqlExec(query, params...)
}

func (m *_PeopleMgr) Count(where string, args ...interface{}) (int32, error) {
	query := "SELECT COUNT(*) FROM [dbo].[People]"
	if where != "" {
		query = query + " WHERE " + where
	}

	rows, err := mssqlQuery(query, args...)
	if err != nil {
		return 0, err
	}

	defer rows.Close()

	var count int32
	if rows.Next() {
		err = rows.Scan(&count)
	}

	return count, err
}

func (m *_PeopleMgr) timeConvToLocal(t *time.Time) *time.Time {
	if t == nil {
		return nil
	}
	localTime := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(),
		t.Second(), t.Nanosecond(), time.Local)
	return &localTime
}
