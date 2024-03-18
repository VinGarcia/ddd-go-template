package main

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/vingarcia/ddd-go-template/v2-domain-adapters-and-helpers/domain"
	tt "github.com/vingarcia/ddd-go-template/v2-domain-adapters-and-helpers/helpers/testtools"
	"github.com/vingarcia/krest"
)

func usersTest(ctx context.Context, t *testing.T, data testData) {
	// These are integration tests, they are starting the API with a real database
	// and then making an HTTP request to this API.
	//
	// This is significantly more complex than unit tests but allows you to
	// fully test your middlewares and the communication with the database
	// which is great.
	//
	// Note that the way these tests were built, after you create a test for
	// a route for the first time, writing more tests becomes a lot simpler.
	t.Run("GET /users/:id", func(t *testing.T) {
		tests := []struct {
			desc               string
			requestUserID      string
			dbUsers            []domain.User
			expectedUser       domain.User
			expectErrToContain []string
		}{
			{
				desc:          "should retrieve a user correctly",
				requestUserID: "42",
				dbUsers: []domain.User{
					{
						ID:   42,
						Name: "FakeName1",
					},
				},
				expectedUser: domain.User{
					ID:   42,
					Name: "FakeName1",
				},
			},
			{
				desc:          "should return 404 if no user is found",
				requestUserID: "42",
				dbUsers: []domain.User{
					{
						ID:   43,
						Name: "NotTheUserYouAreLookingFor",
					},
				},
				expectErrToContain: []string{"NotFound", "42", "404"},
			},
		}

		for _, test := range tests {
			t.Run(test.desc, func(t *testing.T) {
				resetTestState(ctx, data)

				for _, u := range test.dbUsers {
					err := data.db.Insert(ctx, domain.UsersTable, &u)
					tt.AssertNoErr(t, err)
				}

				resp, err := data.http.Get(ctx, data.serverURL+"/users/"+test.requestUserID, krest.RequestData{
					// Any headers necessary for this request to work:
					Headers: map[string]any{
						"Authorization": "Bearer <fake token since we don't have auth yet>",
					},
				})
				if test.expectErrToContain != nil {
					tt.AssertErrContains(t, err, test.expectErrToContain...)
					return
				}
				tt.AssertNoErr(t, err)

				var user domain.User
				err = json.Unmarshal(resp.Body, &user)
				tt.AssertNoErr(t, err)

				tt.AssertEqual(t, user, test.expectedUser)
			})
		}
	})

	t.Run("POST /users", func(t *testing.T) {
		tests := []struct {
			desc               string
			requestBody        map[string]any
			dbUsers            []domain.User
			expectedUser       domain.User
			expectErrToContain []string
		}{
			{
				desc: "should create a user correctly",
				requestBody: map[string]any{
					"name": "NewUserName",
					"age":  22,
				},
				expectedUser: domain.User{
					Name: "NewUserName",
					Age:  22,
				},
			},
			{
				desc: "should update a user correctly",
				requestBody: map[string]any{
					"user_id": 42,
					"name":    "NewUserName",
					"age":     22,
				},
				dbUsers: []domain.User{
					{
						ID:   42,
						Name: "OldUserName",
						Age:  22,
					},
				},
				expectedUser: domain.User{
					ID:   42,
					Name: "NewUserName",
					Age:  22,
				},
			},
		}

		for _, test := range tests {
			t.Run(test.desc, func(t *testing.T) {
				resetTestState(ctx, data)

				for _, u := range test.dbUsers {
					err := data.db.Insert(ctx, domain.UsersTable, &u)
					tt.AssertNoErr(t, err)
				}

				resp, err := data.http.Post(ctx, data.serverURL+"/users", krest.RequestData{
					// Any headers necessary for this request to work:
					Headers: map[string]any{
						"Authorization": "Bearer <fake token since we don't have auth yet>",
					},
					// The krest library will convert the user to JSON for us:
					Body: test.requestBody,
				})
				if test.expectErrToContain != nil {
					tt.AssertErrContains(t, err, test.expectErrToContain...)
					return
				}
				tt.AssertNoErr(t, err)

				var respDTO struct {
					Status string `json:"status"`
					UserID int    `json:"user_id"`
				}
				err = json.Unmarshal(resp.Body, &respDTO)
				tt.AssertNoErr(t, err)

				tt.AssertEqual(t, respDTO.Status, "success")

				var user domain.User
				err = data.db.QueryOne(ctx, &user, "FROM users WHERE id = $1", respDTO.UserID)
				tt.AssertNoErr(t, err)

				// You can overwrite fields that are hard to compare:
				tt.AssertNotEqual(t, user.CreatedAt, time.Time{})
				tt.AssertNotEqual(t, user.UpdatedAt, time.Time{})
				user.CreatedAt = time.Time{}
				user.UpdatedAt = time.Time{}
				test.expectedUser.ID = respDTO.UserID

				tt.AssertEqual(t, user, test.expectedUser)
			})
		}
	})
}
