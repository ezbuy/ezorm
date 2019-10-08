package test

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/ezbuy/ezorm/db"
)

func requestTimeLogger(queryer db.Queryer, query string, args ...interface{}) db.Queryer {
	return func(query string, args ...interface{}) (interface{}, error) {
		start := time.Now()
		defer func() {
			fmt.Printf("query time: %.6f seconds\n", time.Now().Sub(start).Seconds())
		}()

		time.Sleep(time.Millisecond * time.Duration(rand.Int31n(100)))
		return queryer(query, args...)
	}
}

var (
	host     = os.Getenv("MSSQL_HOST")
	userId   = os.Getenv("MSSQL_USER")
	password = os.Getenv("MSSQL_PASSWORD")
	database = os.Getenv("MSSQL_DATABASE")
)

func init() {
	dsn := fmt.Sprintf("server=%s;user id=%s;password=%s;DATABASE=master",
		host, userId, password)
	MssqlSetUp(dsn)

	MssqlSetMaxOpenConns(255)
	MssqlSetMaxIdleConns(255)

	MssqlAddQueryWrapper(requestTimeLogger)

	query, err := ioutil.ReadFile("People.sql")
	if err != nil {
		panic(err)
	}

	if _, err := mssqlExec("CREATE DATABASE test"); err != nil {
		panic(err)
	}

	if _, err := mssqlExec(string(query)); err != nil {
		panic(err)
	}

	dsn = fmt.Sprintf("server=%s;user id=%s;password=%s;DATABASE=%s",
		host, userId, password, database)
	MssqlSetUp(dsn)

}

func savePeople(t *testing.T, name string) (*People, error) {
	now := time.Now()
	p := &People{
		Name:        name,
		Age:         1,
		UniquePart1: rand.Int31n(1000000),
		UniquePart2: rand.Int31n(1000000),
		CreateDate:  &now,
		UpdateDate:  &now,
	}

	res, err := PeopleMgr.Save(p)
	if err != nil {
		t.Fatal(err)
	}

	_, err = res.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}

	return p, err
}

func saveDuplicatedPeople(p1, p2 int32, t *testing.T) (*People, error) {

	now := time.Now()
	p := &People{
		Name:        "dup",
		Age:         1,
		UniquePart1: p1,
		UniquePart2: p2,
		CreateDate:  &now,
		UpdateDate:  &now,
	}

	_, err := PeopleMgr.Save(p)
	if err == nil {
		return nil, errors.New("dup: save expected duplicate error,but got nil")
	}
	return nil, nil
}

func assertPeopleEqual(a, b *People, t *testing.T) {
	if a.Age != b.Age || a.Name != b.Name {
		t.Errorf("%#v != %#v. people not equal", a, b)
	}
}

func newUnsavedPeople() *People {
	now := time.Now()
	return &People{
		Name:        fmt.Sprintf("testname_%d", time.Now().Nanosecond()),
		Age:         rand.Int31n(200),
		NonIndexA:   fmt.Sprintf("testname_%d", time.Now().Nanosecond()),
		NonIndexB:   fmt.Sprintf("testname_%d", time.Now().Nanosecond()),
		IndexAPart1: rand.Int63n(1000000),
		IndexAPart2: rand.Int31n(1000000),
		IndexAPart3: rand.Int31n(1000000),
		UniquePart1: rand.Int31n(1000000),
		UniquePart2: rand.Int31n(1000000),
		CreateDate:  &now,
		UpdateDate:  &now,
	}
}

func TestSaveInsert(t *testing.T) {
	_, err := PeopleMgr.Del("")
	if err != nil {
		t.Errorf("delete error:%s", err.Error())
	}

	p, err := savePeople(t, "testuser")
	if err != nil {
		t.Errorf("save err:%s", err.Error())
	}

	if p.PeopleId != 1 {
		t.Fatalf("1. TestSaveInsert: expect 1 but get %d", p.PeopleId)
	}

	if _, err := saveDuplicatedPeople(p.UniquePart1, p.UniquePart2, t); err != nil {
		t.Fatalf("1. TestSaveDuplicateInsert: %q", err)
	}

}

