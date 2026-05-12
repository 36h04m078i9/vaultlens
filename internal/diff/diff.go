// Package diff provides utilities for comparing two secret snapshots
// and producing a structured change report for audit and display purposes.
package diff

import "sort"

// ChangeKind represents the type of change detected between two snapshots.
type ChangeKind string

const (
	Added   ChangeKind = "added"
	Removed ChangeKind = "removed"
	Changed ChangeKind = "changed"
)

// Change represents a single detected difference between two secret maps.
type Change struct {
	Path string
	Kind ChangeKind
	OldValue string
	NewValue string
}

// Compare computes the diff between two snapshots of secret key-value pairs.
// Each map is keyed by secret path, with the value being the secret data string.
func Compare(before, after map[string]string) []Change {
	var changes []Change

	for path, oldVal := range before {
		newVal, exists := after[path]
		if !exists {
			changes = append(changes, Change{
				Path:     path,
				Kind:     Removed,
				OldValue: oldVal,
			})
			continue
		}
		if oldVal != newVal {
			changes = append(changes, Change{
				Path:     path,
				Kind:     Changed,
				OldValue: oldVal,
				NewValue: newVal,
			})
		}
	}

	for path, newVal := range after {
		if _, exists := before[path]; !exists {
			changes = append(changes, Change{
				Path:     path,
				Kind:     Added,
				NewValue: newVal,
			})
		}
	}

	sort.Slice(changes, func(i, j int) bool {
		return changes[i].Path < changes[j].Path
	})

	return changes
}

// Summary returns counts of each change kind from a slice of changes.
func Summary(changes []Change) map[ChangeKind]int {
	summary := map[ChangeKind]int{
		Added:   0,
		Removed: 0,
		Changed: 0,
	}
	for _, c := range changes {
		summary[c.Kind]++
	}
	return summary
}

// Filter returns only the changes that match the given ChangeKind.
// This is useful for displaying or processing a specific subset of changes.
func Filter(changes []Change, kind ChangeKind) []Change {
	var result []Change
	for _, c := range changes {
		if c.Kind == kind {
			result = append(result, c)
		}
	}
	return result
}
