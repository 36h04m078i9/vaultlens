package snapshot

import (
	"testing"
)

func makeLabelStore(t *testing.T) (*LabelStore, Store) {
	t.Helper()
	s := NewStore()
	ls, err := NewLabelStore(s)
	if err != nil {
		t.Fatalf("NewLabelStore: %v", err)
	}
	return ls, s
}

func seedLabelSnap(t *testing.T, s Store, id string) {
	t.Helper()
	if err := s.Save(&Snapshot{ID: id, Paths: map[string]map[string]string{}});  err != nil {
		t.Fatalf("Save: %v", err)
	}
}

func TestNewLabelStoreNilStoreReturnsError(t *testing.T) {
	_, err := NewLabelStore(nil)
	if err == nil {
		t.Fatal("expected error for nil store")
	}
}

func TestAddLabelSuccess(t *testing.T) {
	ls, s := makeLabelStore(t)
	seedLabelSnap(t, s, "snap1")
	if err := ls.Add("snap1", "production"); err != nil {
		t.Fatalf("Add: %v", err)
	}
	labels := ls.List("snap1")
	if len(labels) != 1 || labels[0].Name != "production" {
		t.Fatalf("expected label 'production', got %v", labels)
	}
}

func TestAddLabelEmptyIDReturnsError(t *testing.T) {
	ls, _ := makeLabelStore(t)
	if err := ls.Add("", "production"); err == nil {
		t.Fatal("expected error for empty id")
	}
}

func TestAddLabelEmptyNameReturnsError(t *testing.T) {
	ls, s := makeLabelStore(t)
	seedLabelSnap(t, s, "snap1")
	if err := ls.Add("snap1", ""); err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestAddLabelMissingSnapshotReturnsError(t *testing.T) {
	ls, _ := makeLabelStore(t)
	if err := ls.Add("ghost", "label"); err == nil {
		t.Fatal("expected error for unknown snapshot")
	}
}

func TestAddLabelIdempotent(t *testing.T) {
	ls, s := makeLabelStore(t)
	seedLabelSnap(t, s, "snap1")
	_ = ls.Add("snap1", "stable")
	_ = ls.Add("snap1", "stable")
	if got := len(ls.List("snap1")); got != 1 {
		t.Fatalf("expected 1 label, got %d", got)
	}
}

func TestRemoveLabelSuccess(t *testing.T) {
	ls, s := makeLabelStore(t)
	seedLabelSnap(t, s, "snap1")
	_ = ls.Add("snap1", "canary")
	if !ls.Remove("snap1", "canary") {
		t.Fatal("expected Remove to return true")
	}
	if len(ls.List("snap1")) != 0 {
		t.Fatal("expected no labels after removal")
	}
}

func TestRemoveMissingLabelReturnsFalse(t *testing.T) {
	ls, s := makeLabelStore(t)
	seedLabelSnap(t, s, "snap1")
	if ls.Remove("snap1", "nonexistent") {
		t.Fatal("expected Remove to return false")
	}
}

func TestFindByLabel(t *testing.T) {
	ls, s := makeLabelStore(t)
	seedLabelSnap(t, s, "a")
	seedLabelSnap(t, s, "b")
	seedLabelSnap(t, s, "c")
	_ = ls.Add("a", "prod")
	_ = ls.Add("b", "prod")
	_ = ls.Add("c", "staging")
	ids := ls.FindByLabel("prod")
	if len(ids) != 2 {
		t.Fatalf("expected 2 ids, got %d: %v", len(ids), ids)
	}
}
