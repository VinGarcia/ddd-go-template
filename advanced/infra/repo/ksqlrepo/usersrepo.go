package ksqlrepo

import (
	"context"
	"time"

	"github.com/vingarcia/ddd-go-template/advanced/domain"
	"github.com/vingarcia/ksql"
)

type UsersRepo struct {
	db ksql.Provider
}

func NewUsersRepo(db ksql.Provider) UsersRepo {
	return UsersRepo{
		db: db,
	}
}

func (u UsersRepo) UpsertUser(ctx context.Context, user domain.User) (userID int, _ error) {
	now := time.Now()
	user.UpdatedAt = &now
	err := u.db.Update(ctx, domain.UsersTable, &user)
	if err == ksql.ErrRecordNotFound {
		user.CreatedAt = &now
		err = u.db.Insert(ctx, domain.UsersTable, &user)
	}
	if err != nil {
		return 0, domain.InternalErr("unexpected error when saving user", map[string]interface{}{
			"user":  user,
			"error": err.Error(),
		})
	}

	return user.ID, nil
}

func (u UsersRepo) GetUser(ctx context.Context, userID int) (domain.User, error) {
	var user domain.User
	err := u.db.QueryOne(ctx, &user, "FROM users WHERE id = $1", userID)
	if err == ksql.ErrRecordNotFound {
		return domain.User{}, domain.NotFoundErr("no user found with provided id", map[string]interface{}{
			"user_id": userID,
		})
	}
	if err != nil {
		return domain.User{}, domain.InternalErr("unexpected error when saving user", map[string]interface{}{
			"user_id": userID,
			"error":   err.Error(),
		})
	}

	return user, nil
}
