// Package history tracks recent secret paths viewed during a session.
package history

import (
	"sync"
	"time"
)

// Entry represents a single history record.
type Entry struct {
	Path      string    `json:"path"`
	AccessedAt time.Time `json:"accessed_at"`
}

// History holds an ordered list of recently accessed secret paths.
type History struct {
	mu      sync.Mutex
	entries []Entry
	maxSize int
}

// New creates a History with the given maximum capacity.
func New(maxSize int) *History {
	if maxSize <= 0 {
		maxSize = 50
	}
	return &History{maxSize: maxSize}
}

// Push records a path access, deduplicating and moving it to the front.
func (h *History) Push(path string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Remove existing entry for the same path.
	filtered := h.entries[:0]
	for _, e := range h.entries {
		if e.Path != path {
			filtered = append(filtered, e)
		}
	}

	// Prepend the new entry.
	newEntry := Entry{Path: path, AccessedAt: time.Now()}
	h.entries = append([]Entry{newEntry}, filtered...)

	// Trim to maxSize.
	if len(h.entries) > h.maxSize {
		h.entries = h.entries[:h.maxSize]
	}
}

// List returns a copy of all history entries, most recent first.
func (h *History) List() []Entry {
	h.mu.Lock()
	defer h.mu.Unlock()

	out := make([]Entry, len(h.entries))
	copy(out, h.entries)
	return out
}

// Clear removes all history entries.
func (h *History) Clear() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.entries = nil
}

// Len returns the current number of entries.
func (h *History) Len() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return len(h.entries)
}
