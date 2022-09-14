package middlewares

import (
	"context"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/vingarcia/ddd-go-template/v2-domain-adapters-and-helpers/adapters/log"
	"github.com/vingarcia/ddd-go-template/v2-domain-adapters-and-helpers/domain"
)

// HandleError is responsible for converting domain errors to HTTP errors
// simplifying error handling overall.
func HandleError(logger log.Provider) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		err := c.Next()
		if err == nil {
			return nil
		}

		req := c.Request()
		status, body := handleDomainErrAsHTTP(
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

func handleDomainErrAsHTTP(ctx context.Context, logger log.Provider, err error, method string, path string) (status int, responseBody []byte) {
	domainErr := domain.AsDomainErr(err)

	response := map[string]interface{}{
		"code":       domainErr.Code,
		"title":      domainErr.Title,
		"request_id": domain.GetRequestIDFromContext(ctx),
	}

	switch domainErr.Code {
	case "InternalErr":
		status = 500

		data := log.Body{
			"route": method + ": " + path,
		}
		for k, v := range domainErr.Data {
			data[k] = v
		}
		logger.Error(ctx, "request-error", data)

	case "BadRequest":
		status = 400
		for k, v := range domainErr.Data {
			response[k] = v
		}

	case "NotFoundErr":
		status = 404
		for k, v := range domainErr.Data {
			response[k] = v
		}
	}

	responseBody, _ = json.Marshal(response)
	return status, responseBody
}
