package user_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/ezbuy/ezorm/v2/db"
	"github.com/ezbuy/ezorm/v2/e2e/mongo/user"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func getConfigFromEnv() *db.MongoConfig {
	return &db.MongoConfig{
		DBName: "ezorm",
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
	user.MgoSetup(getConfigFromEnv())
	os.Exit(m.Run())
}

func TestSave(t *testing.T) {
	ctx := context.TODO()
	u1 := user.Get_UserMgr().NewUser()
	u1.Username = "username_1"

	t.Log("insert user start")
	fmt.Println("insert user start")
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
	fmt.Println("update user start")
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

	u2 := user.Get_UserMgr().NewUser()
	u2.UserId = uid2
	u2.Username = uname2
	u2.Age = globalAge

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
	if err != mongo.ErrNoDocuments {
		t.Fatalf("not found error not match, got: %s, expect: %s", err, mongo.ErrNoDocuments)
	}
}