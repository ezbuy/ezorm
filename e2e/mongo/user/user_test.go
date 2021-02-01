package user_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/ezbuy/ezorm/db"
	"github.com/ezbuy/ezorm/e2e/mongo/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestMain(m *testing.M) {
	user.MgoSetup(&db.MongoConfig{
		MongoDB:   "mongodb://127.0.0.1:27017",
		DBName:    "test",
		PoolLimit: 2,
	})
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

func TestQuery(t *testing.T) {
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

		if l := len(users); 1 != l {
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

		if l := len(users); 2 != l {
			t.Fatalf("unexpected users, got length: %d, expect length: %d", l, 2)
		}
		if uid := users[0].UserId; uid != uid1 {
			t.Fatalf("unexpected users[0], got uid: %d, expect uid: %d", uid, uid1)
		}
		if uid := users[1].UserId; uid != uid2 {
			t.Fatalf("unexpected users[0], got uid: %d, expect uid: %d", uid, uid2)
		}
	}
	if _, err := user.Get_UserMgr().RemoveAll(ctx, nil); err != nil {
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
