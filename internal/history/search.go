package history

import (
	"strings"

	"github.com/your-org/vaultlens/internal/search"
)

// SearchResult holds a history entry paired with its relevance score.
type SearchResult struct {
	Entry Entry
	Score int
}

// Search returns history entries that match the given query using fuzzy
// matching against the path field. Results are ordered by score descending,
// preserving recency order for equal scores.
func Search(h *History, query string) []SearchResult {
	h.mu.RLock()
	entries := make([]Entry, len(h.entries))
	copy(entries, h.entries)
	h.mu.RUnlock()

	if strings.TrimSpace(query) == "" {
		results := make([]SearchResult, len(entries))
		for i, e := range entries {
			results[i] = SearchResult{Entry: e, Score: 0}
		}
		return results
	}

	paths := make([]string, len(entries))
	for i, e := range entries {
		paths[i] = e.Path
	}

	matched := search.Fuzzy(query, paths)

	// Build a score index keyed by path for O(1) lookup.
	scoreByPath := make(map[string]int, len(matched))
	for _, m := range matched {
		scoreByPath[m.Path] = m.Score
	}

	var results []SearchResult
	for _, e := range entries {
		if score, ok := scoreByPath[e.Path]; ok {
			results = append(results, SearchResult{Entry: e, Score: score})
		}
	}

	// Stable sort: higher score first; equal scores preserve recency order.
	for i := 1; i < len(results); i++ {
		for j := i; j > 0 && results[j].Score > results[j-1].Score; j-- {
			results[j], results[j-1] = results[j-1], results[j]
		}
	}

	return results
}
