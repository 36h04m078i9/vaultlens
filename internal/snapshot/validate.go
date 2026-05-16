package snapshot

import (
	"errors"
	"fmt"	
	"strings"
)

// ValidationResult holds the outcome of a snapshot validation check.
type ValidationResult struct {
	SnapshotID string
	Valid      bool
	Issues     []string
}

// Validator checks snapshot integrity against a set of rules.
type Validator struct {
	store *Store
}

// NewValidator creates a Validator backed by the given Store.
// Returns an error if store is nil.
func NewValidator(store *Store) (*Validator, error) {
	if store == nil {
		return nil, errors.New("snapshot: validator requires a non-nil store")
	}
	return &Validator{store: store}, nil
}

// Validate checks a single snapshot by ID and returns a ValidationResult.
// It verifies the snapshot exists, has a non-empty ID, a non-zero timestamp,
// and that all secret paths are non-empty strings.
func (v *Validator) Validate(id string) (ValidationResult, error) {
	if strings.TrimSpace(id) == "" {
		return ValidationResult{}, errors.New("snapshot: validate requires a non-empty id")
	}

	snap, ok := v.store.Get(id)
	if !ok {
		return ValidationResult{}, fmt.Errorf("snapshot: not found: %s", id)
	}

	result := ValidationResult{
		SnapshotID: snap.ID,
		Valid:      true,
	}

	if snap.ID == "" {
		result.Issues = append(result.Issues, "snapshot has an empty ID")
	}
	if snap.CapturedAt.IsZero() {
		result.Issues = append(result.Issues, "snapshot has a zero CapturedAt timestamp")
	}
	for i, entry := range snap.Entries {
		if strings.TrimSpace(entry.Path) == "" {
			result.Issues = append(result.Issues, fmt.Sprintf("entry[%d] has an empty path", i))
		}
	}

	if len(result.Issues) > 0 {
		result.Valid = false
	}
	return result, nil
}

// ValidateAll runs Validate on every snapshot in the store and returns all results.
func (v *Validator) ValidateAll() ([]ValidationResult, error) {
	snaps := v.store.List()
	results := make([]ValidationResult, 0, len(snaps))
	for _, s := range snaps {
		r, err := v.Validate(s.ID)
		if err != nil {
			return nil, err
		}
		results = append(results, r)
	}
	return results, nil
}
