package user_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/ezbuy/ezorm/v2/e2e/mongo/nested"
	"github.com/ezbuy/ezorm/v2/e2e/mongo/user"
	"github.com/ezbuy/ezorm/v2/pkg/db"
	"github.com/ezbuy/ezorm/v2/pkg/orm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getConfigFromEnv(dbName string) *db.MongoConfig {
	return &db.MongoConfig{
		DBName: dbName,
		MongoDB: fmt.Sprintf(
			"mongodb://%s:%s@%s:%s",
			os.Getenv("MONGO_USER"),
			os.Getenv("MONGO_PASSWORD"),
			os.Getenv("MONGO_HOST"),
			os.Getenv("MONGO_PORT"),
		),
	}
}

func TestMain(m *testing.M) {
	user.MgoSetup(getConfigFromEnv("ezorm"))
	expr := 3600
	exprInt32 := int32(expr)
	var exist, created int
	if err := orm.EnsureAllIndex(orm.WithIndexNameHandler(user.UserIndexKey_Age, &options.IndexOptions{
		ExpireAfterSeconds: &exprInt32,
	}), orm.IndexCreateResult(&exist, &created)); err != nil {
		panic(err)
	}
	if exist != 0 && created != len(user.UserIndexes) {
		panic(fmt.Errorf("failed to create index, created: %d, exist: %d", created, exist))
	}
	var exist_s, created_s int
	if err := orm.EnsureAllIndex(orm.IndexCreateResult(&exist_s, &created_s)); err != nil {
		panic(err)
	}
	if created_s != 0 && exist_s != len(user.UserIndexes) {
		panic(fmt.Errorf("failed to create index, created: %d, exist: %d", created_s, exist_s))
	}
	indexMap := make(map[string]struct{})
	indexes, err := user.Col("test_user").Indexes().ListSpecifications(context.Background())
	if err != nil {
		panic(err)
	}

	for _, index := range indexes {
		indexMap[index.Name] = struct{}{}
	}

	for _, index := range user.UserIndexes {
		k := orm.IndexKey(index.Keys)
		if _, ok := indexMap[k]; !ok {
			panic(fmt.Errorf("not all index found in mongo: key: %s", k))
		}
		if k == "Age_1" {
			if *index.Options.ExpireAfterSeconds != 3600 {
				panic(
					fmt.Errorf("index expire after seconds not match,got %d,expect %d",
						*index.Options.ExpireAfterSeconds,
						3600,
					))
			}
		}
	}

	os.Exit(m.Run())
}

func TestOperateMultipleDB(t *testing.T) {
	// fetch user from ezorm
	user.MgoSetup(getConfigFromEnv("ezorm"))
	// fetch user from ezorm_nested
	nested.MgoSetup(getConfigFromEnv("ezorm_nested"))

	u1 := user.Get_UserMgr().NewUser()
	u1.Username = "username_1"
	if _, err := u1.Save(context.TODO()); err != nil {
		t.Fatalf("failed to save user: %s", err)
	}

	c := nested.Get_UserMgr().Count(context.TODO(), bson.M{})
	require.Equalf(t, 0, c, "unexpected count of users, got: %d, expect: %d", c, 0)

	u2 := nested.Get_UserMgr().NewUser()
	u2.Username = "username_2"
	if _, err := u2.Save(context.TODO()); err != nil {
		t.Fatalf("failed to save user: %s", err)
	}

	c2 := nested.Get_UserMgr().Count(context.TODO(), bson.M{})
	require.Equalf(t, 1, c2, "unexpected count of users, got: %d, expect: %d", c2, 1)

	c3 := user.Get_UserMgr().Count(context.TODO(), bson.M{})
	require.Equalf(t, 1, c3, "unexpected count of users, got: %d, expect: %d", c3, 1)

	t.Cleanup(func() {
		if _, err := user.Get_UserMgr().RemoveAll(context.TODO(), nil); err != nil {
			t.Fatalf("failed to remove all users: %s", err)
		}
		if _, err := nested.Get_UserMgr().RemoveAll(context.TODO(), nil); err != nil {
			t.Fatalf("failed to remove all users: %s", err)
		}
	})
}

