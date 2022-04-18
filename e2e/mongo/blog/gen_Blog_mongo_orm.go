package blog

import (
	"time"

	//3rd party libs
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	//Own libs
	"github.com/ezbuy/ezorm/v2/db"
	. "github.com/ezbuy/ezorm/v2/pkg/orm"
)

var _ time.Time

func init() {

	db.SetOnEnsureIndex(initBlogIndex)
	addMgoIndexFunc(initBlogIndex)

	RegisterEzOrmObjByID("blog", "Blog", newBlogFindByID)
	RegisterEzOrmObjRemove("blog", "Blog", BlogMgr.RemoveByID)

}

func initBlogIndex() {
	session, collection := BlogMgr.GetCol()
	defer session.Close()

	if err := collection.EnsureIndex(mgo.Index{
		Key:        []string{"User", "IsPublished"},
		Background: true,
		Sparse:     true,
	}); err != nil {
		panic("ensureIndex test.Blog UserIsPublished error:" + err.Error())
	}

	if err := collection.EnsureIndex(mgo.Index{
		Key:        []string{"Slug"},
		Unique:     true,
		Background: true,
		Sparse:     true,
	}); err != nil {
		panic("ensureIndex test.Blog Slug error:" + err.Error())
	}

	if err := collection.EnsureIndex(mgo.Index{
		Key:        []string{"CreateDate"},
		Background: true,
	}); err != nil {
		panic("ensureIndex test.Blog CreateDate error:" + err.Error())
	}

	if err := collection.EnsureIndex(mgo.Index{
		Key:        []string{"IsPublished"},
		Background: true,
		Sparse:     true,
	}); err != nil {
		panic("ensureIndex test.Blog IsPublished error:" + err.Error())
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

// _BlogSort 将排序字段应用到查询对象中，如果找不到有效的排序字段，则默认使用 `-_id` 作为排序字段
func _BlogSort(q *mgo.Query, sortFields []string) {
	sortFields = XSortFieldsFilter(sortFields)
	if len(sortFields) > 0 {
		q.Sort(sortFields...)
		return
	}

	q.Sort("-_id")
}

// Query 按照查询条件、分页、排序等构建 MongoDB 查询对象，默认情况按照插入倒序返回全量数据
//   - 如果 limit 小于等于 0，则忽略该参数
//   - 如果 offset 小于等于 0，则忽略该参数
//   - 如果 sortFields 为空或全为非法值，则使用 `-_id` 作为排序条件（注意：如果表数据量很大，请显式传递该字段，否则会发生慢查询）
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

// NQuery 按照查询条件、分页、排序等构建 MongoDB 查询对象，如果不指定排序字段，则 MongoDB
// 会按照引擎中的存储顺序返回（Natural-Order）， 不保证返回数据保持插入顺序或插入倒序。
// 建议仅在保证返回数据唯一的情况下使用
// Ref: https://docs.mongodb.com/manual/reference/method/cursor.sort/#return-in-natural-order
//   - 如果 limit 小于等于 0，则忽略该参数
//   - 如果 offset 小于等于 0，则忽略该参数
//   - 如果 sortFields 为空或全为非法值，则忽略该参数
func (o *_BlogMgr) NQuery(query interface{}, limit, offset int, sortFields []string) (*mgo.Session, *mgo.Query) {
	session, col := BlogMgr.GetCol()
	q := col.Find(query)
	if limit > 0 {
		q.Limit(limit)
	}
	if offset > 0 {
		q.Skip(offset)
	}

	if sortFields = XSortFieldsFilter(sortFields); len(sortFields) > 0 {
		q.Sort(sortFields...)
	}

	return session, q
}
func (o *_BlogMgr) FindByUserIsPublished(User int32, IsPublished bool, limit int, offset int, sortFields ...string) (result []*Blog, err error) {
	query := db.M{
		"User":        User,
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
	session, q := BlogMgr.NQuery(query, 1, 0, nil)
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

func (o *_BlogMgr) RemoveBySlug(Slug string) (err error) {
	session, col := BlogMgr.GetCol()
	defer session.Close()

	query := db.M{
		"Slug": Slug,
	}
	return col.Remove(query)
}
func (o *_BlogMgr) FindByCreateDate(CreateDate time.Time, limit int, offset int, sortFields ...string) (result []*Blog, err error) {
	query := db.M{
		"CreateDate": CreateDate,
	}
	session, q := BlogMgr.Query(query, limit, offset, sortFields)
	defer session.Close()
	err = q.All(&result)
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
	result, _ = o.CountE(query)
	return
}

func (o *_BlogMgr) CountE(query interface{}) (result int, err error) {
	session, col := BlogMgr.GetCol()
	defer session.Close()

	result, err = col.Find(query).Count()
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
	if mgoInstances == nil {
		return db.GetCol("test", "test_blog")
	}
	return getCol("test", "test_blog")
}

//Search

func (o *Blog) IsSearchEnabled() bool {

	return false

}

//end search
