package snapshot

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func makeImportStore(t *testing.T) *Store {
	t.Helper()
	s, err := NewStore(t.TempDir())
	if err != nil {
		t.Fatalf("makeImportStore: %v", err)
	}
	return s
}

func encodeSnaps(t *testing.T, snaps ...Snapshot) io.Reader {
	t.Helper()
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	for _, s := range snaps {
		if err := enc.Encode(s); err != nil {
			t.Fatalf("encodeSnaps: %v", err)
		}
	}
	return &buf
}

func TestImportNilReaderReturnsError(t *testing.T) {
	store := makeImportStore(t)
	_, err := Import(nil, store, ImportOptions{})
	if err == nil {
		t.Fatal("expected error for nil reader")
	}
}

func TestImportNilStoreReturnsError(t *testing.T) {
	_, err := Import(strings.NewReader(""), nil, ImportOptions{})
	if err == nil {
		t.Fatal("expected error for nil store")
	}
}

func TestImportStoresSnapshot(t *testing.T) {
	store := makeImportStore(t)
	snap := Snapshot{ID: "snap-1", Prefix: "secret/", CapturedAt: time.Now().UTC()}
	r := encodeSnaps(t, snap)

	results, err := Import(r, store, ImportOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || !results[0].Imported {
		t.Fatalf("expected one imported result, got %+v", results)
	}
	_, ok := store.Get("snap-1")
	if !ok {
		t.Fatal("snapshot not found in store after import")
	}
}

func TestImportSkipsDuplicateByDefault(t *testing.T) {
	store := makeImportStore(t)
	snap := Snapshot{ID: "snap-dup", Prefix: "secret/", CapturedAt: time.Now().UTC()}
	_ = store.Save(&snap)

	r := encodeSnaps(t, snap)
	results, err := Import(r, store, ImportOptions{OverwriteExisting: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || !results[0].Skipped {
		t.Fatalf("expected skipped result, got %+v", results)
	}
}

func TestImportOverwritesWhenEnabled(t *testing.T) {
	store := makeImportStore(t)
	snap := Snapshot{ID: "snap-ow", Prefix: "secret/", CapturedAt: time.Now().UTC()}
	_ = store.Save(&snap)

	r := encodeSnaps(t, snap)
	results, err := Import(r, store, ImportOptions{OverwriteExisting: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || !results[0].Overwrite || !results[0].Imported {
		t.Fatalf("expected overwrite result, got %+v", results)
	}
}

func TestImportAppendsTagsToAdd(t *testing.T) {
	store := makeImportStore(t)
	snap := Snapshot{ID: "snap-tag", Prefix: "secret/", CapturedAt: time.Now().UTC()}
	r := encodeSnaps(t, snap)

	_, err := Import(r, store, ImportOptions{TagsToAdd: []string{"imported", "ci"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, ok := store.Get("snap-tag")
	if !ok {
		t.Fatal("snapshot not found")
	}
	tagSet := make(map[string]bool)
	for _, tag := range got.Tags {
		tagSet[tag] = true
	}
	for _, want := range []string{"imported", "ci"} {
		if !tagSet[want] {
			t.Errorf("expected tag %q to be present", want)
		}
	}
}

func TestImportMissingIDReturnsError(t *testing.T) {
	store := makeImportStore(t)
	snap := Snapshot{Prefix: "secret/"}
	r := encodeSnaps(t, snap)
	_, err := Import(r, store, ImportOptions{})
	if err == nil {
		t.Fatal("expected error for snapshot with empty ID")
	}
}
