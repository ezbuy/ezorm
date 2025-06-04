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

const ColUser = "test_user"

var UserIndexes = []mongo.IndexModel{
	{
		Keys: UserIndexKey_UsernameAge,
	},
	{
		Keys: UserIndexKey_Username,
	},
	{
		Keys: UserIndexKey_Age,
	},
	{
		Keys: UserIndexKey_RegisterDate,
	},
}

var UserIndexesFunc = func() {
	orm.SetupIndexModel(Col(ColUser), UserIndexes)
}
var UserIndexKey_UsernameAge = bson.D{
	{Key: "Username", Value: 1},
	{Key: "Age", Value: 1},
}
var UserIndexKey_Username = bson.D{
	{Key: "Username", Value: 1},
}
var UserIndexKey_Age = bson.D{
	{Key: "Age", Value: 1},
}
var UserIndexKey_RegisterDate = bson.D{
	{Key: "RegisterDate", Value: 1},
}

func init() {
	orm.RegisterEzOrmObjByID("mongo_e2e", "User", newUserFindByID)
	orm.RegisterEzOrmObjRemove("mongo_e2e", "User", newUserRemoveByID)
}

func newUserFindByID(id string) (result orm.EzOrmObj, err error) {
	return UserMgr.FindByID(context.TODO(), id)
}

func newUserRemoveByID(id string) error {
	return UserMgr.RemoveByID(context.TODO(), id)
}

// =====================================
// INSERT METHODS
// =====================================

var (
	insertUserCBs []func(obj orm.EzOrmObj)
	updateUserCBs []func(obj orm.EzOrmObj)
)

func UserAddInsertCallback(cb func(obj orm.EzOrmObj)) {
	insertUserCBs = append(insertUserCBs, cb)
}

func UserAddUpdateCallback(cb func(obj orm.EzOrmObj)) {
	updateUserCBs = append(updateUserCBs, cb)
}

func (o *User) Id() string {
	return o.ID.Hex()
}

// FindOneAndSave try to find one doc by `query`,  and then upsert the result with the current object
func (o *User) FindOneAndSave(ctx context.Context, query interface{}) (*mongo.SingleResult, error) {
	col := UserMgr.GetCol()
	opts := options.FindOneAndUpdate().SetUpsert(true)
	opts.SetReturnDocument(options.After)
	update := bson.M{
		"$set": bson.M{
			UserMgoFieldUserId:       o.UserId,
			UserMgoFieldUsername:     o.Username,
			UserMgoFieldAge:          o.Age,
			UserMgoFieldRegisterDate: o.RegisterDate,
		},
	}
	ret := col.FindOneAndUpdate(ctx, query, update, opts)
	if ret.Err() != nil {
		return nil, ret.Err()
	}
	return ret, nil
}

// Save upserts the document , Save itself is concurrent-safe , but maybe it is not atomic together with other operations, such as `Find`
func (o *User) Save(ctx context.Context) (*mongo.UpdateResult, error) {
	isNew := o.isNew
	if o.ID == primitive.NilObjectID {
		o.ID = primitive.NewObjectID()
	}

	filter := bson.M{"_id": o.ID}
	update := bson.M{
		"$set": bson.M{
			UserMgoFieldUserId:       o.UserId,
			UserMgoFieldUsername:     o.Username,
			UserMgoFieldAge:          o.Age,
			UserMgoFieldRegisterDate: o.RegisterDate,
		},
	}

	opts := options.Update().SetUpsert(true)
	col := UserMgr.GetCol()
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
		UserInsertCallback(o)
	} else {
		UserUpdateCallback(o)
	}
	return ret, err
}

func (o *User) InsertUnique(ctx context.Context, query interface{}) (saved bool, err error) {
	update := bson.M{
		"$setOnInsert": bson.M{
			UserMgoFieldID:           o.ID,
			UserMgoFieldUserId:       o.UserId,
			UserMgoFieldUsername:     o.Username,
			UserMgoFieldAge:          o.Age,
			UserMgoFieldRegisterDate: o.RegisterDate,
		},
	}

	opts := options.Update().SetUpsert(true)
	col := UserMgr.GetCol()
	ret, err := col.UpdateOne(ctx, query, update, opts)
	if err != nil {
		return false, err
	}
	if ret.UpsertedCount != 0 {
		saved = true
	}

	o.isNew = false
	if saved {
		UserInsertCallback(o)
	}
	return saved, nil
}

