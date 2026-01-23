package user

import (
	"context"
	"encoding/json"
	"time"

	"github.com/ezbuy/ezorm/v2/pkg/orm"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// To import `time` package globally to satisfy `time.Time` index in yaml definition
var _ time.Time

// To import `json` package globally to satisfy `json.Marshaler` interface in yaml definition
var _ json.Marshaler

const ColUserNullable = "test_user_nullable"

var UserNullableIndexes = []mongo.IndexModel{
	{
		Keys: UserNullableIndexKey_Username,
	},
}

var UserNullableIndexesFunc = func() {
	orm.SetupIndexModel(Col(ColUserNullable), UserNullableIndexes)
}
var UserNullableIndexKey_Username = bson.D{
	{Key: "Username", Value: 1},
}

func init() {
	orm.RegisterEzOrmObjByID("mongo_e2e", "UserNullable", newUserNullableFindByID)
	orm.RegisterEzOrmObjRemove("mongo_e2e", "UserNullable", newUserNullableRemoveByID)
}

func newUserNullableFindByID(id string) (result orm.EzOrmObj, err error) {
	return UserNullableMgr.FindByID(context.TODO(), id)
}

func newUserNullableRemoveByID(id string) error {
	return UserNullableMgr.RemoveByID(context.TODO(), id)
}

// =====================================
// INSERT METHODS
// =====================================

var (
	insertUserNullableCBs []func(obj orm.EzOrmObj)
	updateUserNullableCBs []func(obj orm.EzOrmObj)
)

func UserNullableAddInsertCallback(cb func(obj orm.EzOrmObj)) {
	insertUserNullableCBs = append(insertUserNullableCBs, cb)
}

func UserNullableAddUpdateCallback(cb func(obj orm.EzOrmObj)) {
	updateUserNullableCBs = append(updateUserNullableCBs, cb)
}

func (o *UserNullable) Id() string {
	return o.ID.Hex()
}

// FindOneAndSave try to find one doc by `query`,  and then upsert the result with the current object
func (o *UserNullable) FindOneAndSave(ctx context.Context, query interface{}) (*mongo.SingleResult, error) {
	col := UserNullableMgr.GetCol()
	opts := options.FindOneAndUpdate().SetUpsert(true)
	opts.SetReturnDocument(options.After)
	setFields := bson.M{}
	setFields[UserNullableMgoFieldUserId] = o.UserId
	setFields[UserNullableMgoFieldUsername] = o.Username
	if o.Age != nil {
		setFields[UserNullableMgoFieldAge] = o.Age
	}
	if o.Nickname != nil {
		setFields[UserNullableMgoFieldNickname] = o.Nickname
	}
	if o.RegisterDate != nil {
		setFields[UserNullableMgoFieldRegisterDate] = o.RegisterDate
	}
	update := bson.M{
		"$set": setFields,
	}
	ret := col.FindOneAndUpdate(ctx, query, update, opts)
	if ret.Err() != nil {
		return nil, ret.Err()
	}
	return ret, nil
}

// Save upserts the document , Save itself is concurrent-safe , but maybe it is not atomic together with other operations, such as `Find`
func (o *UserNullable) Save(ctx context.Context) (*mongo.UpdateResult, error) {
	isNew := o.isNew
	if o.ID == primitive.NilObjectID {
		o.ID = primitive.NewObjectID()
	}

	filter := bson.M{"_id": o.ID}
	setFields := bson.M{}
	setFields[UserNullableMgoFieldUserId] = o.UserId
	setFields[UserNullableMgoFieldUsername] = o.Username
	if o.Age != nil {
		setFields[UserNullableMgoFieldAge] = o.Age
	}
	if o.Nickname != nil {
		setFields[UserNullableMgoFieldNickname] = o.Nickname
	}
	if o.RegisterDate != nil {
		setFields[UserNullableMgoFieldRegisterDate] = o.RegisterDate
	}
	update := bson.M{
		"$set": setFields,
	}

	opts := options.Update().SetUpsert(true)
	col := UserNullableMgr.GetCol()
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
		UserNullableInsertCallback(o)
	} else {
		UserNullableUpdateCallback(o)
	}
	return ret, err
}

