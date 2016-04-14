package people

import (
	"database/sql"
	"fmt"
	"strings"
)

func (m *_PeopleMgr) query(query string, args ...interface{}) ([]*People, error) {
	rows, err := _db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var NonIndexA sql.NullString
	var NonIndexB sql.NullString
	var Age sql.NullInt64
	var Name sql.NullString
	var IndexAPart1 sql.NullInt64
	var IndexAPart2 sql.NullInt64
	var IndexAPart3 sql.NullInt64
	var UniquePart1 sql.NullInt64
	var UniquePart2 sql.NullInt64

	var results []*People
	for rows.Next() {
		var result People
		err := rows.Scan(&NonIndexA, &NonIndexB, &(result.PeopleId), &Age, &Name, &IndexAPart1, &IndexAPart2, &IndexAPart3, &UniquePart1, &UniquePart2)
		if err != nil {
			return nil, err
		}

		result.NonIndexA = NonIndexA.String
		result.NonIndexB = NonIndexB.String
		result.Age = int32(Age.Int64)
		result.Name = Name.String
		result.IndexAPart1 = IndexAPart1.Int64
		result.IndexAPart2 = int32(IndexAPart2.Int64)
		result.IndexAPart3 = int32(IndexAPart3.Int64)
		result.UniquePart1 = int32(UniquePart1.Int64)
		result.UniquePart2 = int32(UniquePart2.Int64)
		results = append(results, &result)
	}
	return results, nil
}

func (m *_PeopleMgr) queryOne(query string, args ...interface{}) (*People, error) {
	rows, err := _db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var NonIndexA sql.NullString
	var NonIndexB sql.NullString
	var Age sql.NullInt64
	var Name sql.NullString
	var IndexAPart1 sql.NullInt64
	var IndexAPart2 sql.NullInt64
	var IndexAPart3 sql.NullInt64
	var UniquePart1 sql.NullInt64
	var UniquePart2 sql.NullInt64

	if !rows.Next() {
		return nil, sql.ErrNoRows
	}

	var result People
	err = rows.Scan(&NonIndexA, &NonIndexB, &(result.PeopleId), &Age, &Name, &IndexAPart1, &IndexAPart2, &IndexAPart3, &UniquePart1, &UniquePart2)
	if err != nil {
		return nil, err
	}

	result.NonIndexA = NonIndexA.String
	result.NonIndexB = NonIndexB.String
	result.Age = int32(Age.Int64)
	result.Name = Name.String
	result.IndexAPart1 = IndexAPart1.Int64
	result.IndexAPart2 = int32(IndexAPart2.Int64)
	result.IndexAPart3 = int32(IndexAPart3.Int64)
	result.UniquePart1 = int32(UniquePart1.Int64)
	result.UniquePart2 = int32(UniquePart2.Int64)
	return &result, nil
}

func (m *_PeopleMgr) Save(obj *People) (sql.Result, error) {
	if obj.PeopleId == 0 {
		return m.saveInsert(obj)
	}
	return m.saveUpdate(obj)
}

func (m *_PeopleMgr) saveInsert(obj *People) (sql.Result, error) {
	query := "INSERT INTO [dbo].[People] (NonIndexA, NonIndexB, Age, Name, IndexAPart1, IndexAPart2, IndexAPart3, UniquePart1, UniquePart2) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"
	result, err := _sqlServer.Exec(query, obj.NonIndexA, obj.NonIndexB, obj.Age, obj.Name, obj.IndexAPart1, obj.IndexAPart2, obj.IndexAPart3, obj.UniquePart1, obj.UniquePart2)
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
	query := "UPDATE [dbo].[People] SET NonIndexA=?, NonIndexB=?, Age=?, Name=?, IndexAPart1=?, IndexAPart2=?, IndexAPart3=?, UniquePart1=?, UniquePart2=? WHERE PeopleId=?"
	return _sqlServer.Exec(query, obj.NonIndexA, obj.NonIndexB, obj.Age, obj.Name, obj.IndexAPart1, obj.IndexAPart2, obj.IndexAPart3, obj.UniquePart1, obj.UniquePart2, obj.PeopleId)
}

func (m *_PeopleMgr) FindByID(id int32) (*People, error) {
	query := "SELECT NonIndexA, NonIndexB, PeopleId, Age, Name, IndexAPart1, IndexAPart2, IndexAPart3, UniquePart1, UniquePart2 FROM [dbo].[People] WHERE PeopleId=?"
	return m.queryOne(query, id)
}

