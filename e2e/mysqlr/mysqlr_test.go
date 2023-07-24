package mysqlr

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	port, err := strconv.ParseInt(os.Getenv("MYSQL_PORT"), 10, 64)
	if err != nil {
		panic(fmt.Errorf("failed to parse MYSQL_PORT: %s", err))
	}
	MySQLSetup(&MySQLConfig{
		Host:     os.Getenv("MYSQL_HOST"),
		Port:     int(port),
		UserName: "ezbuy",
		Password: "ezbuyisthebest",
		Database: "test",
	})

	sqls := []string{"gen.script.mysqlr.blog.sql", "gen.script.mysqlr.auto.blog.sql"}

	for _, s := range sqls {
		sql, err := os.ReadFile(s)
		if err != nil {
			panic(fmt.Errorf("failed to read gen.script.mysqlr.blog.sql: %s", err))
		}
		if _, err := MySQL().Exec(ctx, "CREATE DATABASE IF NOT EXISTS test"); err != nil {
			panic(fmt.Errorf("failed to create database: %s", err))
		}
		for _, q := range strings.Split(string(sql), ";") {
			if len(strings.TrimSpace(q)) == 0 {
				continue
			}
			if _, err := MySQL().Exec(ctx, q); err != nil {
				panic(fmt.Errorf("failed to create table: %s", err))
			}
		}
	}
	os.Exit(m.Run())
}

func TestBlogsCRUD(t *testing.T) {
	ctx := context.Background()
	db := MySQL()
	defer t.Cleanup(func() {
		MySQL().Exec(ctx, "TRUNCATE TABLE blogs")
		MySQL().Exec(ctx, "TRUNCATE TABLE auto_blogs")
	})
	t.Run("Create", func(t *testing.T) {
		bl := &Blog{
			Id:        1,
			UserId:    1,
			Title:     "test",
			Content:   "test",
			Status:    1,
			Readed:    0,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		af, err := BlogDBMgr(db).Create(ctx, bl)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), af)
		count, err := BlogDBMgr(db).SearchConditionsCount(ctx, []string{"user_id = ?"}, 1)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), count)
	})

	t.Run("GetByPrimaryKey", func(t *testing.T) {
		b, err := BlogDBMgr(db).FetchByPrimaryKey(ctx, 1, 1)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), b.Id)
	})

	t.Run("GetByUniqueKey", func(t *testing.T) {
		bs, err := BlogDBMgr(db).FetchByTitles(ctx, []string{"test"})
		assert.NoError(t, err)
		assert.Equal(t, 1, len(bs))
	})

	t.Run("GetUniqueMultiKey", func(t *testing.T) {
		b, err := BlogDBMgr(db).FetchByUserIdTitle(ctx, &UserIdTitleOfBlogUK{
			UserId: 1,
			Title:  "test",
		})
		assert.NoError(t, err)
		assert.Equal(t, int64(1), b.Id)
	})

	t.Run("GetUniqueMultiKeyByBatch", func(t *testing.T) {
		bs, err := BlogDBMgr(db).FetchByUserIdTitles(ctx, []*UserIdTitleOfBlogUK{
			{UserId: 1, Title: "test"},
		})
		assert.NoError(t, err)
		assert.Equal(t, 1, len(bs))
	})

	t.Run("GetByUniqueKeyNoRows", func(t *testing.T) {
		_, err := BlogDBMgr(db).FetchByTitle(ctx, "tet")
		assert.ErrorIs(t, err, sql.ErrNoRows)
	})

	t.Run("GetByIndex", func(t *testing.T) {
		bs, err := BlogDBMgr(db).FindAllByStatus(ctx, 1)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(bs))
	})

	t.Run("Update", func(t *testing.T) {
		b, err := BlogDBMgr(db).FetchByPrimaryKey(ctx, 1, 1)
		assert.NoError(t, err)
		b.Status = 2
		af, err := BlogDBMgr(db).Update(ctx, b)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), af)
	})

	t.Run("Delete", func(t *testing.T) {
		af, err := BlogDBMgr(db).DeleteByPrimaryKey(ctx, 1, 1)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), af)
		count, err := BlogDBMgr(db).SearchConditionsCount(ctx, []string{"user_id = ?"}, 1)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), count)
	})

	t.Run("BatchCreate", func(t *testing.T) {
		blogs := []*AutoBlog{
			{
				UserId:    2,
				Title:     "test2",
				Content:   "test2",
				Status:    1,
				Readed:    0,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				UserId:    3,
				Title:     "test3",
				Content:   "test3",
				Status:    1,
				Readed:    0,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}
		af, err := AutoBlogDBMgr(db).BatchCreate(ctx, blogs)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), af)
		assert.Equal(t, int64(1), blogs[0].Id)
		assert.Equal(t, int64(2), blogs[1].Id)
	})

	t.Run("GetByLimitOffset_NO_CONDITION", func(t *testing.T) {
		resp, err := GetRawQuery().Blog(ctx, &BlogReq{
			Limit:  1,
			Offset: 0,
		}, WithDB(db.DB))
		assert.NoError(t, err)
		assert.Equal(t, 1, len(resp))
	})

	t.Run("GetByLimitOffset_WITH_CONDITION", func(t *testing.T) {
		resp, err := GetRawQuery().Blog(ctx, &BlogReq{
			// not exist id
			Id:     9999,
			Limit:  1,
			Offset: 0,
		}, WithDB(db.DB))
		assert.NoError(t, err)
		assert.Equal(t, 0, len(resp))
	})
}

