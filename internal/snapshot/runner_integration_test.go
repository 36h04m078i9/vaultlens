package snapshot_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/example/vaultlens/internal/snapshot"
)

// TestRunnerIntegrationFullLifecycle exercises Start → Stop → Reload
// using exported types only (black-box integration test).
func TestRunnerIntegrationFullLifecycle(t *testing.T) {
	t.Parallel()

	store := snapshot.NewStore()
	ss, err := snapshot.NewScheduleStore(store)
	if err != nil {
		t.Fatalf("NewScheduleStore: %v", err)
	}

	if err := ss.Add(snapshot.Schedule{ID: "integ", Interval: 20 * time.Millisecond}); err != nil {
		t.Fatalf("Add: %v", err)
	}

	var count int64
	runner, err := snapshot.NewRunner(ss, func(_ context.Context, id string) error {
		atomic.AddInt64(&count, 1)
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

	time.Sleep(65 * time.Millisecond)
	runner.Stop("integ")

	after := atomic.LoadInt64(&count)
	if after < 2 {
		t.Errorf("expected ≥2 captures, got %d", after)
	}

	// Reload should resume capturing.
	if err := runner.Reload(ctx, "integ"); err != nil {
		t.Fatalf("Reload: %v", err)
	}
	time.Sleep(45 * time.Millisecond)
	runner.StopAll()

	final := atomic.LoadInt64(&count)
	if final <= after {
		t.Errorf("expected more captures after Reload, got %d (was %d)", final, after)
	}
}
