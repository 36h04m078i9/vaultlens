package history

import (
	"testing"
	"time"
)

func makeHistory(paths []string) *History {
	h := New(50)
	// Push in reverse so index 0 is the most recent.
	for i := len(paths) - 1; i >= 0; i-- {
		h.Push(Entry{Path: paths[i], AccessedAt: time.Now()})
	}
	return h
}

func TestSearchEmptyQueryReturnsAll(t *testing.T) {
	paths := []string{"secret/a", "secret/b", "secret/c"}
	h := makeHistory(paths)

	results := Search(h, "")
	if len(results) != len(paths) {
		t.Fatalf("expected %d results, got %d", len(paths), len(results))
	}
}

func TestSearchWhitespaceQueryReturnsAll(t *testing.T) {
	paths := []string{"secret/x", "secret/y"}
	h := makeHistory(paths)

	results := Search(h, "   ")
	if len(results) != len(paths) {
		t.Fatalf("expected %d results, got %d", len(paths), len(results))
	}
}

func TestSearchMatchesSubstring(t *testing.T) {
	paths := []string{"secret/alpha", "secret/beta", "secret/alphabet"}
	h := makeHistory(paths)

	results := Search(h, "alpha")
	for _, r := range results {
		if r.Entry.Path == "secret/beta" {
			t.Errorf("unexpected match for 'secret/beta' with query 'alpha'")
		}
	}
	if len(results) < 2 {
		t.Errorf("expected at least 2 matches, got %d", len(results))
	}
}

func TestSearchNoMatchReturnsEmpty(t *testing.T) {
	paths := []string{"secret/foo", "secret/bar"}
	h := makeHistory(paths)

	results := Search(h, "zzzzz")
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestSearchResultsOrderedByScore(t *testing.T) {
	// "db" is a closer match to "database/password" than "secret/dba-creds"
	// but we just verify ordering is descending.
	paths := []string{"secret/dba-creds", "database/password", "kv/unrelated"}
	h := makeHistory(paths)

	results := Search(h, "db")
	for i := 1; i < len(results); i++ {
		if results[i].Score > results[i-1].Score {
			t.Errorf("results not sorted: index %d score %d > index %d score %d",
				i, results[i].Score, i-1, results[i-1].Score)
		}
	}
}

func TestSearchOnEmptyHistoryReturnsEmpty(t *testing.T) {
	h := New(10)
	results := Search(h, "anything")
	if len(results) != 0 {
		t.Errorf("expected 0 results from empty history, got %d", len(results))
	}
}
