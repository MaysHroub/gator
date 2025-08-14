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
}
