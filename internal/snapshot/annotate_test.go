package snapshot

import (
	"testing"
)

func makeAnnotationStore(t *testing.T) (*AnnotationStore, *Store) {
	t.Helper()
	s := NewStore()
	as, err := NewAnnotationStore(s)
	if err != nil {
		t.Fatalf("NewAnnotationStore: %v", err)
	}
	return as, s
}

func TestNewAnnotationStoreNilStoreReturnsError(t *testing.T) {
	_, err := NewAnnotationStore(nil)
	if err == nil {
		t.Fatal("expected error for nil store, got nil")
	}
}

func TestAddAnnotationSuccess(t *testing.T) {
	as, s := makeAnnotationStore(t)
	snap := makeSnap("ann-1", "secret/foo")
	s.Save(snap)

	if err := as.Add("ann-1", "alice", "initial baseline"); err != nil {
		t.Fatalf("Add: %v", err)
	}
	list := as.List("ann-1")
	if len(list) != 1 {
		t.Fatalf("expected 1 annotation, got %d", len(list))
	}
	if list[0].Note != "initial baseline" {
		t.Errorf("unexpected note: %q", list[0].Note)
	}
	if list[0].Author != "alice" {
		t.Errorf("unexpected author: %q", list[0].Author)
	}
}

func TestAddAnnotationEmptyIDReturnsError(t *testing.T) {
	as, _ := makeAnnotationStore(t)
	if err := as.Add("", "alice", "note"); err == nil {
		t.Fatal("expected error for empty id")
	}
}

func TestAddAnnotationEmptyNoteReturnsError(t *testing.T) {
	as, s := makeAnnotationStore(t)
	snap := makeSnap("ann-2", "secret/bar")
	s.Save(snap)
	if err := as.Add("ann-2", "bob", ""); err == nil {
		t.Fatal("expected error for empty note")
	}
}

func TestAddAnnotationMissingSnapshotReturnsError(t *testing.T) {
	as, _ := makeAnnotationStore(t)
	if err := as.Add("nonexistent", "alice", "some note"); err == nil {
		t.Fatal("expected error for missing snapshot")
	}
}

func TestListReturnsEmptyForUnknownID(t *testing.T) {
	as, _ := makeAnnotationStore(t)
	list := as.List("ghost")
	if len(list) != 0 {
		t.Fatalf("expected empty list, got %d", len(list))
	}
}

func TestRemoveAnnotationsSuccess(t *testing.T) {
	as, s := makeAnnotationStore(t)
	snap := makeSnap("ann-3", "secret/baz")
	s.Save(snap)
	as.Add("ann-3", "carol", "note one")
	as.Add("ann-3", "carol", "note two")

	if removed := as.Remove("ann-3"); !removed {
		t.Fatal("expected Remove to return true")
	}
	if list := as.List("ann-3"); len(list) != 0 {
		t.Fatalf("expected empty list after remove, got %d", len(list))
	}
}

func TestRemoveMissingAnnotationsReturnsFalse(t *testing.T) {
	as, _ := makeAnnotationStore(t)
	if removed := as.Remove("no-such-id"); removed {
		t.Fatal("expected Remove to return false for unknown id")
	}
}
