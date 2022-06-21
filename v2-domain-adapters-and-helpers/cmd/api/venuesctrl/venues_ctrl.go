package venuesctrl

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/vingarcia/ddd-go-template/v2-domain-adapters-and-helpers/domain/venues"
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
	venuesService venues.Service
}

func NewController(venuesService venues.Service) Controller {
	return Controller{
		venuesService: venuesService,
	}
}

func (c Controller) GetVenuesByCoordinates(ctx *fiber.Ctx) error {
	latitude := ctx.Params("latitude")
	longitude := ctx.Params("latitude")

	venues, err := c.venuesService.GetVenues(ctx.Context(), latitude, longitude)
	if err != nil {
		return err
	}

	rawJSON, err := json.Marshal(venues)
	if err != nil {
		return fmt.Errorf("error building GET venues response JSON: %s", err)
	}

	return ctx.Send(rawJSON)
}

func (c Controller) GetDetails(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	venue, err := c.venuesService.GetVenue(ctx.Context(), id)
	if err != nil {
		return err
	}

	return ctx.Send(venue)
}
