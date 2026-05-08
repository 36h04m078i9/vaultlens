package diff_test

import (
	"strings"
	"testing"

	"github.com/youorg/vaultlens/internal/diff"
)

func TestFormatNoChanges(t *testing.T) {
	f := diff.Formatter{}
	out := f.Format(nil)
	if out != "no changes detected" {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestFormatAdded(t *testing.T) {
	changes := []diff.Change{
		{Path: "secret/foo", Kind: diff.Added, NewValue: "bar"},
	}
	f := diff.Formatter{}
	out := f.Format(changes)
	if !strings.Contains(out, "+ secret/foo = bar") {
		t.Errorf("expected added line in output, got:\n%s", out)
	}
}

func TestFormatRemoved(t *testing.T) {
	changes := []diff.Change{
		{Path: "secret/old", Kind: diff.Removed, OldValue: "gone"},
	}
	f := diff.Formatter{}
	out := f.Format(changes)
	if !strings.Contains(out, "- secret/old = gone") {
		t.Errorf("expected removed line in output, got:\n%s", out)
	}
}

func TestFormatChanged(t *testing.T) {
	changes := []diff.Change{
		{Path: "secret/x", Kind: diff.Changed, OldValue: "v1", NewValue: "v2"},
	}
	f := diff.Formatter{}
	out := f.Format(changes)
	if !strings.Contains(out, "~ secret/x: v1 -> v2") {
		t.Errorf("expected changed line in output, got:\n%s", out)
	}
}

func TestFormatMasksValues(t *testing.T) {
	changes := []diff.Change{
		{Path: "secret/pw", Kind: diff.Added, NewValue: "supersecret"},
	}
	f := diff.Formatter{MaskValues: true}
	out := f.Format(changes)
	if strings.Contains(out, "supersecret") {
		t.Errorf("expected value to be redacted, got:\n%s", out)
	}
	if !strings.Contains(out, "[redacted]") {
		t.Errorf("expected [redacted] in output, got:\n%s", out)
	}
}

func TestFormatSummaryLine(t *testing.T) {
	changes := []diff.Change{
		{Kind: diff.Added},
		{Kind: diff.Removed},
		{Kind: diff.Changed},
	}
	f := diff.Formatter{}
	out := f.Format(changes)
	if !strings.Contains(out, "summary: +1 added, -1 removed, ~1 changed") {
		t.Errorf("expected summary line, got:\n%s", out)
	}
}