func TestBlogsTx(t *testing.T) {
	t.Run("SucceedTranscation", func(t *testing.T) {
		ctx := context.Background()
		tx, err := MySQL().BeginTx()
		assert.NoError(t, err)
		defer func() {
			err := tx.Close()
			assert.NoError(t, err)
		}()
		af, err := BlogDBMgr(tx).Create(ctx, &Blog{
			Id:        1,
			UserId:    1,
			Title:     "test",
			Content:   "test",
			Status:    1,
			Readed:    0,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
		assert.NoError(t, err)
		assert.Equal(t, int64(1), af)
		// not the same tx
		count, err := BlogDBMgr(MySQL()).SearchConditionsCount(ctx, []string{"user_id = ?"}, 1)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), count)
		// the same tx
		count, err = BlogDBMgr(tx).SearchConditionsCount(ctx, []string{"user_id = ?"}, 1)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), count)

		t.Cleanup(func() {
			MySQL().Exec(ctx, "TRUNCATE TABLE blogs")
		})
	})

	t.Run("BadTranscation", func(t *testing.T) {
		ctx := context.Background()
		tx, err := MySQL().BeginTx()
		assert.NoError(t, err)
		defer func() {
		}()
		af, err := BlogDBMgr(tx).Create(ctx, &Blog{
			Id:        1,
			UserId:    1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
		assert.NoError(t, err)
		assert.Equal(t, int64(1), af)
		// produce the duplicate key error
		af, err = BlogDBMgr(tx).Create(ctx, &Blog{
			Id:     1,
			UserId: 1,
		})
		assert.Error(t, err)
		if err != nil {
			tx.SetError(err)
		}
		assert.Equal(t, int64(0), af)
		// rollbacked
		assert.Equal(t, tx.IsRollback(), true)
		err = tx.Close()
		assert.NoError(t, err)

		count, err := BlogDBMgr(MySQL()).SearchConditionsCount(ctx, []string{"user_id = ?"}, 1)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), count)
		t.Cleanup(func() {
			MySQL().Exec(ctx, "TRUNCATE TABLE blogs")
		})
	})

	t.Run("PanicTransaction", func(t *testing.T) {
		ctx := context.Background()
		id, uid := int64(0o33), int32(0o33)

		defer func() {
			if err := recover(); err == nil {
				t.Fatal("panic error should be recover")
			}

			blogDBMgr := BlogDBMgr(MySQL())
			_, err := blogDBMgr.FetchByPrimaryKey(ctx, id, uid)
			if err == nil || !blogDBMgr.IsErrNotFound(err) {
				t.Fatalf("unexpected error during fetch blog: %s", err)
			}
		}()

		tx, err := MySQL().BeginTx()
		assert.NoError(t, err)
		defer tx.Close()

		_, err = BlogDBMgr(tx).Create(ctx, &Blog{
			Id:        id,
			UserId:    uid,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
		assert.NoError(t, err)

		arr := make([]any, 0)
		_ = arr[100]
	})
}
