package snapshot

import (
	"testing"
	"time"
)

func makeSearchStore(t *testing.T) *Store {
	t.Helper()
	store := NewStore()
	now := time.Now()

	snaps := []Snapshot{
		{ID: "snap-alpha", Prefix: "secret/prod", Tags: []string{"production", "db"}, TakenAt: now.Add(-3 * time.Hour), Secrets: map[string]map[string]interface{}{"secret/prod/db": {"pass": "x"}}},
		{ID: "snap-beta", Prefix: "secret/staging", Tags: []string{"staging"}, TakenAt: now.Add(-2 * time.Hour), Secrets: map[string]map[string]interface{}{"secret/staging/api": {"key": "y"}}},
		{ID: "snap-gamma", Prefix: "secret/prod", Tags: []string{"production", "cache"}, TakenAt: now.Add(-1 * time.Hour), Secrets: map[string]map[string]interface{}{"secret/prod/cache": {"url": "z"}}},
	}

	for i := range snaps {
		if err := store.Save(&snaps[i]); err != nil {
			t.Fatalf("seed save: %v", err)
		}
	}
	return store
}

func TestSearchNilStoreReturnsError(t *testing.T) {
	_, err := Search(nil, SearchOptions{})
	if err == nil {
		t.Fatal("expected error for nil store")
	}
}

func TestSearchEmptyQueryReturnsAll(t *testing.T) {
	store := makeSearchStore(t)
	results, err := Search(store, SearchOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
}

func TestSearchMatchesID(t *testing.T) {
	store := makeSearchStore(t)
	results, err := Search(store, SearchOptions{Query: "alpha"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Snapshot.ID != "snap-alpha" {
		t.Errorf("expected snap-alpha, got %s", results[0].Snapshot.ID)
	}
}

func TestSearchMatchesPrefix(t *testing.T) {
	store := makeSearchStore(t)
	results, err := Search(store, SearchOptions{Query: "staging"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
}

func TestSearchMatchesTag(t *testing.T) {
	store := makeSearchStore(t)
	results, err := Search(store, SearchOptions{Query: "cache"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Snapshot.ID != "snap-gamma" {
		t.Errorf("expected snap-gamma, got %s", results[0].Snapshot.ID)
	}
}

func TestSearchFilterByTags(t *testing.T) {
	store := makeSearchStore(t)
	results, err := Search(store, SearchOptions{Tags: []string{"production"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 production results, got %d", len(results))
	}
}

func TestSearchFilterBySince(t *testing.T) {
	store := makeSearchStore(t)
	cutoff := time.Now().Add(-90 * time.Minute)
	results, err := Search(store, SearchOptions{Since: cutoff})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result since cutoff, got %d", len(results))
	}
	if results[0].Snapshot.ID != "snap-gamma" {
		t.Errorf("expected snap-gamma, got %s", results[0].Snapshot.ID)
	}
}

func TestSearchLimitCapsResults(t *testing.T) {
	store := makeSearchStore(t)
	results, err := Search(store, SearchOptions{Limit: 2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results with limit, got %d", len(results))
	}
}

func TestSearchNoMatchReturnsEmpty(t *testing.T) {
	store := makeSearchStore(t)
	results, err := Search(store, SearchOptions{Query: "zzznomatch"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}