func (m *_PeopleMgr) FindByIndexAPart1IndexAPart2IndexAPart3(IndexAPart1 int64, IndexAPart2 int32, IndexAPart3 int32, offset int, limit int, sortFields ...string) ([]*People, error) {
	orderBy := "ORDER BY %s"
	if len(sortFields) != 0 {
		orderBy = fmt.Sprintf(orderBy, strings.Join(sortFields, ","))
	} else {
		orderBy = fmt.Sprintf(orderBy, "PeopleId")
	}

	query := fmt.Sprintf("SELECT NonIndexA, NonIndexB, PeopleId, Age, Name, IndexAPart1, IndexAPart2, IndexAPart3, UniquePart1, UniquePart2 FROM [dbo].[People] WHERE IndexAPart1=? AND IndexAPart2=? AND IndexAPart3=? %s  OFFSET ? Rows FETCH NEXT ? Rows ONLY", orderBy)

	return m.query(query, IndexAPart1, IndexAPart2, IndexAPart3, offset, limit)
}

func (m *_PeopleMgr) FindOneByUniquePart1UniquePart2(UniquePart1 int32, UniquePart2 int32) (*People, error) {
	query := "SELECT NonIndexA, NonIndexB, PeopleId, Age, Name, IndexAPart1, IndexAPart2, IndexAPart3, UniquePart1, UniquePart2 FROM [dbo].[People] WHERE UniquePart1=? AND UniquePart2=?"
	return m.queryOne(query, UniquePart1, UniquePart2)
}

func (m *_PeopleMgr) FindByAge(Age int32, offset int, limit int, sortFields ...string) ([]*People, error) {
	orderBy := "ORDER BY %s"
	if len(sortFields) != 0 {
		orderBy = fmt.Sprintf(orderBy, strings.Join(sortFields, ","))
	} else {
		orderBy = fmt.Sprintf(orderBy, "PeopleId")
	}

	query := fmt.Sprintf("SELECT NonIndexA, NonIndexB, PeopleId, Age, Name, IndexAPart1, IndexAPart2, IndexAPart3, UniquePart1, UniquePart2 FROM [dbo].[People] WHERE Age=? %s  OFFSET ? Rows FETCH NEXT ? Rows ONLY", orderBy)

	return m.query(query, Age, offset, limit)
}

func (m *_PeopleMgr) FindOneByName(Name string) (*People, error) {
	query := "SELECT NonIndexA, NonIndexB, PeopleId, Age, Name, IndexAPart1, IndexAPart2, IndexAPart3, UniquePart1, UniquePart2 FROM [dbo].[People] WHERE Name=?"
	return m.queryOne(query, Name)
}

func (m *_PeopleMgr) FindOne(where string, args ...interface{}) (*People, error) {
	query := m.getQuerysql(true, where)
	return m.queryOne(query, args...)
}

func (m *_PeopleMgr) Find(where string, args ...interface{}) (results []*People, err error) {
	query := m.getQuerysql(false, where)
	return m.query(query, args...)
}

func (m *_PeopleMgr) FindAll() (results []*People, err error) {
	return m.Find("")
}

func (m *_PeopleMgr) FindWithOffset(where string, offset int, limit int, args ...interface{}) (results []*People, err error) {
	query := m.getQuerysql(false, where)

	if !strings.Contains(strings.ToLower(where), "ORDER BY") {
		where = " ORDER BY Name"
	}
	query = query + where + " OFFSET ? Rows FETCH NEXT ? Rows ONLY"
	args = append(args, offset)
	args = append(args, limit)

	return m.query(query, args...)
}

func (m *_PeopleMgr) getQuerysql(topOne bool, where string) string {
	query := `SELECT `
	if topOne {
		query = query + ` TOP 1 `
	}
	query = query + ` NonIndexA, NonIndexB, PeopleId, Age, Name, IndexAPart1, IndexAPart2, IndexAPart3, UniquePart1, UniquePart2 FROM [dbo].[People] WITH(NOLOCK) `

	if where != "" {
		if strings.Index(strings.Trim(where, " "), "WHERE") == -1 {
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
	return _db.Exec(query, params...)
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
	return _db.Exec(query, params...)
}
