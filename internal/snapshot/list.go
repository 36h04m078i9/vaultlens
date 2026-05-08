package snapshot

import (
	"sort"
	"strings"
	"time"
)

// Meta holds lightweight metadata about a stored snapshot.
type Meta struct {
	ID        string
	Path      string
	CreatedAt time.Time
	KeyCount  int
}

// ListOptions controls filtering and ordering of snapshot listings.
type ListOptions struct {
	// Prefix filters snapshots whose Path starts with the given string.
	Prefix string
	// Since excludes snapshots created before this time (zero means no filter).
	Since time.Time
	// Limit caps the number of results returned (0 means no limit).
	Limit int
}

// List returns metadata for all snapshots in the store that match opts.
// Results are ordered by CreatedAt descending (newest first).
func List(store *Store, opts ListOptions) []Meta {
	if store == nil {
		return nil
	}

	store.mu.RLock()
	defer store.mu.RUnlock()

	var out []Meta
	for id, snap := range store.data {
		if opts.Prefix != "" && !strings.HasPrefix(snap.Path, opts.Prefix) {
			continue
		}
		if !opts.Since.IsZero() && snap.CapturedAt.Before(opts.Since) {
			continue
		}
		out = append(out, Meta{
			ID:        id,
			Path:      snap.Path,
			CreatedAt: snap.CapturedAt,
			KeyCount:  len(snap.Secrets),
		})
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})

	if opts.Limit > 0 && len(out) > opts.Limit {
		out = out[:opts.Limit]
	}
	return out
}