func TestE2ESave(t *testing.T) {
	user.MgoSetup(getConfigFromEnv("ezorm"))

	ctx := context.TODO()
	u1 := user.Get_UserMgr().NewUser()
	u1.Username = "username_1"

	t.Log("insert user start")

	if _, err := u1.Save(ctx); err != nil {
		t.Fatalf("failed to insert user_1 by Save: %s", err)
	}
	if fq, err := user.Get_UserMgr().FindByID(ctx, u1.Id()); err != nil {
		t.Fatalf("failed to find user_1: %s", err)
	} else {
		if fq.Username != u1.Username {
			t.Fatalf("unexpected user_1 name, got: %s, expect: %s",
				fq.Username, u1.Username)
		}
	}

	t.Log("update user start")

	u1.Username = "username_1_new"
	if _, err := u1.Save(ctx); err != nil {
		t.Fatalf("failed to update user_1 by Save: %s", err)
	}
	if fq, err := user.Get_UserMgr().FindByID(ctx, u1.Id()); err != nil {
		t.Fatalf("failed to find user_1: %s", err)
	} else {
		if fq.Username != u1.Username {
			t.Fatalf("unexpected user_1 name, got: %s, expect: %s",
				fq.Username, u1.Username)
		}
	}

	if _, err := user.Get_UserMgr().RemoveAll(ctx, nil); err != nil {
		t.Fatalf("failed to remove all users: %s", err)
	}
}

func TestFindAndSave(t *testing.T) {
	ctx := context.TODO()
	u1 := user.Get_UserMgr().NewUser()
	u1.Username = "username_1"
	_, err := u1.FindOneAndSave(ctx, bson.M{
		user.UserMgoFieldUsername: u1.Username,
	})
	assert.NoError(t, err)
	u1.Username = "username_1_new"
	res, err := u1.FindOneAndSave(ctx, bson.M{
		user.UserMgoFieldUsername: "username_1",
	})
	assert.NoError(t, err)
	oid, err := orm.GetIDFromSingleResult(res)
	assert.NoError(t, err)
	u, err := user.Get_UserMgr().FindOne(ctx, bson.M{
		user.UserMgoFieldUsername: "username_1_new",
	}, nil)
	assert.NoError(t, err)
	assert.Equal(t, u.ID.Hex(), oid)
}

func TestInsertUnique(t *testing.T) {
	const (
		uniqueUsername = "unique_user_1"
	)

	ctx := context.TODO()
	u1 := user.Get_UserMgr().NewUser()
	u1.Age++
	u1.Username = uniqueUsername
	u1.RegisterDate = time.Now()

	saved, err := u1.InsertUnique(ctx, bson.M{user.UserMgoFieldUsername: uniqueUsername})
	if err != nil {
		t.Fatalf("failed to insert unique user: %s", err)
	}
	if !saved {
		t.Fatal("unique_user_1 never initialized but got saved is false")
	}

	u1.Age++
	saved, err = u1.InsertUnique(ctx, bson.M{
		user.UserMgoFieldUsername: uniqueUsername,
	})
	if err != nil {
		t.Fatalf("failed to reinsert unique user: %s", err)
	}
	if saved {
		t.Fatal("unique_user_1 has been initialized but got saved is true")
	}

	if _, err := user.Get_UserMgr().RemoveAll(ctx, nil); err != nil {
		t.Fatalf("failed to remove all users: %s", err)
	}
}

