package snapshot_test

import (
	"testing"
	"time"

	"github.com/youorg/vaultlens/internal/snapshot"
)

func makeDiffStore(t *testing.T) *snapshot.Store {
	t.Helper()
	return snapshot.NewStore()
}

func TestDiffSnapshotsDetectsAdded(t *testing.T) {
	store := makeDiffStore(t)

	snapA := &snapshot.Snapshot{ID: "a", CapturedAt: time.Now(), Secrets: map[string]string{
		"secret/foo": "bar",
	}}
	snapB := &snapshot.Snapshot{ID: "b", CapturedAt: time.Now(), Secrets: map[string]string{
		"secret/foo": "bar",
		"secret/new": "value",
	}}

	_ = store.Save(snapA)
	_ = store.Save(snapB)

	results, err := snapshot.DiffSnapshots(store, "a", "b")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found := false
	for _, r := range results {
		if r.Key == "secret/new" && r.Type == "added" {
			found = true
		}
	}
	if !found {
		t.Error("expected added entry for secret/new")
	}
}

func TestDiffSnapshotsMissingA(t *testing.T) {
	store := makeDiffStore(t)

	snapB := &snapshot.Snapshot{ID: "b", CapturedAt: time.Now(), Secrets: map[string]string{}}
	_ = store.Save(snapB)

	_, err := snapshot.DiffSnapshots(store, "missing", "b")
	if err == nil {
		t.Error("expected error for missing snapshot A")
	}
}

func TestDiffSnapshotsMissingB(t *testing.T) {
	store := makeDiffStore(t)

	snapA := &snapshot.Snapshot{ID: "a", CapturedAt: time.Now(), Secrets: map[string]string{}}
	_ = store.Save(snapA)

	_, err := snapshot.DiffSnapshots(store, "a", "missing")
	if err == nil {
		t.Error("expected error for missing snapshot B")
	}
}

func TestDiffSummaryNoChanges(t *testing.T) {
	store := makeDiffStore(t)

	secrets := map[string]string{"secret/foo": "bar"}
	snapA := &snapshot.Snapshot{ID: "a", CapturedAt: time.Now(), Secrets: secrets}
	snapB := &snapshot.Snapshot{ID: "b", CapturedAt: time.Now(), Secrets: secrets}

	_ = store.Save(snapA)
	_ = store.Save(snapB)

	summary, err := snapshot.DiffSummary(store, "a", "b")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if summary == "" {
		t.Error("expected non-empty summary")
	}
}
