package snapshot

import (
	"errors"
	"testing"
	"time"
)

// mockWriter records WriteSecret calls and optionally returns an error.
type mockWriter struct {
	written map[string]map[string]interface{}
	failOn  string
}

func newMockWriter() *mockWriter {
	return &mockWriter{written: make(map[string]map[string]interface{})}
}

func (m *mockWriter) WriteSecret(path string, data map[string]interface{}) error {
	if m.failOn == path {
		return errors.New("write error")
	}
	m.written[path] = data
	return nil
}

func makeRollbackStore() Store {
	s, _ := NewStore()
	return s
}

func seedRollbackSnap(s Store, id string, data map[string]map[string]interface{}) {
	snap := &Snapshot{ID: id, CreatedAt: time.Now().UTC(), Data: data}
	_ = s.Save(snap)
}

func TestRollbackNilStoreReturnsError(t *testing.T) {
	_, err := Rollback(nil, newMockWriter(), "abc", RollbackOptions{})
	if err == nil {
		t.Fatal("expected error for nil store")
	}
}

func TestRollbackNilWriterReturnsError(t *testing.T) {
	s := makeRollbackStore()
	_, err := Rollback(s, nil, "abc", RollbackOptions{})
	if err == nil {
		t.Fatal("expected error for nil writer")
	}
}

func TestRollbackEmptyIDReturnsError(t *testing.T) {
	s := makeRollbackStore()
	_, err := Rollback(s, newMockWriter(), "", RollbackOptions{})
	if err == nil {
		t.Fatal("expected error for empty id")
	}
}

func TestRollbackMissingSnapshotReturnsError(t *testing.T) {
	s := makeRollbackStore()
	_, err := Rollback(s, newMockWriter(), "missing", RollbackOptions{})
	if err == nil {
		t.Fatal("expected error for missing snapshot")
	}
}

func TestRollbackWritesAllPaths(t *testing.T) {
	s := makeRollbackStore()
	data := map[string]map[string]interface{}{
		"secret/a": {"key": "val1"},
		"secret/b": {"key": "val2"},
	}
	seedRollbackSnap(s, "snap1", data)
	w := newMockWriter()
	recs, err := Rollback(s, w, "snap1", RollbackOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(recs) != 2 {
		t.Fatalf("expected 2 records, got %d", len(recs))
	}
	for _, r := range recs {
		if !r.RolledBack {
			t.Errorf("expected path %q to be rolled back", r.Path)
		}
	}
}

func TestRollbackDryRunDoesNotWrite(t *testing.T) {
	s := makeRollbackStore()
	seedRollbackSnap(s, "snap2", map[string]map[string]interface{}{
		"secret/x": {"k": "v"},
	})
	w := newMockWriter()
	recs, err := Rollback(s, w, "snap2", RollbackOptions{DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(w.written) != 0 {
		t.Error("dry-run should not write any secrets")
	}
	if recs[0].DryRun != true {
		t.Error("record should be marked as dry-run")
	}
}

func TestRollbackPrefixFiltersResults(t *testing.T) {
	s := makeRollbackStore()
	seedRollbackSnap(s, "snap3", map[string]map[string]interface{}{
		"secret/app/db": {"pass": "x"},
		"secret/infra/key": {"val": "y"},
	})
	w := newMockWriter()
	recs, err := Rollback(s, w, "snap3", RollbackOptions{Prefix: "secret/app"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(recs) != 1 {
		t.Fatalf("expected 1 record, got %d", len(recs))
	}
	if recs[0].Path != "secret/app/db" {
		t.Errorf("unexpected path: %s", recs[0].Path)
	}
}

func TestRollbackRecordsWriteError(t *testing.T) {
	s := makeRollbackStore()
	seedRollbackSnap(s, "snap4", map[string]map[string]interface{}{
		"secret/bad": {"k": "v"},
	})
	w := newMockWriter()
	w.failOn = "secret/bad"
	recs, err := Rollback(s, w, "snap4", RollbackOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if recs[0].Err == "" {
		t.Error("expected error to be recorded in RollbackRecord")
	}
	if recs[0].RolledBack {
		t.Error("failed path should not be marked as rolled back")
	}
}
