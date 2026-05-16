package snapshot

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func makeArchiveStore(t *testing.T) *Store {
	t.Helper()
	s, err := NewStore(t.TempDir())
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	return s
}

func seedArchiveSnap(t *testing.T, store *Store, id string, at time.Time, data map[string]string) {
	t.Helper()
	snap := &Snapshot{ID: id, CapturedAt: at, Data: data}
	if err := store.Save(snap); err != nil {
		t.Fatalf("Save(%s): %v", id, err)
	}
}

func TestArchiveNilStoreReturnsError(t *testing.T) {
	err := Archive(nil, &bytes.Buffer{}, ArchiveOptions{})
	if err == nil {
		t.Fatal("expected error for nil store")
	}
}

func TestArchiveNilWriterReturnsError(t *testing.T) {
	store := makeArchiveStore(t)
	err := Archive(store, nil, ArchiveOptions{})
	if err == nil {
		t.Fatal("expected error for nil writer")
	}
}

func TestArchiveNoMatchReturnsError(t *testing.T) {
	store := makeArchiveStore(t)
	now := time.Now()
	seedArchiveSnap(t, store, "snap-1", now.Add(-2*time.Hour), map[string]string{"k": "v"})

	future := now.Add(time.Hour)
	err := Archive(store, &bytes.Buffer{}, ArchiveOptions{Since: future})
	if err == nil {
		t.Fatal("expected error when no snapshots match")
	}
}

func TestArchiveWritesAllSnapshots(t *testing.T) {
	store := makeArchiveStore(t)
	now := time.Now()
	seedArchiveSnap(t, store, "snap-a", now.Add(-1*time.Hour), map[string]string{"secret/foo": "bar"})
	seedArchiveSnap(t, store, "snap-b", now, map[string]string{"secret/baz": "qux"})

	var buf bytes.Buffer
	if err := Archive(store, &buf, ArchiveOptions{}); err != nil {
		t.Fatalf("Archive: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
}

func TestArchiveFilterByID(t *testing.T) {
	store := makeArchiveStore(t)
	now := time.Now()
	seedArchiveSnap(t, store, "snap-x", now, map[string]string{"a": "1"})
	seedArchiveSnap(t, store, "snap-y", now, map[string]string{"b": "2"})

	var buf bytes.Buffer
	if err := Archive(store, &buf, ArchiveOptions{IDs: []string{"snap-x"}}); err != nil {
		t.Fatalf("Archive: %v", err)
	}

	var rec ArchiveRecord
	if err := json.NewDecoder(&buf).Decode(&rec); err != nil {
		t.Fatalf("Decode: %v", err)
	}
	if rec.ID != "snap-x" {
		t.Errorf("expected snap-x, got %s", rec.ID)
	}
}

func TestArchiveFilterBySince(t *testing.T) {
	store := makeArchiveStore(t)
	now := time.Now()
	seedArchiveSnap(t, store, "old", now.Add(-5*time.Hour), map[string]string{"x": "1"})
	seedArchiveSnap(t, store, "new", now, map[string]string{"y": "2"})

	var buf bytes.Buffer
	cutoff := now.Add(-1 * time.Hour)
	if err := Archive(store, &buf, ArchiveOptions{Since: cutoff}); err != nil {
		t.Fatalf("Archive: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}
	var rec ArchiveRecord
	if err := json.Unmarshal([]byte(lines[0]), &rec); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if rec.ID != "new" {
		t.Errorf("expected 'new', got %s", rec.ID)
	}
}
