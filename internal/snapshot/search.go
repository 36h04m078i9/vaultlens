package snapshot

import (
	"strings"
	"time"
)

// SearchOptions controls how snapshots are filtered during a search.
type SearchOptions struct {
	// Query is matched against snapshot ID, prefix, and tags (case-insensitive).
	Query string
	// Since restricts results to snapshots taken at or after this time.
	Since time.Time
	// Tags filters to snapshots that contain ALL of the specified tags.
	Tags []string
	// Limit caps the number of results returned. 0 means no limit.
	Limit int
}

// SearchResult wraps a Snapshot with a relevance score.
type SearchResult struct {
	Snapshot *Snapshot
	// Score indicates how well the snapshot matched the query.
	// Higher is better.
	Score int
}

// Search filters and scores snapshots from the given store according to opts.
// Results are returned in descending score order, newest-first for ties.
func Search(store *Store, opts SearchOptions) ([]SearchResult, error) {
	if store == nil {
		return nil, ErrNilStore
	}

	all, err := store.List()
	if err != nil {
		return nil, err
	}

	var results []SearchResult
	q := strings.ToLower(strings.TrimSpace(opts.Query))

	for i := range all {
		snap := &all[i]

		if !opts.Since.IsZero() && snap.TakenAt.Before(opts.Since) {
			continue
		}

		if len(opts.Tags) > 0 && !hasAllTags(snap, opts.Tags) {
			continue
		}

		score := scoreSnapshot(snap, q)
		if q != "" && score == 0 {
			continue
		}

		results = append(results, SearchResult{Snapshot: snap, Score: score})
	}

	// Sort: higher score first, then newer first.
	sortSearchResults(results)

	if opts.Limit > 0 && len(results) > opts.Limit {
		results = results[:opts.Limit]
	}

	return results, nil
}

func scoreSnapshot(snap *Snapshot, q string) int {
	if q == "" {
		return 1
	}
	score := 0
	if strings.Contains(strings.ToLower(snap.ID), q) {
		score += 10
	}
	if strings.Contains(strings.ToLower(snap.Prefix), q) {
		score += 8
	}
	for _, tag := range snap.Tags {
		if strings.Contains(strings.ToLower(tag), q) {
			score += 5
			break
		}
	}
	return score
}

func hasAllTags(snap *Snapshot, required []string) bool {
	tagSet := make(map[string]struct{}, len(snap.Tags))
	for _, t := range snap.Tags {
		tagSet[strings.ToLower(t)] = struct{}{}
	}
	for _, r := range required {
		if _, ok := tagSet[strings.ToLower(r)]; !ok {
			return false
		}
	}
	return true
}

func sortSearchResults(results []SearchResult) {
	// Insertion sort is fine for the typical small result sets here.
	for i := 1; i < len(results); i++ {
		for j := i; j > 0; j-- {
			a, b := results[j-1], results[j]
			if a.Score < b.Score || (a.Score == b.Score && a.Snapshot.TakenAt.Before(b.Snapshot.TakenAt)) {
				results[j-1], results[j] = results[j], results[j-1]
			} else {
				break
			}
		}
	}
}
