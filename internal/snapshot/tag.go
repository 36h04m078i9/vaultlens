package snapshot

import (
	"errors"
	"strings"
	"time"
)

// Tag associates a human-readable label with a snapshot ID.
type Tag struct {
	Name      string    `json:"name"`
	SnapshotID string   `json:"snapshot_id"`
	CreatedAt time.Time `json:"created_at"`
	Note      string    `json:"note,omitempty"`
}

// TagStore manages snapshot tags backed by a Store.
type TagStore struct {
	store *Store
}

// NewTagStore creates a TagStore that validates tags against the given Store.
func NewTagStore(s *Store) (*TagStore, error) {
	if s == nil {
		return nil, errors.New("snapshot: store must not be nil")
	}
	return &TagStore{store: s}, nil
}

// Add creates a new tag pointing at snapshotID. Returns an error if the
// snapshot does not exist or the tag name is empty.
func (ts *TagStore) Add(name, snapshotID, note string) (*Tag, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("snapshot: tag name must not be empty")
	}
	if snapshotID == "" {
		return nil, errors.New("snapshot: snapshot ID must not be empty")
	}
	if _, ok := ts.store.Get(snapshotID); !ok {
		return nil, errors.New("snapshot: snapshot not found: " + snapshotID)
	}
	t := &Tag{
		Name:       name,
		SnapshotID: snapshotID,
		CreatedAt:  time.Now().UTC(),
		Note:       note,
	}
	ts.store.tags[name] = t
	return t, nil
}

// Get retrieves a tag by name.
func (ts *TagStore) Get(name string) (*Tag, bool) {
	t, ok := ts.store.tags[name]
	return t, ok
}

// Delete removes a tag by name. Returns false if the tag did not exist.
func (ts *TagStore) Delete(name string) bool {
	if _, ok := ts.store.tags[name]; !ok {
		return false
	}
	delete(ts.store.tags, name)
	return true
}

// List returns all tags in the store.
func (ts *TagStore) List() []*Tag {
	out := make([]*Tag, 0, len(ts.store.tags))
	for _, t := range ts.store.tags {
		out = append(out, t)
	}
	return out
}
