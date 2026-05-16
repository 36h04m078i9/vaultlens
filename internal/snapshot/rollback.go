package snapshot

import (
	"errors"
	"fmt"	
	"time"
)

// RollbackRecord describes a single rollback operation.
type RollbackRecord struct {
	SnapshotID string
	Path       string
	RolledBack bool
	DryRun     bool
	At         time.Time
	Err        string
}

// RollbackOptions controls rollback behaviour.
type RollbackOptions struct {
	// Prefix restricts rollback to paths under this prefix.
	Prefix string
	// DryRun reports what would change without writing.
	DryRun bool
}

// Writer is the minimal interface needed to persist secret data.
type Writer interface {
	WriteSecret(path string, data map[string]interface{}) error
}

// Rollback restores all secrets from the given snapshot to the provided writer.
// It returns a slice of RollbackRecords describing each path processed.
func Rollback(store Store, writer Writer, id string, opts RollbackOptions) ([]RollbackRecord, error) {
	if store == nil {
		return nil, errors.New("rollback: store must not be nil")
	}
	if writer == nil {
		return nil, errors.New("rollback: writer must not be nil")
	}
	if id == "" {
		return nil, errors.New("rollback: snapshot id must not be empty")
	}

	snap, ok := store.Get(id)
	if !ok {
		return nil, fmt.Errorf("rollback: snapshot %q not found", id)
	}

	var records []RollbackRecord
	for path, data := range snap.Data {
		if opts.Prefix != "" && !hasPrefix(path, opts.Prefix) {
			continue
		}
		rec := RollbackRecord{
			SnapshotID: id,
			Path:       path,
			DryRun:     opts.DryRun,
			At:         time.Now().UTC(),
		}
		if !opts.DryRun {
			if err := writer.WriteSecret(path, data); err != nil {
				rec.Err = err.Error()
			} else {
				rec.RolledBack = true
			}
		} else {
			rec.RolledBack = false
		}
		records = append(records, rec)
	}
	return records, nil
}
