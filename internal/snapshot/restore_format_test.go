package snapshot

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func baseRecord(path string) RestoreRecord {
	return RestoreRecord{
		Path:       path,
		Data:       map[string]interface{}{"username": "admin", "password": "s3cr3t"},
		RestoredAt: time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
	}
}

func TestFormatRestoreRecordsEmpty(t *testing.T) {
	var buf bytes.Buffer
	if err := FormatRestoreRecords(&buf, nil, RestoreFormatOptions{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No secrets restored") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}

func TestFormatRestoreRecordsContainsPath(t *testing.T) {
	var buf bytes.Buffer
	records := []RestoreRecord{baseRecord("secret/app/db")}
	_ = FormatRestoreRecords(&buf, records, RestoreFormatOptions{})
	if !strings.Contains(buf.String(), "secret/app/db") {
		t.Errorf("expected path in output, got: %s", buf.String())
	}
}

func TestFormatRestoreRecordsDryRunPrefix(t *testing.T) {
	var buf bytes.Buffer
	records := []RestoreRecord{baseRecord("secret/x")}
	_ = FormatRestoreRecords(&buf, records, RestoreFormatOptions{DryRun: true})
	if !strings.Contains(buf.String(), "[DRY-RUN]") {
		t.Errorf("expected [DRY-RUN] prefix, got: %s", buf.String())
	}
}

func TestFormatRestoreRecordsShowDataIncludesKeys(t *testing.T) {
	var buf bytes.Buffer
	records := []RestoreRecord{baseRecord("secret/y")}
	_ = FormatRestoreRecords(&buf, records, RestoreFormatOptions{ShowData: true})
	out := buf.String()
	if !strings.Contains(out, "username") && !strings.Contains(out, "password") {
		t.Errorf("expected data keys in output, got: %s", out)
	}
}

func TestFormatRestoreRecordsHidesDataByDefault(t *testing.T) {
	var buf bytes.Buffer
	records := []RestoreRecord{baseRecord("secret/z")}
	_ = FormatRestoreRecords(&buf, records, RestoreFormatOptions{})
	out := buf.String()
	if strings.Contains(out, "username") || strings.Contains(out, "password") {
		t.Errorf("expected keys hidden by default, got: %s", out)
	}
}

func TestFormatRestoreRecordsTimestamp(t *testing.T) {
	var buf bytes.Buffer
	records := []RestoreRecord{baseRecord("secret/ts")}
	_ = FormatRestoreRecords(&buf, records, RestoreFormatOptions{})
	if !strings.Contains(buf.String(), "2024-06-01T12:00:00Z") {
		t.Errorf("expected formatted timestamp, got: %s", buf.String())
	}
}
