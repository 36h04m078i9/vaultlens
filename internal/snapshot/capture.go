package snapshot

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// PathLister is the minimal interface required to enumerate Vault paths.
type PathLister interface {
	ListSecrets(path string) ([]string, error)
}

// Capture walks the given root path via lister, records every discovered
// path, and stores the resulting Snapshot in store under a generated ID.
func Capture(label, root string, lister PathLister, store *Store) (*Snapshot, error) {
	if lister == nil {
		return nil, fmt.Errorf("snapshot: lister must not be nil")
	}
	if store == nil {
		return nil, fmt.Errorf("snapshot: store must not be nil")
	}

	paths, err := walk(root, lister)
	if err != nil {
		return nil, fmt.Errorf("snapshot: walk failed: %w", err)
	}

	snap := &Snapshot{
		ID:         uuid.NewString(),
		Label:      label,
		CapturedAt: time.Now().UTC(),
		Paths:      paths,
	}

	if err := store.Save(snap); err != nil {
		return nil, err
	}
	return snap, nil
}

// walk recursively lists all leaf paths under root.
func walk(root string, lister PathLister) ([]string, error) {
	entries, err := lister.ListSecrets(root)
	if err != nil {
		return nil, err
	}

	var results []string
	for _, entry := range entries {
		full := root + entry
		if len(entry) > 0 && entry[len(entry)-1] == '/' {
			sub, err := walk(full, lister)
			if err != nil {
				return nil, err
			}
			results = append(results, sub...)
		} else {
			results = append(results, full)
		}
	}
	return results, nil
}
