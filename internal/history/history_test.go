package history_test

import (
	"testing"
	"time"

	"github.com/vaultlens/vaultlens/internal/history"
)

func TestPushAddsEntry(t *testing.T) {
	h := history.New(10)
	h.Push("secret/foo")
	if h.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", h.Len())
	}
}

func TestPushDeduplicatesAndMovesToFront(t *testing.T) {
	h := history.New(10)
	h.Push("secret/foo")
	h.Push("secret/bar")
	h.Push("secret/foo") // re-push foo

	entries := h.List()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries after dedup, got %d", len(entries))
	}
	if entries[0].Path != "secret/foo" {
		t.Errorf("expected foo at front, got %s", entries[0].Path)
	}
}

func TestPushRespectsMaxSize(t *testing.T) {
	h := history.New(3)
	h.Push("a")
	h.Push("b")
	h.Push("c")
	h.Push("d")

	if h.Len() != 3 {
		t.Fatalf("expected 3 entries, got %d", h.Len())
	}
	entries := h.List()
	if entries[0].Path != "d" {
		t.Errorf("expected newest entry first, got %s", entries[0].Path)
	}
}

func TestListReturnsCopy(t *testing.T) {
	h := history.New(10)
	h.Push("secret/x")

	a := h.List()
	a[0].Path = "mutated"

	b := h.List()
	if b[0].Path == "mutated" {
		t.Error("List should return a copy, not a reference")
	}
}

func TestClearRemovesAllEntries(t *testing.T) {
	h := history.New(10)
	h.Push("secret/a")
	h.Push("secret/b")
	h.Clear()

	if h.Len() != 0 {
		t.Errorf("expected 0 entries after Clear, got %d", h.Len())
	}
}

func TestPushSetsAccessedAt(t *testing.T) {
	before := time.Now()
	h := history.New(10)
	h.Push("secret/ts")
	after := time.Now()

	entries := h.List()
	if entries[0].AccessedAt.Before(before) || entries[0].AccessedAt.After(after) {
		t.Error("AccessedAt timestamp is outside expected range")
	}
}

func TestNewDefaultsMaxSizeOnZero(t *testing.T) {
	h := history.New(0)
	for i := 0; i < 60; i++ {
		h.Push(string(rune('a' + i%26)))
	}
	if h.Len() > 50 {
		t.Errorf("expected default max 50, got %d", h.Len())
	}
}
