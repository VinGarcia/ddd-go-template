package repo

import (
	"context"

	"github.com/vingarcia/ddd-go-template/v1-simple-with-short-interface-names/domain"
)

// Users represents the operations we use for
// retrieving a user from a persistent storage
type Users interface {
	GetUser(ctx context.Context, userID int) (domain.User, error)
	UpsertUser(ctx context.Context, user domain.User) (userID int, err error)
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
}
