package history

// DedupeOptions controls how deduplication is applied to a history list.
type DedupeOptions struct {
	// KeepFirst retains the first occurrence of a path instead of the most recent.
	KeepFirst bool
}

// Dedupe removes duplicate path entries from a slice of Entry values.
// By default the last (most-recent) occurrence of each path is kept.
// When KeepFirst is true the first occurrence is kept instead.
func Dedupe(entries []Entry, opts DedupeOptions) []Entry {
	if len(entries) == 0 {
		return nil
	}

	seen := make(map[string]int, len(entries)) // path -> index in result
	result := make([]Entry, 0, len(entries))

	for _, e := range entries {
		if idx, exists := seen[e.Path]; exists {
			if opts.KeepFirst {
				// already have the first occurrence – skip this one
				continue
			}
			// replace the earlier occurrence with this (more-recent) one
			result[idx] = e
			continue
		}
		seen[e.Path] = len(result)
		result = append(result, e)
	}

	return result
}
