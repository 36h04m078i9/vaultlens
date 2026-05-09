package snapshot

import (
	"fmt"
	"sort"
	"time"
)

// CompareOptions controls how two snapshots are compared.
type CompareOptions struct {
	// MaskValues hides secret values in the output.
	MaskValues bool
}

// CompareResult holds the diff between two named snapshots.
type CompareResult struct {
	SnapshotA  string
	SnapshotB  string
	TakenA     time.Time
	TakenB     time.Time
	Added      []string
	Removed    []string
	Changed    []string
	Unchanged  int
}

// Summary returns a human-readable one-line description of the diff.
func (r *CompareResult) Summary() string {
	return fmt.Sprintf(
		"%s → %s: +%d -%d ~%d =%d",
		r.SnapshotA, r.SnapshotB,
		len(r.Added), len(r.Removed), len(r.Changed), r.Unchanged,
	)
}

// Compare loads two snapshots by ID from store and returns a CompareResult.
func Compare(store *Store, idA, idB string, opts CompareOptions) (*CompareResult, error) {
	snA, okA := store.Get(idA)
	if !okA {
		return nil, fmt.Errorf("snapshot not found: %s", idA)
	}
	snB, okB := store.Get(idB)
	if !okB {
		return nil, fmt.Errorf("snapshot not found: %s", idB)
	}

	mapA := make(map[string]string, len(snA.Secrets))
	for _, s := range snA.Secrets {
		mapA[s.Path] = s.Value
	}
	mapB := make(map[string]string, len(snB.Secrets))
	for _, s := range snB.Secrets {
		mapB[s.Path] = s.Value
	}

	result := &CompareResult{
		SnapshotA: idA,
		SnapshotB: idB,
		TakenA:    snA.TakenAt,
		TakenB:    snB.TakenAt,
	}

	for path, valA := range mapA {
		if valB, exists := mapB[path]; !exists {
			result.Removed = append(result.Removed, path)
		} else if valA != valB {
			result.Changed = append(result.Changed, path)
		} else {
			result.Unchanged++
		}
	}
	for path := range mapB {
		if _, exists := mapA[path]; !exists {
			result.Added = append(result.Added, path)
		}
	}

	sort.Strings(result.Added)
	sort.Strings(result.Removed)
	sort.Strings(result.Changed)
	return result, nil
}
