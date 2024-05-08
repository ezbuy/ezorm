package user

import (
	"context"
	"time"

	"github.com/ezbuy/ezorm/v2/pkg/orm"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// To import `time` package globally to satisfy `time.Time` index in yaml definition
var _ time.Time

var UserBlogIndexes = []mongo.IndexModel{
	{
		Keys: UserBlogIndexKey_UserId,
	},
}

var UserBlogIndexesFunc = func() {
	orm.SetupIndexModel(Col("test_user_blog"), UserBlogIndexes)
}
var UserBlogIndexKey_UserId = bson.D{
	{Key: "UserId", Value: 1},
}

func init() {
	orm.RegisterEzOrmObjByID("mongo_e2e", "UserBlog", newUserBlogFindByID)
	orm.RegisterEzOrmObjRemove("mongo_e2e", "UserBlog", newUserBlogRemoveByID)
}

func newUserBlogFindByID(id string) (result orm.EzOrmObj, err error) {
	return UserBlogMgr.FindByID(context.TODO(), id)
}

func newUserBlogRemoveByID(id string) error {
	return UserBlogMgr.RemoveByID(context.TODO(), id)
}

// =====================================
// INSERT METHODS
// =====================================

var (
	insertUserBlogCBs []func(obj orm.EzOrmObj)
	updateUserBlogCBs []func(obj orm.EzOrmObj)
)

func UserBlogAddInsertCallback(cb func(obj orm.EzOrmObj)) {
	insertUserBlogCBs = append(insertUserBlogCBs, cb)
}

func UserBlogAddUpdateCallback(cb func(obj orm.EzOrmObj)) {
	updateUserBlogCBs = append(updateUserBlogCBs, cb)
}

func (o *UserBlog) Id() string {
	return o.ID.Hex()
}

// FindOneAndSave try to find one doc by `query`,  and then upsert the result with the current object
func (o *UserBlog) FindOneAndSave(ctx context.Context, query interface{}) (*mongo.SingleResult, error) {
	col := UserBlogMgr.GetCol()
	opts := options.FindOneAndUpdate().SetUpsert(true)
	opts.SetReturnDocument(options.After)
	update := bson.M{
		"$set": bson.M{
			UserBlogMgoFieldUserId:  o.UserId,
			UserBlogMgoFieldBlogId:  o.BlogId,
			UserBlogMgoFieldContent: o.Content,
		},
	}
	ret := col.FindOneAndUpdate(ctx, query, update, opts)
	if ret.Err() != nil {
		return nil, ret.Err()
	}
	return ret, nil
}

// Save upserts the document , Save itself is concurrent-safe , but maybe it is not atomic together with other operations, such as `Find`
func (o *UserBlog) Save(ctx context.Context) (*mongo.UpdateResult, error) {
	isNew := o.isNew
	if o.ID == primitive.NilObjectID {
		o.ID = primitive.NewObjectID()
	}

	filter := bson.M{"_id": o.ID}
	update := bson.M{
		"$set": bson.M{
			UserBlogMgoFieldUserId:  o.UserId,
			UserBlogMgoFieldBlogId:  o.BlogId,
			UserBlogMgoFieldContent: o.Content,
		},
	}

	opts := options.Update().SetUpsert(true)
	col := UserBlogMgr.GetCol()
	ret, err := col.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return ret, err
	}
	if ret.UpsertedID != nil {
		if id, ok := ret.UpsertedID.(primitive.ObjectID); ok {
			o.ID = id
		}
	}

	o.isNew = false
	if isNew {
		UserBlogInsertCallback(o)
	} else {
		UserBlogUpdateCallback(o)
	}
	return ret, err
}

func (o *UserBlog) InsertUnique(ctx context.Context, query interface{}) (saved bool, err error) {
	update := bson.M{
		"$setOnInsert": bson.M{
			UserBlogMgoFieldID:      o.ID,
			UserBlogMgoFieldUserId:  o.UserId,
			UserBlogMgoFieldBlogId:  o.BlogId,
			UserBlogMgoFieldContent: o.Content,
		},
	}

	opts := options.Update().SetUpsert(true)
	col := UserBlogMgr.GetCol()
	ret, err := col.UpdateOne(ctx, query, update, opts)
	if err != nil {
		return false, err
	}
	if ret.UpsertedCount != 0 {
		saved = true
	}

	o.isNew = false
	if saved {
		UserBlogInsertCallback(o)
	}
	return saved, nil
}

