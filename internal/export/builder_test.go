package export_test

import (
	"testing"
	"time"

	"github.com/yourorg/vaultlens/internal/bookmark"
	"github.com/yourorg/vaultlens/internal/export"
)

func newBuilderWithFixedTime() *export.Builder {
	b := export.NewBuilder()
	// Access via exported field not available; we rely on NewBuilder default
	// and accept that timestamps are close to now — tested structurally.
	_ = b
	return export.NewBuilder()
}

func TestBuilderFromBookmarks(t *testing.T) {
	bms := []bookmark.Bookmark{
		{Path: "secret/a", Note: "alpha"},
		{Path: "secret/b", Note: ""},
	}
	b := export.NewBuilder()
	records := b.FromBookmarks(bms)
	if len(records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(records))
	}
	if records[0].Path != "secret/a" {
		t.Errorf("unexpected path: %s", records[0].Path)
	}
	if records[0].Note != "alpha" {
		t.Errorf("unexpected note: %s", records[0].Note)
	}
	if records[1].Note != "" {
		t.Errorf("expected empty note for second record")
	}
}

func TestBuilderFromBookmarksTimestampSet(t *testing.T) {
	before := time.Now()
	b := export.NewBuilder()
	records := b.FromBookmarks([]bookmark.Bookmark{{Path: "secret/x"}})
	after := time.Now()
	if records[0].ExportedAt.Before(before) || records[0].ExportedAt.After(after) {
		t.Errorf("ExportedAt %v not within expected range", records[0].ExportedAt)
	}
}

func TestBuilderFromPaths(t *testing.T) {
	paths := []string{"secret/foo", "secret/bar"}
	b := export.NewBuilder()
	records := b.FromPaths(paths)
	if len(records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(records))
	}
	if records[0].Path != "secret/foo" {
		t.Errorf("unexpected path: %s", records[0].Path)
	}
	if records[0].Note != "" {
		t.Errorf("expected empty note for path-only record")
	}
}

func TestBuilderFromEmptyPaths(t *testing.T) {
	b := export.NewBuilder()
	records := b.FromPaths([]string{})
	if len(records) != 0 {
		t.Errorf("expected empty records, got %d", len(records))
	}
}
