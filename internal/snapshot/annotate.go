package snapshot

import (
	"errors"
	"sync"
	"time"
)

// Annotation holds a free-form note attached to a snapshot.
type Annotation struct {
	SnapshotID string    `json:"snapshot_id"`
	Author     string    `json:"author"`
	Note       string    `json:"note"`
	CreatedAt  time.Time `json:"created_at"`
}

// AnnotationStore manages annotations keyed by snapshot ID.
type AnnotationStore struct {
	mu          sync.RWMutex
	annotations map[string][]Annotation
	store       *Store
}

// NewAnnotationStore returns an AnnotationStore backed by the given snapshot Store.
// It returns an error if store is nil.
func NewAnnotationStore(store *Store) (*AnnotationStore, error) {
	if store == nil {
		return nil, errors.New("snapshot: annotation store requires a non-nil Store")
	}
	return &AnnotationStore{
		annotations: make(map[string][]Annotation),
		store:       store,
	}, nil
}

// Add attaches an annotation to the snapshot identified by id.
// Returns an error if id is empty, note is empty, or the snapshot does not exist.
func (a *AnnotationStore) Add(id, author, note string) error {
	if id == "" {
		return errors.New("snapshot: annotation id must not be empty")
	}
	if note == "" {
		return errors.New("snapshot: annotation note must not be empty")
	}
	if _, ok := a.store.Get(id); !ok {
		return ErrSnapshotNotFound
	}
	ann := Annotation{
		SnapshotID: id,
		Author:     author,
		Note:       note,
		CreatedAt:  time.Now().UTC(),
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	a.annotations[id] = append(a.annotations[id], ann)
	return nil
}

// List returns all annotations for the given snapshot ID.
// Returns an empty slice when none exist.
func (a *AnnotationStore) List(id string) []Annotation {
	a.mu.RLock()
	defer a.mu.RUnlock()
	result := make([]Annotation, len(a.annotations[id]))
	copy(result, a.annotations[id])
	return result
}

// Count returns the number of annotations for the given snapshot ID.
func (a *AnnotationStore) Count(id string) int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return len(a.annotations[id])
}

// Remove deletes all annotations for the given snapshot ID.
// Returns true if any annotations were removed.
func (a *AnnotationStore) Remove(id string) bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	if _, ok := a.annotations[id]; !ok {
		return false
	}
	delete(a.annotations, id)
	return true
}
