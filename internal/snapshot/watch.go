package snapshot

import (
	"context"
	"errors"
	"sync"
	"time"
)

// WatchEvent describes a change detected between two consecutive snapshots.
type WatchEvent struct {
	PrefixID  string
	SnapshotA string
	SnapshotB string
	Summary   string
	Changes   int
	At        time.Time
}

// WatchOptions configures a Watch run.
type WatchOptions struct {
	// Interval between snapshot captures. Defaults to 60 seconds.
	Interval time.Duration
	// Prefix filters which vault paths are captured.
	Prefix string
	// OnEvent is called for every polling cycle that produces changes.
	OnEvent func(WatchEvent)
}

// Watcher polls vault at a fixed interval, captures snapshots, and emits
// WatchEvents whenever the captured state differs from the previous one.
type Watcher struct {
	store   Storer
	capture CaptureFunc
	opts    WatchOptions
	mu      sync.Mutex
	prevID  string
}

// CaptureFunc is the signature expected by Watcher for snapshot capture.
type CaptureFunc func(ctx context.Context, prefix, id string) (*Snapshot, error)

// NewWatcher creates a Watcher. store and capture must not be nil.
func NewWatcher(store Storer, capture CaptureFunc, opts WatchOptions) (*Watcher, error) {
	if store == nil {
		return nil, errors.New("snapshot: watcher store must not be nil")
	}
	if capture == nil {
		return nil, errors.New("snapshot: watcher capture func must not be nil")
	}
	if opts.Interval <= 0 {
		opts.Interval = 60 * time.Second
	}
	if opts.OnEvent == nil {
		opts.OnEvent = func(WatchEvent) {}
	}
	return &Watcher{store: store, capture: capture, opts: opts}, nil
}

// Start begins the polling loop and blocks until ctx is cancelled.
func (w *Watcher) Start(ctx context.Context) {
	ticker := time.NewTicker(w.opts.Interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case t := <-ticker.C:
			w.poll(ctx, t)
		}
	}
}

func (w *Watcher) poll(ctx context.Context, at time.Time) {
	id := "watch-" + at.UTC().Format("20060102-150405")
	snap, err := w.capture(ctx, w.opts.Prefix, id)
	if err != nil || snap == nil {
		return
	}
	w.mu.Lock()
	prev := w.prevID
	w.prevID = snap.ID
	w.mu.Unlock()

	if prev == "" {
		return
	}
	diffs, err := DiffSnapshots(w.store, prev, snap.ID)
	if err != nil || len(diffs) == 0 {
		return
	}
	w.opts.OnEvent(WatchEvent{
		PrefixID:  w.opts.Prefix,
		SnapshotA: prev,
		SnapshotB: snap.ID,
		Summary:   DiffSummary(diffs),
		Changes:   len(diffs),
		At:        at,
	})
}