func TestSaveUpdate(t *testing.T) {
	_, err := PeopleMgr.Del("")
	if err != nil {
		t.Errorf("delete error:%s", err.Error())
	}

	p, err := savePeople(t, "testuser")
	if err != nil {
		t.Errorf("save err:%s", err.Error())
	}

	if p.PeopleId != 2 {
		t.Fatalf("2. TestSaveUpdate: expect 2 but get %d", p.PeopleId)
	}

	p.Age = p.Age + 1
	result, err := PeopleMgr.Save(p)
	if err != nil {
		t.Errorf("save err:%s", err.Error())
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		t.Errorf("save err:%s", err.Error())
	}
	if affectedRows != 1 {
		t.Errorf("save err:affectedRows[%d]!=1", affectedRows)
	}

	pFound, err := PeopleMgr.FindOne("")
	if err != nil {
		t.Errorf("find one error:%s", err.Error())
	}
	assertPeopleEqual(p, pFound, t)
}

func TestFindOne(t *testing.T) {
	_, err := PeopleMgr.Del("")
	if err != nil {
		t.Errorf("delete error:%s", err.Error())
	}

	p, err := savePeople(t, "testuser")
	if err != nil {
		t.Errorf("save err:%s", err.Error())
	}

	if p.PeopleId != 3 {
		t.Fatalf("3. TestFindOne: expect 3 but get %d", p.PeopleId)
	}

	pFound, err := PeopleMgr.FindOne("")
	if err != nil {
		t.Errorf("find one error:%s", err.Error())
	}
	assertPeopleEqual(p, pFound, t)
}

func TestFind(t *testing.T) {
	_, err := PeopleMgr.Del("")
	if err != nil {
		t.Errorf("delete error:%s", err.Error())
	}

	p, err := savePeople(t, "testuser")
	if err != nil {
		t.Errorf("save err:%s", err.Error())
	}

	if p.PeopleId != 4 {
		t.Fatalf("4. TestFindOne: expect 4 but get %d", p.PeopleId)
	}

	pSlice, err := PeopleMgr.Find("")
	if err != nil {
		t.Errorf("find error:%s", err.Error())
	}

	if len(pSlice) == 0 {
		t.Errorf("fail to find people")
	}
	assertPeopleEqual(p, pSlice[0], t)
}

func TestFindAll(t *testing.T) {
	_, err := PeopleMgr.Del("")
	if err != nil {
		t.Errorf("delete error:%s", err.Error())
	}

	p, err := savePeople(t, "testuser1")
	if err != nil {
		t.Errorf("save err:%s", err.Error())
	}

	if p.PeopleId != 5 {
		t.Fatalf("5. TestFindAll: expect 5 but get %d", p.PeopleId)
	}

	p, err = savePeople(t, "testuser2")
	if err != nil {
		t.Errorf("save err:%s", err.Error())
	}

	if p.PeopleId != 6 {
		t.Fatalf("6. TestFindAll: expect 6 but get %d", p.PeopleId)
	}

	pSlice, err := PeopleMgr.FindAll()
	if err != nil {
		t.Errorf("find all error:%s", err.Error())
	}

	if len(pSlice) != 2 {
		t.Errorf("FindAll result incorrect, len(result)=%d, not 2", len(pSlice))
	}
}

func TestFindWithOffset(t *testing.T) {
	_, err := PeopleMgr.Del("")
	if err != nil {
		t.Errorf("delete error:%s", err.Error())
	}

	for i := 0; i < 5; i++ {
		_, err = savePeople(t, fmt.Sprint(i))
		if err != nil {
			t.Errorf("save err:%s", err.Error())
		}
	}

	pSlice, err := PeopleMgr.FindWithOffset("ORDER BY PeopleId", 0, 4)
	if err != nil {
		t.Errorf("FindWithOffset error:%s", err.Error())
	}

	if len(pSlice) != 4 {
		t.Errorf("FindWithOffset result incorrect, len(result)[%d]!=4", len(pSlice))
	}

	pSlice, err = PeopleMgr.FindWithOffset("ORDER BY PeopleId", 3, 4)
	if err != nil {
		t.Errorf("FindWithOffset error:%s", err.Error())
	}

	if len(pSlice) != 2 {
		t.Errorf("FindWithOffset result incorrect, len(result)[%d]!=2", len(pSlice))
	}
}

