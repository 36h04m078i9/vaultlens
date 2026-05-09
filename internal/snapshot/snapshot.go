package snapshot

import (
	"errors"
	"sort"
	"strings"
	"sync"
	"time"
)

// SecretEntry is a single path/value pair captured in a snapshot.
type SecretEntry struct {
	Path  string
	Value string
}

// Snapshot represents a point-in-time capture of a set of Vault secrets.
type Snapshot struct {
	ID      string
	TakenAt time.Time
	Prefix  string
	Secrets []SecretEntry
}

// Store holds snapshots in memory keyed by ID.
type Store struct {
	mu    sync.RWMutex
	items map[string]*Snapshot
}

// NewStore creates an empty snapshot store.
func NewStore() *Store {
	return &Store{items: make(map[string]*Snapshot)}
}

// Save persists a snapshot. Returns an error if sn is nil or has an empty ID.
func (s *Store) Save(sn *Snapshot) error {
	if sn == nil {
		return errors.New("snapshot must not be nil")
	}
	if sn.ID == "" {
		return errors.New("snapshot ID must not be empty")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[sn.ID] = sn
	return nil
}

// Get retrieves a snapshot by ID.
func (s *Store) Get(id string) (*Snapshot, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sn, ok := s.items[id]
	return sn, ok
}

// Delete removes a snapshot by ID and reports whether it existed.
func (s *Store) Delete(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.items[id]
	delete(s.items, id)
	return ok
}

// All returns all snapshots sorted by TakenAt descending.
func (s *Store) All() []*Snapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*Snapshot, 0, len(s.items))
	for _, sn := range s.items {
		out = append(out, sn)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].TakenAt.After(out[j].TakenAt)
	})
	return out
}

// FilterByPrefix returns snapshots whose Prefix starts with the given string.
func (s *Store) FilterByPrefix(prefix string) []*Snapshot {
	all := s.All()
	if prefix == "" {
		return all
	}
	out := all[:0:0]
	for _, sn := range all {
		if strings.HasPrefix(sn.Prefix, prefix) {
			out = append(out, sn)
		}
	}
	return out
}
