package pgrepo

import (
	"context"

	"github.com/vingarcia/ddd-go-template/v2-domain-adapters-and-helpers/adapters/log"
	"github.com/vingarcia/ddd-go-template/v2-domain-adapters-and-helpers/domain"
	"github.com/vingarcia/ksql"
	"github.com/vingarcia/ksql/adapters/kpgx"
)

// Repo implements the repo.Users interface by using the ksql database.
type Repo struct {
	db ksql.Provider
}

// New instantiates a new Repo
func New(ctx context.Context, postgresURL string) (Repo, error) {
	db, err := kpgx.New(ctx, postgresURL, ksql.Config{})
	if err != nil {
		return Repo{}, domain.InternalErr("unable to start database", log.Body{
			"error": err.Error(),
		})
	}

	return Repo{
		db: db,
	}, nil
}

// ChangeUserEmail implements the repo.Users interface
func (u Repo) ChangeUserEmail(ctx context.Context, userID int, newEmail string) error {
	return changeUserEmail(ctx, u.db, userID, newEmail)
}

// UpsertUser implements the repo.Users interface
func (u Repo) UpsertUser(ctx context.Context, user domain.User) (userID int, _ error) {
	return upsertUser(ctx, u.db, user)
}

// GetUser implements the repo.Users interface
func (u Repo) GetUser(ctx context.Context, userID int) (domain.User, error) {
	return getUser(ctx, u.db, userID)
}

// GetUserByEmail implements the repo.Users interface
func (u Repo) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	return getUserByEmail(ctx, u.db, email)
}
