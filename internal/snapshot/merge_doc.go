// Package snapshot provides facilities for capturing, storing, comparing,
// and managing point-in-time snapshots of Vault secret paths.
//
// # Merge
//
// The Merge function combines two existing snapshots into a new one:
//
//	result, err := snapshot.Merge(store, "snap-2024-01", "snap-2024-02", snapshot.MergeOptions{
//		PreferA: true,
//		NewID:   "snap-merged",
//	})
//
// When the same path/key appears in both snapshots with different values, a
// MergeConflict is recorded. MergeOptions.PreferA controls which value wins:
//
//   - PreferA = true  → keep the value from snapshot A
//   - PreferA = false → keep the value from snapshot B (default)
//
// The merged snapshot is saved to the store and its ID is returned in
// MergeResult.ID. Source snapshot IDs are preserved in MergeResult.SourceIDs
// for auditing purposes.
package snapshot
