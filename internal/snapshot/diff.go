package snapshot

import (
	"fmt"

	"github.com/youorg/vaultlens/internal/diff"
)

// DiffSnapshots compares two stored snapshots by ID and returns the diff results.
// Returns an error if either snapshot is not found.
func DiffSnapshots(store *Store, idA, idB string) ([]diff.Result, error) {
	snapA, ok := store.Get(idA)
	if !ok {
		return nil, fmt.Errorf("snapshot %q not found", idA)
	}

	snapB, ok := store.Get(idB)
	if !ok {
		return nil, fmt.Errorf("snapshot %q not found", idB)
	}

	return diff.Compare(snapA.Secrets, snapB.Secrets), nil
}

// DiffSummary returns a human-readable summary of changes between two snapshots.
func DiffSummary(store *Store, idA, idB string) (string, error) {
	results, err := DiffSnapshots(store, idA, idB)
	if err != nil {
		return "", err
	}

	return diff.Summary(results), nil
}
