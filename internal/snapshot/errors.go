package snapshot

import "errors"

// Sentinel errors shared across the snapshot package.
var (
	// ErrNilStore is returned when a nil store is provided.
	ErrNilStore = errors.New("snapshot: store must not be nil")

	// ErrNotFound is returned when a requested snapshot does not exist.
	ErrNotFound = errors.New("snapshot: not found")

	// ErrEmptyID is returned when an empty ID is supplied.
	ErrEmptyID = errors.New("snapshot: id must not be empty")

	// ErrZeroInterval is returned when a schedule interval is zero.
	ErrZeroInterval = errors.New("snapshot: interval must be greater than zero")

	errNilCapture = errors.New("snapshot: capture func must not be nil")
)
