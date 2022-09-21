package users

import (
	"context"
	"testing"

	"github.com/vingarcia/ddd-go-template/v2-domain-adapters-and-helpers/adapters/log"
	"github.com/vingarcia/ddd-go-template/v2-domain-adapters-and-helpers/adapters/repo"
	"github.com/vingarcia/ddd-go-template/v2-domain-adapters-and-helpers/domain"
	tt "github.com/vingarcia/ddd-go-template/v2-domain-adapters-and-helpers/helpers/testtools"
)

func TestNewService(t *testing.T) {
	t.Run("should build the struct correctly", func(t *testing.T) {
		logger := log.Mock{}
		usersRepo := repo.UsersMock{}

		svc := NewService(&logger, &usersRepo)

		tt.AssertEqual(t, svc.logger, &logger)
		tt.AssertEqual(t, svc.usersRepo, &usersRepo)
	})
}

func TestUpsertUser(t *testing.T) {
	ctx := context.Background()

	t.Run("should upsert a user correctly", func(t *testing.T) {
		var userArg domain.User

		svc := NewService(
			log.Mock{},
			repo.UsersMock{
				UpsertUserFn: func(ctx context.Context, user domain.User) (userID int, _ error) {
					// Collect any arguments relevant to this test
					// so we can assert its value at the end of the test:
					userArg = user

					// Return fake values to provoke the behavior we want on this test:
					return 42, nil
				},
			},
		)

		userID, err := svc.UpsertUser(ctx, domain.User{
			Name:  "fakeName",
			Email: "fake@email.com",
			Age:   24,
		})
		tt.AssertNoErr(t, err)

		tt.AssertEqual(t, userID, 42)
		tt.AssertEqual(t, userArg, domain.User{
			Name:  "fakeName",
			Email: "fake@email.com",
			Age:   24,
		})
	})
}
