package snapshot

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// ImportOptions controls how snapshots are imported.
type ImportOptions struct {
	// OverwriteExisting replaces a snapshot if one with the same ID already exists.
	OverwriteExisting bool
	// TagsToAdd are additional tags applied to every imported snapshot.
	TagsToAdd []string
}

// ImportResult describes the outcome of a single snapshot import.
type ImportResult struct {
	ID        string
	Imported  bool
	Skipped   bool
	Overwrite bool
	Err       error
}

// Import reads newline-delimited JSON snapshots from r and stores each one
// using store. Results are returned in the same order as the input records.
func Import(r io.Reader, store *Store, opts ImportOptions) ([]ImportResult, error) {
	if r == nil {
		return nil, fmt.Errorf("import: reader must not be nil")
	}
	if store == nil {
		return nil, fmt.Errorf("import: store must not be nil")
	}

	var results []ImportResult
	dec := json.NewDecoder(r)

	for dec.More() {
		var snap Snapshot
		if err := dec.Decode(&snap); err != nil {
			return results, fmt.Errorf("import: decode error: %w", err)
		}

		if snap.ID == "" {
			return results, fmt.Errorf("import: snapshot missing ID")
		}

		result := ImportResult{ID: snap.ID}

		_, exists := store.Get(snap.ID)
		if exists && !opts.OverwriteExisting {
			result.Skipped = true
			results = append(results, result)
			continue
		}

		if exists {
			result.Overwrite = true
		}

		for _, tag := range opts.TagsToAdd {
			snap.Tags = appendUnique(snap.Tags, tag)
		}

		if snap.CapturedAt.IsZero() {
			snap.CapturedAt = time.Now().UTC()
		}

		if err := store.Save(&snap); err != nil {
			result.Err = err
			results = append(results, result)
			continue
		}

		result.Imported = true
		results = append(results, result)
	}

	return results, nil
}

func appendUnique(slice []string, val string) []string {
	for _, v := range slice {
		if v == val {
			return slice
		}
	}
	return append(slice, val)
}
