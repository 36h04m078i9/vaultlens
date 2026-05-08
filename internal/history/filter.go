package history

import (
	"strings"
	"time"
)

// FilterOptions controls which history entries are returned.
type FilterOptions struct {
	// Since filters entries to those visited at or after this time.
	// Zero value means no lower bound.
	Since time.Time

	// Prefix filters entries whose path starts with the given string.
	// Empty string means no prefix filter.
	Prefix string

	// Limit caps the number of returned entries. 0 means no limit.
	Limit int
}

// Filter returns history entries that match all criteria in opts.
// Entries are returned in the same order as h.List() (most-recent first).
func Filter(h *History, opts FilterOptions) []Entry {
	all := h.List()
	result := make([]Entry, 0, len(all))

	for _, e := range all {
		if !opts.Since.IsZero() && e.VisitedAt.Before(opts.Since) {
			continue
		}
		if opts.Prefix != "" && !strings.HasPrefix(e.Path, opts.Prefix) {
			continue
		}
		result = append(result, e)
		if opts.Limit > 0 && len(result) >= opts.Limit {
			break
		}
	}

	return result
}
