package snapshot

import (
	"strings"
	"testing"
	"time"
)

func makeCloneStore(t *testing.T) *Store {
	t.Helper()
	s, err := NewStore(t.TempDir())
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	return s
}

func seedCloneSnap(t *testing.T, s *Store, id string) *Snapshot {
	t.Helper()
	snap := &Snapshot{
		ID:        id,
		CreatedAt: time.Now().UTC(),
		Paths: map[string]map[string]string{
			"secret/app": {"key": "value"},
		},
		Tags: []string{"prod"},
	}
	if err := s.Save(snap); err != nil {
		t.Fatalf("Save: %v", err)
	}
	return snap
}

func TestCloneNilStoreReturnsError(t *testing.T) {
	_, err := Clone(nil, "abc", CloneOptions{})
	if err == nil {
		t.Fatal("expected error for nil store")
	}
}

func TestCloneEmptySourceIDReturnsError(t *testing.T) {
	s := makeCloneStore(t)
	_, err := Clone(s, "", CloneOptions{})
	if err == nil {
		t.Fatal("expected error for empty sourceID")
	}
}

func TestCloneMissingSourceReturnsError(t *testing.T) {
	s := makeCloneStore(t)
	_, err := Clone(s, "nonexistent", CloneOptions{})
	if err == nil {
		t.Fatal("expected error for missing source")
	}
}

func TestCloneCreatesNewSnapshot(t *testing.T) {
	s := makeCloneStore(t)
	seedCloneSnap(t, s, "snap-1")

	result, err := Clone(s, "snap-1", CloneOptions{NewID: "snap-1-copy"})
	if err != nil {
		t.Fatalf("Clone: %v", err)
	}
	if result.CloneID != "snap-1-copy" {
		t.Errorf("expected CloneID snap-1-copy, got %q", result.CloneID)
	}
	if result.SourceID != "snap-1" {
		t.Errorf("expected SourceID snap-1, got %q", result.SourceID)
	}

	clone, ok := s.Get("snap-1-copy")
	if !ok {
		t.Fatal("cloned snapshot not found in store")
	}
	if v := clone.Paths["secret/app"]["key"]; v != "value" {
		t.Errorf("expected path value 'value', got %q", v)
	}
}

func TestCloneInheritsAndExtendsTagss(t *testing.T) {
	s := makeCloneStore(t)
	seedCloneSnap(t, s, "snap-2")

	_, err := Clone(s, "snap-2", CloneOptions{
		NewID:     "snap-2-clone",
		ExtraTags: []string{"clone", "staging"},
	})
	if err != nil {
		t.Fatalf("Clone: %v", err)
	}

	clone, _ := s.Get("snap-2-clone")
	tagSet := map[string]bool{}
	for _, tg := range clone.Tags {
		tagSet[tg] = true
	}
	for _, want := range []string{"prod", "clone", "staging"} {
		if !tagSet[want] {
			t.Errorf("expected tag %q in clone", want)
		}
	}
}

func TestCloneGeneratesIDWhenEmpty(t *testing.T) {
	s := makeCloneStore(t)
	seedCloneSnap(t, s, "snap-3")

	result, err := Clone(s, "snap-3", CloneOptions{})
	if err != nil {
		t.Fatalf("Clone: %v", err)
	}
	if !strings.HasPrefix(result.CloneID, "snap-3-clone-") {
		t.Errorf("auto-generated ID should start with 'snap-3-clone-', got %q", result.CloneID)
	}
}

func TestCloneIsIndependentCopy(t *testing.T) {
	s := makeCloneStore(t)
	original := seedCloneSnap(t, s, "snap-4")

	_, err := Clone(s, "snap-4", CloneOptions{NewID: "snap-4-clone"})
	if err != nil {
		t.Fatalf("Clone: %v", err)
	}

	// Mutate original in memory — clone should be unaffected.
	original.Paths["secret/app"]["key"] = "mutated"

	clone, _ := s.Get("snap-4-clone")
	if v := clone.Paths["secret/app"]["key"]; v != "value" {
		t.Errorf("clone should be independent; got mutated value %q", v)
	}
}
