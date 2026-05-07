package bookmark

import "strings"

// SearchResult pairs a bookmark with a relevance score.
type SearchResult struct {
	Bookmark Bookmark
	Score    int
}

// Search filters bookmarks whose path or note contains query (case-insensitive).
// Results with path matches are ranked higher than note-only matches.
func Search(items []Bookmark, query string) []SearchResult {
	if query == "" {
		results := make([]SearchResult, len(items))
		for i, b := range items {
			results[i] = SearchResult{Bookmark: b, Score: 0}
		}
		return results
	}

	lower := strings.ToLower(query)
	var results []SearchResult
	for _, b := range items {
		score := 0
		if strings.Contains(strings.ToLower(b.Path), lower) {
			score += 2
		}
		if strings.Contains(strings.ToLower(b.Note), lower) {
			score += 1
		}
		if score > 0 {
			results = append(results, SearchResult{Bookmark: b, Score: score})
		}
	}
	return results
}
