package bookmark_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/vaultlens/vaultlens/internal/bookmark"
)

func tempStore(t *testing.T) *bookmark.Store {
	t.Helper()
	f := filepath.Join(t.TempDir(), "bookmarks.json")
	s, err := bookmark.NewStore(f)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	return s
}

func TestAddAndList(t *testing.T) {
	s := tempStore(t)
	if err := s.Add("secret/foo", "my note"); err != nil {
		t.Fatalf("Add: %v", err)
	}
	list := s.List()
	if len(list) != 1 {
		t.Fatalf("expected 1 bookmark, got %d", len(list))
	}
	if list[0].Path != "secret/foo" {
		t.Errorf("unexpected path: %s", list[0].Path)
	}
	if list[0].Note != "my note" {
		t.Errorf("unexpected note: %s", list[0].Note)
	}
}

func TestAddDuplicateIsIdempotent(t *testing.T) {
	s := tempStore(t)
	_ = s.Add("secret/bar", "")
	_ = s.Add("secret/bar", "duplicate")
	if n := len(s.List()); n != 1 {
		t.Errorf("expected 1, got %d", n)
	}
}

func TestRemoveExisting(t *testing.T) {
	s := tempStore(t)
	_ = s.Add("secret/baz", "")
	ok, err := s.Remove("secret/baz")
	if err != nil {
		t.Fatalf("Remove: %v", err)
	}
	if !ok {
		t.Error("expected Remove to return true")
	}
	if n := len(s.List()); n != 0 {
		t.Errorf("expected empty list, got %d", n)
	}
}

func TestRemoveMissingReturnsFalse(t *testing.T) {
	s := tempStore(t)
	ok, err := s.Remove("secret/nonexistent")
	if err != nil {
		t.Fatalf("Remove: %v", err)
	}
	if ok {
		t.Error("expected Remove to return false for missing path")
	}
}

func TestPersistenceAcrossReloads(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "bookmarks.json")
	s1, _ := bookmark.NewStore(f)
	_ = s1.Add("secret/persist", "persisted")

	s2, err := bookmark.NewStore(f)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	list := s2.List()
	if len(list) != 1 || list[0].Path != "secret/persist" {
		t.Errorf("data not persisted: %+v", list)
	}
}

func TestNewStoreCreatesFileOnFirstAdd(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "bookmarks.json")
	s, _ := bookmark.NewStore(f)
	_ = s.Add("secret/new", "")
	if _, err := os.Stat(f); err != nil {
		t.Errorf("file not created: %v", err)
	}
}
