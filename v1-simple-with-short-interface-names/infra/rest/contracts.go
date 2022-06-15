package rest

import "context"

// Provider provides the functions to perform
// REST requests automatically marshalling the input body as JSON.
//
// It returns error only if it was not possible to complete the request
// either because of a marshal error or a network error.
//
// Otherwise the statusCode should be used to check if the request
// was processed successfully.
type Provider interface {
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
