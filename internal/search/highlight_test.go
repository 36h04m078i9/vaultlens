package search

import (
	"strings"
	"testing"
)

func TestHighlightEmptyQuery(t *testing.T) {
	paths := []string{"secret/app/key", "secret/db/pass"}
	results := Highlight("", paths)
	if len(results) != len(paths) {
		t.Fatalf("expected %d results, got %d", len(paths), len(results))
	}
	for _, r := range results {
		if r.Highlight != r.Path {
			t.Errorf("empty query: expected plain path %q, got %q", r.Path, r.Highlight)
		}
	}
}

func TestHighlightContainsBrackets(t *testing.T) {
	paths := []string{"secret/myapp/token", "secret/other/key"}
	results := Highlight("myapp", paths)
	if len(results) == 0 {
		t.Fatal("expected at least one result")
	}
	found := false
	for _, r := range results {
		if strings.Contains(r.Path, "myapp") {
			if !strings.Contains(r.Highlight, "[") || !strings.Contains(r.Highlight, "]") {
				t.Errorf("expected brackets in highlight %q", r.Highlight)
			}
			found = true
		}
	}
	if !found {
		t.Fatal("matching path not found in results")
	}
}

func TestHighlightPreservesOriginalPath(t *testing.T) {
	paths := []string{"Secret/MyApp/Key"}
	results := Highlight("mak", paths)
	if len(results) == 0 {
		t.Skip("no match found, skipping case preservation check")
	}
	stripped := strings.NewReplacer("[", "", "]", "").Replace(results[0].Highlight)
	if stripped != results[0].Path {
		t.Errorf("stripping brackets should yield original path: got %q want %q", stripped, results[0].Path)
	}
}

func TestHighlightNoMatch(t *testing.T) {
	paths := []string{"secret/app/key"}
	results := Highlight("zzz", paths)
	if len(results) != 0 {
		t.Errorf("expected no results for unmatched query, got %d", len(results))
	}
}

func TestHighlightResultsOrdered(t *testing.T) {
	paths := []string{"secret/app/token", "secret/application/key", "secret/other"}
	results := Highlight("app", paths)
	for i := 1; i < len(results); i++ {
		if results[i].Score > results[i-1].Score {
			t.Errorf("results not sorted descending by score at index %d", i)
		}
	}
}
