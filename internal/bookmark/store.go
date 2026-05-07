package bookmark

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"
)

// Bookmark represents a saved Vault secret path.
type Bookmark struct {
	Path      string    `json:"path"`
	Note      string    `json:"note,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// Store manages a persisted list of bookmarked paths.
type Store struct {
	mu       sync.RWMutex
	filePath string
	items    []Bookmark
}

// NewStore creates a Store backed by the given file path.
func NewStore(filePath string) (*Store, error) {
	s := &Store{filePath: filePath}
	if err := s.load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	return s, nil
}

// Add inserts a new bookmark. Duplicate paths are silently ignored.
func (s *Store) Add(path, note string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, b := range s.items {
		if b.Path == path {
			return nil
		}
	}
	s.items = append(s.items, Bookmark{Path: path, Note: note, CreatedAt: time.Now().UTC()})
	return s.save()
}

// Remove deletes a bookmark by path. Returns false if not found.
func (s *Store) Remove(path string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, b := range s.items {
		if b.Path == path {
			s.items = append(s.items[:i], s.items[i+1:]...)
			return true, s.save()
		}
	}
	return false, nil
}

// List returns a copy of all bookmarks.
func (s *Store) List() []Bookmark {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Bookmark, len(s.items))
	copy(out, s.items)
	return out
}

func (s *Store) load() error {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &s.items)
}

func (s *Store) save() error {
	data, err := json.MarshalIndent(s.items, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.filePath, data, 0o600)
}
