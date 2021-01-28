package user_test

import (
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
	u1 := user.Get_UserMgr().NewUser()
	u1.Username = "username_1"

	t.Log("insert user start")
	fmt.Println("insert user start")
	if _, err := u1.Save(); err != nil {
		t.Fatalf("failed to insert user_1 by Save: %s", err)
	}
	if fq, err := user.Get_UserMgr().FindByID(u1.Id()); err != nil {
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
	if _, err := u1.Save(); err != nil {
		t.Fatalf("failed to update user_1 by Save: %s", err)
	}
	if fq, err := user.Get_UserMgr().FindByID(u1.Id()); err != nil {
		t.Fatalf("failed to find user_1: %s", err)
	} else {
		if fq.Username != u1.Username {
			t.Fatalf("unexpected user_1 name, got: %s, expect: %s",
				fq.Username, u1.Username)
		}
	}

	if _, err := user.Get_UserMgr().RemoveAll(nil); err != nil {
		t.Fatalf("failed to remove all users: %s", err)
	}
}

func TestInsertUnique(t *testing.T) {
	const (
		unique_username = "unique_user_1"
	)

	u1 := user.Get_UserMgr().NewUser()
	u1.Age++
	u1.Username = unique_username
	u1.RegisterDate = time.Now()

	saved, err := u1.InsertUnique(bson.M{user.UserMgoFieldUsername: unique_username})
	if err != nil {
		t.Fatalf("failed to insert unique user: %s", err)
	}
	if !saved {
		t.Fatal("unique_user_1 never initialized but got saved is false")
	}

	u1.Age++
	saved, err = u1.InsertUnique(bson.M{
		user.UserMgoFieldUsername: unique_username,
	})
	if err != nil {
		t.Fatalf("failed to reinsert unique user: %s", err)
	}
	if saved {
		t.Fatal("unique_user_1 has been initialized but got saved is true")
	}

	if _, err := user.Get_UserMgr().RemoveAll(nil); err != nil {
		t.Fatalf("failed to remove all users: %s", err)
	}
}

func TestErrNotFound(t *testing.T) {
	_, err := user.Get_UserMgr().FindOne(bson.M{
		user.UserMgoFieldUsername: "user-not-found",
	}, user.UserMgoSortField_WRP{user.UserMgoSortFieldIDAsc})
	if err == nil {
		t.Fatalf("should not return empty error cause user not exist: %s", err)
	}
	if err != mongo.ErrNoDocuments {
		t.Fatalf("not found error not match, got: %s, expect: %s", err, mongo.ErrNoDocuments)
	}
}
