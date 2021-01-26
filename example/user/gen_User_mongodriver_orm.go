package user

import (
	"context"

	"github.com/ezbuy/ezorm/orm"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	orm.RegisterEzOrmObjByID("user", "User", newUserFindByID)
	orm.RegisterEzOrmObjRemove("user", "User", UserMgr.RemoveByID)
}

func newUserFindByID(id string) (result orm.EzOrmObj, err error) {
	return UserMgr.FindByID(id)
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

func (o *User) Save() (*mongo.UpdateResult, error) {
	isNew := o.isNew

	filter := bson.M{"_id": o.ID}
	update := bson.M{
		"$set": bson.M{
			UserMgoFieldID:           o.ID,
			UserMgoFieldUserId:       o.UserId,
			UserMgoFieldUsername:     o.Username,
			UserMgoFieldAge:          o.Age,
			UserMgoFieldRegisterDate: o.RegisterDate,
		},
	}

	opts := options.Update().SetUpsert(true)
	col := UserMgr.GetCol()
	ret, err := col.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		return ret, err
	}

	o.isNew = false
	if isNew {
		UserInsertCallback(o)
	} else {
		UserUpdateCallback(o)
	}
	return ret, err
}

func (o *User) InsertUnique(query interface{}) (saved bool, err error) {
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
	ret, err := col.UpdateOne(context.TODO(), query, update, opts)
	if err != nil {
		return false, err
	}
	if ret.UpsertedCount == 0 {
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

func (o *_UserMgr) FindOne(query interface{}, sortFields interface{}) (result *User, err error) {
	col := o.GetCol()
	opts := options.FindOne()

	if sortFields != nil {
		opts.SetSort(sortFields)
	}

	ret := col.FindOne(context.TODO(), query, opts)
	if err = ret.Err(); err != nil {
		return nil, err
	}
	err = ret.Decode(&result)
	return
}

func (o *_UserMgr) Query(query interface{}, limit, offset int, sortFields interface{}) (*mongo.Cursor, error) {
	col := o.GetCol()
	opts := options.Find()

	if limit > 0 {
		opts.SetLimit(int64(limit))
	}
	if offset > 0 {
		opts.SetLimit(int64(offset))
	}
	if sortFields != nil {
		opts.SetSort(sortFields)
	}

	return col.Find(context.TODO(), query, opts)
}

func (o *_UserMgr) FindByUsernameAge(Username string, Age int32, limit int, offset int, sortFields interface{}) (result []*User, err error) {
	query := bson.M{
		"Username": Username,
		"Age":      Age,
	}
	cursor, err := UserMgr.Query(query, limit, offset, sortFields)
	if err != nil {
		return nil, err
	}
	err = cursor.All(context.TODO(), &result)
	return
}

func (o *_UserMgr) FindByUsername(Username string, limit int, offset int, sortFields interface{}) (result []*User, err error) {
	query := bson.M{
		"Username": Username,
	}
	cursor, err := UserMgr.Query(query, limit, offset, sortFields)
	if err != nil {
		return nil, err
	}
	err = cursor.All(context.TODO(), &result)
	return
}

func (o *_UserMgr) FindByAge(Age int32, limit int, offset int, sortFields interface{}) (result []*User, err error) {
	query := bson.M{
		"Age": Age,
	}
	cursor, err := UserMgr.Query(query, limit, offset, sortFields)
	if err != nil {
		return nil, err
	}
	err = cursor.All(context.TODO(), &result)
	return
}

func (o *_UserMgr) Find(query interface{}, limit int, offset int, sortFields interface{}) (result []*User, err error) {
	cursor, err := UserMgr.Query(query, limit, offset, sortFields)
	if err != nil {
		return nil, err
	}
	err = cursor.All(context.TODO(), &result)
	return
}

func (o *_UserMgr) FindAll(query interface{}, sortFields interface{}) (result []*User, err error) {
	cursor, err := UserMgr.Query(query, -1, -1, sortFields)
	if err != nil {
		return nil, err
	}
	err = cursor.All(context.TODO(), &result)
	return
}

func (o *_UserMgr) Has(query interface{}) bool {
	count, err := o.CountE(query)
	if err != nil || count == 0 {
		return false
	}
	return true
}

func (o *_UserMgr) Count(query interface{}) int {
	count, _ := o.CountE(query)
	return count
}

func (o *_UserMgr) CountE(query interface{}) (int, error) {
	col := o.GetCol()
	count, err := col.CountDocuments(context.TODO(), query)
	return int(count), err
}

func (o *_UserMgr) FindByIDs(id []string, sortFields interface{}) (result []*User, err error) {
	ids := make([]primitive.ObjectID, 0, len(id))
	for _, i := range id {
		if oid, err := primitive.ObjectIDFromHex(i); err == nil {
			ids = append(ids, oid)
		}
	}
	return o.FindAll(bson.M{"_id": bson.M{"$in": ids}}, sortFields)
}

func (o *_UserMgr) FindByID(id string) (result *User, err error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, mongo.ErrNoDocuments
	}

	col := o.GetCol()
	ret := col.FindOne(context.TODO(), bson.M{"_id": oid})
	if err = ret.Err(); err != nil {
		return nil, err
	}
	err = ret.Decode(&result)
	return
}

func (o *_UserMgr) RemoveAll(query interface{}) (int64, error) {
	if query == nil {
		query = bson.M{}
	}

	col := o.GetCol()
	ret, err := col.DeleteMany(context.TODO(), query)
	if err != nil {
		return 0, err
	}
	return ret.DeletedCount, nil
}

func (o *_UserMgr) RemoveByID(id string) (err error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return mongo.ErrNoDocuments
	}

	col := o.GetCol()
	_, err = col.DeleteOne(context.TODO(), bson.M{"_id": oid})
	return err
}

func (m *_UserMgr) GetCol() *mongo.Collection {
	return Col("test_user")
}
