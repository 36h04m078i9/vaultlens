package snapshot

import (
	"strings"
	"testing"
	"time"
)

func baseRollbackRecord(path string) RollbackRecord {
	return RollbackRecord{
		SnapshotID: "snap1",
		Path:       path,
		RolledBack: true,
		DryRun:     false,
		At:         time.Now().UTC(),
	}
}

func TestFormatRollbackRecordsEmpty(t *testing.T) {
	out := FormatRollbackRecords(nil, false)
	if !strings.Contains(out, "no paths matched") {
		t.Errorf("expected empty message, got: %s", out)
	}
}

func TestFormatRollbackRecordsContainsPath(t *testing.T) {
	rec := baseRollbackRecord("secret/app/db")
	out := FormatRollbackRecords([]RollbackRecord{rec}, false)
	if !strings.Contains(out, "secret/app/db") {
		t.Errorf("expected path in output, got: %s", out)
	}
}

func TestFormatRollbackRecordsDryRunPrefix(t *testing.T) {
	rec := baseRollbackRecord("secret/app/token")
	rec.DryRun = true
	rec.RolledBack = false
	out := FormatRollbackRecords([]RollbackRecord{rec}, false)
	if !strings.Contains(out, "dry-run") {
		t.Errorf("expected dry-run label, got: %s", out)
	}
}

func TestFormatRollbackRecordsShowsErrorWhenEnabled(t *testing.T) {
	rec := baseRollbackRecord("secret/broken")
	rec.RolledBack = false
	rec.Err = "permission denied"
	out := FormatRollbackRecords([]RollbackRecord{rec}, true)
	if !strings.Contains(out, "permission denied") {
		t.Errorf("expected error in output, got: %s", out)
	}
}

func TestFormatRollbackRecordsHidesErrorWhenDisabled(t *testing.T) {
	rec := baseRollbackRecord("secret/broken")
	rec.RolledBack = false
	rec.Err = "permission denied"
	out := FormatRollbackRecords([]RollbackRecord{rec}, false)
	if strings.Contains(out, "permission denied") {
		t.Errorf("error should be hidden when showErrors=false, got: %s", out)
	}
}

func TestFormatRollbackRecordsSummaryLine(t *testing.T) {
	recs := []RollbackRecord{
		baseRollbackRecord("secret/a"),
		baseRollbackRecord("secret/b"),
	}
	out := FormatRollbackRecords(recs, false)
	if !strings.Contains(out, "summary:") {
		t.Errorf("expected summary line, got: %s", out)
	}
	if !strings.Contains(out, "2 restored") {
		t.Errorf("expected '2 restored' in summary, got: %s", out)
	}
}
