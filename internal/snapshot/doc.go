// Package snapshot provides facilities for capturing, storing, comparing,
// listing, pruning, tagging, exporting, and scheduling periodic snapshots
// of Vault secret paths.
//
// # Capturing
//
// Use [Capture] with a [PathLister] and a [Store] to walk a Vault prefix and
// persist a point-in-time snapshot of all discovered paths.
//
// # Comparing
//
// [Compare] and [DiffSnapshots] detect added, removed, and changed keys
// between two stored snapshots, returning structured [diff.Result] slices.
//
// # Scheduling
//
// [ScheduleStore] manages recurring capture policies. Call [ScheduleStore.Due]
// with the current time to obtain schedules ready to run, then [ScheduleStore.MarkRun]
// after each successful capture to advance the schedule cursor.
//
// # Pruning
//
// [Prune] trims a snapshot store by retention count or age, keeping the
// most recent N snapshots or removing those older than a given threshold.
package snapshot
