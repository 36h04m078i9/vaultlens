package snapshot

import (
	"fmt"
	"sort"
	"time"
)

// PruneOptions controls which snapshots are removed.
type PruneOptions struct {
	// KeepLast retains the N most recent snapshots. 0 means keep all.
	KeepLast int
	// OlderThan removes snapshots created before this time. Zero value is ignored.
	OlderThan time.Time
	// Prefix restricts pruning to snapshots whose ID starts with this prefix.
	Prefix string
}

// PruneResult summarises what was removed.
type PruneResult struct {
	Removed []string
	Retained int
}

// Prune deletes snapshots from store according to opts.
// It returns a PruneResult describing what was removed, or an error.
func Prune(store *Store, opts PruneOptions) (PruneResult, error) {
	if store == nil {
		return PruneResult{}, fmt.Errorf("snapshot: store must not be nil")
	}

	all, err := store.List()
	if err != nil {
		return PruneResult{}, fmt.Errorf("snapshot: list failed: %w", err)
	}

	// Filter by prefix when requested.
	var candidates []*Snapshot
	for _, s := range all {
		if opts.Prefix == "" || len(s.ID) >= len(opts.Prefix) && s.ID[:len(opts.Prefix)] == opts.Prefix {
			candidates = append(candidates, s)
		}
	}

	// Sort newest-first so KeepLast is straightforward.
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].CreatedAt.After(candidates[j].CreatedAt)
	})

	toDelete := make(map[string]bool)

	// Mark by age.
	if !opts.OlderThan.IsZero() {
		for _, s := range candidates {
			if s.CreatedAt.Before(opts.OlderThan) {
				toDelete[s.ID] = true
			}
		}
	}

	// Mark by count — everything beyond KeepLast.
	if opts.KeepLast > 0 && len(candidates) > opts.KeepLast {
		for _, s := range candidates[opts.KeepLast:] {
			toDelete[s.ID] = true
		}
	}

	var result PruneResult
	for id := range toDelete {
		if err := store.Delete(id); err != nil {
			return result, fmt.Errorf("snapshot: delete %q failed: %w", id, err)
		}
		result.Removed = append(result.Removed, id)
	}
	sort.Strings(result.Removed)
	result.Retained = len(candidates) - len(result.Removed)
	return result, nil
}