func TestFindOne(t *testing.T) {
	const (
		uid1      = 1
		uid2      = 2
		globalAge = 30
	)

	ctx := context.TODO()
	{
		u1 := user.Get_UserMgr().NewUser()
		u1.UserId = uid1
		u1.Age = globalAge
		if _, err := u1.Save(ctx); err != nil {
			t.Fatalf("failed to save uid=1 user: %s", err)
		}
	}
	{
		u2 := user.Get_UserMgr().NewUser()
		u2.UserId = uid2
		u2.Age = globalAge
		if _, err := u2.Save(ctx); err != nil {
			t.Fatalf("failed to save uid=1 user: %s", err)
		}
	}

	gotUser, err := user.Get_UserMgr().FindOne(ctx, bson.M{
		user.UserMgoFieldAge: globalAge,
	}, user.UserMgoSortField_WRP{user.UserMgoSortFieldIDDesc})
	if err != nil {
		t.Fatalf("failed to find one user sort by _id: %s", err)
	}
	if uid := gotUser.UserId; uid != uid2 {
		t.Fatalf("unexpected user id, got: %d, expect: %d\n", uid, uid2)
	}

	if _, err := user.Get_UserMgr().RemoveAll(ctx, nil); err != nil {
		t.Fatalf("failed to remove all users: %s", err)
	}
}

type removeFn func(ctx context.Context) error

func initUsersHelper(t *testing.T, users ...*user.User) (removeFn removeFn) {
	ctx := context.Background()
	for _, u := range users {
		if _, e := u.Save(ctx); e != nil {
			t.Fatalf("failed to create user: %s", e)
		}
	}

	return func(ctx context.Context) error {
		if _, err := user.Get_UserMgr().RemoveAll(ctx, nil); err != nil {
			t.Fatalf("failed to remove all users: %s", err)
		}
		return nil
	}
}

func TestQuery(t *testing.T) {
	const (
		uid1      = 1
		uid2      = 2
		globalAge = 30
	)

	ctx := context.TODO()
	u1 := user.Get_UserMgr().NewUser()
	u1.UserId = uid1
	u1.Age = globalAge

	u2 := user.Get_UserMgr().NewUser()
	u2.UserId = uid2
	u2.Age = globalAge

	cleanFn := initUsersHelper(t, u1, u2)

	{
		cursor, err := user.Get_UserMgr().Query(ctx, bson.M{
			user.UserMgoFieldAge: globalAge,
		}, 1, 1, user.UserMgoSortField_WRP{user.UserMgoSortFieldIDDesc})
		if err != nil {
			t.Fatalf("failed to query users: %s", err)
		}

		var users []*user.User
		if err = cursor.All(ctx, &users); err != nil {
			t.Fatalf("failed to decode users: %s", err)
		}
		if err = cursor.Close(ctx); err != nil {
			t.Fatalf("failed to close cursor: %s", err)
		}

		if l := len(users); l != 1 {
			t.Fatalf("unexpected users, got length: %d, expect length: %d", l, 1)
		}
		if uid := users[0].UserId; uid != uid1 {
			t.Fatalf("unexpected users, got uid: %d, expect uid: %d", uid, uid1)
		}
	}
	{
		cursor, err := user.Get_UserMgr().Query(ctx, bson.M{
			user.UserMgoFieldAge: globalAge,
		}, 10, 0, user.UserMgoSortField_WRP{user.UserMgoSortFieldIDAsc})
		if err != nil {
			t.Fatalf("failed to query users: %s", err)
		}

		var users []*user.User
		if err = cursor.All(ctx, &users); err != nil {
			t.Fatalf("failed to decode users: %s", err)
		}
		if err = cursor.Close(ctx); err != nil {
			t.Fatalf("failed to close cursor: %s", err)
		}

		if l := len(users); l != 2 {
			t.Fatalf("unexpected users, got length: %d, expect length: %d", l, 2)
		}
		if uid := users[0].UserId; uid != uid1 {
			t.Fatalf("unexpected users[0], got uid: %d, expect uid: %d", uid, uid1)
		}
		if uid := users[1].UserId; uid != uid2 {
			t.Fatalf("unexpected users[0], got uid: %d, expect uid: %d", uid, uid2)
		}
	}

	if err := cleanFn(ctx); err != nil {
		t.Fatalf("failed to remove all users: %s", err)
	}
}

