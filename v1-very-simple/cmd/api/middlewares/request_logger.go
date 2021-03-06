package middlewares

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/vingarcia/ddd-go-template/v1-very-simple/domain"
)

// RequestLogger will log every request including the request payload,
// status code, duration and information about any errors.
func RequestLogger(logger domain.LogProvider) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		startTime := time.Now()

		err := c.Next()
		if err != nil {
			logger.Error(c.Context(), "request-failed", domain.LogBody{
				"error":        err.Error(),
				"route":        c.Method() + " " + c.Path(),
				"request_body": string(c.Body()),
				"duration_ms":  time.Since(startTime).Milliseconds(),
			})
			return err
		}

		logger.Info(c.Context(), "request-completed", domain.LogBody{
			"route":        c.Method() + " " + c.Path(),
			"request_body": string(c.Body()),
			"duration_ms":  time.Since(startTime).Milliseconds(),
		})

		return nil
	}
}
