package domain

import (
	"fmt"
	"strings"
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

func ErrCode(err error) string {
	domainErr, ok := err.(DomainErr)
	if !ok {
		return "InternalErr"
	}
	return domainErr.code
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

func NotFoundErr(title string, data map[string]interface{}) DomainErr {
	return DomainErr{
		code:  "NotFoundErr",
		title: title,
		data:  data,
	}
}
