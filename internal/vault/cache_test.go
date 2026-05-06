package vault

import (
	"testing"
	"time"
)

func TestCacheSetAndGet(t *testing.T) {
	c := NewCache(5 * time.Second)
	c.Set("secret/foo", "bar")

	val, ok := c.Get("secret/foo")
	if !ok {
		t.Fatal("expected cache hit, got miss")
	}
	if val != "bar" {
		t.Fatalf("expected 'bar', got %v", val)
	}
}

func TestCacheMissOnUnknownKey(t *testing.T) {
	c := NewCache(5 * time.Second)

	_, ok := c.Get("nonexistent")
	if ok {
		t.Fatal("expected cache miss, got hit")
	}
}

func TestCacheExpiry(t *testing.T) {
	c := NewCache(50 * time.Millisecond)
	c.Set("key", "value")

	time.Sleep(100 * time.Millisecond)

	_, ok := c.Get("key")
	if ok {
		t.Fatal("expected cache miss after TTL expiry, got hit")
	}
}

func TestCacheInvalidate(t *testing.T) {
	c := NewCache(5 * time.Second)
	c.Set("key", "value")
	c.Invalidate("key")

	_, ok := c.Get("key")
	if ok {
		t.Fatal("expected cache miss after invalidation")
	}
}

func TestCacheFlush(t *testing.T) {
	c := NewCache(5 * time.Second)
	c.Set("a", 1)
	c.Set("b", 2)
	c.Set("c", 3)

	c.Flush()

	if c.Size() != 0 {
		t.Fatalf("expected size 0 after flush, got %d", c.Size())
	}
}

func TestCacheSize(t *testing.T) {
	c := NewCache(5 * time.Second)
	if c.Size() != 0 {
		t.Fatalf("expected initial size 0, got %d", c.Size())
	}
	c.Set("x", true)
	c.Set("y", true)
	if c.Size() != 2 {
		t.Fatalf("expected size 2, got %d", c.Size())
	}
}

func TestCacheOverwrite(t *testing.T) {
	c := NewCache(5 * time.Second)
	c.Set("key", "first")
	c.Set("key", "second")

	val, ok := c.Get("key")
	if !ok {
		t.Fatal("expected cache hit")
	}
	if val != "second" {
		t.Fatalf("expected 'second', got %v", val)
	}
}
