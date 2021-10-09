package domain

import (
	"context"

	"github.com/google/uuid"
)

const requestIDKey = "request_id_key"

func GenerateRequestID() (requestIDKey string, id string) {
	return "request_id_key", uuid.NewString()
}

func GetRequestIDFromContext(ctx context.Context) string {
	requestID, _ := ctx.Value("request_id_key").(string)
	return requestID
}
