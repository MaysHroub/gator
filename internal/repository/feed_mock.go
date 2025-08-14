package repository

import (
	"context"
	"github/MaysHroub/gator/internal/database"

	"github.com/stretchr/testify/mock"
)

type MockFeedStore struct {
	mock.Mock
}

func (mk *MockFeedStore) CreateFeed(ctx context.Context, arg database.CreateFeedParams) (database.Feed, error) {
	args := mk.Called(ctx, arg)
	return args.Get(0).(database.Feed), args.Error(1)
}
