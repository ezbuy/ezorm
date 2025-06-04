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

type UserSuite struct {
	suite.Suite
	ctx context.Context
}

func (s *UserSuite) SetupTest() {
	s.ctx = context.TODO()
}

// TearDownTest runs after each test in the suite.
func (s *UserSuite) TearDownTest() {
}

// TestSave tests the Save method with various scenarios
func (s *UserSuite) TestSave() {
	cases := []struct {
		name    string
		data    User
		wantErr bool
	}{
		{
			name: "save_success",
			data: User{
				// ID will be auto-generated
				UserId:       1,
				Username:     "test_username",
				Age:          1,
				RegisterDate: time.Now(),
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
func (s *UserSuite) TestFindOneAndSave() {
	cases := []struct {
		name    string
		data    User
		query   bson.M
		wantErr bool
	}{
		{
			name: "findOneAndSave_success",
			data: User{
				// ID will be auto-generated
				UserId:       2,
				Username:     "test_username",
				Age:          2,
				RegisterDate: time.Now(),
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
						{Key: UserMgoFieldUserId, Value: c.data.UserId},
						{Key: UserMgoFieldUsername, Value: c.data.Username},
						{Key: UserMgoFieldAge, Value: c.data.Age},
						{Key: UserMgoFieldRegisterDate, Value: c.data.RegisterDate},
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
func (s *UserSuite) TestInsertUnique() {
	cases := []struct {
		name    string
		data    User
		query   bson.M
		wantErr bool
	}{
		{
			name: "insertUnique_success",
			data: User{
				// ID will be auto-generated
				UserId:       3,
				Username:     "unique_username",
				Age:          3,
				RegisterDate: time.Now(),
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
func (s *UserSuite) TestFindOne() {
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
			sortFields: UserMgoSortField_WRP{UserMgoSortFieldIDAsc},
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
				findResp := mtest.CreateCursorResponse(1, "test.user", mtest.FirstBatch, bson.D{
					{Key: "_id", Value: primitive.NewObjectID()},
					{Key: UserMgoFieldUserId, Value: uint64(1)},
					{Key: UserMgoFieldUsername, Value: "test_username"},
					{Key: UserMgoFieldAge, Value: int32(1)},
					{Key: UserMgoFieldRegisterDate, Value: time.Now()},
				})
				t.AddMockResponses(findResp)
				result, err := Get_UserMgr().FindOne(s.ctx, c.query, c.sortFields)
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
func (s *UserSuite) TestQuery() {
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
			sortFields: UserMgoSortField_WRP{UserMgoSortFieldIDAsc},
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
				first := mtest.CreateCursorResponse(1, "test.user", mtest.FirstBatch, bson.D{
					{Key: "_id", Value: primitive.NewObjectID()},
					{Key: UserMgoFieldUserId, Value: uint64(1)},
					{Key: UserMgoFieldUsername, Value: "test_username_1"},
					{Key: UserMgoFieldAge, Value: int32(1)},
					{Key: UserMgoFieldRegisterDate, Value: time.Now()},
				})
				getMore := mtest.CreateCursorResponse(1, "test.user", mtest.NextBatch, bson.D{
					{Key: "_id", Value: primitive.NewObjectID()},
					{Key: UserMgoFieldUserId, Value: uint64(2)},
					{Key: UserMgoFieldUsername, Value: "test_username_2"},
					{Key: UserMgoFieldAge, Value: int32(2)},
					{Key: UserMgoFieldRegisterDate, Value: time.Now()},
				})
				killCursors := mtest.CreateCursorResponse(0, "test.user", mtest.NextBatch)
				t.AddMockResponses(first, getMore, killCursors)
				cursor, err := Get_UserMgr().Query(s.ctx, c.query, c.limit, c.offset, c.sortFields)
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

// TestFindByUsernameAge tests the FindByUsernameAge method
func (s *UserSuite) TestFindByUsernameAge() {
	cases := []struct {
		name       string
		username   string
		age        int32
		limit      int
		offset     int
		sortFields interface{}
		wantErr    bool
	}{
		{
			name:       "findbyusernameage_success",
			username:   "test_username",
			age:        int32(1),
			limit:      10,
			offset:     0,
			sortFields: nil,
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
				first := mtest.CreateCursorResponse(1, "test.user", mtest.FirstBatch, bson.D{
					{Key: "_id", Value: primitive.NewObjectID()},
					{Key: UserMgoFieldUserId, Value: uint64(1)},
					{Key: UserMgoFieldUsername, Value: "test_username"},
					{Key: UserMgoFieldAge, Value: int32(1)},
					{Key: UserMgoFieldRegisterDate, Value: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)},
				})
				killCursors := mtest.CreateCursorResponse(0, "test.user", mtest.NextBatch)
				t.AddMockResponses(first, killCursors)
				result, err := Get_UserMgr().FindByUsernameAge(s.ctx, c.username, c.age, c.limit, c.offset, c.sortFields)
				if c.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					assert.NotNil(t, result)
					assert.Len(t, result, 1)
					assert.Equal(t, c.username, result[0].Username)
					assert.Equal(t, c.age, result[0].Age)
				}
			})
		})
	}
}

// TestFindByUsername tests the FindByUsername method
func (s *UserSuite) TestFindByUsername() {
	cases := []struct {
		name       string
		username   string
		limit      int
		offset     int
		sortFields interface{}
		wantErr    bool
	}{
		{
			name:       "findbyusername_success",
			username:   "test_username",
			limit:      10,
			offset:     0,
			sortFields: nil,
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
				first := mtest.CreateCursorResponse(1, "test.user", mtest.FirstBatch, bson.D{
					{Key: "_id", Value: primitive.NewObjectID()},
					{Key: UserMgoFieldUserId, Value: uint64(1)},
					{Key: UserMgoFieldUsername, Value: "test_username"},
					{Key: UserMgoFieldAge, Value: int32(1)},
					{Key: UserMgoFieldRegisterDate, Value: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)},
				})
				killCursors := mtest.CreateCursorResponse(0, "test.user", mtest.NextBatch)
				t.AddMockResponses(first, killCursors)
				result, err := Get_UserMgr().FindByUsername(s.ctx, c.username, c.limit, c.offset, c.sortFields)
				if c.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					assert.NotNil(t, result)
					assert.Len(t, result, 1)
					assert.Equal(t, c.username, result[0].Username)
				}
			})
		})
	}
}

// TestFindByAge tests the FindByAge method
func (s *UserSuite) TestFindByAge() {
	cases := []struct {
		name       string
		age        int32
		limit      int
		offset     int
		sortFields interface{}
		wantErr    bool
	}{
		{
			name:       "findbyage_success",
			age:        int32(1),
			limit:      10,
			offset:     0,
			sortFields: nil,
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
				first := mtest.CreateCursorResponse(1, "test.user", mtest.FirstBatch, bson.D{
					{Key: "_id", Value: primitive.NewObjectID()},
					{Key: UserMgoFieldUserId, Value: uint64(1)},
					{Key: UserMgoFieldUsername, Value: "test_username"},
					{Key: UserMgoFieldAge, Value: int32(1)},
					{Key: UserMgoFieldRegisterDate, Value: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)},
				})
				killCursors := mtest.CreateCursorResponse(0, "test.user", mtest.NextBatch)
				t.AddMockResponses(first, killCursors)
				result, err := Get_UserMgr().FindByAge(s.ctx, c.age, c.limit, c.offset, c.sortFields)
				if c.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					assert.NotNil(t, result)
					assert.Len(t, result, 1)
					assert.Equal(t, c.age, result[0].Age)
				}
			})
		})
	}
}

