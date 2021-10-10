package main

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/vingarcia/ddd-go-layout/domain"
	"github.com/vingarcia/ddd-go-layout/domain/venues"
	"github.com/vingarcia/ddd-go-layout/infra/env"
	"github.com/vingarcia/ddd-go-layout/infra/jsonlogs"
	"github.com/vingarcia/ddd-go-layout/infra/memorycache"
	"github.com/vingarcia/ddd-go-layout/infra/redis"
	"github.com/vingarcia/ddd-go-layout/infra/rest"

	adapter "github.com/vingarcia/go-adapter/fiber/v2"
)

func main() {
	ctx := context.Background()

	// Read all configs at once so its easy to spot all of them:
	port := env.GetString("PORT", "80")
	logLevel := env.GetString("LOG_LEVEL", "INFO")
	foursquareBaseURL := env.MustGetString("FOURSQUARE_BASE_URL")
	foursquareClientID := env.MustGetString("FOURSQUARE_CLIENT_ID")
	foursquareSecret := env.MustGetString("FOURSQUARE_SECRET")
	redisURL := env.GetString("REDIS_URL", "")
	redisPassword := env.GetString("REDIS_PASSWORD", "")

	// Dependency Injection goes here:
	logger := jsonlogs.New(logLevel)

	restClient := rest.New(30 * time.Second)

	var cacheClient domain.CacheProvider
	if redisURL != "" {
		cacheClient = redis.New(redisURL, redisPassword, 24*time.Hour)
	} else {
		cacheClient = memorycache.New(24*time.Hour, 10*time.Minute)
	}

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
	//
	// It should be kept here because the main package is the only one that is allowed
	// to depend on anything, and also because this logic is unique to this endpoint,
	// so you won't reuse it anywhere else.
	app := fiber.New()

	app.Use(handleRequestID())
	app.Use(handleError(logger))

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})
	app.Get("/venues/<latitude>,<longitude>", adapter.Adapt(venuesController.GetVenuesByCoordinates))
	app.Get("/venues/details/<id>", adapter.Adapt(venuesController.GetDetails))

	logger.Info(ctx, "server-starting-up", domain.LogBody{
		"port": port,
	})

	if err := app.Listen(":" + port); err != nil {
		logger.Error(ctx, "server-stopped-with-an-error", domain.LogBody{
			"error": err.Error(),
		})
	}
}
