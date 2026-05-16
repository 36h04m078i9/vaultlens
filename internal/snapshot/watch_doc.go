// Package snapshot provides facilities for capturing, storing, comparing,
// and watching HashiCorp Vault secret trees.
//
// # Watch
//
// The Watcher type polls Vault at a configurable interval, captures a new
// snapshot on every tick, and compares it against the previous capture.
// When differences are detected a WatchEvent is delivered to the caller's
// OnEvent callback.
//
// Basic usage:
//
//	store := snapshot.NewStore()
//	watcher, err := snapshot.NewWatcher(store, myCaptureFn, snapshot.WatchOptions{
//		Interval: 5 * time.Minute,
//		Prefix:   "secret/production",
//		OnEvent: func(e snapshot.WatchEvent) {
//			fmt.Printf("[%s] %d change(s): %s\n", e.At.Format(time.RFC3339), e.Changes, e.Summary)
//		},
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	watcher.Start(ctx)
//
// The Watcher is read-only: it never writes back to Vault.
package snapshot
