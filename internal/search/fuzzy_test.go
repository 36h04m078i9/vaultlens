package search

import (
	"testing"
)

var testPaths = []string{
	"secret/database/prod/password",
	"secret/database/staging/password",
	"secret/api/prod/token",
	"secret/api/staging/token",
	"secret/infra/tls/cert",
}

func TestFuzzyEmptyQuery(t *testing.T) {
	results := Fuzzy(testPaths, "")
	if len(results) != len(testPaths) {
		t.Fatalf("expected %d results for empty query, got %d", len(testPaths), len(results))
	}
}

func TestFuzzyExactSubstring(t *testing.T) {
	results := Fuzzy(testPaths, "prod")
	if len(results) != 2 {
		t.Fatalf("expected 2 results for 'prod', got %d", len(results))
	}
	for _, r := range results {
		if r.Score <= 50 {
			t.Errorf("expected high score for exact substring match, got %d for %s", r.Score, r.Path)
		}
	}
}

func TestFuzzySubsequenceMatch(t *testing.T) {
	// 'd', 'b', 'p' appear in order in "database/prod/password"
	results := Fuzzy(testPaths, "dbp")
	if len(results) == 0 {
		t.Fatal("expected at least one subsequence match for 'dbp'")
	}
}

func TestFuzzyNoMatch(t *testing.T) {
	results := Fuzzy(testPaths, "zzzzz")
	if len(results) != 0 {
		t.Fatalf("expected 0 results for non-matching query, got %d", len(results))
	}
}

func TestFuzzyResultsAreSortedDescending(t *testing.T) {
	results := Fuzzy(testPaths, "prod")
	for i := 1; i < len(results); i++ {
		if results[i].Score > results[i-1].Score {
			t.Errorf("results not sorted: index %d has higher score than %d", i, i-1)
		}
	}
}

func TestFuzzyEmptyPaths(t *testing.T) {
	results := Fuzzy([]string{}, "prod")
	if len(results) != 0 {
		t.Fatalf("expected 0 results for empty paths slice, got %d", len(results))
	}
}
