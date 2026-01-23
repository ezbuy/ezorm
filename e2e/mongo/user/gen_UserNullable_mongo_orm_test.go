package user

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

var _ time.Time

type UserNullableSuite struct {
	suite.Suite
	ctx context.Context
}

func (s *UserNullableSuite) SetupTest() {
	s.ctx = context.TODO()
}

// TearDownTest runs after each test in the suite.
func (s *UserNullableSuite) TearDownTest() {
}

// TestSave tests the Save method with various scenarios
func (s *UserNullableSuite) TestSave() {
	cases := []struct {
		name    string
		data    UserNullable
		wantErr bool
	}{
		{
			name: "save_success",
			data: UserNullable{
				// ID will be auto-generated
				UserId:   1,
				Username: "test_username",
			},
			wantErr: false,
		},
	}
	mt := mtest.New(
		s.T(),
		mtest.NewOptions().ClientType(mtest.Mock),
	)
	defer mt.Close()
	for _, c := range cases {
		s.Run(c.name, func() {
			mt.Run(c.name, func(t *mtest.T) {
				if t.Client == nil {
					panic("t.Client is nil - mtest not properly initialized")
				}
				MgoSetup(nil, WithMockStub(t))
				updateResp := mtest.CreateSuccessResponse(
					bson.D{
						{Key: "acknowledged", Value: true},
						{Key: "matchedCount", Value: 1},
						{Key: "modifiedCount", Value: 1},
					}...,
				)
				t.AddMockResponses(updateResp)
				obj := c.data
				_, err := obj.Save(s.ctx)
				if c.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		})
	}
}

// TestFindOneAndSave tests the FindOneAndSave method
func (s *UserNullableSuite) TestFindOneAndSave() {
	cases := []struct {
		name    string
		data    UserNullable
		query   bson.M
		wantErr bool
	}{
		{
			name: "findOneAndSave_success",
			data: UserNullable{
				// ID will be auto-generated
				UserId:   2,
				Username: "test_username",
			},
			query: bson.M{
				"Username": "test_username",
			},
			wantErr: false,
		},
	}
	mt := mtest.New(
		s.T(),
		mtest.NewOptions().ClientType(mtest.Mock),
	)
	defer mt.Close()
	for _, c := range cases {
		s.Run(c.name, func() {
			mt.Run(c.name, func(t *mtest.T) {
				if t.Client == nil {
					panic("t.Client is nil - mtest not properly initialized")
				}
				MgoSetup(nil, WithMockStub(t))
				// FindOneAndUpdate returns the actual document
				objectID := primitive.NewObjectID()
				findOneAndUpdateResp := mtest.CreateSuccessResponse(bson.D{
					{Key: "value", Value: bson.D{
						{Key: "_id", Value: objectID},
						{Key: UserNullableMgoFieldUserId, Value: c.data.UserId},
						{Key: UserNullableMgoFieldUsername, Value: c.data.Username},
						{Key: UserNullableMgoFieldAge, Value: c.data.Age},
						{Key: UserNullableMgoFieldNickname, Value: c.data.Nickname},
						{Key: UserNullableMgoFieldRegisterDate, Value: c.data.RegisterDate},
					}},
					{Key: "ok", Value: 1},
				}...)
				t.AddMockResponses(findOneAndUpdateResp)
				obj := c.data
				_, err := obj.FindOneAndSave(s.ctx, c.query)
				if c.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		})
	}
}

// TestInsertUnique tests the InsertUnique method
func (s *UserNullableSuite) TestInsertUnique() {
	cases := []struct {
		name    string
		data    UserNullable
		query   bson.M
		wantErr bool
	}{
		{
			name: "insertUnique_success",
			data: UserNullable{
				// ID will be auto-generated
				UserId:   3,
				Username: "unique_username",
			},
			query: bson.M{
				"Username": "unique_username",
			},
			wantErr: false,
		},
	}
	mt := mtest.New(
		s.T(),
		mtest.NewOptions().ClientType(mtest.Mock),
	)
	defer mt.Close()
	for _, c := range cases {
		s.Run(c.name, func() {
			mt.Run(c.name, func(t *mtest.T) {
				if t.Client == nil {
					panic("t.Client is nil - mtest not properly initialized")
				}
				MgoSetup(nil, WithMockStub(t))

				updateResp := mtest.CreateSuccessResponse(
					bson.D{
						{Key: "ok", Value: 1},
						{Key: "n", Value: 1},
						{Key: "nModified", Value: 0},
						{Key: "upserted", Value: bson.A{
							bson.D{
								{Key: "index", Value: 0},
								{Key: "_id", Value: primitive.NewObjectID()},
							},
						}},
					}...,
				)
				t.AddMockResponses(updateResp)
				obj := c.data
				saved, err := obj.InsertUnique(s.ctx, c.query)
				if c.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					assert.True(t, saved)
				}
			})
		})
	}
}

// TestFindOne tests the FindOne method
func (s *UserNullableSuite) TestFindOne() {
	cases := []struct {
		name       string
		query      bson.M
		sortFields interface{}
		wantErr    bool
	}{
		{
			name: "findOne_success",
			query: bson.M{
				"Username": "test_username",
			},
			sortFields: UserNullableMgoSortField_WRP{UserNullableMgoSortFieldIDAsc},
			wantErr:    false,
		},
	}
	mt := mtest.New(
		s.T(),
		mtest.NewOptions().ClientType(mtest.Mock),
	)
	defer mt.Close()
	for _, c := range cases {
		s.Run(c.name, func() {
			mt.Run(c.name, func(t *mtest.T) {
				if t.Client == nil {
					panic("t.Client is nil - mtest not properly initialized")
				}
				MgoSetup(nil, WithMockStub(t))
				findResp := mtest.CreateCursorResponse(1, "test.usernullable", mtest.FirstBatch, bson.D{
					{Key: "_id", Value: primitive.NewObjectID()},
					{Key: UserNullableMgoFieldUserId, Value: uint64(1)},
					{Key: UserNullableMgoFieldUsername, Value: "test_username"},
					{Key: UserNullableMgoFieldAge, Value: int32(1)},
					{Key: UserNullableMgoFieldNickname, Value: "test_nickname"},
					{Key: UserNullableMgoFieldRegisterDate, Value: time.Now()},
				})
				t.AddMockResponses(findResp)
				result, err := Get_UserNullableMgr().FindOne(s.ctx, c.query, c.sortFields)
				if c.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					assert.NotNil(t, result)
					assert.Equal(t, "test_username", result.Username)
				}
			})
		})
	}
}

// TestQuery tests the Query method
func (s *UserNullableSuite) TestQuery() {
	cases := []struct {
		name       string
		query      bson.M
		limit      int
		offset     int
		sortFields interface{}
		wantErr    bool
	}{
		{
			name: "query_success",
			query: bson.M{
				"UserId": uint64(1),
			},
			limit:      10,
			offset:     0,
			sortFields: UserNullableMgoSortField_WRP{UserNullableMgoSortFieldIDAsc},
			wantErr:    false,
		},
	}
	mt := mtest.New(
		s.T(),
		mtest.NewOptions().ClientType(mtest.Mock),
	)
	defer mt.Close()
	for _, c := range cases {
		s.Run(c.name, func() {
			mt.Run(c.name, func(t *mtest.T) {
				if t.Client == nil {
					panic("t.Client is nil - mtest not properly initialized")
				}
				MgoSetup(nil, WithMockStub(t))
				first := mtest.CreateCursorResponse(1, "test.usernullable", mtest.FirstBatch, bson.D{
					{Key: "_id", Value: primitive.NewObjectID()},
					{Key: UserNullableMgoFieldUserId, Value: uint64(1)},
					{Key: UserNullableMgoFieldUsername, Value: "test_username_1"},
					{Key: UserNullableMgoFieldAge, Value: int32(1)},
					{Key: UserNullableMgoFieldNickname, Value: "test_nickname_1"},
					{Key: UserNullableMgoFieldRegisterDate, Value: time.Now()},
				})
				getMore := mtest.CreateCursorResponse(1, "test.usernullable", mtest.NextBatch, bson.D{
					{Key: "_id", Value: primitive.NewObjectID()},
					{Key: UserNullableMgoFieldUserId, Value: uint64(2)},
					{Key: UserNullableMgoFieldUsername, Value: "test_username_2"},
					{Key: UserNullableMgoFieldAge, Value: int32(2)},
					{Key: UserNullableMgoFieldNickname, Value: "test_nickname_2"},
					{Key: UserNullableMgoFieldRegisterDate, Value: time.Now()},
				})
				killCursors := mtest.CreateCursorResponse(0, "test.usernullable", mtest.NextBatch)
				t.AddMockResponses(first, getMore, killCursors)
				cursor, err := Get_UserNullableMgr().Query(s.ctx, c.query, c.limit, c.offset, c.sortFields)
				if c.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					assert.NotNil(t, cursor)
					// Close cursor to clean up
					cursor.Close(s.ctx)
				}
			})
		})
	}
}

// TestFindByUsername tests the FindByUsername method
func (s *UserNullableSuite) TestFindByUsername() {
	cases := []struct {
		name           string
		field_username string
		limit          int
		offset         int
		sortFields     interface{}
		wantErr        bool
	}{
		{
			name:           "findbyusername_success",
			field_username: "test_username",
			limit:          10,
			offset:         0,
			sortFields:     nil,
			wantErr:        false,
		},
	}
	mt := mtest.New(
		s.T(),
		mtest.NewOptions().ClientType(mtest.Mock),
	)
	defer mt.Close()
	for _, c := range cases {
		s.Run(c.name, func() {
			mt.Run(c.name, func(t *mtest.T) {
				if t.Client == nil {
					panic("t.Client is nil - mtest not properly initialized")
				}
				MgoSetup(nil, WithMockStub(t))
				first := mtest.CreateCursorResponse(1, "test.usernullable", mtest.FirstBatch, bson.D{
					{Key: "_id", Value: primitive.NewObjectID()},
					{Key: UserNullableMgoFieldUserId, Value: uint64(1)},
					{Key: UserNullableMgoFieldUsername, Value: "test_username"},
					{Key: UserNullableMgoFieldAge, Value: int32(1)},
					{Key: UserNullableMgoFieldNickname, Value: "test_nickname"},
					{Key: UserNullableMgoFieldRegisterDate, Value: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)},
				})
				killCursors := mtest.CreateCursorResponse(0, "test.usernullable", mtest.NextBatch)
				t.AddMockResponses(first, killCursors)
				result, err := Get_UserNullableMgr().FindByUsername(s.ctx, c.field_username, c.limit, c.offset, c.sortFields)
				if c.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					assert.NotNil(t, result)
					assert.Len(t, result, 1)
					assert.Equal(t, c.field_username, result[0].Username)
				}
			})
		})
	}
}

func TestUserNullable(t *testing.T) {
	suite.Run(t, &UserNullableSuite{})
}
