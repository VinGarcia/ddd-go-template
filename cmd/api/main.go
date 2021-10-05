package main

import (
	"context"
	"time"

	"github.com/vingarcia/my-ddd-go-layout/domain"
	"github.com/vingarcia/my-ddd-go-layout/domain/venues"
	"github.com/vingarcia/my-ddd-go-layout/infra/env"
	"github.com/vingarcia/my-ddd-go-layout/infra/jsonlogs"
	"github.com/vingarcia/my-ddd-go-layout/infra/memorycache"
	"github.com/vingarcia/my-ddd-go-layout/infra/rest"

	routing "github.com/jackwhelpton/fasthttp-routing/v2"
	"github.com/valyala/fasthttp"

	adapter "github.com/vingarcia/go-adapter"
)

func main() {
	ctx := context.Background()

	// Read all configs at once so its easy to spot all of them:
	port := env.GetString("PORT", "80")
	logLevel := env.GetString("LOG_LEVEL", "INFO")
	foursquareBaseURL := env.MustGetString("FOURSQUARE_BASE_URL")
	foursquareClientID := env.MustGetString("FOURSQUARE_CLIENT_ID")
	foursquareSecret := env.MustGetString("FOURSQUARE_SECRET")

	// Dependency Injection goes here:
	logger := jsonlogs.New(logLevel)

	restClient := rest.New(30 * time.Second)

	cacheClient := memorycache.New(24*time.Hour, 10*time.Minute)

	venuesService := venues.New(
		logger,
		restClient,
		cacheClient,
		foursquareBaseURL,
		foursquareClientID,
		foursquareSecret,
	)

	// The controllers handle HTTP stuff so the services can be kept as simple as possible
	// only working on top of the domain language, i.e. types and interfaces from the domain/ package
	venuesController := venues.NewController(venuesService)

	// Any framework you need for serving HTTP or GRPC goes in the main package,
	// since it is allowed to know and reference everything, it is then, also ok for it to
	// have direct external dependencies such as the fasthttp framework and router:
	router := routing.New()
	router.Get("/ping", func(ctx *routing.Context) error {
		ctx.SetBody([]byte("pong"))
		return nil
	})
	router.Get("/venues/<latitude>,<longitude>", adapter.Adapt(venuesController.GetVenuesByCoordinates))
	router.Get("/venues/details/<id>", adapter.Adapt(venuesController.GetDetails))

	logger.Info(ctx, "server-starting-up", domain.LogBody{
		"port": port,
	})
	if err := fasthttp.ListenAndServe(":"+port, router.HandleRequest); err != nil {
		logger.Error(ctx, "server-stopped-with-an-error", domain.LogBody{
			"error": err.Error(),
		})
	}
}
