package snapshot

import (
	"fmt"
	"testing"
	"time"
)

// buildPruneStore creates a Store populated with n snapshots spaced one hour apart.
func buildPruneStore(t *testing.T, n int) (*Store, time.Time) {
	t.Helper()
	store := NewStore()
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < n; i++ {
		s := &Snapshot{
			ID:        fmt.Sprintf("snap-%02d", i),
			CreatedAt: base.Add(time.Duration(i) * time.Hour),
			Secrets:   map[string]string{},
		}
		if err := store.Save(s); err != nil {
			t.Fatalf("Save: %v", err)
		}
	}
	return store, base
}

func TestPruneNilStoreReturnsError(t *testing.T) {
	_, err := Prune(nil, PruneOptions{})
	if err == nil {
		t.Fatal("expected error for nil store")
	}
}

func TestPruneNoOptionsKeepsAll(t *testing.T) {
	store, _ := buildPruneStore(t, 5)
	res, err := Prune(store, PruneOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Removed) != 0 {
		t.Errorf("expected 0 removed, got %d", len(res.Removed))
	}
	if res.Retained != 5 {
		t.Errorf("expected 5 retained, got %d", res.Retained)
	}
}

func TestPruneKeepLast(t *testing.T) {
	store, _ := buildPruneStore(t, 6)
	res, err := Prune(store, PruneOptions{KeepLast: 2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Removed) != 4 {
		t.Errorf("expected 4 removed, got %d", len(res.Removed))
	}
	if res.Retained != 2 {
		t.Errorf("expected 2 retained, got %d", res.Retained)
	}
}

func TestPruneOlderThan(t *testing.T) {
	store, base := buildPruneStore(t, 5)
	// Remove snapshots created before hour 3 (snap-00, snap-01, snap-02).
	cutoff := base.Add(3 * time.Hour)
	res, err := Prune(store, PruneOptions{OlderThan: cutoff})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Removed) != 3 {
		t.Errorf("expected 3 removed, got %d: %v", len(res.Removed), res.Removed)
	}
}

func TestPrunePrefixFilters(t *testing.T) {
	store := NewStore()
	base := time.Now()
	for _, id := range []string{"prod-01", "prod-02", "staging-01"} {
		_ = store.Save(&Snapshot{ID: id, CreatedAt: base, Secrets: map[string]string{}})
		base = base.Add(time.Hour)
	}
	res, err := Prune(store, PruneOptions{KeepLast: 1, Prefix: "prod"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Only prod snapshots are candidates; keep newest 1, remove 1.
	if len(res.Removed) != 1 {
		t.Errorf("expected 1 removed, got %d", len(res.Removed))
	}
	// staging-01 must still exist.
	if _, ok := store.Get("staging-01"); !ok {
		t.Error("staging-01 should not have been pruned")
	}
}
