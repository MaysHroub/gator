package repository

import (
	"context"
	"github/MaysHroub/gator/internal/database"

	"github.com/stretchr/testify/mock"
)

type MockUserStore struct {
	mock.Mock
}

func (mk *MockUserStore) CreateUser(ctx context.Context, arg database.CreateUserParams) (database.User, error) {
	args := mk.Called(ctx, arg)
	return args.Get(0).(database.User), args.Error(1)
}