func (o *UserNullable) InsertUnique(ctx context.Context, query interface{}) (saved bool, err error) {
	setFields := bson.M{}
	setFields[UserNullableMgoFieldID] = o.ID
	setFields[UserNullableMgoFieldUserId] = o.UserId
	setFields[UserNullableMgoFieldUsername] = o.Username
	if o.Age != nil {
		setFields[UserNullableMgoFieldAge] = o.Age
	}
	if o.Nickname != nil {
		setFields[UserNullableMgoFieldNickname] = o.Nickname
	}
	if o.RegisterDate != nil {
		setFields[UserNullableMgoFieldRegisterDate] = o.RegisterDate
	}
	update := bson.M{
		"$setOnInsert": setFields,
	}

	opts := options.Update().SetUpsert(true)
	col := UserNullableMgr.GetCol()
	ret, err := col.UpdateOne(ctx, query, update, opts)
	if err != nil {
		return false, err
	}
	if ret.UpsertedCount != 0 {
		saved = true
	}

	o.isNew = false
	if saved {
		UserNullableInsertCallback(o)
	}
	return saved, nil
}

func UserNullableInsertCallback(o *UserNullable) {
	for _, cb := range insertUserNullableCBs {
		cb(o)
	}
}

func UserNullableUpdateCallback(o *UserNullable) {
	for _, cb := range updateUserNullableCBs {
		cb(o)
	}
}

// =====================================
// FOREIGN KEYS
// =====================================

// =====================================
// COLLECTION
// =====================================

func (o *_UserNullableMgr) FindOne(ctx context.Context, query interface{}, sortFields interface{}) (result *UserNullable, err error) {
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

func (o *_UserNullableMgr) Query(ctx context.Context, query interface{}, limit, offset int, sortFields interface{}) (*mongo.Cursor, error) {
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

func (o *_UserNullableMgr) FindByUsername(ctx context.Context, Username string, limit int, offset int, sortFields interface{}) (result []*UserNullable, err error) {
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

func (o *_UserNullableMgr) Find(ctx context.Context, query interface{}, limit int, offset int, sortFields interface{}) (result []*UserNullable, err error) {
	cursor, err := o.Query(ctx, query, limit, offset, sortFields)
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &result)
	return
}

func (o *_UserNullableMgr) FindAll(ctx context.Context, query interface{}, sortFields interface{}) (result []*UserNullable, err error) {
	cursor, err := o.Query(ctx, query, -1, -1, sortFields)
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &result)
	return
}

func (o *_UserNullableMgr) Has(ctx context.Context, query interface{}) bool {
	count, err := o.CountE(ctx, query)
	if err != nil || count == 0 {
		return false
	}
	return true
}

func (o *_UserNullableMgr) Count(ctx context.Context, query interface{}) int {
	count, _ := o.CountE(ctx, query)
	return count
}

func (o *_UserNullableMgr) CountE(ctx context.Context, query interface{}) (int, error) {
	col := o.GetCol()
	count, err := col.CountDocuments(ctx, query)
	return int(count), err
}

func (o *_UserNullableMgr) FindByIDs(ctx context.Context, id []string, sortFields interface{}) (result []*UserNullable, err error) {
	ids := make([]primitive.ObjectID, 0, len(id))
	for _, i := range id {
		if oid, err := primitive.ObjectIDFromHex(i); err == nil {
			ids = append(ids, oid)
		}
	}
	return o.FindAll(ctx, bson.M{"_id": bson.M{"$in": ids}}, sortFields)
}

func (o *_UserNullableMgr) FindByID(ctx context.Context, id string) (result *UserNullable, err error) {
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

func (o *_UserNullableMgr) RemoveAll(ctx context.Context, query interface{}) (int64, error) {
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

func (o *_UserNullableMgr) RemoveByID(ctx context.Context, id string) (err error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return mongo.ErrNoDocuments
	}

	col := o.GetCol()
	_, err = col.DeleteOne(ctx, bson.M{"_id": oid})
	return err
}

func (m *_UserNullableMgr) GetCol() *mongo.Collection {
	return Col("test_user_nullable")
}
