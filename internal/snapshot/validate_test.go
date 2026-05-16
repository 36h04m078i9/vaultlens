package snapshot

import (
	"testing"
	"time"
)

func makeValidateStore(t *testing.T) *Store {
	t.Helper()
	s, err := NewStore()
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	return s
}

func TestNewValidatorNilStoreReturnsError(t *testing.T) {
	_, err := NewValidator(nil)
	if err == nil {
		t.Fatal("expected error for nil store, got nil")
	}
}

func TestValidateEmptyIDReturnsError(t *testing.T) {
	store := makeValidateStore(t)
	v, _ := NewValidator(store)
	_, err := v.Validate("")
	if err == nil {
		t.Fatal("expected error for empty id, got nil")
	}
}

func TestValidateMissingSnapshotReturnsError(t *testing.T) {
	store := makeValidateStore(t)
	v, _ := NewValidator(store)
	_, err := v.Validate("ghost-id")
	if err == nil {
		t.Fatal("expected error for missing snapshot, got nil")
	}
}

func TestValidateHealthySnapshot(t *testing.T) {
	store := makeValidateStore(t)
	snap := &Snapshot{
		ID:          "snap-ok",
		CapturedAt:  time.Now(),
		Entries:     []Entry{{Path: "secret/foo", Data: map[string]string{"k": "v"}}},
	}
	if err := store.Save(snap); err != nil {
		t.Fatalf("Save: %v", err)
	}
	v, _ := NewValidator(store)
	result, err := v.Validate("snap-ok")
	if err != nil {
		t.Fatalf("Validate: %v", err)
	}
	if !result.Valid {
		t.Errorf("expected valid snapshot, got issues: %v", result.Issues)
	}
}

func TestValidateSnapshotWithEmptyEntryPath(t *testing.T) {
	store := makeValidateStore(t)
	snap := &Snapshot{
		ID:         "snap-bad-entry",
		CapturedAt: time.Now(),
		Entries:    []Entry{{Path: "", Data: map[string]string{"k": "v"}}},
	}
	if err := store.Save(snap); err != nil {
		t.Fatalf("Save: %v", err)
	}
	v, _ := NewValidator(store)
	result, err := v.Validate("snap-bad-entry")
	if err != nil {
		t.Fatalf("Validate: %v", err)
	}
	if result.Valid {
		t.Error("expected invalid result for empty entry path")
	}
	if len(result.Issues) == 0 {
		t.Error("expected at least one issue reported")
	}
}

func TestValidateAllReturnsResultsForEachSnapshot(t *testing.T) {
	store := makeValidateStore(t)
	for _, id := range []string{"s1", "s2", "s3"} {
		snap := &Snapshot{ID: id, CapturedAt: time.Now()}
		if err := store.Save(snap); err != nil {
			t.Fatalf("Save %s: %v", id, err)
		}
	}
	v, _ := NewValidator(store)
	results, err := v.ValidateAll()
	if err != nil {
		t.Fatalf("ValidateAll: %v", err)
	}
	if len(results) != 3 {
		t.Errorf("expected 3 results, got %d", len(results))
	}
}
