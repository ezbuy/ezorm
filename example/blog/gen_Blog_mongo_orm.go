package blog

import (
	//3rd party libs
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	//Own libs
	"github.com/ezbuy/ezorm/db"
	. "github.com/ezbuy/ezorm/orm"

	
)

func init() {
	
	db.SetOnFinishInit(initBlogIndex)
	
	
	RegisterEzOrmObjByID("blog", "Blog", newBlogFindByID)
	RegisterEzOrmObjRemove("blog", "Blog", BlogMgr.RemoveByID)
	
}




func initBlogIndex() {
	session, collection := BlogMgr.GetCol()
	defer session.Close()
	
reEnsureUserIsPublished:
	if err := collection.EnsureIndex(mgo.Index{
		Key: []string{"User","IsPublished"},
		Sparse: true,
	}); err != nil {
		println("error ensureIndex Blog UserIsPublished", err.Error())
		err = collection.DropIndex("UserIsPublished")
		if err != nil {
			panic(err)
		}
		goto reEnsureUserIsPublished

	}
	
reEnsureSlug:
	if err := collection.EnsureIndex(mgo.Index{
		Key: []string{"Slug"},
		Unique: true,
		Sparse: true,
	}); err != nil {
		println("error ensureIndex Blog Slug", err.Error())
		err = collection.DropIndex("Slug")
		if err != nil {
			panic(err)
		}
		goto reEnsureSlug

	}
	
reEnsureIsPublished:
	if err := collection.EnsureIndex(mgo.Index{
		Key: []string{"IsPublished"},
		Sparse: true,
	}); err != nil {
		println("error ensureIndex Blog IsPublished", err.Error())
		err = collection.DropIndex("IsPublished")
		if err != nil {
			panic(err)
		}
		goto reEnsureIsPublished

	}
	
}


func newBlogFindByID(id string) (result EzOrmObj, err error) {
	return BlogMgr.FindByID(id)
}



//mongo methods
var (
	insertCB_Blog []func(obj EzOrmObj)
	updateCB_Blog []func(obj EzOrmObj)
)

func BlogAddInsertCallback(cb func(obj EzOrmObj)) {
	insertCB_Blog = append(insertCB_Blog, cb)
}

func BlogAddUpdateCallback(cb func(obj EzOrmObj)) {
	updateCB_Blog = append(updateCB_Blog, cb)
}

func (o *Blog) Id() string {
	return o.ID.Hex()
}

func (o *Blog) Save() (info *mgo.ChangeInfo, err error) {
	session, col := BlogMgr.GetCol()
	defer session.Close()

	isNew := o.isNew

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	

	info, err = col.UpsertId(o.ID, o)
	o.isNew = false

	

	if isNew {
		BlogInsertCallback(o)
	} else {
		BlogUpdateCallback(o)
	}

	return
}

func (o *Blog) InsertUnique(query interface{}) (saved bool, err error) {
	session, col := BlogMgr.GetCol()
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
		BlogInsertCallback(o)
	}
	return
}

func BlogInsertCallback(o *Blog) {
	for _, cb := range insertCB_Blog {
		cb(o)
	}
}

func BlogUpdateCallback(o *Blog) {
	for _, cb := range updateCB_Blog {
		cb(o)
	}
}





//foreigh keys


	

	

	

	

	

	

	




//Collection Manage methods

func (o *_BlogMgr) FindOne(query interface{}, sortFields ...string) (result *Blog, err error) {
	session, col := BlogMgr.GetCol()
	defer session.Close()

	q := col.Find(query)

	_BlogSort(q, sortFields)

	err = q.One(&result)
	return
}

func _BlogSort(q *mgo.Query, sortFields []string) {
	sortFields = XSortFieldsFilter(sortFields)
	if len(sortFields) > 0 {
		q.Sort(sortFields...)
		return
	}

	q.Sort("-_id")
}

