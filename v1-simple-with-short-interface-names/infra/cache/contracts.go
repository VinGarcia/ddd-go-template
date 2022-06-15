package cache

import "context"

// Provider implements a simple type-safe cache
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
type Provider interface {
	Get(ctx context.Context, key string, record interface{}) error
	Set(ctx context.Context, key string, record interface{}) error
}
