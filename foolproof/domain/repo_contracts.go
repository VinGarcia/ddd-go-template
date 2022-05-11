package domain

import "context"

// UsersRepo represents the operations we use for
// retrieving a user from a persistent storage
type UsersRepo interface {
	GetUser(ctx context.Context, userID int) (User, error)
	UpsertUser(ctx context.Context, user User) (userID int, err error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
}
