package snapshot

import (
	"testing"
	"time"
)

func makeMergeStore() *Store {
	s, _ := NewStore()
	return s
}

func seedMergeSnap(t *testing.T, store *Store, id string, data map[string]map[string]string) {
	t.Helper()
	err := store.Save(&Snapshot{ID: id, CreatedAt: time.Now().UTC(), Data: data})
	if err != nil {
		t.Fatalf("seed snapshot %q: %v", id, err)
	}
}

func TestMergeNilStoreReturnsError(t *testing.T) {
	_, err := Merge(nil, "a", "b", MergeOptions{})
	if err == nil {
		t.Fatal("expected error for nil store")
	}
}

func TestMergeEmptyIDReturnsError(t *testing.T) {
	store := makeMergeStore()
	_, err := Merge(store, "", "b", MergeOptions{})
	if err == nil {
		t.Fatal("expected error for empty idA")
	}
	_, err = Merge(store, "a", "", MergeOptions{})
	if err == nil {
		t.Fatal("expected error for empty idB")
	}
}

func TestMergeMissingSnapshotReturnsError(t *testing.T) {
	store := makeMergeStore()
	seedMergeSnap(t, store, "a", map[string]map[string]string{
		"secret/foo": {"key": "val"},
	})
	_, err := Merge(store, "a", "missing", MergeOptions{})
	if err == nil {
		t.Fatal("expected error for missing snapshot B")
	}
}

func TestMergeDisjointPathsMergesAll(t *testing.T) {
	store := makeMergeStore()
	seedMergeSnap(t, store, "a", map[string]map[string]string{
		"secret/alpha": {"x": "1"},
	})
	seedMergeSnap(t, store, "b", map[string]map[string]string{
		"secret/beta": {"y": "2"},
	})
	res, err := Merge(store, "a", "b", MergeOptions{NewID: "merged-ab"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Paths) != 2 {
		t.Fatalf("expected 2 paths, got %d", len(res.Paths))
	}
	if len(res.Conflicts) != 0 {
		t.Fatalf("expected no conflicts, got %d", len(res.Conflicts))
	}
	if res.ID != "merged-ab" {
		t.Fatalf("expected ID merged-ab, got %s", res.ID)
	}
}

func TestMergeConflictPreferA(t *testing.T) {
	store := makeMergeStore()
	seedMergeSnap(t, store, "a", map[string]map[string]string{
		"secret/shared": {"token": "from-a"},
	})
	seedMergeSnap(t, store, "b", map[string]map[string]string{
		"secret/shared": {"token": "from-b"},
	})
	res, err := Merge(store, "a", "b", MergeOptions{PreferA: true, NewID: "merged-prefer-a"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Conflicts) != 1 {
		t.Fatalf("expected 1 conflict, got %d", len(res.Conflicts))
	}
	if res.Paths["secret/shared"]["token"] != "from-a" {
		t.Fatalf("expected from-a, got %s", res.Paths["secret/shared"]["token"])
	}
}

func TestMergeConflictPreferB(t *testing.T) {
	store := makeMergeStore()
	seedMergeSnap(t, store, "a", map[string]map[string]string{
		"secret/shared": {"token": "from-a"},
	})
	seedMergeSnap(t, store, "b", map[string]map[string]string{
		"secret/shared": {"token": "from-b"},
	})
	res, err := Merge(store, "a", "b", MergeOptions{PreferA: false, NewID: "merged-prefer-b"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Paths["secret/shared"]["token"] != "from-b" {
		t.Fatalf("expected from-b, got %s", res.Paths["secret/shared"]["token"])
	}
}

func TestMergeDefaultIDGenerated(t *testing.T) {
	store := makeMergeStore()
	seedMergeSnap(t, store, "snap1", map[string]map[string]string{
		"secret/a": {"k": "v"},
	})
	seedMergeSnap(t, store, "snap2", map[string]map[string]string{
		"secret/b": {"k": "v"},
	})
	res, err := Merge(store, "snap1", "snap2", MergeOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.ID != "merged-snap1-snap2" {
		t.Fatalf("expected merged-snap1-snap2, got %s", res.ID)
	}
}

func TestMergeResultPersistedInStore(t *testing.T) {
	store := makeMergeStore()
	seedMergeSnap(t, store, "x", map[string]map[string]string{
		"secret/x": {"k": "1"},
	})
	seedMergeSnap(t, store, "y", map[string]map[string]string{
		"secret/y": {"k": "2"},
	})
	res, err := Merge(store, "x", "y", MergeOptions{NewID: "merged-xy"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, ok := store.Get(res.ID)
	if !ok {
		t.Fatal("merged snapshot not persisted in store")
	}
}
