package history

import (
	"testing"
	"time"
)

func makeHistoryWithTimes(paths []string, base time.Time) *History {
	h := New(50)
	for i, p := range paths {
		// Simulate ascending access times so index 0 is oldest.
		h.Push(p)
		// Manually adjust the stored time for deterministic tests.
		h.mu.Lock()
		for j := range h.entries {
			if h.entries[j].Path == p {
				h.entries[j].AccessedAt = base.Add(time.Duration(i) * time.Minute)
				break
			}
		}
		h.mu.Unlock()
	}
	return h
}

func TestToExportRecordsNoOptions(t *testing.T) {
	h := New(10)
	h.Push("secret/a")
	h.Push("secret/b")

	records := ToExportRecords(h, ExportOptions{})
	if len(records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(records))
	}
}

func TestToExportRecordsLimit(t *testing.T) {
	h := New(10)
	for _, p := range []string{"a", "b", "c", "d"} {
		h.Push(p)
	}

	records := ToExportRecords(h, ExportOptions{Limit: 2})
	if len(records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(records))
	}
}

func TestToExportRecordsSinceFilter(t *testing.T) {
	base := time.Now().Add(-10 * time.Minute)
	paths := []string{"old/path", "mid/path", "new/path"}
	h := makeHistoryWithTimes(paths, base)

	// Only entries accessed at or after base+1m should be included.
	cutoff := base.Add(1 * time.Minute)
	records := ToExportRecords(h, ExportOptions{Since: cutoff})
	if len(records) != 2 {
		t.Fatalf("expected 2 records after cutoff, got %d", len(records))
	}
}

func TestToExportRecordsPathPreserved(t *testing.T) {
	h := New(10)
	h.Push("secret/ops/token")

	records := ToExportRecords(h, ExportOptions{})
	if len(records) == 0 {
		t.Fatal("expected at least one record")
	}
	if records[0].Path != "secret/ops/token" {
		t.Errorf("expected path %q, got %q", "secret/ops/token", records[0].Path)
	}
}

func TestToExportRecordsTimestampSet(t *testing.T) {
	h := New(10)
	h.Push("secret/check")

	records := ToExportRecords(h, ExportOptions{})
	if records[0].Timestamp.IsZero() {
		t.Error("expected non-zero timestamp on exported record")
	}
}
