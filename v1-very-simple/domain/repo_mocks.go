package domain

import "context"

type UsersRepoMock struct {
	GetUserFn        func(ctx context.Context, userID int) (User, error)
	UpsertUserFn     func(ctx context.Context, user User) (userID int, err error)
	GetUserByEmailFn func(ctx context.Context, email string) (User, error)
}

func (m UsersRepoMock) GetUser(ctx context.Context, userID int) (User, error) {
	return m.GetUserFn(ctx, userID)
}

func (m UsersRepoMock) UpsertUser(ctx context.Context, user User) (userID int, err error) {
	return m.UpsertUserFn(ctx, user)
}

func (m UsersRepoMock) GetUserByEmail(ctx context.Context, email string) (User, error) {
	return m.GetUserByEmailFn(ctx, email)
}
