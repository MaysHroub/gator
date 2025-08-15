// Package repository provides interfaces and implementations for user data storage and retrieval.
package repository

import (
	"context"
	"github/MaysHroub/gator/internal/database"
)

type Repository interface {
	CreateUser(ctx context.Context, arg database.CreateUserParams) (database.User, error)
	GetUserByName(ctx context.Context, name string) (database.User, error)
	DeleteAllUsers(ctx context.Context) error
	GetNamesOfAllUsers(ctx context.Context) ([]string, error)
	CreateFeed(ctx context.Context, arg database.CreateFeedParams) (database.Feed, error)
	GetAllFeeds(ctx context.Context) ([]database.GetAllFeedsRow, error)
	GetFeedByURL(ctx context.Context, url string) (database.Feed, error)
	CreateFeedFollow(ctx context.Context, arg database.CreateFeedFollowParams) ([]database.CreateFeedFollowRow, error)
	GetFeedFollowsForUser(ctx context.Context, name string) ([]database.GetFeedFollowsForUserRow, error)
	DeleteFeedFollowByUserAndURL(ctx context.Context, arg database.DeleteFeedFollowByUserAndURLParams) error 
}
