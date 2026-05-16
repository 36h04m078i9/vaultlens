package snapshot

import (
	"errors"
	"sync"
	"time"
)

// Label represents a named marker attached to a snapshot.
type Label struct {
	SnapshotID string    `json:"snapshot_id"`
	Name       string    `json:"name"`
	CreatedAt  time.Time `json:"created_at"`
}

// LabelStore manages labels assigned to snapshots.
type LabelStore struct {
	mu     sync.RWMutex
	store  Store
	labels map[string][]Label // keyed by snapshot ID
}

// NewLabelStore creates a LabelStore backed by the given snapshot Store.
func NewLabelStore(s Store) (*LabelStore, error) {
	if s == nil {
		return nil, errors.New("snapshot store must not be nil")
	}
	return &LabelStore{
		store:  s,
		labels: make(map[string][]Label),
	}, nil
}

// Add attaches a label to the snapshot identified by id.
func (ls *LabelStore) Add(id, name string) error {
	if id == "" {
		return errors.New("snapshot id must not be empty")
	}
	if name == "" {
		return errors.New("label name must not be empty")
	}
	if _, ok := ls.store.Get(id); !ok {
		return errors.New("snapshot not found")
	}
	ls.mu.Lock()
	defer ls.mu.Unlock()
	for _, l := range ls.labels[id] {
		if l.Name == name {
			return nil // idempotent
		}
	}
	ls.labels[id] = append(ls.labels[id], Label{
		SnapshotID: id,
		Name:       name,
		CreatedAt:  time.Now().UTC(),
	})
	return nil
}

// Remove detaches a label from a snapshot. Returns false if the label was not present.
func (ls *LabelStore) Remove(id, name string) bool {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	list := ls.labels[id]
	for i, l := range list {
		if l.Name == name {
			ls.labels[id] = append(list[:i], list[i+1:]...)
			return true
		}
	}
	return false
}

// List returns all labels attached to the given snapshot.
func (ls *LabelStore) List(id string) []Label {
	ls.mu.RLock()
	defer ls.mu.RUnlock()
	out := make([]Label, len(ls.labels[id]))
	copy(out, ls.labels[id])
	return out
}

// FindByLabel returns snapshot IDs that carry the given label name.
func (ls *LabelStore) FindByLabel(name string) []string {
	ls.mu.RLock()
	defer ls.mu.RUnlock()
	var ids []string
	for id, labels := range ls.labels {
		for _, l := range labels {
			if l.Name == name {
				ids = append(ids, id)
				break
			}
		}
	}
	return ids
}
