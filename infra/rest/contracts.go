package rest

import (
	"context"
	"io"
	"net/http"
	"time"
)

// Provider provides the functions to perform
// REST requests automatically marshalling the input body as JSON.
//
// It returns error only if it was not possible to complete the request
// either because of a marshal error or a network error.
//
// Otherwise the statusCode should be used to check if the request
// was processed successfully.
type Provider interface {
	Get(url string, data RequestData) (resp Response, err error)
	Post(url string, data RequestData) (resp Response, err error)
	Put(url string, data RequestData) (resp Response, err error)
	Patch(url string, data RequestData) (resp Response, err error)
	Delete(url string, data RequestData) (resp Response, err error)
}

// RequestData describes the optional arguments for all
// the http methods of this client.
type RequestData struct {
	// The body accepts any struct that can
	// be marshaled into JSON
	Body interface{}

	Headers map[string]string

	// Optional
	Context context.Context

	// Set this option to true if you
	// expect to receive big bodies of data
	// and you don't want this library to
	// to load all of it in memory.
	//
	// When using this option the resp.Body
	// field will be set to null and you'll
	// need to use the response struct as an io.ReadCloser
	// for streaming the data wherever you need to.
	//
	// Don't forget to close it afterwards.
	//
	// Also note that there is no need to call resp.Close()
	// if you are not using the Stream option or if the call
	// returns an error.
	Stream bool

	// It's the max number of retries, if 0 it defaults 1
	MaxRetries int

	// The start and max delay for the exponential backoff strategy
	// if unset they default to 300ms and 32s respectively
	BaseRetryDelay time.Duration
	MaxRetryDelay  time.Duration

	// Set this attribute if you want to personalize the retry behavior
	// if nil it defaults to `rest.DefaultRetryRule()`
	RetryRule func(resp *http.Response, err error) bool
}

// SetDefaultsIfNecessary sets the default values
// for the RequestData structure
func (r *RequestData) SetDefaultsIfNecessary() {
	if r.MaxRetries == 0 {
		r.MaxRetries = 1
	}
	if r.BaseRetryDelay == 0 {
		r.BaseRetryDelay = 300 * time.Millisecond
	}
	if r.MaxRetryDelay == 0 {
		r.MaxRetryDelay = 32 * time.Second
	}
	if r.RetryRule == nil {
		r.RetryRule = DefaultRetryRule
	}
}

// Response describes the expected attributes
// on the response for a REST request
type Response struct {
	io.ReadCloser

	Body       []byte
	Header     map[string]string
	StatusCode int
}

func DefaultRetryRule(resp *http.Response, err error) bool {
	retriableStatus := map[int]bool{
		http.StatusLocked:          true,
		http.StatusTooEarly:        true,
		http.StatusTooManyRequests: true,
	}
	return err != nil || retriableStatus[resp.StatusCode] || resp.StatusCode > 500
}
