package snapshot

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func makeRunnerStore(t *testing.T) *ScheduleStore {
	t.Helper()
	ss, err := makeScheduleStore(t)
	if err != nil {
		t.Fatalf("makeRunnerStore: %v", err)
	}
	return ss
}

func TestNewRunnerNilStoreReturnsError(t *testing.T) {
	_, err := NewRunner(nil, func(_ context.Context, _ string) error { return nil })
	if err == nil {
		t.Fatal("expected error for nil store")
	}
}

func TestNewRunnerNilCaptureReturnsError(t *testing.T) {
	ss := makeRunnerStore(t)
	_, err := NewRunner(ss, nil)
	if err == nil {
		t.Fatal("expected error for nil capture func")
	}
}

func TestRunnerStartInvokesCaptureOnTick(t *testing.T) {
	ss := makeRunnerStore(t)
	if err := ss.Add(Schedule{ID: "s1", Interval: 20 * time.Millisecond}); err != nil {
		t.Fatalf("Add: %v", err)
	}

	var calls int64
	runner, err := NewRunner(ss, func(_ context.Context, id string) error {
		atomic.AddInt64(&calls, 1)
		return nil
	})
	if err != nil {
		t.Fatalf("NewRunner: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := runner.Start(ctx); err != nil {
		t.Fatalf("Start: %v", err)
	}

	time.Sleep(70 * time.Millisecond)
	cancel()

	got := atomic.LoadInt64(&calls)
	if got < 2 {
		t.Errorf("expected at least 2 captures, got %d", got)
	}
}

func TestRunnerStopHaltsGoroutine(t *testing.T) {
	ss := makeRunnerStore(t)
	if err := ss.Add(Schedule{ID: "s2", Interval: 15 * time.Millisecond}); err != nil {
		t.Fatalf("Add: %v", err)
	}

	var calls int64
	runner, _ := NewRunner(ss, func(_ context.Context, _ string) error {
		atomic.AddInt64(&calls, 1)
		return nil
	})

	ctx := context.Background()
	_ = runner.Start(ctx)
	time.Sleep(35 * time.Millisecond)
	runner.Stop("s2")
	snapshot := atomic.LoadInt64(&calls)
	time.Sleep(40 * time.Millisecond)

	if after := atomic.LoadInt64(&calls); after != snapshot {
		t.Errorf("expected no more captures after Stop, got %d extra", after-snapshot)
	}
}

func TestRunnerStopAllHaltsAll(t *testing.T) {
	ss := makeRunnerStore(t)
	for _, id := range []string{"a", "b"} {
		if err := ss.Add(Schedule{ID: id, Interval: 10 * time.Millisecond}); err != nil {
			t.Fatalf("Add %s: %v", id, err)
		}
	}

	var calls int64
	runner, _ := NewRunner(ss, func(_ context.Context, _ string) error {
		atomic.AddInt64(&calls, 1)
		return nil
	})

	_ = runner.Start(context.Background())
	time.Sleep(35 * time.Millisecond)
	runner.StopAll()
	snap := atomic.LoadInt64(&calls)
	time.Sleep(30 * time.Millisecond)

	if after := atomic.LoadInt64(&calls); after != snap {
		t.Errorf("expected no captures after StopAll, got %d extra", after-snap)
	}
}

func TestRunnerReloadPicksUpNewSchedule(t *testing.T) {
	ss := makeRunnerStore(t)
	if err := ss.Add(Schedule{ID: "r1", Interval: 20 * time.Millisecond}); err != nil {
		t.Fatalf("Add: %v", err)
	}

	var calls int64
	runner, _ := NewRunner(ss, func(_ context.Context, _ string) error {
		atomic.AddInt64(&calls, 1)
		return nil
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := runner.Reload(ctx, "r1"); err != nil {
		t.Fatalf("Reload: %v", err)
	}
	time.Sleep(55 * time.Millisecond)
	cancel()

	if got := atomic.LoadInt64(&calls); got < 1 {
		t.Errorf("expected at least 1 capture after Reload, got %d", got)
	}
}
