package snapshot_test

import (
	"testing"

	"github.com/your-org/vaultlens/internal/snapshot"
)

func makePinStore(t *testing.T) (*snapshot.PinStore, snapshot.Store) {
	t.Helper()
	s := snapshot.NewMemStore()
	ps, err := snapshot.NewPinStore(s)
	if err != nil {
		t.Fatalf("NewPinStore: %v", err)
	}
	return ps, s
}

func seedPinSnap(t *testing.T, s snapshot.Store, id string) {
	t.Helper()
	snap := &snapshot.Snapshot{ID: id, Prefix: "secret/"}
	if err := s.Save(snap); err != nil {
		t.Fatalf("Save: %v", err)
	}
}

func TestNewPinStoreNilStoreReturnsError(t *testing.T) {
	_, err := snapshot.NewPinStore(nil)
	if err == nil {
		t.Fatal("expected error for nil store")
	}
}

func TestPinSuccessAndIsPinned(t *testing.T) {
	ps, s := makePinStore(t)
	seedPinSnap(t, s, "snap-1")

	if err := ps.Pin("snap-1", "do not prune"); err != nil {
		t.Fatalf("Pin: %v", err)
	}
	if !ps.IsPinned("snap-1") {
		t.Fatal("expected snap-1 to be pinned")
	}
}

func TestPinEmptyIDReturnsError(t *testing.T) {
	ps, _ := makePinStore(t)
	if err := ps.Pin("", "reason"); err == nil {
		t.Fatal("expected error for empty id")
	}
}

func TestPinMissingSnapshotReturnsError(t *testing.T) {
	ps, _ := makePinStore(t)
	if err := ps.Pin("nonexistent", "reason"); err == nil {
		t.Fatal("expected error for missing snapshot")
	}
}

func TestUnpinExistingReturnsTrueAndRemovesPin(t *testing.T) {
	ps, s := makePinStore(t)
	seedPinSnap(t, s, "snap-2")
	_ = ps.Pin("snap-2", "keep")

	if !ps.Unpin("snap-2") {
		t.Fatal("expected Unpin to return true")
	}
	if ps.IsPinned("snap-2") {
		t.Fatal("expected snap-2 to be unpinned")
	}
}

func TestUnpinMissingReturnsFalse(t *testing.T) {
	ps, _ := makePinStore(t)
	if ps.Unpin("does-not-exist") {
		t.Fatal("expected Unpin to return false for unknown id")
	}
}

func TestListReturnsPinnedSnapshots(t *testing.T) {
	ps, s := makePinStore(t)
	seedPinSnap(t, s, "snap-a")
	seedPinSnap(t, s, "snap-b")
	_ = ps.Pin("snap-a", "first")
	_ = ps.Pin("snap-b", "second")

	pins := ps.List()
	if len(pins) != 2 {
		t.Fatalf("expected 2 pins, got %d", len(pins))
	}
}

func TestListReturnsEmptyWhenNoPins(t *testing.T) {
	ps, _ := makePinStore(t)
	if len(ps.List()) != 0 {
		t.Fatal("expected empty list")
	}
}

func TestPinnedSnapshotHasReason(t *testing.T) {
	ps, s := makePinStore(t)
	seedPinSnap(t, s, "snap-r")
	_ = ps.Pin("snap-r", "critical baseline")

	pins := ps.List()
	if len(pins) != 1 {
		t.Fatalf("expected 1 pin")
	}
	if pins[0].Reason != "critical baseline" {
		t.Errorf("expected reason 'critical baseline', got %q", pins[0].Reason)
	}
}
