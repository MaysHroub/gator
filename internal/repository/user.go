// Package repository provides interfaces and implementations for user data storage and retrieval.
package repository

import (
	"context"
	"github/MaysHroub/gator/internal/database"
)

type UserStore interface {
	CreateUser(ctx context.Context, arg database.CreateUserParams) (database.User, error)
	GetUser(ctx context.Context, name string) (database.User, error)
	DeleteAllUsers(ctx context.Context) error
	GetUsersNames(ctx context.Context) ([]string, error)
}