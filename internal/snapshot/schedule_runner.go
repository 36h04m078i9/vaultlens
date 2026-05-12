package snapshot

import (
	"context"
	"log"
	"sync"
	"time"
)

// Runner executes scheduled snapshots at their configured intervals.
type Runner struct {
	schedules *ScheduleStore
	capture  func(ctx context.Context, id string) error
	mu       sync.Mutex
	cancel   map[string]context.CancelFunc
}

// NewRunner creates a Runner backed by the given ScheduleStore.
// The capture func is called with the schedule ID when a snapshot is due.
func NewRunner(schedules *ScheduleStore, capture func(ctx context.Context, id string) error) (*Runner, error) {
	if schedules == nil {
		return nil, ErrNilStore
	}
	if capture == nil {
		return nil, errNilCapture
	}
	return &Runner{
		schedules: schedules,
		capture:   capture,
		cancel:    make(map[string]context.CancelFunc),
	}, nil
}

// Start launches a goroutine for every active schedule.
func (r *Runner) Start(ctx context.Context) error {
	list, err := r.schedules.List()
	if err != nil {
		return err
	}
	for _, s := range list {
		r.startOne(ctx, s)
	}
	return nil
}

// Reload stops any running goroutine for id and starts a fresh one.
func (r *Runner) Reload(ctx context.Context, id string) error {
	r.Stop(id)
	s, ok, err := r.schedules.Get(id)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}
	r.startOne(ctx, s)
	return nil
}

// Stop cancels the goroutine for the given schedule id.
func (r *Runner) Stop(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if cancel, ok := r.cancel[id]; ok {
		cancel()
		delete(r.cancel, id)
	}
}

// StopAll cancels every running goroutine.
func (r *Runner) StopAll() {
	r.mu.Lock()
	defer r.mu.Unlock()
	for id, cancel := range r.cancel {
		cancel()
		delete(r.cancel, id)
	}
}

func (r *Runner) startOne(ctx context.Context, s Schedule) {
	child, cancel := context.WithCancel(ctx)
	r.mu.Lock()
	r.cancel[s.ID] = cancel
	r.mu.Unlock()

	go func() {
		ticker := time.NewTicker(s.Interval)
		defer ticker.Stop()
		for {
			select {
			case <-child.Done():
				return
			case <-ticker.C:
				if err := r.capture(child, s.ID); err != nil {
					log.Printf("snapshot runner: capture %q: %v", s.ID, err)
				}
			}
		}
	}()
}
