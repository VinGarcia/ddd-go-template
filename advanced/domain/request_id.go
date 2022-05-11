package domain

import (
	"context"

	"github.com/google/uuid"
)

// RequestIDKey needs to be a string, since its used on Fiber `.Locals()`
const RequestIDKey string = "request_id_key"

func GenerateRequestID() string {
	return uuid.NewString()
}

func GetRequestIDFromContext(ctx context.Context) string {
	requestID, _ := ctx.Value(RequestIDKey).(string)
	return requestID
}