func TestFindByIndexes(t *testing.T) {
	const (
		uid1   = 1
		uname1 = "John"

		uid2   = 2
		uname2 = "Mike"

		globalAge = 30
	)

	ctx := context.TODO()
	u1 := user.Get_UserMgr().NewUser()
	u1.UserId = uid1
	u1.Username = uname1
	u1.Age = globalAge
	u1.RegisterDate = time.Now()

	u2 := user.Get_UserMgr().NewUser()
	u2.UserId = uid2
	u2.Username = uname2
	u2.Age = globalAge
	u2.RegisterDate = time.Now().Add(time.Hour)

	cleanFn := initUsersHelper(t, u1, u2)
	{
		users, err := user.Get_UserMgr().FindByUsernameAge(ctx, uname1, globalAge, 10, 0, nil)
		if err != nil {
			t.Fatalf("failed to find by username and age: %s", err)
		}
		if l := len(users); l != 1 {
			t.Fatalf("unexpected length of users, got %d, expect: %d", l, 1)
		} else if uid := users[0].UserId; uid != uid1 {
			t.Fatalf("unexpected uid of users, got: %d, expect: %d", uid, uid1)
		}
	}
	{
		users, err := user.Get_UserMgr().FindByUsernameAge(ctx, uname2, 0, 10, 0, nil)
		if err != nil {
			t.Fatalf("failed to find by username and age: %s", err)
		}
		if l := len(users); l != 0 {
			t.Fatalf("unexpected length of users, got %d, expect: %d", l, 0)
		}
	}
	{
		users, err := user.Get_UserMgr().FindAll(ctx, bson.M{}, user.UserMgoSortField_WRP{user.UserMgoSortFieldRegisterDateDesc})
		if err != nil {
			t.Fatalf("failed to find by username and age: %s", err)
		}
		if l := len(users); l != 2 {
			t.Fatalf("unexpected length of users, got %d, expect: %d", l, 2)
		}
		if uid := users[0].UserId; uid != uid2 {
			t.Fatalf("unexpected uid of users, got: %d, expect: %d", uid, uid2)
		}
	}

	if err := cleanFn(ctx); err != nil {
		t.Fatalf("failed to remove all users: %s", err)
	}
}

func TestCount(t *testing.T) {
	const (
		uid1   = 1
		uname1 = "John"

		uid2   = 2
		uname2 = "Mike"

		globalAge = 30
	)

	ctx := context.TODO()
	u1 := user.Get_UserMgr().NewUser()
	u1.UserId = uid1
	u1.Username = uname1
	u1.Age = globalAge

	u2 := user.Get_UserMgr().NewUser()
	u2.UserId = uid2
	u2.Username = uname2
	u2.Age = globalAge

	cleanFn := initUsersHelper(t, u1, u2)

	{
		count, err := user.Get_UserMgr().CountE(ctx, bson.M{
			user.UserMgoFieldAge: globalAge,
		})
		if err != nil {
			t.Fatalf("failed to count users: %s", err)
		}
		if count != 2 {
			t.Fatalf("unexpected count of users, got: %d, expect: %d", count, 2)
		}
	}
	{
		count, err := user.Get_UserMgr().CountE(ctx, bson.M{
			user.UserMgoFieldUsername: "",
		})
		if err != nil {
			t.Fatalf("failed to count users: %s", err)
		}
		if count != 0 {
			t.Fatalf("unexpected count of users, got: %d, expect: %d", count, 0)
		}
	}

	if err := cleanFn(ctx); err != nil {
		t.Fatalf("failed to remove all users: %s", err)
	}
}

func TestErrNotFound(t *testing.T) {
	_, err := user.Get_UserMgr().FindOne(context.TODO(), bson.M{
		user.UserMgoFieldUsername: "user-not-found",
	}, user.UserMgoSortField_WRP{user.UserMgoSortFieldIDAsc})
	if err == nil {
		t.Fatalf("should not return empty error cause user not exist: %s", err)
	}
	if !orm.IsErrNotFound(err) {
		t.Fatalf("not found error not match, got: %s, expect: %s", err, mongo.ErrNoDocuments)
	}
}
