package snapshot

import (
	"errors"
	"time"
)

// RetentionPolicy defines rules for how long snapshots are kept.
type RetentionPolicy struct {
	// KeepLast retains the N most recent snapshots. Zero means no limit.
	KeepLast int
	// OlderThan removes snapshots older than this duration. Zero means no limit.
	OlderThan time.Duration
	// Prefix restricts the policy to snapshots whose ID starts with this string.
	Prefix string
}

// Validate returns an error if the policy has no meaningful rules.
func (p RetentionPolicy) Validate() error {
	if p.KeepLast == 0 && p.OlderThan == 0 {
		return errors.New("retention policy must specify at least one of KeepLast or OlderThan")
	}
	if p.KeepLast < 0 {
		return errors.New("KeepLast must be non-negative")
	}
	return nil
}

// PolicyStore manages named retention policies.
type PolicyStore struct {
	policies map[string]RetentionPolicy
}

// NewPolicyStore returns an initialised PolicyStore.
func NewPolicyStore() *PolicyStore {
	return &PolicyStore{policies: make(map[string]RetentionPolicy)}
}

// Add registers a named policy, replacing any existing one with the same name.
func (s *PolicyStore) Add(name string, p RetentionPolicy) error {
	if name == "" {
		return errors.New("policy name must not be empty")
	}
	if err := p.Validate(); err != nil {
		return err
	}
	s.policies[name] = p
	return nil
}

// Get retrieves a policy by name.
func (s *PolicyStore) Get(name string) (RetentionPolicy, bool) {
	p, ok := s.policies[name]
	return p, ok
}

// Remove deletes a named policy. Returns false if it did not exist.
func (s *PolicyStore) Remove(name string) bool {
	_, ok := s.policies[name]
	if ok {
		delete(s.policies, name)
	}
	return ok
}

// List returns all policy names in an unspecified order.
func (s *PolicyStore) List() []string {
	names := make([]string, 0, len(s.policies))
	for n := range s.policies {
		names = append(names, n)
	}
	return names
}
