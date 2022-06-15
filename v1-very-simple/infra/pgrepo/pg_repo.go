package pgrepo

import (
	"context"

	"github.com/vingarcia/ddd-go-template/v1-very-simple/domain"
	"github.com/vingarcia/ksql"
	"github.com/vingarcia/ksql/adapters/kpgx"
)

// Client implements the repo.Users interface by using the ksql database.
type Client struct {
	db ksql.Provider
}

// NewClient instantiates a new Client
func NewClient(ctx context.Context, postgresURL string) (Client, error) {
	db, err := kpgx.New(ctx, postgresURL, ksql.Config{})
	if err != nil {
		return Client{}, domain.InternalErr("unable to start database", domain.LogBody{
			"postgres_url": postgresURL,
			"error":        err.Error(),
		})
	}

	return Client{
		db: db,
	}, nil
}

// ChangeUserEmail implements the repo.Users interface
func (c Client) ChangeUserEmail(ctx context.Context, userID int, newEmail string) error {
	return changeUserEmail(ctx, c.db, userID, newEmail)
}

// UpsertUser implements the repo.Users interface
func (c Client) UpsertUser(ctx context.Context, user domain.User) (userID int, _ error) {
	return upsertUser(ctx, c.db, user)
}

// GetUser implements the repo.Users interface
func (c Client) GetUser(ctx context.Context, userID int) (domain.User, error) {
	return getUser(ctx, c.db, userID)
}

// GetUserByEmail implements the repo.Users interface
func (c Client) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	return getUserByEmail(ctx, c.db, email)
}
