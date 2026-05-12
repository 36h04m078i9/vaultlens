package snapshot

import (
	"errors"
	"time"
)

// Schedule defines a recurring snapshot capture policy.
type Schedule struct {
	ID       string        `json:"id"`
	Prefix   string        `json:"prefix"`
	Interval time.Duration `json:"interval"`
	KeepLast int           `json:"keep_last"`
	Enabled  bool          `json:"enabled"`
	LastRun  time.Time     `json:"last_run,omitempty"`
}

// ScheduleStore manages snapshot schedules.
type ScheduleStore struct {
	schedules map[string]*Schedule
}

// NewScheduleStore creates an empty ScheduleStore.
func NewScheduleStore() *ScheduleStore {
	return &ScheduleStore{schedules: make(map[string]*Schedule)}
}

// Add inserts or replaces a schedule. ID and Interval are required.
func (s *ScheduleStore) Add(sc Schedule) error {
	if sc.ID == "" {
		return errors.New("schedule id must not be empty")
	}
	if sc.Interval <= 0 {
		return errors.New("schedule interval must be positive")
	}
	copy := sc
	s.schedules[sc.ID] = &copy
	return nil
}

// Get returns the schedule with the given ID.
func (s *ScheduleStore) Get(id string) (*Schedule, bool) {
	sc, ok := s.schedules[id]
	if !ok {
		return nil, false
	}
	copy := *sc
	return &copy, true
}

// Remove deletes a schedule by ID. Returns false if not found.
func (s *ScheduleStore) Remove(id string) bool {
	_, ok := s.schedules[id]
	if !ok {
		return false
	}
	delete(s.schedules, id)
	return true
}

// Due returns all enabled schedules whose next run time has passed.
func (s *ScheduleStore) Due(now time.Time) []*Schedule {
	var due []*Schedule
	for _, sc := range s.schedules {
		if !sc.Enabled {
			continue
		}
		next := sc.LastRun.Add(sc.Interval)
		if now.Equal(next) || now.After(next) {
			copy := *sc
			due = append(due, &copy)
		}
	}
	return due
}

// MarkRun updates the LastRun timestamp for the given schedule ID.
func (s *ScheduleStore) MarkRun(id string, at time.Time) error {
	sc, ok := s.schedules[id]
	if !ok {
		return errors.New("schedule not found: " + id)
	}
	sc.LastRun = at
	return nil
}
