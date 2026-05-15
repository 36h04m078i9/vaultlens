// Package snapshot provides facilities for capturing, storing, comparing,
// and managing point-in-time snapshots of Vault secret paths.
//
// # Clone
//
// The Clone function duplicates an existing snapshot under a new identifier.
// This is useful when an operator wants to branch from a known-good baseline
// before making changes, or to create a labelled copy for archival purposes.
//
// Basic usage:
//
//	result, err := snapshot.Clone(store, "snap-20240601", snapshot.CloneOptions{
//		NewID:     "snap-20240601-staging",
//		ExtraTags: []string{"staging", "pre-release"},
//		Note:      "Cloned before staging deployment",
//	})
//
// If NewID is left empty, a unique identifier is generated automatically
// using the source ID and a nanosecond timestamp suffix.
//
// The clone receives a deep copy of all secret paths from the source snapshot,
// so mutations to the original in memory do not affect the clone and vice versa.
//
// Extra tags are appended to any tags already present on the source snapshot.
package snapshot
