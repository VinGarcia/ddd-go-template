package memorycache

import (
	"time"

	cache "github.com/patrickmn/go-cache"
)

type Client struct {
	cache *cache.Cache
}

func New(defaultExpiration time.Duration, cleanupInterval time.Duration) Client {
	return Client{
		cache: cache.New(defaultExpiration, cleanupInterval),
	}
}

func (c Client) Get(key string) interface{} {
	value, found := c.cache.Get(key)
	if !found {
		return nil
	}
	return value
}

func (c Client) Set(key string, value interface{}) {
	c.cache.Set(key, value, 0)
}
