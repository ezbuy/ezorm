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

type UserBlogSuite struct {
	suite.Suite
	ctx context.Context
}

func (s *UserBlogSuite) SetupTest() {
	s.ctx = context.TODO()
}

// TearDownTest runs after each test in the suite.
func (s *UserBlogSuite) TearDownTest() {
}

// TestSave tests the Save method with various scenarios
func (s *UserBlogSuite) TestSave() {
	cases := []struct {
		name    string
		data    UserBlog
		wantErr bool
	}{
		{
			name: "save_success",
			data: UserBlog{
				// ID will be auto-generated
				UserId:  1,
				BlogId:  1,
				Content: "test_content",
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
func (s *UserBlogSuite) TestFindOneAndSave() {
	cases := []struct {
		name    string
		data    UserBlog
		query   bson.M
		wantErr bool
	}{
		{
			name: "findOneAndSave_success",
			data: UserBlog{
				// ID will be auto-generated
				UserId:  2,
				BlogId:  2,
				Content: "test_content",
			},
			query: bson.M{
				"Content": "test_content",
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
						{Key: UserBlogMgoFieldUserId, Value: c.data.UserId},
						{Key: UserBlogMgoFieldBlogId, Value: c.data.BlogId},
						{Key: UserBlogMgoFieldContent, Value: c.data.Content},
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
func (s *UserBlogSuite) TestInsertUnique() {
	cases := []struct {
		name    string
		data    UserBlog
		query   bson.M
		wantErr bool
	}{
		{
			name: "insertUnique_success",
			data: UserBlog{
				// ID will be auto-generated
				UserId:  3,
				BlogId:  3,
				Content: "unique_content",
			},
			query: bson.M{
				"Content": "unique_content",
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
func (s *UserBlogSuite) TestFindOne() {
	cases := []struct {
		name       string
		query      bson.M
		sortFields interface{}
		wantErr    bool
	}{
		{
			name: "findOne_success",
			query: bson.M{
				"Content": "test_content",
			},
			sortFields: UserBlogMgoSortField_WRP{UserBlogMgoSortFieldIDAsc},
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
				findResp := mtest.CreateCursorResponse(1, "test.userblog", mtest.FirstBatch, bson.D{
					{Key: "_id", Value: primitive.NewObjectID()},
					{Key: UserBlogMgoFieldUserId, Value: uint64(1)},
					{Key: UserBlogMgoFieldBlogId, Value: uint64(1)},
					{Key: UserBlogMgoFieldContent, Value: "test_content"},
				})
				t.AddMockResponses(findResp)
				result, err := Get_UserBlogMgr().FindOne(s.ctx, c.query, c.sortFields)
				if c.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					assert.NotNil(t, result)
					assert.Equal(t, "test_content", result.Content)
				}
			})
		})
	}
}

// TestQuery tests the Query method
func (s *UserBlogSuite) TestQuery() {
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
			sortFields: UserBlogMgoSortField_WRP{UserBlogMgoSortFieldIDAsc},
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
				first := mtest.CreateCursorResponse(1, "test.userblog", mtest.FirstBatch, bson.D{
					{Key: "_id", Value: primitive.NewObjectID()},
					{Key: UserBlogMgoFieldUserId, Value: uint64(1)},
					{Key: UserBlogMgoFieldBlogId, Value: uint64(1)},
					{Key: UserBlogMgoFieldContent, Value: "test_content_1"},
				})
				getMore := mtest.CreateCursorResponse(1, "test.userblog", mtest.NextBatch, bson.D{
					{Key: "_id", Value: primitive.NewObjectID()},
					{Key: UserBlogMgoFieldUserId, Value: uint64(2)},
					{Key: UserBlogMgoFieldBlogId, Value: uint64(2)},
					{Key: UserBlogMgoFieldContent, Value: "test_content_2"},
				})
				killCursors := mtest.CreateCursorResponse(0, "test.userblog", mtest.NextBatch)
				t.AddMockResponses(first, getMore, killCursors)
				cursor, err := Get_UserBlogMgr().Query(s.ctx, c.query, c.limit, c.offset, c.sortFields)
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

// TestFindByUserId tests the FindByUserId method
func (s *UserBlogSuite) TestFindByUserId() {
	cases := []struct {
		name         string
		field_userid uint64
		limit        int
		offset       int
		sortFields   interface{}
		wantErr      bool
	}{
		{
			name:         "findbyuserid_success",
			field_userid: uint64(1),
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
				first := mtest.CreateCursorResponse(1, "test.userblog", mtest.FirstBatch, bson.D{
					{Key: "_id", Value: primitive.NewObjectID()},
					{Key: UserBlogMgoFieldUserId, Value: uint64(1)},
					{Key: UserBlogMgoFieldBlogId, Value: uint64(1)},
					{Key: UserBlogMgoFieldContent, Value: "test_content"},
				})
				killCursors := mtest.CreateCursorResponse(0, "test.userblog", mtest.NextBatch)
				t.AddMockResponses(first, killCursors)
				result, err := Get_UserBlogMgr().FindByUserId(s.ctx, c.field_userid, c.limit, c.offset, c.sortFields)
				if c.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					assert.NotNil(t, result)
					assert.Len(t, result, 1)
					assert.Equal(t, c.field_userid, result[0].UserId)
				}
			})
		})
	}
}

func TestUserBlog(t *testing.T) {
	suite.Run(t, &UserBlogSuite{})
}
