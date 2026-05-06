package search

import (
	"strings"
)

// Result holds a matched secret path and its relevance score.
type Result struct {
	Path  string
	Score int
}

// Fuzzy performs a fuzzy search over a list of secret paths using the given
// query. It returns results sorted by descending score, omitting non-matches.
func Fuzzy(paths []string, query string) []Result {
	if query == "" {
		results := make([]Result, len(paths))
		for i, p := range paths {
			results[i] = Result{Path: p, Score: 0}
		}
		return results
	}

	var results []Result
	q := strings.ToLower(query)

	for _, path := range paths {
		if score := score(strings.ToLower(path), q); score > 0 {
			results = append(results, Result{Path: path, Score: score})
		}
	}

	sortResults(results)
	return results
}

// score computes a simple fuzzy match score between a candidate string and a
// query. Returns 0 if the query characters are not a subsequence of candidate.
func score(candidate, query string) int {
	if strings.Contains(candidate, query) {
		// Exact substring match gets a high bonus.
		return 100 + (100 - len(candidate))
	}

	// Subsequence match.
	ci, qi := 0, 0
	for ci < len(candidate) && qi < len(query) {
		if candidate[ci] == query[qi] {
			qi++
		}
		ci++
	}
	if qi < len(query) {
		return 0 // query is not a subsequence
	}
	return 50 - ci // fewer skipped chars = higher score
}

// sortResults sorts results in place by descending Score.
func sortResults(results []Result) {
	for i := 1; i < len(results); i++ {
		for j := i; j > 0 && results[j].Score > results[j-1].Score; j-- {
			results[j], results[j-1] = results[j-1], results[j]
		}
	}
}
