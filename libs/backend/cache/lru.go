package cache

import (
	"context"
	"libs/backend/boot"
	"log/slog"

	lru "github.com/hashicorp/golang-lru"
)

// LRUCache is a shared store for APQ and query AST caching
// using LRU cache.
// Implements the graphql.Cache interface
type LRUCache[T any] struct {
	cache  *lru.Cache
	logger boot.Logger
}

// NewLRUCache creates a new cache with the given size.
func NewLRUCache[T any](size int, logger boot.Logger) (LRUCache[T], error) {
	cache, err := lru.New(size)
	if err != nil {
		return LRUCache[T]{}, err
	}

	return LRUCache[T]{cache: cache, logger: logger}, nil
}

// Get looks up a key's value from the cache and asserts the type
func (c LRUCache[T]) Get(ctx context.Context, key string) (value T, ok bool) {
	// Get looks up a key's value from the cache.
	item, ok := c.cache.Get(key)
	if !ok {
		return value, false
	}

	// Type assertion that value is type T
	properValue, ok := item.(T)
	if !ok {
		c.logger.Warn("Type assertion failed for item in LRU cache", slog.Any("item", item))
		return value, false
	}

	return properValue, true
}

// Add adds a value to the cache.
func (c LRUCache[T]) Add(ctx context.Context, key string, value T) {
	c.cache.Add(key, value)
}
