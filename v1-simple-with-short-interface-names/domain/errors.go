package domain

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/vingarcia/ddd-go-template/v1-simple-with-short-interface-names/infra/log"
)

type DomainErr struct {
	code  string
	title string
	data  map[string]interface{}
}

func (e DomainErr) Error() string {
	fields := []string{
		e.code + ": " + e.title,
	}
	for k, v := range e.data {
		fields = append(fields, fmt.Sprintf("%s = %#v", k, v))
	}

	return strings.Join(fields, "; ")
}

func AsDomainErr(err error) DomainErr {
	domainErr, ok := err.(DomainErr)
	if ok {
		return domainErr
	}
	return DomainErr{
		code:  "InternalErr",
		title: err.Error(),
	}
}

func InternalErr(title string, data map[string]interface{}) DomainErr {
	return DomainErr{
		code:  "InternalErr",
		title: title,
		data:  data,
	}
}

func BadRequestErr(title string, data map[string]interface{}) DomainErr {
	return DomainErr{
		code:  "BadRequestErr",
		title: title,
		data:  data,
	}
}

func UnauthorizedErr(title string, data map[string]interface{}) DomainErr {
	return DomainErr{
		code:  "UnauthorizedErr",
		title: title,
		data:  data,
	}
}

func NotFoundErr(title string, data map[string]interface{}) DomainErr {
	return DomainErr{
		code:  "NotFoundErr",
		title: title,
		data:  data,
	}
}

func HandleDomainErrAsHTTP(ctx context.Context, logger log.Provider, err error, method string, path string) (status int, responseBody []byte) {
	domainErr := AsDomainErr(err)

	response := map[string]interface{}{
		"code":       domainErr.code,
		"title":      domainErr.title,
		"request_id": GetRequestIDFromContext(ctx),
	}

	switch domainErr.code {
	case "InternalErr":
		status = 500

		data := log.Body{
			"route": method + ": " + path,
		}
		for k, v := range domainErr.data {
			data[k] = v
		}
		logger.Error(ctx, "request-error", data)

	case "BadRequest":
		status = 400
		for k, v := range domainErr.data {
			response[k] = v
		}

	case "NotFoundErr":
		status = 404
		for k, v := range domainErr.data {
			response[k] = v
		}
	}

	responseBody, _ = json.Marshal(response)
	return status, responseBody
}
