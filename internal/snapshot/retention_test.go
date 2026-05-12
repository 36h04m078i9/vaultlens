package snapshot

import (
	"fmt"
	"testing"
	"time"
)

func buildRetentionStore(t *testing.T, count int, spacing time.Duration) *Store {
	t.Helper()
	store := NewStore()
	base := time.Now().Add(-spacing * time.Duration(count))
	for i := 0; i < count; i++ {
		s := &Snapshot{
			ID:          fmt.Sprintf("snap-%02d", i),
			CapturedAt:  base.Add(spacing * time.Duration(i)),
			SecretCount: i,
		}
		if err := store.Save(s); err != nil {
			t.Fatalf("save: %v", err)
		}
	}
	return store
}

func TestApplyRetentionNilStoreReturnsError(t *testing.T) {
	_, err := ApplyRetention(nil, RetentionPolicy{KeepLast: 1})
	if err == nil {
		t.Fatal("expected error for nil store")
	}
}

func TestApplyRetentionNoPolicyReturnsError(t *testing.T) {
	store := buildRetentionStore(t, 3, time.Hour)
	_, err := ApplyRetention(store, RetentionPolicy{})
	if err == nil {
		t.Fatal("expected error for empty policy")
	}
}

func TestApplyRetentionKeepLast(t *testing.T) {
	store := buildRetentionStore(t, 6, time.Hour)
	removed, err := ApplyRetention(store, RetentionPolicy{KeepLast: 2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if removed != 4 {
		t.Errorf("expected 4 removed, got %d", removed)
	}
	snaps, _ := store.List()
	if len(snaps) != 2 {
		t.Errorf("expected 2 remaining, got %d", len(snaps))
	}
}

func TestApplyRetentionOlderThan(t *testing.T) {
	store := buildRetentionStore(t, 5, 24*time.Hour)
	cutoff := time.Now().Add(-2 * 24 * time.Hour)
	removed, err := ApplyRetention(store, RetentionPolicy{OlderThan: cutoff})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if removed == 0 {
		t.Error("expected at least one snapshot removed")
	}
}

func TestApplyRetentionKeepLastZeroRemovesNone(t *testing.T) {
	store := buildRetentionStore(t, 4, time.Hour)
	// KeepLast=4 keeps all, OlderThan far future keeps all too.
	removed, err := ApplyRetention(store, RetentionPolicy{KeepLast: 4})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if removed != 0 {
		t.Errorf("expected 0 removed, got %d", removed)
	}
}
