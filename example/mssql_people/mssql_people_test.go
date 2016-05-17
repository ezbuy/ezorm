package people

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func init() {
	dsn := fmt.Sprintf("server=%s;user id=%s;password=%s;DATABASE=%s",
		host, userId, password, database)
	MssqlSetUp(dsn)

	MssqlSetMaxOpenConns(255)
	MssqlSetMaxIdleConns(255)
}

func savePeople(name string) (*People, error) {
	p := &People{
		Name:        name,
		Age:         1,
		UniquePart1: rand.Int31n(1000000),
		UniquePart2: rand.Int31n(1000000),
	}

	_, err := PeopleMgr.Save(p)
	return p, err
}

func TestSaveInsert(t *testing.T) {
	_, err := PeopleMgr.Del("")
	if err != nil {
		t.Errorf("delete error:%s", err.Error())
	}

	_, err = savePeople("testuser")
	if err != nil {
		t.Errorf("save err:%s", err.Error())
	}
}

func TestSaveUpdate(t *testing.T) {
	_, err := PeopleMgr.Del("")
	if err != nil {
		t.Errorf("delete error:%s", err.Error())
	}

	p, err := savePeople("testuser")
	if err != nil {
		t.Errorf("save err:%s", err.Error())
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

func assertPeopleEqual(a, b *People, t *testing.T) {
	if a.Age != b.Age || a.Name != b.Name {
		t.Errorf("%#v != %#v. people not equal", a, b)
	}
}

func TestFindOne(t *testing.T) {
	_, err := PeopleMgr.Del("")
	if err != nil {
		t.Errorf("delete error:%s", err.Error())
	}

	p, err := savePeople("testuser")
	if err != nil {
		t.Errorf("save err:%s", err.Error())
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

	p, err := savePeople("testuser")
	if err != nil {
		t.Errorf("save err:%s", err.Error())
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

	_, err = savePeople("testuser1")
	if err != nil {
		t.Errorf("save err:%s", err.Error())
	}

	_, err = savePeople("testuser2")
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
}

func TestFindWithOffset(t *testing.T) {
	_, err := PeopleMgr.Del("")
	if err != nil {
		t.Errorf("delete error:%s", err.Error())
	}

	for i := 0; i < 5; i++ {
		_, err = savePeople(fmt.Sprint(i))
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

	p1, err := savePeople("testuser1")
	if err != nil {
		t.Errorf("save err:%s", err.Error())
	}

	p2, err := savePeople("testuser2")
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

	p, err := savePeople("testuser")
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

	p1, err := savePeople("testuser")
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

	p1, err := savePeople("testuser")
	if err != nil {
		t.Errorf("save err:%s", err.Error())
	}

	p2, err := PeopleMgr.FindOneByName(p1.Name)
	if err != nil {
		t.Errorf("FindOneByName err:%s", err.Error())
	}
	assertPeopleEqual(p1, p2, t)
}

func newUnsavedPeople() *People {
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
	}
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
