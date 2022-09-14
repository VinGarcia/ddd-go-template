package domain

import (
	"fmt"
	"strings"
)

type DomainErr struct {
	Code  string
	Title string
	Data  map[string]interface{}
}

func (e DomainErr) Error() string {
	fields := []string{
		e.Code + ": " + e.Title,
	}
	for k, v := range e.Data {
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
		Code:  "InternalErr",
		Title: err.Error(),
	}
}

func InternalErr(title string, data map[string]interface{}) DomainErr {
	return DomainErr{
		Code:  "InternalErr",
		Title: title,
		Data:  data,
	}
}

func BadRequestErr(title string, data map[string]interface{}) DomainErr {
	return DomainErr{
		Code:  "BadRequestErr",
		Title: title,
		Data:  data,
	}
}

func UnauthorizedErr(title string, data map[string]interface{}) DomainErr {
	return DomainErr{
		Code:  "UnauthorizedErr",
		Title: title,
		Data:  data,
	}
}

func NotFoundErr(title string, data map[string]interface{}) DomainErr {
	return DomainErr{
		Code:  "NotFoundErr",
		Title: title,
		Data:  data,
	}
}
