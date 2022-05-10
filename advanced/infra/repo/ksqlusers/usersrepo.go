package ksqlusers

import (
	"context"

	"github.com/vingarcia/ddd-go-template/advanced/domain"
	"github.com/vingarcia/ksql"
)

// UsersRepo implements the repo.Users interface by using the ksql database.
type UsersRepo struct {
	db ksql.Provider
}

// New instantiates a new UsersRepo
func New(db ksql.Provider) UsersRepo {
	return UsersRepo{
		db: db,
	}
}

// ChangeUserEmail implements the repo.Users interface
func (u UsersRepo) ChangeUserEmail(ctx context.Context, userID int, newEmail string) error {
	return changeUserEmail(ctx, u.db, userID, newEmail)
}

// UpsertUser implements the repo.Users interface
func (u UsersRepo) UpsertUser(ctx context.Context, user domain.User) (userID int, _ error) {
	return upsertUser(ctx, u.db, user)
}

// GetUser implements the repo.Users interface
func (u UsersRepo) GetUser(ctx context.Context, userID int) (domain.User, error) {
	return getUser(ctx, u.db, userID)
}

// GetUserByEmail implements the repo.Users interface
func (u UsersRepo) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	return getUserByEmail(ctx, u.db, email)
}
