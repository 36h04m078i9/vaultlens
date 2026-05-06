package vault

import (
	"fmt"
	"time"
)

// CachedClient wraps a Client and transparently caches list and read results.
type CachedClient struct {
	client *Client
	cache  *Cache
}

// NewCachedClient creates a CachedClient with the given underlying client and TTL.
func NewCachedClient(client *Client, ttl time.Duration) *CachedClient {
	return &CachedClient{
		client: client,
		cache:  NewCache(ttl),
	}
}

// ListSecrets returns cached secret keys for path, or fetches and caches them.
func (cc *CachedClient) ListSecrets(mountPath, secretPath string) ([]string, error) {
	key := fmt.Sprintf("list:%s/%s", mountPath, secretPath)

	if cached, ok := cc.cache.Get(key); ok {
		if keys, ok := cached.([]string); ok {
			return keys, nil
		}
	}

	keys, err := cc.client.ListSecrets(mountPath, secretPath)
	if err != nil {
		return nil, err
	}
	cc.cache.Set(key, keys)
	return keys, nil
}

// ReadSecret returns a cached secret map for path, or fetches and caches it.
func (cc *CachedClient) ReadSecret(mountPath, secretPath string) (map[string]interface{}, error) {
	key := fmt.Sprintf("read:%s/%s", mountPath, secretPath)

	if cached, ok := cc.cache.Get(key); ok {
		if data, ok := cached.(map[string]interface{}); ok {
			return data, nil
		}
	}

	data, err := cc.client.ReadSecret(mountPath, secretPath)
	if err != nil {
		return nil, err
	}
	cc.cache.Set(key, data)
	return data, nil
}

// InvalidatePath removes cached entries for a specific mount+path combination.
func (cc *CachedClient) InvalidatePath(mountPath, secretPath string) {
	cc.cache.Invalidate(fmt.Sprintf("list:%s/%s", mountPath, secretPath))
	cc.cache.Invalidate(fmt.Sprintf("read:%s/%s", mountPath, secretPath))
}

// FlushCache clears all cached entries.
func (cc *CachedClient) FlushCache() {
	cc.cache.Flush()
}
