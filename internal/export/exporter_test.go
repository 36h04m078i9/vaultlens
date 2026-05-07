package export_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/vaultlens/internal/export"
)

var fixedTime = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

func sampleRecords() []export.Record {
	return []export.Record{
		{Path: "secret/prod/db", Note: "main db", ExportedAt: fixedTime},
		{Path: "secret/prod/api", Note: "", ExportedAt: fixedTime},
	}
}

func TestExportJSON(t *testing.T) {
	var buf bytes.Buffer
	exp := export.NewExporter(&buf, export.FormatJSON)
	if err := exp.Export(sampleRecords()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var records []export.Record
	if err := json.Unmarshal(buf.Bytes(), &records); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if len(records) != 2 {
		t.Errorf("expected 2 records, got %d", len(records))
	}
	if records[0].Path != "secret/prod/db" {
		t.Errorf("unexpected path: %s", records[0].Path)
	}
}

func TestExportCSV(t *testing.T) {
	var buf bytes.Buffer
	exp := export.NewExporter(&buf, export.FormatCSV)
	if err := exp.Export(sampleRecords()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 { // header + 2 rows
		t.Errorf("expected 3 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "path,note") {
		t.Errorf("expected CSV header, got: %s", lines[0])
	}
	if !strings.Contains(lines[1], "secret/prod/db") {
		t.Errorf("expected path in row, got: %s", lines[1])
	}
}

func TestExportCSVContainsTimestamp(t *testing.T) {
	var buf bytes.Buffer
	exp := export.NewExporter(&buf, export.FormatCSV)
	_ = exp.Export(sampleRecords())
	if !strings.Contains(buf.String(), "2024-06-01T12:00:00Z") {
		t.Errorf("expected RFC3339 timestamp in output")
	}
}

func TestExportUnsupportedFormat(t *testing.T) {
	var buf bytes.Buffer
	exp := export.NewExporter(&buf, export.Format("xml"))
	if err := exp.Export(sampleRecords()); err == nil {
		t.Error("expected error for unsupported format, got nil")
	}
}

func TestExportEmptyRecords(t *testing.T) {
	var buf bytes.Buffer
	exp := export.NewExporter(&buf, export.FormatJSON)
	if err := exp.Export([]export.Record{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var records []export.Record
	if err := json.Unmarshal(buf.Bytes(), &records); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if len(records) != 0 {
		t.Errorf("expected empty slice, got %d records", len(records))
	}
}
