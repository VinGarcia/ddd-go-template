package domain

import "context"

// UserRepo represents the operations we use for
// retrieving a user from a persistent storage
type UsersRepo interface {
	GetUser(ctx context.Context, userID int) (User, error)
	UpsertUser(ctx context.Context, user User) (userID int, err error)
}
