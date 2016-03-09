package mssql_people

import (
	"fmt"
	"testing"
	"time"

	"github.com/ezbuy/ezorm/db"
)

func init() {
	conf := &db.SqlDbConfig{
		SqlConnStr: "server=localhost;user id=testuser;password=888888;DATABASE=test",
	}
	db.SetDBConfig(conf)
}

func savePeole(t *testing.T) (*People, error) {
	p := &People{
		Name: fmt.Sprintf("test_%d", time.Now().Nanosecond()),
		Age:  1,
	}

	_, err := PeopleMgr.Save(p)
	return p, err
}

func TestSaveInsert(t *testing.T) {
	_, err := PeopleMgr.Del("")
	if err != nil {
		t.Errorf("delete error:%s", err.Error())
	}

	_, err = savePeole(t)
	if err != nil {
		t.Errorf("save err:%s", err.Error())
	}
}

func TestSaveUpdate(t *testing.T) {
	_, err := PeopleMgr.Del("")
	if err != nil {
		t.Errorf("delete error:%s", err.Error())
	}

	p, err := savePeole(t)
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

	p, err := savePeole(t)
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

	p, err := savePeole(t)
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

	_, err = savePeole(t)
	if err != nil {
		t.Errorf("save err:%s", err.Error())
	}

	_, err = savePeole(t)
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
		_, err = savePeole(t)
		if err != nil {
			t.Errorf("save err:%s", err.Error())
		}
	}

	pSlice, err := PeopleMgr.FindWithOffset("", 0, 4)
	if err != nil {
		t.Errorf("FindWithOffset error:%s", err.Error())
	}

	if len(pSlice) != 4 {
		t.Errorf("FindWithOffset result incorrect, len(result)[%d]!=4", len(pSlice))
	}

	pSlice, err = PeopleMgr.FindWithOffset("", 3, 4)
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

	p1, err := savePeole(t)
	if err != nil {
		t.Errorf("save err:%s", err.Error())
	}

	p2, err := savePeole(t)
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

	p, err := savePeole(t)
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
