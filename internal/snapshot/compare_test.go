package snapshot

import (
	"testing"
	"time"
)

func makeCompareStore(t *testing.T) *Store {
	t.Helper()
	return NewStore()
}

func saveSnap(t *testing.T, st *Store, id string, secrets []SecretEntry) {
	t.Helper()
	sn := &Snapshot{ID: id, TakenAt: time.Now(), Secrets: secrets}
	if err := st.Save(sn); err != nil {
		t.Fatalf("save: %v", err)
	}
}

func TestCompareNoChanges(t *testing.T) {
	st := makeCompareStore(t)
	entries := []SecretEntry{{Path: "secret/a", Value: "v1"}}
	saveSnap(t, st, "snap1", entries)
	saveSnap(t, st, "snap2", entries)

	res, err := Compare(st, "snap1", "snap2", CompareOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Added)+len(res.Removed)+len(res.Changed) != 0 {
		t.Errorf("expected no diff, got %+v", res)
	}
	if res.Unchanged != 1 {
		t.Errorf("expected 1 unchanged, got %d", res.Unchanged)
	}
}

func TestCompareDetectsAdded(t *testing.T) {
	st := makeCompareStore(t)
	saveSnap(t, st, "a", []SecretEntry{})
	saveSnap(t, st, "b", []SecretEntry{{Path: "secret/new", Value: "x"}})

	res, err := Compare(st, "a", "b", CompareOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Added) != 1 || res.Added[0] != "secret/new" {
		t.Errorf("expected added secret/new, got %v", res.Added)
	}
}

func TestCompareDetectsRemoved(t *testing.T) {
	st := makeCompareStore(t)
	saveSnap(t, st, "a", []SecretEntry{{Path: "secret/gone", Value: "v"}})
	saveSnap(t, st, "b", []SecretEntry{})

	res, err := Compare(st, "a", "b", CompareOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Removed) != 1 || res.Removed[0] != "secret/gone" {
		t.Errorf("expected removed secret/gone, got %v", res.Removed)
	}
}

func TestCompareDetectsChanged(t *testing.T) {
	st := makeCompareStore(t)
	saveSnap(t, st, "a", []SecretEntry{{Path: "secret/x", Value: "old"}})
	saveSnap(t, st, "b", []SecretEntry{{Path: "secret/x", Value: "new"}})

	res, err := Compare(st, "a", "b", CompareOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Changed) != 1 || res.Changed[0] != "secret/x" {
		t.Errorf("expected changed secret/x, got %v", res.Changed)
	}
}

func TestCompareMissingSnapshotReturnsError(t *testing.T) {
	st := makeCompareStore(t)
	saveSnap(t, st, "only", []SecretEntry{})

	_, err := Compare(st, "only", "missing", CompareOptions{})
	if err == nil {
		t.Error("expected error for missing snapshot")
	}
}

func TestCompareSummaryFormat(t *testing.T) {
	st := makeCompareStore(t)
	saveSnap(t, st, "s1", []SecretEntry{{Path: "p", Value: "v"}})
	saveSnap(t, st, "s2", []SecretEntry{{Path: "p", Value: "v"}})

	res, _ := Compare(st, "s1", "s2", CompareOptions{})
	got := res.Summary()
	if got == "" {
		t.Error("expected non-empty summary")
	}
}
