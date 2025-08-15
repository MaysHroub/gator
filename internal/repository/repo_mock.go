package repository

import (
	"context"
	"github/MaysHroub/gator/internal/database"

	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (mk *MockRepository) CreateUser(ctx context.Context, arg database.CreateUserParams) (database.User, error) {
	args := mk.Called(ctx, arg)
	return args.Get(0).(database.User), args.Error(1)
}

func (mk *MockRepository) GetUserByName(ctx context.Context, name string) (database.User, error) {
	args := mk.Called(ctx, name)
	return args.Get(0).(database.User), args.Error(1)
}

func (mk *MockRepository) DeleteAllUsers(ctx context.Context) error {
	args := mk.Called(ctx)
	return args.Error(0)
}

func (mk *MockRepository) GetNamesOfAllUsers(ctx context.Context) ([]string, error) {
	args := mk.Called(ctx)
	return args.Get(0).([]string), args.Error(1)
}

func (mk *MockRepository) CreateFeed(ctx context.Context, arg database.CreateFeedParams) (database.Feed, error) {
	args := mk.Called(ctx, arg)
	return args.Get(0).(database.Feed), args.Error(1)
}

func (mk *MockRepository) GetAllFeeds(ctx context.Context) ([]database.GetAllFeedsRow, error) {
	args := mk.Called(ctx)
	return args.Get(0).([]database.GetAllFeedsRow), args.Error(1)
}

func (mk *MockRepository) GetFeedByURL(ctx context.Context, url string) (database.Feed, error) {
	args := mk.Called(ctx, url)
	return args.Get(0).(database.Feed), args.Error(1)
}

func (mk *MockRepository) CreateFeedFollow(ctx context.Context, arg database.CreateFeedFollowParams) ([]database.CreateFeedFollowRow, error) {
	args := mk.Called(ctx, arg)
	return args.Get(0).([]database.CreateFeedFollowRow), args.Error(1)
}