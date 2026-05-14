package snapshot_test

import (
	"testing"
	"time"

	"github.com/your-org/vaultlens/internal/snapshot"
)

func TestPolicyValidateRequiresRule(t *testing.T) {
	p := snapshot.RetentionPolicy{}
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for empty policy")
	}
}

func TestPolicyValidateNegativeKeepLast(t *testing.T) {
	p := snapshot.RetentionPolicy{KeepLast: -1}
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for negative KeepLast")
	}
}

func TestPolicyValidateOlderThanOnly(t *testing.T) {
	p := snapshot.RetentionPolicy{OlderThan: time.Hour}
	if err := p.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPolicyStoreAddAndGet(t *testing.T) {
	s := snapshot.NewPolicyStore()
	p := snapshot.RetentionPolicy{KeepLast: 5}
	if err := s.Add("weekly", p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, ok := s.Get("weekly")
	if !ok {
		t.Fatal("expected policy to exist")
	}
	if got.KeepLast != 5 {
		t.Errorf("expected KeepLast=5, got %d", got.KeepLast)
	}
}

func TestPolicyStoreAddEmptyNameReturnsError(t *testing.T) {
	s := snapshot.NewPolicyStore()
	err := s.Add("", snapshot.RetentionPolicy{KeepLast: 1})
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestPolicyStoreAddInvalidPolicyReturnsError(t *testing.T) {
	s := snapshot.NewPolicyStore()
	err := s.Add("bad", snapshot.RetentionPolicy{})
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestPolicyStoreRemoveExisting(t *testing.T) {
	s := snapshot.NewPolicyStore()
	_ = s.Add("daily", snapshot.RetentionPolicy{KeepLast: 7})
	if !s.Remove("daily") {
		t.Fatal("expected Remove to return true")
	}
	_, ok := s.Get("daily")
	if ok {
		t.Fatal("expected policy to be deleted")
	}
}

func TestPolicyStoreRemoveMissingReturnsFalse(t *testing.T) {
	s := snapshot.NewPolicyStore()
	if s.Remove("ghost") {
		t.Fatal("expected false for missing policy")
	}
}

func TestPolicyStoreList(t *testing.T) {
	s := snapshot.NewPolicyStore()
	_ = s.Add("a", snapshot.RetentionPolicy{KeepLast: 1})
	_ = s.Add("b", snapshot.RetentionPolicy{KeepLast: 2})
	names := s.List()
	if len(names) != 2 {
		t.Errorf("expected 2 policies, got %d", len(names))
	}
}
