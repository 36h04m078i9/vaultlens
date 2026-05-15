package snapshot

import (
	"errors"
	"fmt"
	"time"
)

// CloneOptions controls how a snapshot is cloned.
type CloneOptions struct {
	// NewID is the identifier for the cloned snapshot.
	// If empty, a timestamped ID is generated.
	NewID string
	// ExtraTags are additional tags applied to the clone.
	ExtraTags []string
	// Note is an optional annotation added to the clone.
	Note string
}

// CloneResult describes the outcome of a clone operation.
type CloneResult struct {
	SourceID string
	CloneID  string
	CreatedAt time.Time
}

// Clone duplicates an existing snapshot under a new ID.
// The cloned snapshot shares the same secret paths and values as the source.
func Clone(store *Store, sourceID string, opts CloneOptions) (*CloneResult, error) {
	if store == nil {
		return nil, errors.New("snapshot: Clone requires a non-nil store")
	}
	if sourceID == "" {
		return nil, errors.New("snapshot: Clone requires a non-empty sourceID")
	}

	src, ok := store.Get(sourceID)
	if !ok {
		return nil, fmt.Errorf("snapshot: source snapshot %q not found", sourceID)
	}

	newID := opts.NewID
	if newID == "" {
		newID = fmt.Sprintf("%s-clone-%d", sourceID, time.Now().UnixNano())
	}

	clone := &Snapshot{
		ID:        newID,
		CreatedAt: time.Now().UTC(),
		Paths:     make(map[string]map[string]string, len(src.Paths)),
		Tags:      append([]string(nil), src.Tags...),
	}

	for path, secrets := range src.Paths {
		copied := make(map[string]string, len(secrets))
		for k, v := range secrets {
			copied[k] = v
		}
		clone.Paths[path] = copied
	}

	clone.Tags = append(clone.Tags, opts.ExtraTags...)

	if err := store.Save(clone); err != nil {
		return nil, fmt.Errorf("snapshot: failed to save clone: %w", err)
	}

	return &CloneResult{
		SourceID:  sourceID,
		CloneID:   newID,
		CreatedAt: clone.CreatedAt,
	}, nil
}
