package snapshot

import (
	"testing"
	"time"
)

func makeExportStore(t *testing.T) *Store {
	t.Helper()
	st := NewStore()

	now := time.Now()

	st.Save(&Snapshot{
		ID:      "snap-a",
		Tag:     "baseline",
		TakenAt: now.Add(-2 * time.Hour),
		Secrets: map[string]map[string]string{
			"secret/db": {"password": "s3cr3t"},
		},
	})
	st.Save(&Snapshot{
		ID:      "snap-b",
		Tag:     "after-rotation",
		TakenAt: now.Add(-1 * time.Hour),
		Secrets: map[string]map[string]string{
			"secret/api": {"key": "abc123"},
			"secret/db":  {"password": "newpass"},
		},
	})
	return st
}

func TestToExportRecordsNilStoreReturnsError(t *testing.T) {
	_, err := ToExportRecords(nil, ExportOptions{})
	if err == nil {
		t.Fatal("expected error for nil store, got nil")
	}
}

func TestToExportRecordsAllSnapshots(t *testing.T) {
	st := makeExportStore(t)
	records, err := ToExportRecords(st, ExportOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// snap-a has 1 path, snap-b has 2 paths → 3 records total
	if len(records) != 3 {
		t.Errorf("expected 3 records, got %d", len(records))
	}
}

func TestToExportRecordsFilterByID(t *testing.T) {
	st := makeExportStore(t)
	records, err := ToExportRecords(st, ExportOptions{IDs: []string{"snap-a"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(records) != 1 {
		t.Errorf("expected 1 record, got %d", len(records))
	}
	if records[0].Path != "secret/db" {
		t.Errorf("unexpected path: %s", records[0].Path)
	}
}

func TestToExportRecordsFilterBySince(t *testing.T) {
	st := makeExportStore(t)
	since := time.Now().Add(-90 * time.Minute)
	records, err := ToExportRecords(st, ExportOptions{Since: since})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Only snap-b (taken 1h ago) should pass the 90-min filter
	if len(records) != 2 {
		t.Errorf("expected 2 records, got %d", len(records))
	}
}

func TestToExportRecordsIncludeSecrets(t *testing.T) {
	st := makeExportStore(t)
	records, err := ToExportRecords(st, ExportOptions{IDs: []string{"snap-a"}, IncludeSecrets: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("expected 1 record")
	}
	if records[0].Data == nil {
		t.Error("expected Data to be populated when IncludeSecrets is true")
	}
}

func TestToExportRecordsSnapshotIDInMeta(t *testing.T) {
	st := makeExportStore(t)
	records, err := ToExportRecords(st, ExportOptions{IDs: []string{"snap-b"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range records {
		if r.Meta["snapshot_id"] != "snap-b" {
			t.Errorf("expected snapshot_id=snap-b in meta, got %q", r.Meta["snapshot_id"])
		}
	}
}
