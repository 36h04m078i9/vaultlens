package snapshot

import (
	"testing"
	"time"
)

func makeSnap(id, label string, t time.Time) *Snapshot {
	return &Snapshot{ID: id, Label: label, CapturedAt: t, Paths: []string{"secret/a"}}
}

func TestSaveAndGet(t *testing.T) {
	s := NewStore()
	snap := makeSnap("id1", "first", time.Now())
	if err := s.Save(snap); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, ok := s.Get("id1")
	if !ok {
		t.Fatal("expected snapshot to be found")
	}
	if got.Label != "first" {
		t.Errorf("label mismatch: got %q", got.Label)
	}
}

func TestSaveNilReturnsError(t *testing.T) {
	s := NewStore()
	if err := s.Save(nil); err == nil {
		t.Fatal("expected error for nil snapshot")
	}
}

func TestSaveEmptyIDReturnsError(t *testing.T) {
	s := NewStore()
	snap := &Snapshot{ID: "", Label: "x"}
	if err := s.Save(snap); err == nil {
		t.Fatal("expected error for empty ID")
	}
}

func TestGetMissReturnsFalse(t *testing.T) {
	s := NewStore()
	_, ok := s.Get("missing")
	if ok {
		t.Fatal("expected miss")
	}
}

func TestListOrderedByTime(t *testing.T) {
	s := NewStore()
	now := time.Now()
	_ = s.Save(makeSnap("b", "second", now.Add(time.Second)))
	_ = s.Save(makeSnap("a", "first", now))
	list := s.List()
	if len(list) != 2 {
		t.Fatalf("expected 2, got %d", len(list))
	}
	if list[0].ID != "a" {
		t.Errorf("expected 'a' first, got %q", list[0].ID)
	}
}

func TestDeleteExisting(t *testing.T) {
	s := NewStore()
	_ = s.Save(makeSnap("id1", "x", time.Now()))
	if !s.Delete("id1") {
		t.Fatal("expected Delete to return true")
	}
	_, ok := s.Get("id1")
	if ok {
		t.Fatal("expected snapshot to be gone")
	}
}

func TestDeleteMissingReturnsFalse(t *testing.T) {
	s := NewStore()
	if s.Delete("nope") {
		t.Fatal("expected false for missing ID")
	}
}
