package search

import "strings"

// HighlightResult holds a secret path with match positions highlighted.
type HighlightResult struct {
	Path      string
	Highlight string
	Score     int
}

// Highlight wraps Fuzzy results and annotates each matching path with
// bracket-style highlighting around the matched characters, e.g.
// "secret/[my][ap]p/key".
func Highlight(query string, paths []string) []HighlightResult {
	results := Fuzzy(query, paths)
	highlighted := make([]HighlightResult, 0, len(results))

	for _, r := range results {
		hl := highlightPath(query, r.Path)
		highlighted = append(highlighted, HighlightResult{
			Path:      r.Path,
			Highlight: hl,
			Score:     r.Score,
		})
	}
	return highlighted
}

// highlightPath wraps each character in path that participates in a
// subsequence match of query with square brackets.
func highlightPath(query, path string) string {
	if query == "" {
		return path
	}

	lower := strings.ToLower(path)
	q := strings.ToLower(query)

	// Collect matched indices via greedy subsequence walk.
	matched := make([]bool, len(path))
	qi := 0
	for pi := 0; pi < len(lower) && qi < len(q); pi++ {
		if lower[pi] == q[qi] {
			matched[pi] = true
			qi++
		}
	}

	if qi < len(q) {
		// No full subsequence match — return plain path.
		return path
	}

	var sb strings.Builder
	inMatch := false
	for i, ch := range path {
		if matched[i] && !inMatch {
			sb.WriteByte('[')
			inMatch = true
		} else if !matched[i] && inMatch {
			sb.WriteByte(']')
			inMatch = false
		}
		sb.WriteRune(ch)
	}
	if inMatch {
		sb.WriteByte(']')
	}
	return sb.String()
}
