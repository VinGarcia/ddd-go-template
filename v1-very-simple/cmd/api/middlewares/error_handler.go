package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vingarcia/ddd-go-template/v1-very-simple/domain"
)

// HandleError is responsible for converting domain errors to HTTP errors
// simplifying error handling overall.
func HandleError(logger domain.LogProvider) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		err := c.Next()
		if err == nil {
			return nil
		}

		req := c.Request()
		status, body := domain.HandleDomainErrAsHTTP(
			c.Context(),
			logger,
			err,
			string(req.Header.Method()),
			string(req.RequestURI()),
		)
		c.Status(status).Send(body)
		return nil
	}
}
