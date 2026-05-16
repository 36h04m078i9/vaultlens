package snapshot

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// ArchiveRecord represents a single snapshot serialized for export/archive.
type ArchiveRecord struct {
	ID        string            `json:"id"`
	CapturedAt time.Time        `json:"captured_at"`
	Tags      []string          `json:"tags,omitempty"`
	Paths     map[string]string `json:"paths"`
}

// ArchiveOptions controls which snapshots are included in the archive.
type ArchiveOptions struct {
	// IDs restricts the archive to specific snapshot IDs. If empty, all snapshots are included.
	IDs []string
	// Since, if non-zero, excludes snapshots captured before this time.
	Since time.Time
}

// Archive writes all matching snapshots from store as a newline-delimited JSON
// stream to w. Each line is a valid JSON object representing one ArchiveRecord.
func Archive(store *Store, w io.Writer, opts ArchiveOptions) error {
	if store == nil {
		return fmt.Errorf("snapshot: archive: store must not be nil")
	}
	if w == nil {
		return fmt.Errorf("snapshot: archive: writer must not be nil")
	}

	filter := func(s *Snapshot) bool {
		if !opts.Since.IsZero() && s.CapturedAt.Before(opts.Since) {
			return false
		}
		if len(opts.IDs) > 0 {
			for _, id := range opts.IDs {
				if id == s.ID {
					return true
				}
			}
			return false
		}
		return true
	}

	snaps := store.All()
	enc := json.NewEncoder(w)
	written := 0
	for _, s := range snaps {
		if !filter(s) {
			continue
		}
		rec := ArchiveRecord{
			ID:         s.ID,
			CapturedAt: s.CapturedAt,
			Tags:       s.Tags,
			Paths:      s.Data,
		}
		if err := enc.Encode(rec); err != nil {
			return fmt.Errorf("snapshot: archive: encode id=%s: %w", s.ID, err)
		}
		written++
	}
	if written == 0 {
		return fmt.Errorf("snapshot: archive: no snapshots matched the given options")
	}
	return nil
}
