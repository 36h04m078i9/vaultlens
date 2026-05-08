package snapshot

import (
	"testing"
	"time"
)

func makeListStore(t *testing.T) *Store {
	t.Helper()
	s := NewStore()

	now := time.Now()
	snaps := []Snapshot{
		{ID: "a", Path: "secret/prod/db", CapturedAt: now.Add(-3 * time.Hour), Secrets: map[string]string{"k": "v"}},
		{ID: "b", Path: "secret/prod/api", CapturedAt: now.Add(-2 * time.Hour), Secrets: map[string]string{"x": "1", "y": "2"}},
		{ID: "c", Path: "secret/staging/db", CapturedAt: now.Add(-1 * time.Hour), Secrets: map[string]string{}},
		{ID: "d", Path: "secret/prod/cache", CapturedAt: now, Secrets: map[string]string{"a": "b", "c": "d", "e": "f"}},
	}
	for _, snap := range snaps {
		if err := s.Save(snap); err != nil {
			t.Fatalf("Save: %v", err)
		}
	}
	return s
}

func TestListNoOptions(t *testing.T) {
	s := makeListStore(t)
	results := List(s, ListOptions{})
	if len(results) != 4 {
		t.Fatalf("expected 4 results, got %d", len(results))
	}
}

func TestListSortedNewestFirst(t *testing.T) {
	s := makeListStore(t)
	results := List(s, ListOptions{})
	for i := 1; i < len(results); i++ {
		if results[i].CreatedAt.After(results[i-1].CreatedAt) {
			t.Errorf("results not sorted descending at index %d", i)
		}
	}
}

func TestListFilterByPrefix(t *testing.T) {
	s := makeListStore(t)
	results := List(s, ListOptions{Prefix: "secret/prod"})
	if len(results) != 3 {
		t.Fatalf("expected 3 prod results, got %d", len(results))
	}
	for _, r := range results {
		if len(r.Path) < len("secret/prod") || r.Path[:len("secret/prod")] != "secret/prod" {
			t.Errorf("unexpected path %q in prod-filtered results", r.Path)
		}
	}
}

func TestListFilterBySince(t *testing.T) {
	s := makeListStore(t)
	cutoff := time.Now().Add(-90 * time.Minute)
	results := List(s, ListOptions{Since: cutoff})
	if len(results) != 2 {
		t.Fatalf("expected 2 recent results, got %d", len(results))
	}
}

func TestListLimit(t *testing.T) {
	s := makeListStore(t)
	results := List(s, ListOptions{Limit: 2})
	if len(results) != 2 {
		t.Fatalf("expected 2 results with limit, got %d", len(results))
	}
}

func TestListKeyCount(t *testing.T) {
	s := makeListStore(t)
	results := List(s, ListOptions{Prefix: "secret/prod/cache"})
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].KeyCount != 3 {
		t.Errorf("expected KeyCount=3, got %d", results[0].KeyCount)
	}
}

func TestListNilStoreReturnsNil(t *testing.T) {
	results := List(nil, ListOptions{})
	if results != nil {
		t.Errorf("expected nil for nil store, got %v", results)
	}
}
