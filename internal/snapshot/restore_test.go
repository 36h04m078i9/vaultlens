package snapshot

import (
	"errors"
	"testing"
	"time"
)

// fakeWriter records WriteSecret calls.
type fakeWriter struct {
	written map[string]map[string]interface{}
	errOn   string
}

func (f *fakeWriter) WriteSecret(path string, data map[string]interface{}) error {
	if f.errOn == path {
		return errors.New("write error")
	}
	if f.written == nil {
		f.written = make(map[string]map[string]interface{})
	}
	f.written[path] = data
	return nil
}

func makeRestoreStore(t *testing.T) *Store {
	t.Helper()
	s, err := NewStore(t.TempDir())
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	return s
}

func TestNewRestorerNilStoreReturnsError(t *testing.T) {
	_, err := NewRestorer(nil, &fakeWriter{})
	if err == nil {
		t.Fatal("expected error for nil store")
	}
}

func TestNewRestorerNilWriterReturnsError(t *testing.T) {
	s := makeRestoreStore(t)
	_, err := NewRestorer(s, nil)
	if err == nil {
		t.Fatal("expected error for nil writer")
	}
}

func TestRestoreEmptyIDReturnsError(t *testing.T) {
	s := makeRestoreStore(t)
	r, _ := NewRestorer(s, &fakeWriter{})
	_, err := r.Restore("", RestoreOptions{})
	if err == nil {
		t.Fatal("expected error for empty id")
	}
}

func TestRestoreMissingSnapshotReturnsError(t *testing.T) {
	s := makeRestoreStore(t)
	r, _ := NewRestorer(s, &fakeWriter{})
	_, err := r.Restore("nonexistent", RestoreOptions{})
	if err == nil {
		t.Fatal("expected error for missing snapshot")
	}
}

func TestRestoreWritesAllPaths(t *testing.T) {
	s := makeRestoreStore(t)
	snap := &Snapshot{
		ID:        "snap1",
		CreatedAt: time.Now(),
		Secrets: map[string]map[string]interface{}{
			"secret/a": {"key": "val1"},
			"secret/b": {"key": "val2"},
		},
	}
	if err := s.Save(snap); err != nil {
		t.Fatalf("Save: %v", err)
	}
	w := &fakeWriter{}
	r, _ := NewRestorer(s, w)
	records, err := r.Restore("snap1", RestoreOptions{})
	if err != nil {
		t.Fatalf("Restore: %v", err)
	}
	if len(records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(records))
	}
	if len(w.written) != 2 {
		t.Fatalf("expected 2 writes, got %d", len(w.written))
	}
}

func TestRestoreDryRunDoesNotWrite(t *testing.T) {
	s := makeRestoreStore(t)
	snap := &Snapshot{
		ID:        "snap2",
		CreatedAt: time.Now(),
		Secrets: map[string]map[string]interface{}{
			"secret/a": {"key": "val"},
		},
	}
	_ = s.Save(snap)
	w := &fakeWriter{}
	r, _ := NewRestorer(s, w)
	records, err := r.Restore("snap2", RestoreOptions{DryRun: true})
	if err != nil {
		t.Fatalf("Restore: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}
	if len(w.written) != 0 {
		t.Fatal("expected no writes in dry-run mode")
	}
}

func TestRestorePathPrefixFilters(t *testing.T) {
	s := makeRestoreStore(t)
	snap := &Snapshot{
		ID:        "snap3",
		CreatedAt: time.Now(),
		Secrets: map[string]map[string]interface{}{
			"secret/app/a": {"k": "v"},
			"secret/db/b":  {"k": "v"},
		},
	}
	_ = s.Save(snap)
	w := &fakeWriter{}
	r, _ := NewRestorer(s, w)
	records, err := r.Restore("snap3", RestoreOptions{PathPrefix: "secret/app"})
	if err != nil {
		t.Fatalf("Restore: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}
}

func TestRestoreWriterErrorPropagates(t *testing.T) {
	s := makeRestoreStore(t)
	snap := &Snapshot{
		ID:        "snap4",
		CreatedAt: time.Now(),
		Secrets: map[string]map[string]interface{}{
			"secret/x": {"k": "v"},
		},
	}
	_ = s.Save(snap)
	w := &fakeWriter{errOn: "secret/x"}
	r, _ := NewRestorer(s, w)
	_, err := r.Restore("snap4", RestoreOptions{})
	if err == nil {
		t.Fatal("expected error from writer")
	}
}
