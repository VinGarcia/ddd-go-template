package usersctrl

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/vingarcia/ddd-go-layout/domain"
	"github.com/vingarcia/ddd-go-layout/domain/users"
)

// The Controller can depend directly on the service with no intermediary interface.
// This is ok because I don't usually test controllers (i.e. no need for mocking the service)
//
// And I don't test it for two reasons:
//
// 1. Controllers are expected to be simple, so unit testing it should not be that useful.
// 2. There is a more complete way of testing the whole program including the controllers
// which is to write api tests on `cmd/api/main_test.go` and literally starting the API and
// running requests on it, so if testing the controller is important I prefer writing API tests.
type Controller struct {
	usersService users.Service
}

func NewController(usersService users.Service) Controller {
	return Controller{
		usersService: usersService,
	}
}

func (c Controller) UpsertUser(ctx *fiber.Ctx) error {
	// Intermediary structure so I don't expose my internal
	// user representation to the outside world:
	var user struct {
		UserID int    `json:"user_id"`
		Name   string `json:"name"`
		Age    int    `json:"age"`
	}
	err := json.Unmarshal(ctx.Body(), &user)
	if err != nil {
		return domain.BadRequestErr("unable to parse payload as JSON", map[string]interface{}{
			"payload": string(ctx.Body()),
			"error":   err.Error(),
		})
	}

	userID, err := c.usersService.UpsertUser(ctx.Context(), domain.User{
		// Showcasing that my internal model might differ from the API,
		// in this case the internal name for the ID attribute is just ID not `UserID`:
		ID:   user.UserID,
		Name: user.Name,
		Age:  user.Age,
	})
	if err != nil {
		return err
	}

	return ctx.JSON(map[string]interface{}{
		"status":  "success",
		"user_id": userID,
	})
}

func (c Controller) GetUser(ctx *fiber.Ctx) error {
	userID, err := ctx.ParamsInt("id")
	if err != nil {
		return domain.BadRequestErr("the input user id is not a valid integer", map[string]interface{}{
			"received_id": ctx.Params(":id"),
		})
	}

	user, err := c.usersService.GetUser(ctx.Context(), userID)
	if err != nil {
		return err
	}

	// Again using intermediary structs (or a map) is useful for decoupling
	// the internal entities from what is exposed on the web:
	return ctx.JSON(map[string]interface{}{
		"id":   userID,
		"name": user.Name,
		"age":  user.Age,
	})
}
