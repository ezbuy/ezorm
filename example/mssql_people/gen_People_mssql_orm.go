package people

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/ezbuy/ezorm/db"
)

func (m *_PeopleMgr) Save(obj *People) (sql.Result, error) {
	if obj.PeopleId == 0 {
		return m.saveInsert(obj)
	}
	return m.saveUpdate(obj)
}

func (m *_PeopleMgr) saveInsert(obj *People) (sql.Result, error) {
	query := "insert into dbo.[People] (NonIndexA, NonIndexB, Age, Name, IndexAPart1, IndexAPart2, IndexAPart3, UniquePart1, UniquePart2) values (?, ?, ?, ?, ?, ?, ?, ?, ?)"
	result, err := db.Exec(query, obj.NonIndexA, obj.NonIndexB, obj.Age, obj.Name, obj.IndexAPart1, obj.IndexAPart2, obj.IndexAPart3, obj.UniquePart1, obj.UniquePart2)
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
	query := "update dbo.[People] set NonIndexA=?, NonIndexB=?, Age=?, Name=?, IndexAPart1=?, IndexAPart2=?, IndexAPart3=?, UniquePart1=?, UniquePart2=? where PeopleId=?"
	return db.Exec(query, obj.NonIndexA, obj.NonIndexB, obj.Age, obj.Name, obj.IndexAPart1, obj.IndexAPart2, obj.IndexAPart3, obj.UniquePart1, obj.UniquePart2, obj.PeopleId)
}

func (m *_PeopleMgr) FindByID(id int32) (*People, error) {
	query := "SELECT NonIndexA, NonIndexB, PeopleId, Age, Name, IndexAPart1, IndexAPart2, IndexAPart3, UniquePart1, UniquePart2 FROM People WHERE PeopleId=?"
	var obj People
	err := db.Query(&obj, query, id)
	return &obj, err
}

func (m *_PeopleMgr) FindByIndexAPart1IndexAPart2IndexAPart3(IndexAPart1 int32, IndexAPart2 int32, IndexAPart3 int32, offset int, limit int, sortFields ...string) (objs []*People, err error) {
	orderBy := "ORDER BY %s"
	if len(sortFields) != 0 {
		orderBy = fmt.Sprintf(orderBy, strings.Join(sortFields, ","))
	} else {
		orderBy = fmt.Sprintf(orderBy, "PeopleId")
	}

	query := fmt.Sprintf("SELECT NonIndexA, NonIndexB, PeopleId, Age, Name, IndexAPart1, IndexAPart2, IndexAPart3, UniquePart1, UniquePart2 FROM People WHERE IndexAPart1=?  AND  IndexAPart2=?  AND  IndexAPart3=? %s  OFFSET ? Rows FETCH NEXT ? Rows ONLY", orderBy)

	err = db.Query(&objs, query, IndexAPart1, IndexAPart2, IndexAPart3, offset, limit)
	return
}

func (m *_PeopleMgr) FindOneByUniquePart1UniquePart2(UniquePart1 int32, UniquePart2 int32) (*People, error) {
	query := "SELECT NonIndexA, NonIndexB, PeopleId, Age, Name, IndexAPart1, IndexAPart2, IndexAPart3, UniquePart1, UniquePart2 FROM People WHERE UniquePart1=? AND UniquePart2=?"
	var obj People
	err := db.Query(&obj, query, UniquePart1, UniquePart2)
	return &obj, err
}

func (m *_PeopleMgr) FindByAge(Age int32, offset int, limit int, sortFields ...string) (objs []*People, err error) {
	orderBy := "ORDER BY %s"
	if len(sortFields) != 0 {
		orderBy = fmt.Sprintf(orderBy, strings.Join(sortFields, ","))
	} else {
		orderBy = fmt.Sprintf(orderBy, "PeopleId")
	}

	query := fmt.Sprintf("SELECT NonIndexA, NonIndexB, PeopleId, Age, Name, IndexAPart1, IndexAPart2, IndexAPart3, UniquePart1, UniquePart2 FROM People WHERE Age=? %s  OFFSET ? Rows FETCH NEXT ? Rows ONLY", orderBy)

	err = db.Query(&objs, query, Age, offset, limit)
	return
}

func (m *_PeopleMgr) FindOneByName(Name string) (*People, error) {
	query := "SELECT NonIndexA, NonIndexB, PeopleId, Age, Name, IndexAPart1, IndexAPart2, IndexAPart3, UniquePart1, UniquePart2 FROM People WHERE Name=?"
	var obj People
	err := db.Query(&obj, query, Name)
	return &obj, err
}

func (m *_PeopleMgr) FindOne(where string, args ...interface{}) (*People, error) {
	query := getQuerysql(true, where)
	var obj People
	err := db.Query(&obj, query, args...)
	return &obj, err
}

func (m *_PeopleMgr) Find(where string, args ...interface{}) (results []*People, err error) {
	query := getQuerysql(false, where)
	err = db.Query(&results, query, args...)
	return
}

func (m *_PeopleMgr) FindAll() (results []*People, err error) {
	return m.Find("")
}

func (m *_PeopleMgr) FindWithOffset(where string, offset int, limit int, args ...interface{}) (results []*People, err error) {
	query := getQuerysql(false, where)

	if !strings.Contains(strings.ToLower(where), "ORDER BY") {
		where = " ORDER BY Name"
	}
	query = query + where + " OFFSET ? Rows FETCH NEXT ? Rows ONLY"
	args = append(args, offset)
	args = append(args, limit)

	err = db.Query(&results, query, args...)
	return
}

func getQuerysql(topOne bool, where string) string {
	query := `SELECT `
	if topOne {
		query = query + ` TOP 1 `
	}
	query = query + ` NonIndexA, NonIndexB, PeopleId, Age, Name, IndexAPart1, IndexAPart2, IndexAPart3, UniquePart1, UniquePart2 FROM dbo.[People] WITH(NOLOCK) `

	if where != "" {
		if strings.Index(strings.Trim(where, " "), "WHERE") == -1 {
			where = " WHERE " + where
		}
		query = query + where
	}
	return query
}

func (m *_PeopleMgr) Del(where string, params ...interface{}) (sql.Result, error) {
	query := "delete from People"
	if where != "" {
		query = fmt.Sprintf("delete from People where " + where)
	}
	return db.Exec(query, params...)
}

// argument example:
// set:"a=?, b=?"
// where:"c=? and d=?"
// params:[]interface{}{"a", "b", "c", "d"}...
func (m *_PeopleMgr) Update(set, where string, params ...interface{}) (sql.Result, error) {
	query := fmt.Sprintf("update People set %s", set)
	if where != "" {
		query = fmt.Sprintf("update People set %s where %s", set, where)
	}
	return db.Exec(query, params...)
}