func UserBlogInsertCallback(o *UserBlog) {
	for _, cb := range insertUserBlogCBs {
		cb(o)
	}
}

func UserBlogUpdateCallback(o *UserBlog) {
	for _, cb := range updateUserBlogCBs {
		cb(o)
	}
}

// =====================================
// FOREIGN KEYS
// =====================================

// =====================================
// COLLECTION
// =====================================

func (o *_UserBlogMgr) FindOne(ctx context.Context, query interface{}, sortFields interface{}) (result *UserBlog, err error) {
	col := o.GetCol()
	opts := options.FindOne()

	if sortFields != nil {
		opts.SetSort(sortFields)
	}

	ret := col.FindOne(ctx, query, opts)
	if err = ret.Err(); err != nil {
		return nil, err
	}
	err = ret.Decode(&result)
	return
}

func (o *_UserBlogMgr) Query(ctx context.Context, query interface{}, limit, offset int, sortFields interface{}) (*mongo.Cursor, error) {
	col := o.GetCol()
	opts := options.Find()

	if limit > 0 {
		opts.SetLimit(int64(limit))
	}
	if offset > 0 {
		opts.SetSkip(int64(offset))
	}
	if sortFields != nil {
		opts.SetSort(sortFields)
	}

	return col.Find(ctx, query, opts)
}

func (o *_UserBlogMgr) FindByUserId(ctx context.Context, UserId uint64, limit int, offset int, sortFields interface{}) (result []*UserBlog, err error) {
	query := bson.M{
		"UserId": UserId,
	}
	cursor, err := o.Query(ctx, query, limit, offset, sortFields)
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &result)
	return
}

func (o *_UserBlogMgr) Find(ctx context.Context, query interface{}, limit int, offset int, sortFields interface{}) (result []*UserBlog, err error) {
	cursor, err := o.Query(ctx, query, limit, offset, sortFields)
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &result)
	return
}

func (o *_UserBlogMgr) FindAll(ctx context.Context, query interface{}, sortFields interface{}) (result []*UserBlog, err error) {
	cursor, err := o.Query(ctx, query, -1, -1, sortFields)
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &result)
	return
}

func (o *_UserBlogMgr) Has(ctx context.Context, query interface{}) bool {
	count, err := o.CountE(ctx, query)
	if err != nil || count == 0 {
		return false
	}
	return true
}

func (o *_UserBlogMgr) Count(ctx context.Context, query interface{}) int {
	count, _ := o.CountE(ctx, query)
	return count
}

func (o *_UserBlogMgr) CountE(ctx context.Context, query interface{}) (int, error) {
	col := o.GetCol()
	count, err := col.CountDocuments(ctx, query)
	return int(count), err
}

func (o *_UserBlogMgr) FindByIDs(ctx context.Context, id []string, sortFields interface{}) (result []*UserBlog, err error) {
	ids := make([]primitive.ObjectID, 0, len(id))
	for _, i := range id {
		if oid, err := primitive.ObjectIDFromHex(i); err == nil {
			ids = append(ids, oid)
		}
	}
	return o.FindAll(ctx, bson.M{"_id": bson.M{"$in": ids}}, sortFields)
}

func (o *_UserBlogMgr) FindByID(ctx context.Context, id string) (result *UserBlog, err error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, mongo.ErrNoDocuments
	}

	col := o.GetCol()
	ret := col.FindOne(ctx, bson.M{"_id": oid})
	if err = ret.Err(); err != nil {
		return nil, err
	}
	err = ret.Decode(&result)
	return
}

func (o *_UserBlogMgr) RemoveAll(ctx context.Context, query interface{}) (int64, error) {
	if query == nil {
		query = bson.M{}
	}

	col := o.GetCol()
	ret, err := col.DeleteMany(ctx, query)
	if err != nil {
		return 0, err
	}
	return ret.DeletedCount, nil
}

func (o *_UserBlogMgr) RemoveByID(ctx context.Context, id string) (err error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return mongo.ErrNoDocuments
	}

	col := o.GetCol()
	_, err = col.DeleteOne(ctx, bson.M{"_id": oid})
	return err
}

func (m *_UserBlogMgr) GetCol() *mongo.Collection {
	return Col("test_user_blog")
}
