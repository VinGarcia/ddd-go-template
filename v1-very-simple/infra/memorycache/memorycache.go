package memorycache

import (
	"context"
	"encoding/json"
	"time"

	cache "github.com/patrickmn/go-cache"
	"github.com/vingarcia/ddd-go-template/v1-very-simple/domain"
)

type Client struct {
	cache *cache.Cache
}

func NewClient(defaultExpiration time.Duration, cleanupInterval time.Duration) Client {
	return Client{
		cache: cache.New(defaultExpiration, cleanupInterval),
	}
}

func (c Client) Get(ctx context.Context, key string, record interface{}) error {
	value, _ := c.cache.Get(key)
	rawJSON, ok := value.([]byte)
	if !ok {
		return domain.NotFoundErr("record-not-found", map[string]interface{}{
			"func":      "memorycache.Client.Get",
			"input_key": key,
		})
	}
	return json.Unmarshal(rawJSON, record)
}

func (c Client) Set(ctx context.Context, key string, record interface{}) error {
	rawJSON, err := json.Marshal(record)
	if err != nil {
		return domain.InternalErr("unable-to-marshal-record-as-json", map[string]interface{}{
			"func":         "memorycache.Client.Set",
			"error":        err.Error(),
			"input_record": record,
		})
	}

	c.cache.Set(key, rawJSON, 0)
	return nil
}
