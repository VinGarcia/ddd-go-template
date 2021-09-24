package domain

import "context"

type CacheProvider interface {
	Get(key string) interface{}
	Set(key string, value interface{})
}

type LogProvider interface {
	Debug(ctx context.Context, title string, valueMaps ...LogBody)
	Info(ctx context.Context, title string, valueMaps ...LogBody)
	Warn(ctx context.Context, title string, valueMaps ...LogBody)
	Error(ctx context.Context, title string, valueMaps ...LogBody)
	Fatal(ctx context.Context, title string, valueMaps ...LogBody)
}

// Body is the log body containing the keys and values
// used to build the structured logs
type LogBody map[string]interface{}
