package domain

import "context"

// CacheProvider implements a simple type-safe cache
//
// Usage example:
//
// err := cache.Set("some_key", Foo{Name: "example object"})
// if err != nil {
//   return err
// }
//
// var result Foo
// err := cache.Get("some_key", &result)
// if err != nil {
//   return err
// }
type CacheProvider interface {
	Get(ctx context.Context, key string, record interface{}) error
	Set(ctx context.Context, key string, record interface{}) error
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
type LogBody = map[string]interface{}

// Provider provides the functions to perform
// REST requests automatically marshalling the input body as JSON.
//
// It returns error only if it was not possible to complete the request
// either because of a marshal error or a network error.
//
// Otherwise the statusCode should be used to check if the request
// was processed successfully.
type RestProvider interface {
	Get(ctx context.Context, url string, data RequestData) (resp Response, err error)
	Post(ctx context.Context, url string, data RequestData) (resp Response, err error)
	Put(ctx context.Context, url string, data RequestData) (resp Response, err error)
	Patch(ctx context.Context, url string, data RequestData) (resp Response, err error)
	Delete(ctx context.Context, url string, data RequestData) (resp Response, err error)
}

// RequestData describes the optional arguments for all
// the http methods of this client.
type RequestData struct {
	// The body accepts any struct that can
	// be marshaled into JSON
	Body interface{}

	Headers map[string]string
}

// Response describes the expected attributes
// on the response for a REST request
type Response struct {
	Body       []byte
	Header     map[string]string
	StatusCode int
}
