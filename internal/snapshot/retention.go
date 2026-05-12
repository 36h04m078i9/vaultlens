package snapshot

import (
	"errors"
	"sort"
	"time"
)

// RetentionPolicy defines rules for how many snapshots to keep.
type RetentionPolicy struct {
	// KeepLast retains the N most recent snapshots. Zero means no limit.
	KeepLast int
	// KeepDaily retains the most recent snapshot per calendar day for N days.
	KeepDaily int
	// OlderThan removes snapshots taken before this time. Zero value is ignored.
	OlderThan time.Time
}

// ApplyRetention applies the given policy to the store, deleting snapshots
// that fall outside the retention rules. At least one rule must be set.
func ApplyRetention(store *Store, policy RetentionPolicy) (int, error) {
	if store == nil {
		return 0, errors.New("snapshot: store must not be nil")
	}
	if policy.KeepLast == 0 && policy.KeepDaily == 0 && policy.OlderThan.IsZero() {
		return 0, errors.New("snapshot: retention policy has no rules set")
	}

	snaps, err := store.List()
	if err != nil {
		return 0, err
	}

	// Sort newest first.
	sort.Slice(snaps, func(i, j int) bool {
		return snaps[i].CapturedAt.After(snaps[j].CapturedAt)
	})

	keep := make(map[string]struct{})

	if policy.KeepLast > 0 {
		for i := 0; i < policy.KeepLast && i < len(snaps); i++ {
			keep[snaps[i].ID] = struct{}{}
		}
	}

	if policy.KeepDaily > 0 {
		seenDays := map[string]struct{}{}
		cutoff := time.Now().AddDate(0, 0, -policy.KeepDaily)
		for _, s := range snaps {
			if s.CapturedAt.Before(cutoff) {
				break
			}
			day := s.CapturedAt.Format("2006-01-02")
			if _, seen := seenDays[day]; !seen {
				seenDays[day] = struct{}{}
				keep[s.ID] = struct{}{}
			}
		}
	}

	removed := 0
	for _, s := range snaps {
		if _, ok := keep[s.ID]; ok {
			continue
		}
		if !policy.OlderThan.IsZero() && s.CapturedAt.After(policy.OlderThan) {
			continue
		}
		if err := store.Delete(s.ID); err != nil {
			return removed, err
		}
		removed++
	}
	return removed, nil
}
