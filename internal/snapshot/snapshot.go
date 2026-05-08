// Package snapshot provides point-in-time captures of Vault secret paths
// for comparison, auditing, and drift detection.
package snapshot

import (
	"fmt"
	"time"
)

// Snapshot represents a captured set of Vault secret paths at a point in time.
type Snapshot struct {
	ID        string            `json:"id"`
	Label     string            `json:"label"`
	CapturedAt time.Time        `json:"captured_at"`
	Paths     []string          `json:"paths"`
	Meta      map[string]string `json:"meta,omitempty"`
}

// Store holds named snapshots in memory.
type Store struct {
	snaps map[string]*Snapshot
}

// NewStore returns an initialised snapshot Store.
func NewStore() *Store {
	return &Store{snaps: make(map[string]*Snapshot)}
}

// Save stores a snapshot, keyed by its ID.
func (s *Store) Save(snap *Snapshot) error {
	if snap == nil {
		return fmt.Errorf("snapshot: cannot save nil snapshot")
	}
	if snap.ID == "" {
		return fmt.Errorf("snapshot: ID must not be empty")
	}
	s.snaps[snap.ID] = snap
	return nil
}

// Get retrieves a snapshot by ID. Returns nil, false when not found.
func (s *Store) Get(id string) (*Snapshot, bool) {
	snap, ok := s.snaps[id]
	return snap, ok
}

// List returns all stored snapshots ordered by capture time (oldest first).
func (s *Store) List() []*Snapshot {
	out := make([]*Snapshot, 0, len(s.snaps))
	for _, snap := range s.snaps {
		out = append(out, snap)
	}
	// simple insertion-sort by CapturedAt
	for i := 1; i < len(out); i++ {
		for j := i; j > 0 && out[j].CapturedAt.Before(out[j-1].CapturedAt); j-- {
			out[j], out[j-1] = out[j-1], out[j]
		}
	}
	return out
}

// Delete removes a snapshot by ID. Returns false when the ID was not found.
func (s *Store) Delete(id string) bool {
	if _, ok := s.snaps[id]; !ok {
		return false
	}
	delete(s.snaps, id)
	return true
}
