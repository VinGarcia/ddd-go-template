package repo

import (
	"context"

	"github.com/vingarcia/ddd-go-template/v1-simple-with-short-interface-names/domain"
)

type UsersMock struct {
	GetUserFn        func(ctx context.Context, userID int) (domain.User, error)
	UpsertUserFn     func(ctx context.Context, user domain.User) (userID int, err error)
	GetUserByEmailFn func(ctx context.Context, email string) (domain.User, error)
}

func (m UsersMock) GetUser(ctx context.Context, userID int) (domain.User, error) {
	return m.GetUserFn(ctx, userID)
}

func (m UsersMock) UpsertUser(ctx context.Context, user domain.User) (userID int, err error) {
	return m.UpsertUserFn(ctx, user)
}

func (m UsersMock) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	return m.GetUserByEmailFn(ctx, email)
}
