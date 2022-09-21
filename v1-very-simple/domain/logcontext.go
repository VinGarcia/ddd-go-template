package domain

// This file is part of the domain language, and it allows different packages
// to save values to the context in a way that our log provider implementation
// is able to read and log later on.
//
// In the domain package modules that have function should be kept to a minimum,
// and should only be put here if they really need to be used by several different
// packages.

import (
	"context"

	"github.com/vingarcia/ddd-go-template/v1-very-simple/infra/maps"
)

// Declaring a unique private type for the ctx key
// guarantees that no key colision will ever happen:
type logCtxKey struct{}

// CtxWithValues merges received values with log body currently stored
// on the input ctx.
func CtxWithValues(ctx context.Context, values LogBody) context.Context {
	m, _ := ctx.Value(logCtxKey{}).(LogBody)
	return context.WithValue(ctx, logCtxKey{}, mergeMaps(m, values))
}

// GetCtxValues extracts the log body currently stored on the input ctx.
func GetCtxValues(ctx context.Context) LogBody {
	m, _ := ctx.Value(logCtxKey{}).(LogBody)
	if m == nil {
		m = LogBody{}
	}
	m["request_id"] = GetRequestIDFromContext(ctx)
	return m
}

func mergeMaps(bodies ...LogBody) LogBody {
	body := LogBody{}
	maps.Merge(&body, bodies...)
	return body
}
