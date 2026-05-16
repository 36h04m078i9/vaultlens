package snapshot

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

// stubCapture returns a CaptureFunc that always returns the provided snapshot.
func stubCapture(snap *Snapshot, err error) CaptureFunc {
	return func(_ context.Context, _, id string) (*Snapshot, error) {
		if snap != nil {
			snap.ID = id
		}
		return snap, err
	}
}

func TestNewWatcherNilStoreReturnsError(t *testing.T) {
	_, err := NewWatcher(nil, stubCapture(&Snapshot{}, nil), WatchOptions{})
	if err == nil {
		t.Fatal("expected error for nil store")
	}
}

func TestNewWatcherNilCaptureReturnsError(t *testing.T) {
	store := makeSnap(t)
	_, err := NewWatcher(store, nil, WatchOptions{})
	if err == nil {
		t.Fatal("expected error for nil capture")
	}
}

func TestNewWatcherDefaultsInterval(t *testing.T) {
	store := makeSnap(t)
	w, err := NewWatcher(store, stubCapture(&Snapshot{}, nil), WatchOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if w.opts.Interval != 60*time.Second {
		t.Fatalf("expected 60s default, got %v", w.opts.Interval)
	}
}

func TestWatcherEmitsEventOnChange(t *testing.T) {
	store := makeSnap(t)

	// Seed two different snapshots so DiffSnapshots finds changes.
	snapA := &Snapshot{ID: "a", Paths: map[string]map[string]string{"secret/x": {"k": "v1"}}}
	snapB := &Snapshot{ID: "b", Paths: map[string]map[string]string{"secret/x": {"k": "v2"}}}
	if err := store.Save(snapA); err != nil {
		t.Fatal(err)
	}
	if err := store.Save(snapB); err != nil {
		t.Fatal(err)
	}

	call := 0
	captureFn := func(_ context.Context, _, id string) (*Snapshot, error) {
		call++
		if call == 1 {
			snapA.ID = id
			_ = store.Save(snapA)
			return snapA, nil
		}
		snapB.ID = id
		_ = store.Save(snapB)
		return snapB, nil
	}

	var mu sync.Mutex
	var events []WatchEvent
	opts := WatchOptions{
		Interval: 10 * time.Millisecond,
		OnEvent: func(e WatchEvent) {
			mu.Lock()
			events = append(events, e)
			mu.Unlock()
		},
	}
	w, err := NewWatcher(store, captureFn, opts)
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 80*time.Millisecond)
	defer cancel()
	w.Start(ctx)

	mu.Lock()
	defer mu.Unlock()
	if len(events) == 0 {
		t.Fatal("expected at least one watch event")
	}
	if events[0].Changes == 0 {
		t.Fatal("expected non-zero changes in event")
	}
}

func TestWatcherCaptureErrorSkipsCycle(t *testing.T) {
	store := makeSnap(t)
	captureFn := stubCapture(nil, errors.New("vault unavailable"))

	var events []WatchEvent
	opts := WatchOptions{
		Interval: 10 * time.Millisecond,
		OnEvent: func(e WatchEvent) { events = append(events, e) },
	}
	w, err := NewWatcher(store, captureFn, opts)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	w.Start(ctx)

	if len(events) != 0 {
		t.Fatalf("expected no events on capture error, got %d", len(events))
	}
}
