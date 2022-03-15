package redis

import (
	"context"
	"encoding/json"
	"time"

	redis "github.com/go-redis/redis/v8"
	"github.com/vingarcia/ddd-go-template/foolproof/domain"
)

type Client struct {
	redis             *redis.Client
	defaultExpiration time.Duration
}

func NewClient(
	connectionURL string,
	password string,
	defaultExpiration time.Duration,
) Client {
	return Client{
		redis: redis.NewClient(&redis.Options{
			Addr:     connectionURL, // e.g. localhost:6379
			Password: password,      // set it to empty for no password
		}),
		defaultExpiration: defaultExpiration, // set to 0 for no expiration
	}
}

// Get implements the domain.CacheProvider interface
func (c Client) Get(ctx context.Context, key string, record interface{}) error {
	rawJSONStr, err := c.redis.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return domain.NotFoundErr("record-not-found-on-redis", map[string]interface{}{
				"func":      "redis.Client.Get",
				"input_key": key,
			})
		}
		return domain.InternalErr("error-fetching-record-from-redis", map[string]interface{}{
			"func":      "redis.Client.Get",
			"error":     err.Error(),
			"input_key": key,
		})
	}

	err = json.Unmarshal([]byte(rawJSONStr), &record)
	if err != nil {
		return domain.InternalErr("error-decoding-record-from-redis-as-json", map[string]interface{}{
			"func":        "redis.Client.Get",
			"error":       err.Error(),
			"input_key":   key,
			"record_json": rawJSONStr,
		})
	}

	return nil
}

// Set implements the domain.CacheProvider interface
func (c Client) Set(ctx context.Context, key string, record interface{}) error {
	rawJSON, err := json.Marshal(record)
	if err != nil {
		return domain.InternalErr("error-marshalling-record", map[string]interface{}{
			"func":         "redis.Client.Get",
			"error":        err.Error(),
			"input_key":    key,
			"input_record": record,
		})
	}

	err = c.redis.Set(ctx, key, string(rawJSON), c.defaultExpiration).Err()
	if err != nil {
		return domain.InternalErr("error-saving-record-on-redis", map[string]interface{}{
			"func":         "redis.Client.Get",
			"error":        err.Error(),
			"input_key":    key,
			"input_record": record,
		})
	}
	return nil
}