func TestDel(t *testing.T) {
	_, err := PeopleMgr.Del("")
	if err != nil {
		t.Errorf("delete error:%s", err.Error())
	}

	p1, err := savePeople(t, "testuser1")
	if err != nil {
		t.Errorf("save err:%s", err.Error())
	}

	p2, err := savePeople(t, "testuser2")
	if err != nil {
		t.Errorf("save err:%s", err.Error())
	}

	pSlice, err := PeopleMgr.FindAll()
	if err != nil {
		t.Errorf("find all error:%s", err.Error())
	}

	if len(pSlice) != 2 {
		t.Errorf("FindAll result incorrect, len(result)=%d, not 2", len(pSlice))
	}

	_, err = PeopleMgr.Del("Name=?", p1.Name)
	if err != nil {
		t.Errorf("delete error:%s", err.Error())
	}

	pSlice, err = PeopleMgr.FindAll()
	if err != nil {
		t.Errorf("find all error:%s", err.Error())
	}

	if len(pSlice) != 1 {
		t.Errorf("FindAll result incorrect, len(result)[%d]!=1", len(pSlice))
	}
	assertPeopleEqual(p2, pSlice[0], t)
}

func TestUpdate(t *testing.T) {
	_, err := PeopleMgr.Del("")
	if err != nil {
		t.Errorf("delete error:%s", err.Error())
	}

	p, err := savePeople(t, "testuser")
	if err != nil {
		t.Errorf("save err:%s", err.Error())
	}

	newName := fmt.Sprintf("newname_%d", time.Now().Nanosecond())

	result, err := PeopleMgr.Update("Name=?", "Name=?", newName, p.Name)
	if err != nil {
		t.Errorf("update error:%s", err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Errorf("update error:%s", err.Error())
	}
	if rowsAffected != 1 {
		t.Errorf("update error:rowsAffected[%d]!=1", rowsAffected)
	}

	pNew, err := PeopleMgr.FindOne("")
	if err != nil {
		t.Errorf("find one error:%s", err.Error())
	}

	if pNew.Name != newName {
		t.Errorf("update error:%s!=%s, pNew.Name!=newName", pNew.Name, newName)
	}
}

func TestFindByID(t *testing.T) {
	_, err := PeopleMgr.Del("")
	if err != nil {
		t.Errorf("delete error:%s", err.Error())
	}

	p1, err := savePeople(t, "testuser")
	if err != nil {
		t.Errorf("save err:%s", err.Error())
	}

	p2, err := PeopleMgr.FindByID(p1.PeopleId)
	if err != nil {
		t.Errorf("findByID err:%s", err.Error())
	}
	assertPeopleEqual(p1, p2, t)
}

func TestFindOneByName(t *testing.T) {
	_, err := PeopleMgr.Del("")
	if err != nil {
		t.Errorf("delete error:%s", err.Error())
	}

	p1, err := savePeople(t, "testuser")
	if err != nil {
		t.Errorf("save err:%s", err.Error())
	}

	p2, err := PeopleMgr.FindOneByName(p1.Name)
	if err != nil {
		t.Errorf("FindOneByName err:%s", err.Error())
	}
	assertPeopleEqual(p1, p2, t)
}

func TestFindByAge(t *testing.T) {
	_, err := PeopleMgr.Del("")
	if err != nil {
		t.Errorf("delete error:%s", err.Error())
	}

	rand.Seed(time.Now().UnixNano())

	p1 := newUnsavedPeople()

	_, err = PeopleMgr.Save(p1)
	if err != nil {
		t.Errorf("save error:%s", err.Error())
	}

	p2 := newUnsavedPeople()
	p2.Age = p1.Age + 1

	_, err = PeopleMgr.Save(p2)
	if err != nil {
		t.Errorf("save error:%s", err.Error())
	}

	ps, err := PeopleMgr.FindByAge(p1.Age, 0, 100)
	if err != nil {
		t.Errorf("FindByAge err:%s", err.Error())
	} else if len(ps) != 1 {
		t.Errorf("FindByAge incorrect:len(ps)[%d]!=1", len(ps))
	}

	assertPeopleEqual(p1, ps[0], t)

	ps, err = PeopleMgr.FindByAge(p2.Age, 0, 100)
	if err != nil {
		t.Errorf("FindByAge err:%s", err.Error())
	} else if len(ps) != 1 {
		t.Errorf("FindByAge incorrect:len(ps)[%d]!=1", len(ps))
	}

	assertPeopleEqual(p2, ps[0], t)
}

func TestFindByIndexAPart1IndexAPart2IndexAPart3(t *testing.T) {
	_, err := PeopleMgr.Del("")
	if err != nil {
		t.Errorf("delete error:%s", err.Error())
	}

	rand.Seed(time.Now().UnixNano())

	p1 := newUnsavedPeople()

	_, err = PeopleMgr.Save(p1)
	if err != nil {
		t.Errorf("save error:%s", err.Error())
	}

	p2 := newUnsavedPeople()

	_, err = PeopleMgr.Save(p2)
	if err != nil {
		t.Errorf("save error:%s", err.Error())
	}

	ps, err := PeopleMgr.FindByIndexAPart1IndexAPart2IndexAPart3(p1.IndexAPart1, p1.IndexAPart2, p1.IndexAPart3, 0, 100)
	if err != nil {
		t.Errorf("FindByIndexAPart1IndexAPart2IndexAPart3 err:%s", err.Error())
	} else if len(ps) != 1 {
		t.Errorf("FindByIndexAPart1IndexAPart2IndexAPart3 incorrect:len(ps)[%d]!=1", len(ps))
	}

	assertPeopleEqual(p1, ps[0], t)

	ps, err = PeopleMgr.FindByIndexAPart1IndexAPart2IndexAPart3(p2.IndexAPart1, p2.IndexAPart2, p2.IndexAPart3, 0, 100)
	if err != nil {
		t.Errorf("FindByIndexAPart1IndexAPart2IndexAPart3 err:%s", err.Error())
	} else if len(ps) != 1 {
		t.Errorf("FindByIndexAPart1IndexAPart2IndexAPart3 incorrect:len(ps)[%d]!=1", len(ps))
	}

	assertPeopleEqual(p2, ps[0], t)
}

func TestFindByIDs(t *testing.T) {
	_, err := PeopleMgr.Del("")
	if err != nil {
		t.Errorf("delete error:%s", err.Error())
	}

	rand.Seed(time.Now().UnixNano())

	p1 := newUnsavedPeople()

	_, err = PeopleMgr.Save(p1)
	if err != nil {
		t.Errorf("save error:%s", err.Error())
	}

	p2 := newUnsavedPeople()

	_, err = PeopleMgr.Save(p2)
	if err != nil {
		t.Errorf("save error:%s", err.Error())
	}

	ids := []int32{p1.PeopleId, p2.PeopleId}

	ps, err := PeopleMgr.FindByIDs(ids)
	if err != nil {
		t.Errorf("FindByIds err:%s", err.Error())
	} else if len(ps) != 2 {
		t.Errorf("FindByIds incorrect:len(ps)[%d]!=2", len(ps))
	}
}

func TestCount(t *testing.T) {
	_, err := PeopleMgr.Del("")
	if err != nil {
		t.Errorf("delete error:%s", err.Error())
	}

	rand.Seed(time.Now().UnixNano())

	p1 := newUnsavedPeople()

	_, err = PeopleMgr.Save(p1)
	if err != nil {
		t.Errorf("save error:%s", err.Error())
	}

	p2 := newUnsavedPeople()

	_, err = PeopleMgr.Save(p2)
	if err != nil {
		t.Errorf("save error:%s", err.Error())
	}

	count, err := PeopleMgr.Count("")
	if err != nil {
		t.Errorf("TestCount err:%s", err.Error())
	} else if count != 2 {
		t.Errorf("TestCount incorrect:count[%d]!=2", count)
	}

	count, err = PeopleMgr.Count("PeopleId=?", p1.PeopleId)
	if err != nil {
		t.Errorf("TestCount err:%s", err.Error())
	} else if count != 1 {
		t.Errorf("TestCount incorrect:count[%d]!=1", count)
	}
}

func TestInsertBatch(t *testing.T) {
	_, err := PeopleMgr.Del("")
	if err != nil {
		t.Errorf("delete error:%s", err.Error())
	}

	p1 := newUnsavedPeople()

	p2 := newUnsavedPeople()

	_, err = PeopleMgr.InsertBatch([]*People{p1, p2})
	if err != nil {
		t.Errorf("InsertBatch err:%v", err)
	}

}

func TestErrNoRows(t *testing.T) {
	_, err := PeopleMgr.Del("")
	if err != nil {
		t.Errorf("delete error:%s", err.Error())
	}

	_, err = PeopleMgr.FindByID(0)
	if err != sql.ErrNoRows {
		t.Errorf("error:[%v] not sql.ErrNoRows", err)
	}
}

func TestFindByIds(t *testing.T) {
	_, err := PeopleMgr.Del("")
	if err != nil {
		t.Errorf("delete error:%s", err.Error())
	}

	maxNum := 11000
	ids := make([]int32, 0, maxNum)
	for i := 0; i < maxNum; i++ {
		ids = append(ids, int32(i))
	}

	vals, err := PeopleMgr.FindByIDs(ids)
	if err != nil {
		t.Errorf("error:[%v]", err)
	}

	fmt.Printf("vals len:%d", len(vals))
}

func TestGetId2Obj(t *testing.T) {
	_, err := PeopleMgr.Del("")
	if err != nil {
		t.Errorf("delete error:%s", err.Error())
	}

	rand.Seed(time.Now().UnixNano())

	p1 := newUnsavedPeople()

	_, err = PeopleMgr.Save(p1)
	if err != nil {
		t.Errorf("save error:%s", err.Error())
	}

	p2 := newUnsavedPeople()

	_, err = PeopleMgr.Save(p2)
	if err != nil {
		t.Errorf("save error:%s", err.Error())
	}

	ids := []int32{p1.PeopleId, p2.PeopleId}

	ps, err := PeopleMgr.FindByIDs(ids)
	if err != nil {
		t.Errorf("FindByIds err:%s", err.Error())
	} else if len(ps) != 2 {
		t.Errorf("FindByIds incorrect:len(ps)[%d]!=2", len(ps))
	}

	id2obj := PeopleMgr.GetId2Obj(ps)
	for _, each := range ps {
		if id2obj[each.PeopleId] != each {
			t.Error("id2obj[each.PeopleId] != each")
		}
	}
}

func TestGetIds(t *testing.T) {
	_, err := PeopleMgr.Del("")
	if err != nil {
		t.Errorf("delete error:%s", err.Error())
	}

	rand.Seed(time.Now().UnixNano())

	p1 := newUnsavedPeople()

	_, err = PeopleMgr.Save(p1)
	if err != nil {
		t.Errorf("save error:%s", err.Error())
	}

	p2 := newUnsavedPeople()

	_, err = PeopleMgr.Save(p2)
	if err != nil {
		t.Errorf("save error:%s", err.Error())
	}

	ids := []int32{p1.PeopleId, p2.PeopleId}

	ps, err := PeopleMgr.FindByIDs(ids)
	if err != nil {
		t.Errorf("FindByIds err:%s", err.Error())
	} else if len(ps) != 2 {
		t.Errorf("FindByIds incorrect:len(ps)[%d]!=2", len(ps))
	}

	ids = PeopleMgr.GetIds(ps)
	for i, each := range ps {
		if ids[i] != each.PeopleId {
			t.Error("ids[i] != each.PeopleId")
		}
	}
}