func UserInsertCallback(o *User) {
	for _, cb := range insertUserCBs {
		cb(o)
	}
}

func UserUpdateCallback(o *User) {
	for _, cb := range updateUserCBs {
		cb(o)
	}
}

// =====================================
// FOREIGN KEYS
// =====================================

// =====================================
// COLLECTION
// =====================================

func (o *_UserMgr) FindOne(ctx context.Context, query interface{}, sortFields interface{}) (result *User, err error) {
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

func (o *_UserMgr) Query(ctx context.Context, query interface{}, limit, offset int, sortFields interface{}) (*mongo.Cursor, error) {
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

func (o *_UserMgr) FindByUsernameAge(ctx context.Context, Username string, Age int32, limit int, offset int, sortFields interface{}) (result []*User, err error) {
	query := bson.M{
		"Username": Username,
		"Age":      Age,
	}
	cursor, err := o.Query(ctx, query, limit, offset, sortFields)
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &result)
	return
}

func (o *_UserMgr) FindByUsername(ctx context.Context, Username string, limit int, offset int, sortFields interface{}) (result []*User, err error) {
	query := bson.M{
		"Username": Username,
	}
	cursor, err := o.Query(ctx, query, limit, offset, sortFields)
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &result)
	return
}

func (o *_UserMgr) FindByAge(ctx context.Context, Age int32, limit int, offset int, sortFields interface{}) (result []*User, err error) {
	query := bson.M{
		"Age": Age,
	}
	cursor, err := o.Query(ctx, query, limit, offset, sortFields)
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &result)
	return
}

func (o *_UserMgr) FindByRegisterDate(ctx context.Context, RegisterDate time.Time, limit int, offset int, sortFields interface{}) (result []*User, err error) {
	query := bson.M{
		"RegisterDate": RegisterDate,
	}
	cursor, err := o.Query(ctx, query, limit, offset, sortFields)
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &result)
	return
}

func (o *_UserMgr) Find(ctx context.Context, query interface{}, limit int, offset int, sortFields interface{}) (result []*User, err error) {
	cursor, err := o.Query(ctx, query, limit, offset, sortFields)
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &result)
	return
}

func (o *_UserMgr) FindAll(ctx context.Context, query interface{}, sortFields interface{}) (result []*User, err error) {
	cursor, err := o.Query(ctx, query, -1, -1, sortFields)
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &result)
	return
}

func (o *_UserMgr) Has(ctx context.Context, query interface{}) bool {
	count, err := o.CountE(ctx, query)
	if err != nil || count == 0 {
		return false
	}
	return true
}

func (o *_UserMgr) Count(ctx context.Context, query interface{}) int {
	count, _ := o.CountE(ctx, query)
	return count
}

func (o *_UserMgr) CountE(ctx context.Context, query interface{}) (int, error) {
	col := o.GetCol()
	count, err := col.CountDocuments(ctx, query)
	return int(count), err
}

func (o *_UserMgr) FindByIDs(ctx context.Context, id []string, sortFields interface{}) (result []*User, err error) {
	ids := make([]primitive.ObjectID, 0, len(id))
	for _, i := range id {
		if oid, err := primitive.ObjectIDFromHex(i); err == nil {
			ids = append(ids, oid)
		}
	}
	return o.FindAll(ctx, bson.M{"_id": bson.M{"$in": ids}}, sortFields)
}

func (o *_UserMgr) FindByID(ctx context.Context, id string) (result *User, err error) {
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

func (o *_UserMgr) RemoveAll(ctx context.Context, query interface{}) (int64, error) {
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

func (o *_UserMgr) RemoveByID(ctx context.Context, id string) (err error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return mongo.ErrNoDocuments
	}

	col := o.GetCol()
	_, err = col.DeleteOne(ctx, bson.M{"_id": oid})
	return err
}

func (m *_UserMgr) GetCol() *mongo.Collection {
	return Col("test_user")
}
