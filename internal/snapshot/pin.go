package snapshot

import (
	"errors"
	"sync"
	"time"
)

// PinnedSnapshot represents a snapshot that has been pinned to prevent pruning.
type PinnedSnapshot struct {
	ID        string    `json:"id"`
	Reason    string    `json:"reason"`
	PinnedAt  time.Time `json:"pinned_at"`
}

// PinStore manages pinned snapshots.
type PinStore struct {
	mu      sync.RWMutex
	pins    map[string]PinnedSnapshot
	snapshots Store
}

// NewPinStore creates a new PinStore backed by the given snapshot Store.
func NewPinStore(s Store) (*PinStore, error) {
	if s == nil {
		return nil, errors.New("snapshot store must not be nil")
	}
	return &PinStore{
		pins:      make(map[string]PinnedSnapshot),
		snapshots: s,
	}, nil
}

// Pin marks a snapshot as pinned with an optional reason.
func (p *PinStore) Pin(id, reason string) error {
	if id == "" {
		return errors.New("snapshot id must not be empty")
	}
	_, ok := p.snapshots.Get(id)
	if !ok {
		return errors.New("snapshot not found: " + id)
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.pins[id] = PinnedSnapshot{
		ID:       id,
		Reason:   reason,
		PinnedAt: time.Now().UTC(),
	}
	return nil
}

// Unpin removes the pin from a snapshot.
func (p *PinStore) Unpin(id string) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	_, exists := p.pins[id]
	if !exists {
		return false
	}
	delete(p.pins, id)
	return true
}

// IsPinned reports whether a snapshot is currently pinned.
func (p *PinStore) IsPinned(id string) bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	_, ok := p.pins[id]
	return ok
}

// List returns all currently pinned snapshots.
func (p *PinStore) List() []PinnedSnapshot {
	p.mu.RLock()
	defer p.mu.RUnlock()
	out := make([]PinnedSnapshot, 0, len(p.pins))
	for _, pin := range p.pins {
		out = append(out, pin)
	}
	return out
}
