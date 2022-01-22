package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vingarcia/ddd-go-layout/domain"
)

func HandleRequestID() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		requestID := c.Get("request-id")
		if requestID == "" {
			requestID = domain.GenerateRequestID()
		}

		c.Locals(domain.RequestIDKey, requestID)
		return c.Next()
	}
}
