package mysql

import (
	"context"
	"fmt"
	"testing"

	mock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetUser(t *testing.T) {
	mockDB, mock, err := mock.New()
	if !assert.NoError(t, err) {
		return
	}

	req := &GetUserReq{
		Name: "me",
	}

	sql := fmt.Sprintf(_GetUserSQL, req.Condition())

	mock.ExpectQuery(sql).WithArgs("me").WillReturnRows(
		mock.NewRows([]string{"name"}).AddRow("me"),
	)

	ctx := context.TODO()
	resp, err := GetRawQuery().GetUser(ctx, req, WithDB(mockDB))
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
