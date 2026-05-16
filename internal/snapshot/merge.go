package snapshot

import (
	"errors"
	"fmt"
	"time"
)

// MergeResult holds the outcome of merging two snapshots.
type MergeResult struct {
	ID        string
	CreatedAt time.Time
	SourceIDs []string
	Paths     map[string]map[string]string
	Conflicts []MergeConflict
}

// MergeConflict describes a key that existed in both snapshots with different values.
type MergeConflict struct {
	Path  string
	Key   string
	ValueA string
	ValueB string
}

// MergeOptions controls how conflicts are resolved during a merge.
type MergeOptions struct {
	// PreferA keeps values from snapshot A on conflict; otherwise B wins.
	PreferA bool
	// NewID is the ID assigned to the merged snapshot. Defaults to "merged-<A>-<B>".
	NewID string
}

// Merge combines two snapshots from the store into a new snapshot.
// Conflicting keys are recorded in MergeResult.Conflicts and resolved
// according to MergeOptions.
func Merge(store *Store, idA, idB string, opts MergeOptions) (*MergeResult, error) {
	if store == nil {
		return nil, errors.New("merge: store must not be nil")
	}
	if idA == "" || idB == "" {
		return nil, errors.New("merge: source snapshot IDs must not be empty")
	}

	snapA, okA := store.Get(idA)
	if !okA {
		return nil, fmt.Errorf("merge: snapshot %q not found", idA)
	}
	snapB, okB := store.Get(idB)
	if !okB {
		return nil, fmt.Errorf("merge: snapshot %q not found", idB)
	}

	newID := opts.NewID
	if newID == "" {
		newID = fmt.Sprintf("merged-%s-%s", idA, idB)
	}

	merged := make(map[string]map[string]string)
	var conflicts []MergeConflict

	for path, data := range snapA.Data {
		merged[path] = copyMap(data)
	}

	for path, dataB := range snapB.Data {
		dataA, exists := merged[path]
		if !exists {
			merged[path] = copyMap(dataB)
			continue
		}
		for k, vB := range dataB {
			vA, keyExists := dataA[k]
			if !keyExists {
				dataA[k] = vB
				continue
			}
			if vA != vB {
				conflicts = append(conflicts, MergeConflict{Path: path, Key: k, ValueA: vA, ValueB: vB})
				if !opts.PreferA {
					dataA[k] = vB
				}
			}
		}
	}

	now := time.Now().UTC()
	newSnap := &Snapshot{
		ID:        newID,
		CreatedAt: now,
		Data:      merged,
	}
	if err := store.Save(newSnap); err != nil {
		return nil, fmt.Errorf("merge: saving merged snapshot: %w", err)
	}

	return &MergeResult{
		ID:        newID,
		CreatedAt: now,
		SourceIDs: []string{idA, idB},
		Paths:     merged,
		Conflicts: conflicts,
	}, nil
}

func copyMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
