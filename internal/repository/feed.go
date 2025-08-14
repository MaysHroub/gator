package repository

import (
	"context"
	"github/MaysHroub/gator/internal/database"
)

type FeedStore interface {
	CreateFeed(ctx context.Context, arg database.CreateFeedParams) (database.Feed, error)
}