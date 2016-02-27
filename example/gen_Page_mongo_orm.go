package page

import (
	//3rd party libs
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	//Own libs
	"github.com/ezbuy/ezorm/db"
	. "github.com/ezbuy/ezorm/orm"

	
)

func init() {
	
	
	RegisterEzOrmObjByID("page", "Page", newPageFindByID)
	RegisterEzOrmObjRemove("page", "Page", PageMgr.RemoveByID)
	
}





func newPageFindByID(id string) (result EzOrmObj, err error) {
	return PageMgr.FindByID(id)
}



//mongo methods
var (
	insertCB_Page []func(obj EzOrmObj)
	updateCB_Page []func(obj EzOrmObj)
)

func PageAddInsertCallback(cb func(obj EzOrmObj)) {
	insertCB_Page = append(insertCB_Page, cb)
}

func PageAddUpdateCallback(cb func(obj EzOrmObj)) {
	updateCB_Page = append(updateCB_Page, cb)
}

func (o *Page) Id() string {
	return o.ID.Hex()
}

func (o *Page) Save() (info *mgo.ChangeInfo, err error) {
	session, col := PageMgr.GetCol()
	defer session.Close()

	isNew := o.isNew

	
	
	
	
	
	
	
	
	
	
	
	
	

	info, err = col.UpsertId(o.ID, o)
	o.isNew = false

	

	if isNew {
		PageInsertCallback(o)
	} else {
		PageUpdateCallback(o)
	}

	return
}

func (o *Page) InsertUnique(query interface{}) (saved bool, err error) {
	session, col := PageMgr.GetCol()
	defer session.Close()

	info, err := col.Upsert(query, db.M{"$setOnInsert": o})
	if err != nil {
		return
	}
	if info.Updated == 0 {
		saved = true
	}
	o.isNew = false
	if saved {
		PageInsertCallback(o)
	}
	return
}

func PageInsertCallback(o *Page) {
	for _, cb := range insertCB_Page {
		cb(o)
	}
}

func PageUpdateCallback(o *Page) {
	for _, cb := range updateCB_Page {
		cb(o)
	}
}





//foreigh keys


	

	

	

	

	

	




//Collection Manage methods

func (o *_PageMgr) FindOne(query interface{}, sortFields ...string) (result *Page, err error) {
	session, col := PageMgr.GetCol()
	defer session.Close()

	q := col.Find(query)

	_PageSort(q, sortFields)

	err = q.One(&result)
	return
}

func _PageSort(q *mgo.Query, sortFields []string) {
	sortFields = XSortFieldsFilter(sortFields)
	if len(sortFields) > 0 {
		q.Sort(sortFields...)
		return
	}

	q.Sort("-_id")
}

func (o *_PageMgr) Query(query interface{}, limit, offset int, sortFields []string) (*mgo.Session, *mgo.Query) {
	session, col := PageMgr.GetCol()
	q := col.Find(query)
	if limit > 0 {
		q.Limit(limit)
	}
	if offset > 0 {
		q.Skip(offset)
	}
	_PageSort(q, sortFields)
	return session, q
}









func (o *_PageMgr) FindBySlug(Slug string) (result *Page, err error) {
	query := db.M{
		"Slug": Slug,
	}
	session, q := PageMgr.Query(query, 1, 0, nil)
	defer session.Close()
	err = q.One(&result)
	return
}

func (o *_PageMgr) MustGetBySlug(Slug string) (result *Page) {
	result, _ = o.FindBySlug(Slug)
	if result == nil {
		result = PageMgr.NewPage()
		result.Slug = Slug
		result.Save()
	}
	return
}







func (o *_PageMgr) Find(query interface{}, limit int, offset int, sortFields ...string) (result []*Page, err error) {
	session, q := PageMgr.Query(query, limit, offset, sortFields)
	defer session.Close()
	err = q.All(&result)
	return
}

func (o *_PageMgr) FindAll(query interface{}, sortFields ...string) (result []*Page, err error) {
	session, q := PageMgr.Query(query, -1, -1, sortFields)
	defer session.Close()
	err = q.All(&result)
	return
}

func (o *_PageMgr) Has(query interface{}) bool {
	session, col := PageMgr.GetCol()
	defer session.Close()

	var ret interface{}
	err := col.Find(query).One(&ret)
	if err != nil || ret == nil {
		return false
	}
	return true
}

func (o *_PageMgr) Count(query interface{}) (result int) {
	session, col := PageMgr.GetCol()
	defer session.Close()

	result, _ = col.Find(query).Count()
	return
}

func (o *_PageMgr) FindByIDs(id []string, sortFields ...string) (result []*Page, err error) {
	ids := make([]bson.ObjectId, 0, len(id))
	for _, i := range id {
		if bson.IsObjectIdHex(i) {
			ids = append(ids, bson.ObjectIdHex(i))
		}
	}
	return PageMgr.FindAll(db.M{"_id": db.M{"$in": ids}}, sortFields...)
}

func (m *_PageMgr) FindByID(id string) (result *Page, err error) {
	session, col := PageMgr.GetCol()
	defer session.Close()

	if !bson.IsObjectIdHex(id) {
		err = mgo.ErrNotFound
		return
	}
	err = col.FindId(bson.ObjectIdHex(id)).One(&result)
	return
}

func (m *_PageMgr) RemoveAll(query interface{}) (info *mgo.ChangeInfo, err error) {
	session, col := PageMgr.GetCol()
	defer session.Close()

	return col.RemoveAll(query)
}

func (m *_PageMgr) RemoveByID(id string) (err error) {
	session, col := PageMgr.GetCol()
	defer session.Close()

	if !bson.IsObjectIdHex(id) {
		err = mgo.ErrNotFound
		return
	}
	err = col.RemoveId(bson.ObjectIdHex(id))
	
	return
}

func (m *_PageMgr) GetCol() (session *mgo.Session, col *mgo.Collection) {
	return db.GetCol("page.Page")
}






//Search


func (o *Page) IsSearchEnabled() bool {

	return false

}

//end search


