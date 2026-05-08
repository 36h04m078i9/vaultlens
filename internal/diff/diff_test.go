package diff_test

import (
	"testing"

	"github.com/youorg/vaultlens/internal/diff"
)

func TestCompareNoChanges(t *testing.T) {
	before := map[string]string{"secret/a": "val1", "secret/b": "val2"}
	after := map[string]string{"secret/a": "val1", "secret/b": "val2"}
	changes := diff.Compare(before, after)
	if len(changes) != 0 {
		t.Fatalf("expected 0 changes, got %d", len(changes))
	}
}

func TestCompareDetectsAdded(t *testing.T) {
	before := map[string]string{}
	after := map[string]string{"secret/new": "newval"}
	changes := diff.Compare(before, after)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Kind != diff.Added {
		t.Errorf("expected Added, got %s", changes[0].Kind)
	}
	if changes[0].NewValue != "newval" {
		t.Errorf("unexpected NewValue: %s", changes[0].NewValue)
	}
}

func TestCompareDetectsRemoved(t *testing.T) {
	before := map[string]string{"secret/gone": "oldval"}
	after := map[string]string{}
	changes := diff.Compare(before, after)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Kind != diff.Removed {
		t.Errorf("expected Removed, got %s", changes[0].Kind)
	}
	if changes[0].OldValue != "oldval" {
		t.Errorf("unexpected OldValue: %s", changes[0].OldValue)
	}
}

func TestCompareDetectsChanged(t *testing.T) {
	before := map[string]string{"secret/x": "v1"}
	after := map[string]string{"secret/x": "v2"}
	changes := diff.Compare(before, after)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Kind != diff.Changed {
		t.Errorf("expected Changed, got %s", changes[0].Kind)
	}
}

func TestCompareResultsAreSorted(t *testing.T) {
	before := map[string]string{"secret/z": "v", "secret/a": "v"}
	after := map[string]string{"secret/m": "v"}
	changes := diff.Compare(before, after)
	for i := 1; i < len(changes); i++ {
		if changes[i].Path < changes[i-1].Path {
			t.Errorf("results not sorted at index %d: %s < %s", i, changes[i].Path, changes[i-1].Path)
		}
	}
}

func TestSummary(t *testing.T) {
	changes := []diff.Change{
		{Kind: diff.Added},
		{Kind: diff.Added},
		{Kind: diff.Removed},
		{Kind: diff.Changed},
	}
	s := diff.Summary(changes)
	if s[diff.Added] != 2 {
		t.Errorf("expected 2 added, got %d", s[diff.Added])
	}
	if s[diff.Removed] != 1 {
		t.Errorf("expected 1 removed, got %d", s[diff.Removed])
	}
	if s[diff.Changed] != 1 {
		t.Errorf("expected 1 changed, got %d", s[diff.Changed])
	}
}
