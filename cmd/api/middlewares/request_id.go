package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vingarcia/ddd-go-layout/domain"
)

func HandleRequestID() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		key, requestID := domain.GenerateRequestID()
		c.Locals(key, requestID)
		return c.Next()
	}
}
