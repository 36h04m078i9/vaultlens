// Package snapshot provides facilities for capturing, storing, comparing,
// listing, pruning, tagging, and exporting point-in-time views of Vault
// secret paths.
//
// # Core types
//
//   - Store – in-memory snapshot repository (thread-safe).
//   - ScheduleStore – persists recurring capture schedules.
//   - Runner – executes scheduled captures at their configured intervals.
//   - TagStore – associates human-readable tags with snapshot IDs.
//
// # Typical workflow
//
//  1. Create a Store and a ScheduleStore.
//  2. Register schedules with ScheduleStore.Add.
//  3. Instantiate a Runner with a capture func that calls Capture.
//  4. Call Runner.Start to begin background captures.
//  5. Use Compare or DiffSnapshots to detect secret drift.
//  6. Export results via ToExportRecords for audit purposes.
//
// All exported functions are safe for concurrent use unless stated otherwise.
package snapshot
