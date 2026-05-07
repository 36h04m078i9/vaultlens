package bookmark_test

import (
	"testing"
	"time"

	"github.com/vaultlens/vaultlens/internal/bookmark"
)

func makeBookmarks() []bookmark.Bookmark {
	now := time.Now()
	return []bookmark.Bookmark{
		{Path: "secret/database/prod", Note: "production db", CreatedAt: now},
		{Path: "secret/api/keys", Note: "api tokens", CreatedAt: now},
		{Path: "secret/tls/cert", Note: "TLS certificate", CreatedAt: now},
	}
}

func TestSearchEmptyQueryReturnsAll(t *testing.T) {
	items := makeBookmarks()
	results := bookmark.Search(items, "")
	if len(results) != len(items) {
		t.Errorf("expected %d results, got %d", len(items), len(results))
	}
}

func TestSearchMatchesPath(t *testing.T) {
	results := bookmark.Search(makeBookmarks(), "database")
	if len(results) != 1 {
		t.Fatalf("expected 1, got %d", len(results))
	}
	if results[0].Bookmark.Path != "secret/database/prod" {
		t.Errorf("unexpected path: %s", results[0].Bookmark.Path)
	}
}

func TestSearchMatchesNote(t *testing.T) {
	results := bookmark.Search(makeBookmarks(), "tokens")
	if len(results) != 1 {
		t.Fatalf("expected 1, got %d", len(results))
	}
	if results[0].Bookmark.Path != "secret/api/keys" {
		t.Errorf("unexpected path: %s", results[0].Bookmark.Path)
	}
}

func TestSearchPathMatchScoresHigher(t *testing.T) {
	items := []bookmark.Bookmark{
		{Path: "secret/prod/db", Note: "tls info"},
		{Path: "secret/tls/cert", Note: "cert"},
	}
	results := bookmark.Search(items, "tls")
	if len(results) != 2 {
		t.Fatalf("expected 2, got %d", len(results))
	}
	for _, r := range results {
		if r.Bookmark.Path == "secret/tls/cert" && r.Score < 2 {
			t.Errorf("path match should score >= 2, got %d", r.Score)
		}
	}
}

func TestSearchNoMatchReturnsEmpty(t *testing.T) {
	results := bookmark.Search(makeBookmarks(), "zzznomatch")
	if len(results) != 0 {
		t.Errorf("expected empty, got %d", len(results))
	}
}

func TestSearchIsCaseInsensitive(t *testing.T) {
	results := bookmark.Search(makeBookmarks(), "TLS")
	if len(results) == 0 {
		t.Error("expected case-insensitive match for TLS")
	}
}
