package vault_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCachedClientListSecretsHitsCache(t *testing.T) {
	var callCount atomic.Int32
	srv := newMockVaultServer(t, func(path string) []string {
		callCount.Add(1)
		return []string{"alpha", "beta"}
	}, nil)

	base := mustNewClient(t, srv.URL)
	cache := NewCache(5 * time.Minute)
	client := NewCachedClient(base, cache)

	ctx := context.Background()

	keys1, err := client.ListSecrets(ctx, "secret/")
	require.NoError(t, err)
	assert.Equal(t, []string{"alpha", "beta"}, keys1)

	keys2, err := client.ListSecrets(ctx, "secret/")
	require.NoError(t, err)
	assert.Equal(t, keys1, keys2)

	assert.Equal(t, int32(1), callCount.Load(), "upstream should be called only once")
}

func TestCachedClientReadSecretHitsCache(t *testing.T) {
	var callCount atomic.Int32
	srv := newMockVaultServer(t, nil, func(path string) map[string]interface{} {
		callCount.Add(1)
		return map[string]interface{}{"key": "value"}
	})

	base := mustNewClient(t, srv.URL)
	cache := NewCache(5 * time.Minute)
	client := NewCachedClient(base, cache)

	ctx := context.Background()

	secret1, err := client.ReadSecret(ctx, "secret/data/foo")
	require.NoError(t, err)
	assert.Equal(t, map[string]interface{}{"key": "value"}, secret1)

	secret2, err := client.ReadSecret(ctx, "secret/data/foo")
	require.NoError(t, err)
	assert.Equal(t, secret1, secret2)

	assert.Equal(t, int32(1), callCount.Load(), "upstream should be called only once")
}

func TestCachedClientInvalidateFlushesEntry(t *testing.T) {
	var callCount atomic.Int32
	srv := newMockVaultServer(t, func(path string) []string {
		callCount.Add(1)
		return []string{"x"}
	}, nil)

	base := mustNewClient(t, srv.URL)
	cache := NewCache(5 * time.Minute)
	client := NewCachedClient(base, cache)

	ctx := context.Background()

	_, err := client.ListSecrets(ctx, "secret/")
	require.NoError(t, err)

	client.Invalidate("list:secret/")

	_, err = client.ListSecrets(ctx, "secret/")
	require.NoError(t, err)

	assert.Equal(t, int32(2), callCount.Load(), "upstream should be called again after invalidation")
}