// TestFindByRegisterDate tests the FindByRegisterDate method
func (s *UserSuite) TestFindByRegisterDate() {
	cases := []struct {
		name         string
		registerdate time.Time
		limit        int
		offset       int
		sortFields   interface{}
		wantErr      bool
	}{
		{
			name:         "findbyregisterdate_success",
			registerdate: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
			limit:        10,
			offset:       0,
			sortFields:   nil,
			wantErr:      false,
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
				first := mtest.CreateCursorResponse(1, "test.user", mtest.FirstBatch, bson.D{
					{Key: "_id", Value: primitive.NewObjectID()},
					{Key: UserMgoFieldUserId, Value: uint64(1)},
					{Key: UserMgoFieldUsername, Value: "test_username"},
					{Key: UserMgoFieldAge, Value: int32(1)},
					{Key: UserMgoFieldRegisterDate, Value: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)},
				})
				killCursors := mtest.CreateCursorResponse(0, "test.user", mtest.NextBatch)
				t.AddMockResponses(first, killCursors)
				result, err := Get_UserMgr().FindByRegisterDate(s.ctx, c.registerdate, c.limit, c.offset, c.sortFields)
				if c.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					assert.NotNil(t, result)
					assert.Len(t, result, 1)
					assert.Equal(t, c.registerdate.Unix(), result[0].RegisterDate.Unix())
				}
			})
		})
	}
}

func TestUser(t *testing.T) {
	suite.Run(t, &UserSuite{})
}
