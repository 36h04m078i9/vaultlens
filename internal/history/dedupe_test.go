package history

import (
	"testing"
	"time"
)

func makeDedupEntries() []Entry {
	now := time.Now()
	return []Entry{
		{Path: "secret/a", AccessedAt: now.Add(-4 * time.Minute)},
		{Path: "secret/b", AccessedAt: now.Add(-3 * time.Minute)},
		{Path: "secret/a", AccessedAt: now.Add(-2 * time.Minute)}, // duplicate of first
		{Path: "secret/c", AccessedAt: now.Add(-1 * time.Minute)},
		{Path: "secret/b", AccessedAt: now},                        // duplicate of second
	}
}

func TestDedupeEmpty(t *testing.T) {
	result := Dedupe(nil, DedupeOptions{})
	if result != nil {
		t.Fatalf("expected nil for empty input, got %v", result)
	}
}

func TestDedupeKeepLast(t *testing.T) {
	entries := makeDedupEntries()
	result := Dedupe(entries, DedupeOptions{KeepFirst: false})

	if len(result) != 3 {
		t.Fatalf("expected 3 unique entries, got %d", len(result))
	}

	// secret/a at index 0 should be the later timestamp
	if result[0].Path != "secret/a" {
		t.Errorf("expected secret/a at index 0, got %s", result[0].Path)
	}
	if result[0].AccessedAt != entries[2].AccessedAt {
		t.Errorf("expected updated timestamp for secret/a")
	}

	// secret/b at index 1 should be the later timestamp
	if result[1].Path != "secret/b" {
		t.Errorf("expected secret/b at index 1, got %s", result[1].Path)
	}
	if result[1].AccessedAt != entries[4].AccessedAt {
		t.Errorf("expected updated timestamp for secret/b")
	}
}

func TestDedupeKeepFirst(t *testing.T) {
	entries := makeDedupEntries()
	result := Dedupe(entries, DedupeOptions{KeepFirst: true})

	if len(result) != 3 {
		t.Fatalf("expected 3 unique entries, got %d", len(result))
	}

	// secret/a should retain the original (earliest) timestamp
	if result[0].AccessedAt != entries[0].AccessedAt {
		t.Errorf("expected original timestamp for secret/a")
	}

	// secret/b should retain the original timestamp
	if result[1].AccessedAt != entries[1].AccessedAt {
		t.Errorf("expected original timestamp for secret/b")
	}
}

func TestDedupeNoDuplicates(t *testing.T) {
	now := time.Now()
	entries := []Entry{
		{Path: "secret/x", AccessedAt: now},
		{Path: "secret/y", AccessedAt: now},
	}
	result := Dedupe(entries, DedupeOptions{})
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
}
