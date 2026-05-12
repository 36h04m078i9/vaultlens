package snapshot

import (
	"testing"
	"time"
)

func makeTagStore(t *testing.T) (*TagStore, *Store) {
	t.Helper()
	s := NewStore()
	ts, err := NewTagStore(s)
	if err != nil {
		t.Fatalf("NewTagStore: %v", err)
	}
	return ts, s
}

func seedSnap(t *testing.T, s *Store, id string) *Snapshot {
	t.Helper()
	snap := &Snapshot{ID: id, CapturedAt: time.Now().UTC(), Paths: []string{"secret/foo"}}
	if err := s.Save(snap); err != nil {
		t.Fatalf("Save: %v", err)
	}
	return snap
}

func TestNewTagStoreNilStoreReturnsError(t *testing.T) {
	_, err := NewTagStore(nil)
	if err == nil {
		t.Fatal("expected error for nil store")
	}
}

func TestAddTagSuccess(t *testing.T) {
	ts, s := makeTagStore(t)
	seedSnap(t, s, "snap-1")

	tag, err := ts.Add("release-v1", "snap-1", "first release")
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	if tag.Name != "release-v1" {
		t.Errorf("got name %q, want %q", tag.Name, "release-v1")
	}
	if tag.SnapshotID != "snap-1" {
		t.Errorf("got snapshot_id %q, want %q", tag.SnapshotID, "snap-1")
	}
}

func TestAddTagEmptyNameReturnsError(t *testing.T) {
	ts, s := makeTagStore(t)
	seedSnap(t, s, "snap-1")

	_, err := ts.Add("", "snap-1", "")
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestAddTagMissingSnapshotReturnsError(t *testing.T) {
	ts, _ := makeTagStore(t)

	_, err := ts.Add("my-tag", "nonexistent", "")
	if err == nil {
		t.Fatal("expected error for missing snapshot")
	}
}

func TestGetTagFound(t *testing.T) {
	ts, s := makeTagStore(t)
	seedSnap(t, s, "snap-2")
	ts.Add("v2", "snap-2", "")

	tag, ok := ts.Get("v2")
	if !ok {
		t.Fatal("expected tag to be found")
	}
	if tag.SnapshotID != "snap-2" {
		t.Errorf("got snapshot_id %q, want %q", tag.SnapshotID, "snap-2")
	}
}

func TestDeleteTagRemovesEntry(t *testing.T) {
	ts, s := makeTagStore(t)
	seedSnap(t, s, "snap-3")
	ts.Add("to-delete", "snap-3", "")

	if !ts.Delete("to-delete") {
		t.Fatal("expected Delete to return true")
	}
	if _, ok := ts.Get("to-delete"); ok {
		t.Fatal("expected tag to be gone after deletion")
	}
}

func TestDeleteMissingTagReturnsFalse(t *testing.T) {
	ts, _ := makeTagStore(t)
	if ts.Delete("ghost") {
		t.Fatal("expected Delete to return false for missing tag")
	}
}

func TestListReturnAllTags(t *testing.T) {
	ts, s := makeTagStore(t)
	seedSnap(t, s, "snap-a")
	seedSnap(t, s, "snap-b")
	ts.Add("alpha", "snap-a", "")
	ts.Add("beta", "snap-b", "")

	tags := ts.List()
	if len(tags) != 2 {
		t.Errorf("got %d tags, want 2", len(tags))
	}
}
