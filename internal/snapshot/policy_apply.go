package snapshot

import (
	"errors"
	"fmt"
)

// ApplyPolicy applies a named policy from the given PolicyStore to the
// SnapshotStore, pruning snapshots that violate the policy rules.
// It delegates to ApplyRetention under the hood.
func ApplyPolicy(name string, ps *PolicyStore, ss *Store) error {
	if ps == nil {
		return errors.New("policy store must not be nil")
	}
	if ss == nil {
		return errors.New("snapshot store must not be nil")
	}
	pol, ok := ps.Get(name)
	if !ok {
		return fmt.Errorf("policy %q not found", name)
	}
	return ApplyRetention(ss, &RetentionOptions{
		KeepLast:  pol.KeepLast,
		OlderThan: pol.OlderThan,
		Prefix:    pol.Prefix,
	})
}

// ApplyAllPolicies iterates every policy in the store and applies each one.
// Errors are collected and returned as a combined error after all policies
// have been attempted.
func ApplyAllPolicies(ps *PolicyStore, ss *Store) error {
	if ps == nil {
		return errors.New("policy store must not be nil")
	}
	var errs []error
	for _, name := range ps.List() {
		if err := ApplyPolicy(name, ps, ss); err != nil {
			errs = append(errs, fmt.Errorf("policy %q: %w", name, err))
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("%d policy error(s): %v", len(errs), errs)
	}
	return nil
}
