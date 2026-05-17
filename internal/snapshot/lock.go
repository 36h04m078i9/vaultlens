package snapshot

import (
	"errors"
	"sync"
	"time"
)

// ErrAlreadyLocked is returned when a snapshot is already locked by another holder.
var ErrAlreadyLocked = errors.New("snapshot: already locked")

// ErrNotLocked is returned when attempting to release a lock that does not exist.
var ErrNotLocked = errors.New("snapshot: not locked")

// Lock represents an advisory lock held on a snapshot.
type Lock struct {
	SnapshotID string
	Holder     string
	AcquiredAt time.Time
}

// LockStore manages advisory locks on snapshots.
type LockStore struct {
	mu    sync.Mutex
	locks map[string]Lock
	store *Store
}

// NewLockStore creates a LockStore backed by the given snapshot Store.
// Returns an error if store is nil.
func NewLockStore(store *Store) (*LockStore, error) {
	if store == nil {
		return nil, errors.New("snapshot: lock store requires a non-nil store")
	}
	return &LockStore{
		locks: make(map[string]Lock),
		store: store,
	}, nil
}

// Acquire attempts to place an advisory lock on the snapshot identified by id
// on behalf of holder. Returns ErrAlreadyLocked if the snapshot is already
// locked by a different holder, or ErrNotFound if the snapshot does not exist.
func (ls *LockStore) Acquire(id, holder string) (Lock, error) {
	if id == "" {
		return Lock{}, errors.New("snapshot: lock id must not be empty")
	}
	if holder == "" {
		return Lock{}, errors.New("snapshot: lock holder must not be empty")
	}
	if _, ok := ls.store.Get(id); !ok {
		return Lock{}, ErrNotFound
	}
	ls.mu.Lock()
	defer ls.mu.Unlock()
	if existing, ok := ls.locks[id]; ok && existing.Holder != holder {
		return Lock{}, ErrAlreadyLocked
	}
	l := Lock{SnapshotID: id, Holder: holder, AcquiredAt: time.Now().UTC()}
	ls.locks[id] = l
	return l, nil
}

// Release removes the lock on the snapshot identified by id.
// Returns ErrNotLocked if no lock is currently held.
func (ls *LockStore) Release(id string) error {
	if id == "" {
		return errors.New("snapshot: lock id must not be empty")
	}
	ls.mu.Lock()
	defer ls.mu.Unlock()
	if _, ok := ls.locks[id]; !ok {
		return ErrNotLocked
	}
	delete(ls.locks, id)
	return nil
}

// IsLocked reports whether the snapshot identified by id is currently locked.
func (ls *LockStore) IsLocked(id string) bool {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	_, ok := ls.locks[id]
	return ok
}

// Get returns the current Lock for a snapshot, or false if none exists.
func (ls *LockStore) Get(id string) (Lock, bool) {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	l, ok := ls.locks[id]
	return l, ok
}
