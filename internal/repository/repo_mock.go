package repository

import (
	"context"
	"github/MaysHroub/gator/internal/database"

	"github.com/google/uuid"
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

func (mk *MockRepository) GetFeedFollowsForUser(ctx context.Context, name string) ([]database.GetFeedFollowsForUserRow, error) {
	args := mk.Called(ctx, name)
	return args.Get(0).([]database.GetFeedFollowsForUserRow), args.Error(1)
}

func (mk *MockRepository) DeleteFeedFollowByUserAndURL(ctx context.Context, arg database.DeleteFeedFollowByUserAndURLParams) error {
	args := mk.Called(ctx, arg)
	return args.Error(0)
}

func (mk *MockRepository) GetNextFeedToFetch(ctx context.Context) (database.GetNextFeedToFetchRow, error) {
	args := mk.Called(ctx)
	return args.Get(0).(database.GetNextFeedToFetchRow), args.Error(1)
}

func (mk *MockRepository) MarkFeedFetched(ctx context.Context, id uuid.UUID) error {
	args := mk.Called(ctx, id)
	return args.Error(0)
}

func (mk *MockRepository) CreatePost(ctx context.Context, arg database.CreatePostParams) (database.Post, error) {
	args := mk.Called(ctx, arg)
	return args.Get(0).(database.Post), args.Error(1)
}

func (mk *MockRepository) GetPostsForUser(ctx context.Context, arg database.GetPostsForUserParams) ([]database.Post, error) {
	args := mk.Called(ctx, arg)
	return args.Get(0).([]database.Post), args.Error(1)
}
