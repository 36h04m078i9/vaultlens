package snapshot

import (
	"testing"
	"time"
)

func makeScheduleStore(t *testing.T) *ScheduleStore {
	t.Helper()
	return NewScheduleStore()
}

func TestAddAndGetSchedule(t *testing.T) {
	ss := makeScheduleStore(t)
	sc := Schedule{ID: "daily", Prefix: "secret/", Interval: 24 * time.Hour, Enabled: true}
	if err := ss.Add(sc); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, ok := ss.Get("daily")
	if !ok {
		t.Fatal("expected schedule to exist")
	}
	if got.Prefix != "secret/" {
		t.Errorf("prefix mismatch: got %q", got.Prefix)
	}
}

func TestAddScheduleEmptyIDReturnsError(t *testing.T) {
	ss := makeScheduleStore(t)
	err := ss.Add(Schedule{Interval: time.Hour})
	if err == nil {
		t.Fatal("expected error for empty ID")
	}
}

func TestAddScheduleZeroIntervalReturnsError(t *testing.T) {
	ss := makeScheduleStore(t)
	err := ss.Add(Schedule{ID: "bad", Interval: 0})
	if err == nil {
		t.Fatal("expected error for zero interval")
	}
}

func TestRemoveSchedule(t *testing.T) {
	ss := makeScheduleStore(t)
	_ = ss.Add(Schedule{ID: "s1", Interval: time.Hour})
	if !ss.Remove("s1") {
		t.Fatal("expected Remove to return true")
	}
	_, ok := ss.Get("s1")
	if ok {
		t.Fatal("expected schedule to be gone")
	}
}

func TestRemoveMissingScheduleReturnsFalse(t *testing.T) {
	ss := makeScheduleStore(t)
	if ss.Remove("nonexistent") {
		t.Fatal("expected false for missing schedule")
	}
}

func TestDueReturnsOverdueEnabledSchedules(t *testing.T) {
	ss := makeScheduleStore(t)
	now := time.Now()
	_ = ss.Add(Schedule{ID: "overdue", Interval: time.Hour, Enabled: true, LastRun: now.Add(-2 * time.Hour)})
	_ = ss.Add(Schedule{ID: "future", Interval: time.Hour, Enabled: true, LastRun: now})
	_ = ss.Add(Schedule{ID: "disabled", Interval: time.Hour, Enabled: false, LastRun: now.Add(-2 * time.Hour)})
	due := ss.Due(now)
	if len(due) != 1 {
		t.Fatalf("expected 1 due schedule, got %d", len(due))
	}
	if due[0].ID != "overdue" {
		t.Errorf("expected 'overdue', got %q", due[0].ID)
	}
}

func TestMarkRunUpdatesLastRun(t *testing.T) {
	ss := makeScheduleStore(t)
	_ = ss.Add(Schedule{ID: "s1", Interval: time.Hour, Enabled: true})
	at := time.Now()
	if err := ss.MarkRun("s1", at); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, _ := ss.Get("s1")
	if !got.LastRun.Equal(at) {
		t.Errorf("LastRun not updated: got %v, want %v", got.LastRun, at)
	}
}

func TestMarkRunMissingReturnsError(t *testing.T) {
	ss := makeScheduleStore(t)
	err := ss.MarkRun("ghost", time.Now())
	if err == nil {
		t.Fatal("expected error for missing schedule")
	}
}
