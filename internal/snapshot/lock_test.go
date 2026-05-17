package snapshot

import (
	"testing"
	"time"
)

func makeLockStore(t *testing.T) (*LockStore, *Store) {
	t.Helper()
	s := &Store{data: make(map[string]Snapshot)}
	ls, err := NewLockStore(s)
	if err != nil {
		t.Fatalf("NewLockStore: %v", err)
	}
	return ls, s
}

func seedLockSnap(t *testing.T, s *Store, id string) {
	t.Helper()
	snap := Snapshot{ID: id, CapturedAt: time.Now().UTC(), Paths: map[string]map[string]string{}}
	if err := s.Save(snap); err != nil {
		t.Fatalf("Save: %v", err)
	}
}

func TestNewLockStoreNilStoreReturnsError(t *testing.T) {
	_, err := NewLockStore(nil)
	if err == nil {
		t.Fatal("expected error for nil store")
	}
}

func TestAcquireLockSuccess(t *testing.T) {
	ls, s := makeLockStore(t)
	seedLockSnap(t, s, "snap-1")

	l, err := ls.Acquire("snap-1", "alice")
	if err != nil {
		t.Fatalf("Acquire: %v", err)
	}
	if l.Holder != "alice" {
		t.Errorf("expected holder alice, got %s", l.Holder)
	}
	if l.SnapshotID != "snap-1" {
		t.Errorf("expected snapshot id snap-1, got %s", l.SnapshotID)
	}
	if l.AcquiredAt.IsZero() {
		t.Error("expected non-zero AcquiredAt")
	}
}

func TestAcquireLockAlreadyLockedByOther(t *testing.T) {
	ls, s := makeLockStore(t)
	seedLockSnap(t, s, "snap-2")

	if _, err := ls.Acquire("snap-2", "alice"); err != nil {
		t.Fatalf("first Acquire: %v", err)
	}
	_, err := ls.Acquire("snap-2", "bob")
	if err != ErrAlreadyLocked {
		t.Errorf("expected ErrAlreadyLocked, got %v", err)
	}
}

func TestAcquireLockIdempotentForSameHolder(t *testing.T) {
	ls, s := makeLockStore(t)
	seedLockSnap(t, s, "snap-3")

	if _, err := ls.Acquire("snap-3", "alice"); err != nil {
		t.Fatalf("first Acquire: %v", err)
	}
	if _, err := ls.Acquire("snap-3", "alice"); err != nil {
		t.Errorf("re-acquire by same holder should succeed, got: %v", err)
	}
}

func TestAcquireMissingSnapshotReturnsError(t *testing.T) {
	ls, _ := makeLockStore(t)
	_, err := ls.Acquire("no-such-snap", "alice")
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestReleaseExistingLock(t *testing.T) {
	ls, s := makeLockStore(t)
	seedLockSnap(t, s, "snap-4")

	if _, err := ls.Acquire("snap-4", "alice"); err != nil {
		t.Fatalf("Acquire: %v", err)
	}
	if err := ls.Release("snap-4"); err != nil {
		t.Errorf("Release: %v", err)
	}
	if ls.IsLocked("snap-4") {
		t.Error("expected snapshot to be unlocked after Release")
	}
}

func TestReleaseNotLockedReturnsError(t *testing.T) {
	ls, _ := makeLockStore(t)
	if err := ls.Release("snap-x"); err != ErrNotLocked {
		t.Errorf("expected ErrNotLocked, got %v", err)
	}
}

func TestGetReturnsLock(t *testing.T) {
	ls, s := makeLockStore(t)
	seedLockSnap(t, s, "snap-5")

	if _, err := ls.Acquire("snap-5", "carol"); err != nil {
		t.Fatalf("Acquire: %v", err)
	}
	l, ok := ls.Get("snap-5")
	if !ok {
		t.Fatal("expected lock to be present")
	}
	if l.Holder != "carol" {
		t.Errorf("expected holder carol, got %s", l.Holder)
	}
}
