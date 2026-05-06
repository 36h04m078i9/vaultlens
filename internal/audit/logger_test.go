package audit

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestLogWritesValidJSON(t *testing.T) {
	var buf bytes.Buffer
	l := NewLogger(&buf)

	if err := l.Log("read", "secret/data/foo", "alice", true, ""); err != nil {
		t.Fatalf("Log() error: %v", err)
	}

	var entry Entry
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("failed to parse log output as JSON: %v", err)
	}

	if entry.Operation != "read" {
		t.Errorf("expected operation=read, got %q", entry.Operation)
	}
	if entry.Path != "secret/data/foo" {
		t.Errorf("expected path=secret/data/foo, got %q", entry.Path)
	}
	if entry.User != "alice" {
		t.Errorf("expected user=alice, got %q", entry.User)
	}
	if !entry.Success {
		t.Error("expected success=true")
	}
	if entry.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
	if time.Since(entry.Timestamp) > 5*time.Second {
		t.Error("timestamp appears stale")
	}
}

func TestLogReadFailure(t *testing.T) {
	var buf bytes.Buffer
	l := NewLogger(&buf)

	if err := l.LogRead("secret/data/bar", "bob", false); err != nil {
		t.Fatalf("LogRead() error: %v", err)
	}

	var entry Entry
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("failed to parse log output: %v", err)
	}

	if entry.Success {
		t.Error("expected success=false")
	}
	if !strings.Contains(entry.Message, "failed") {
		t.Errorf("expected failure message, got %q", entry.Message)
	}
}

func TestLogListSuccess(t *testing.T) {
	var buf bytes.Buffer
	l := NewLogger(&buf)

	if err := l.LogList("secret/metadata/", "carol", true); err != nil {
		t.Fatalf("LogList() error: %v", err)
	}

	var entry Entry
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("failed to parse log output: %v", err)
	}

	if entry.Operation != "list" {
		t.Errorf("expected operation=list, got %q", entry.Operation)
	}
	if entry.Message != "" {
		t.Errorf("expected empty message on success, got %q", entry.Message)
	}
}

func TestNewLoggerDefaultsToStdout(t *testing.T) {
	// Just ensure no panic when nil writer is passed.
	l := NewLogger(nil)
	if l == nil {
		t.Fatal("expected non-nil logger")
	}
}
