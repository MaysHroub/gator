package repository

import (
	"context"
	"github/MaysHroub/gator/internal/database"
)

type UserStore interface {
	CreateUser(ctx context.Context, arg database.CreateUserParams) (database.User, error)
}