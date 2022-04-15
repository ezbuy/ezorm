package mysql

import (
	"context"
	"testing"

	mock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetUser(t *testing.T) {
	mockDB, mock, err := mock.New()
	if !assert.NoError(t, err) {
		return
	}
	mock.ExpectQuery(_GetUserSQL).WithArgs("me").WillReturnRows(
		mock.NewRows([]string{"name"}).AddRow("me"),
	)

	ctx := context.TODO()
	resp, err := SQL.GetUser(ctx, &GetUserReq{
		Name: "me",
	}, WithDB(mockDB))
	if !assert.NoError(t, err) {
		return
	}
	if !assert.Equal(t, 1, len(resp)) {
		return
	}
	if !assert.Equal(t, "me", resp[0].Name) {
		return
	}
}
