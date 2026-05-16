// Package snapshot provides facilities for capturing, storing, and managing
// Vault secret path snapshots.
//
// # Pin
//
// The pin sub-feature allows operators to mark specific snapshots as protected
// so that automated pruning and retention policies leave them untouched.
//
// Usage:
//
//	store, _ := snapshot.NewMemStore()
//	pins, _ := snapshot.NewPinStore(store)
//
//	// Protect a known-good snapshot before running a large migration.
//	_ = pins.Pin("snap-20240101-baseline", "pre-migration baseline")
//
//	// Check protection status before pruning.
//	if pins.IsPinned(id) {
//	    // skip pruning
//	}
//
//	// Release the pin once it is no longer needed.
//	pins.Unpin("snap-20240101-baseline")
//
// PinStore is safe for concurrent use.
package snapshot
