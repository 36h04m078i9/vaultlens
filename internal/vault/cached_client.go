package vault

import (
	"context"
	"fmt"
)

// CachedClient wraps a Client with a transparent read-through Cache.
type CachedClient struct {
	base  *Client
	cache *Cache
}

// NewCachedClient creates a CachedClient that stores results in the given Cache.
func NewCachedClient(base *Client, cache *Cache) *CachedClient {
	return &CachedClient{base: base, cache: cache}
}

// ListSecrets returns the cached key list for path, fetching from Vault on a miss.
func (c *CachedClient) ListSecrets(ctx context.Context, path string) ([]string, error) {
	key := fmt.Sprintf("list:%s", path)

	if raw, ok := c.cache.Get(key); ok {
		if keys, ok := raw.([]string); ok {
			return keys, nil
		}
	}

	keys, err := c.base.ListSecrets(ctx, path)
	if err != nil {
		return nil, err
	}

	c.cache.Set(key, keys)
	return keys, nil
}

// ReadSecret returns the cached secret at path, fetching from Vault on a miss.
func (c *CachedClient) ReadSecret(ctx context.Context, path string) (map[string]interface{}, error) {
	key := fmt.Sprintf("read:%s", path)

	if raw, ok := c.cache.Get(key); ok {
		if secret, ok := raw.(map[string]interface{}); ok {
			return secret, nil
		}
	}

	secret, err := c.base.ReadSecret(ctx, path)
	if err != nil {
		return nil, err
	}

	c.cache.Set(key, secret)
	return secret, nil
}

// Invalidate removes a single cache entry by its exact cache key.
// Use the "list:<path>" or "read:<path>" prefix conventions.
func (c *CachedClient) Invalidate(key string) {
	c.cache.Invalidate(key)
}

// Flush clears the entire underlying cache.
func (c *CachedClient) Flush() {
	c.cache.Flush()
}
