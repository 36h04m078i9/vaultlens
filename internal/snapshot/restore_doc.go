// Package snapshot provides snapshot capture, storage, diffing, tagging,
// scheduling, policy management, and restoration for Vault secret paths.
//
// # Restore
//
// The Restorer type reads a previously captured snapshot from the Store and
// replays its secrets through a SecretWriter, which is typically backed by the
// Vault client.
//
// Basic usage:
//
//	restorer, err := snapshot.NewRestorer(store, vaultClient)
//	if err != nil { ... }
//
//	records, err := restorer.Restore("snap-20240601", snapshot.RestoreOptions{
//		PathPrefix: "secret/app",
//		DryRun:     true, // preview without writing
//	})
//
// Use FormatRestoreRecords to render results for human consumption:
//
//	snapshot.FormatRestoreRecords(os.Stdout, records, snapshot.RestoreFormatOptions{
//		DryRun:   true,
//		ShowData: false, // never print secret values
//	})
package snapshot
