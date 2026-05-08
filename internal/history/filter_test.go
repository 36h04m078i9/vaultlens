package history

import (
	"testing"
	"time"
)

func makeFilterHistory(t *testing.T) *History {
	t.Helper()
	h := New(20)
	now := time.Now()
	// Push in reverse order so most-recent ends up first after dedup.
	entries := []struct {
		path string
		offset time.Duration
	}{
		{"secret/prod/db", -1 * time.Minute},
		{"secret/prod/api", -5 * time.Minute},
		{"secret/staging/db", -10 * time.Minute},
		{"kv/old/token", -30 * time.Minute},
	}
	for _, e := range entries {
		h.Push(e.path)
		// Manually backdate the first entry to simulate real timestamps.
		h.entries[0].VisitedAt = now.Add(e.offset)
	}
	return h
}

func TestFilterNoOptions(t *testing.T) {
	h := makeFilterHistory(t)
	got := Filter(h, FilterOptions{})
	if len(got) != 4 {
		t.Fatalf("expected 4 entries, got %d", len(got))
	}
}

func TestFilterByPrefix(t *testing.T) {
	h := makeFilterHistory(t)
	got := Filter(h, FilterOptions{Prefix: "secret/prod"})
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
	for _, e := range got {
		if e.Path[:11] != "secret/prod" {
			t.Errorf("unexpected path %q", e.Path)
		}
	}
}

func TestFilterBySince(t *testing.T) {
	h := makeFilterHistory(t)
	cutoff := time.Now().Add(-6 * time.Minute)
	got := Filter(h, FilterOptions{Since: cutoff})
	// Only entries within the last 6 minutes should appear.
	for _, e := range got {
		if e.VisitedAt.Before(cutoff) {
			t.Errorf("entry %q is older than cutoff", e.Path)
		}
	}
}

func TestFilterLimit(t *testing.T) {
	h := makeFilterHistory(t)
	got := Filter(h, FilterOptions{Limit: 2})
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
}

func TestFilterCombinedPrefixAndLimit(t *testing.T) {
	h := makeFilterHistory(t)
	got := Filter(h, FilterOptions{Prefix: "secret/", Limit: 1})
	if len(got) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(got))
	}
	if got[0].Path[:7] != "secret/" {
		t.Errorf("unexpected path %q", got[0].Path)
	}
}

func TestFilterEmptyHistoryReturnsEmpty(t *testing.T) {
	h := New(10)
	got := Filter(h, FilterOptions{})
	if len(got) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(got))
	}
}