func (o *_BlogMgr) Query(query interface{}, limit, offset int, sortFields []string) (*mgo.Session, *mgo.Query) {
	session, col := BlogMgr.GetCol()
	q := col.Find(query)
	if limit > 0 {
		q.Limit(limit)
	}
	if offset > 0 {
		q.Skip(offset)
	}
	_BlogSort(q, sortFields)
	return session, q
}
func (o *_BlogMgr) FindByUserIsPublished(User int32, IsPublished bool, limit int, offset int, sortFields ...string) (result []*Blog, err error) {
	query := db.M{
		"User": User,
		"IsPublished": IsPublished,
	}
	session, q := BlogMgr.Query(query, limit, offset, sortFields)
	defer session.Close()
	err = q.All(&result)
	return
}
func (o *_BlogMgr) FindOneBySlug(Slug string) (result *Blog, err error) {
	query := db.M{
		"Slug": Slug,
	}
	session, q := BlogMgr.Query(query, 1, 0, nil)
	defer session.Close()
	err = q.One(&result)
	return
}

func (o *_BlogMgr) MustFindOneBySlug(Slug string) (result *Blog) {
	result, _ = o.FindOneBySlug(Slug)
	if result == nil {
		result = BlogMgr.NewBlog()
		result.Slug = Slug
		result.Save()
	}
	return
}
func (o *_BlogMgr) FindByIsPublished(IsPublished bool, limit int, offset int, sortFields ...string) (result []*Blog, err error) {
	query := db.M{
		"IsPublished": IsPublished,
	}
	session, q := BlogMgr.Query(query, limit, offset, sortFields)
	defer session.Close()
	err = q.All(&result)
	return
}

func (o *_BlogMgr) Find(query interface{}, limit int, offset int, sortFields ...string) (result []*Blog, err error) {
	session, q := BlogMgr.Query(query, limit, offset, sortFields)
	defer session.Close()
	err = q.All(&result)
	return
}

func (o *_BlogMgr) FindAll(query interface{}, sortFields ...string) (result []*Blog, err error) {
	session, q := BlogMgr.Query(query, -1, -1, sortFields)
	defer session.Close()
	err = q.All(&result)
	return
}

func (o *_BlogMgr) Has(query interface{}) bool {
	session, col := BlogMgr.GetCol()
	defer session.Close()

	var ret interface{}
	err := col.Find(query).One(&ret)
	if err != nil || ret == nil {
		return false
	}
	return true
}

func (o *_BlogMgr) Count(query interface{}) (result int) {
	session, col := BlogMgr.GetCol()
	defer session.Close()

	result, _ = col.Find(query).Count()
	return
}

func (o *_BlogMgr) FindByIDs(id []string, sortFields ...string) (result []*Blog, err error) {
	ids := make([]bson.ObjectId, 0, len(id))
	for _, i := range id {
		if bson.IsObjectIdHex(i) {
			ids = append(ids, bson.ObjectIdHex(i))
		}
	}
	return BlogMgr.FindAll(db.M{"_id": db.M{"$in": ids}}, sortFields...)
}

func (m *_BlogMgr) FindByID(id string) (result *Blog, err error) {
	session, col := BlogMgr.GetCol()
	defer session.Close()

	if !bson.IsObjectIdHex(id) {
		err = mgo.ErrNotFound
		return
	}
	err = col.FindId(bson.ObjectIdHex(id)).One(&result)
	return
}

func (m *_BlogMgr) RemoveAll(query interface{}) (info *mgo.ChangeInfo, err error) {
	session, col := BlogMgr.GetCol()
	defer session.Close()

	return col.RemoveAll(query)
}

func (m *_BlogMgr) RemoveByID(id string) (err error) {
	session, col := BlogMgr.GetCol()
	defer session.Close()

	if !bson.IsObjectIdHex(id) {
		err = mgo.ErrNotFound
		return
	}
	err = col.RemoveId(bson.ObjectIdHex(id))
	
	return
}

func (m *_BlogMgr) GetCol() (session *mgo.Session, col *mgo.Collection) {
	return db.GetCol("test_blog")
}






//Search


func (o *Blog) IsSearchEnabled() bool {

	return false

}

//end search


