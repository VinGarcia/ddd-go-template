package users

import (
	"context"

	"github.com/vingarcia/ddd-go-template/v1-simple-with-short-interface-names/domain"
	"github.com/vingarcia/ddd-go-template/v1-simple-with-short-interface-names/infra/log"
	"github.com/vingarcia/ddd-go-template/v1-simple-with-short-interface-names/infra/repo"
)

// Usually its here where the business logic complexity builds up,
// but since this is just an example both these functions are actually
// very simple, but in real world scenarios you would want to make
// these contain all your business logic.

type Service struct {
	logger    log.Provider
	usersRepo repo.Users
}

func NewService(
	logger log.Provider,
	usersRepo repo.Users,
) Service {
	return Service{
		logger:    logger,
		usersRepo: usersRepo,
	}
}

func (s Service) UpsertUser(ctx context.Context, user domain.User) (userID int, _ error) {
	userID, err := s.usersRepo.UpsertUser(ctx, user)
	if err != nil {
		return 0, err
	}

	s.logger.Info(ctx, "user created", log.Body{
		"user_id": userID,
	})
	return userID, nil
}

func (s Service) GetUser(ctx context.Context, userID int) (domain.User, error) {
	user, err := s.usersRepo.GetUser(ctx, userID)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}
