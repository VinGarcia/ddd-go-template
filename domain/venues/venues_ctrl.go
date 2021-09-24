package venues

import (
	"context"
	"encoding/json"
	"fmt"

	routing "github.com/jackwhelpton/fasthttp-routing/v2"
)

// The Controller can depend directly on the service with no intermediary interface.
// This is ok because I don't usually test controllers (i.e. no need for mocking the service)
//
// I don't test it for two reasons:
//
// 1. Controllers are expected to be simple, so unit testing it should not be that useful.
// 2. There is a more complete way of testing that test the whole program including the controllers
// which is to write api tests on `cmd/api/main_test.go`
type Controller struct {
	venuesService Service
}

func NewController(venuesService Service) Controller {
	return Controller{
		venuesService: venuesService,
	}
}

func (c Controller) GetVenuesByCoordinates(ctx *routing.Context, args struct {
	Latitude  string `path:"latitude"`
	Longitude string `path:"longitude"`
}) error {
	goCtx := context.Background()

	venues, err := c.venuesService.GetVenues(goCtx, args.Latitude, args.Longitude)
	if err != nil {
		return err
	}

	rawJSON, err := json.Marshal(venues)
	if err != nil {
		return fmt.Errorf("error building GET venues response JSON: %s", err)
	}
	ctx.SetBody(rawJSON)

	return nil
}

func (c Controller) GetDetails(ctx *routing.Context, args struct {
	ID string `path:"id"`
}) error {
	goCtx := context.Background()

	venue, err := c.venuesService.GetVenue(goCtx, args.ID)
	if err != nil {
		return err
	}

	ctx.SetBody(venue)

	return nil
}
