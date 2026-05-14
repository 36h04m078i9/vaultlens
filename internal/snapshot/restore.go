package snapshot

import (
	"errors"
	"fmt"	
	"time"
)

// RestoreRecord describes a single secret path restored from a snapshot.
type RestoreRecord struct {
	Path      string
	Data      map[string]interface{}
	RestoredAt time.Time
}

// RestoreOptions controls which paths are restored from a snapshot.
type RestoreOptions struct {
	// PathPrefix limits restoration to paths with this prefix. Empty means all.
	PathPrefix string
	// DryRun reports what would be restored without writing.
	DryRun bool
}

// Restorer reads secrets from a snapshot and replays them via a writer.
type Restorer struct {
	store  *Store
	writer SecretWriter
}

// SecretWriter is the interface used to write secrets back to Vault.
type SecretWriter interface {
	WriteSecret(path string, data map[string]interface{}) error
}

// NewRestorer creates a Restorer backed by the given snapshot store and writer.
func NewRestorer(store *Store, writer SecretWriter) (*Restorer, error) {
	if store == nil {
		return nil, errors.New("snapshot: store must not be nil")
	}
	if writer == nil {
		return nil, errors.New("snapshot: writer must not be nil")
	}
	return &Restorer{store: store, writer: writer}, nil
}

// Restore reads snapshot id and writes matching secrets via the writer.
// It returns the list of RestoreRecords that were processed.
func (r *Restorer) Restore(id string, opts RestoreOptions) ([]RestoreRecord, error) {
	if id == "" {
		return nil, errors.New("snapshot: id must not be empty")
	}
	snap, ok := r.store.Get(id)
	if !ok {
		return nil, fmt.Errorf("snapshot: %q not found", id)
	}

	var records []RestoreRecord
	for path, data := range snap.Secrets {
		if opts.PathPrefix != "" && !hasPrefix(path, opts.PathPrefix) {
			continue
		}
		if !opts.DryRun {
			if err := r.writer.WriteSecret(path, data); err != nil {
				return records, fmt.Errorf("snapshot: write %q: %w", path, err)
			}
		}
		records = append(records, RestoreRecord{
			Path:       path,
			Data:       data,
			RestoredAt: time.Now().UTC(),
		})
	}
	return records, nil
}

func hasPrefix(s, prefix string) bool {
	if len(prefix) == 0 {
		return true
	}
	if len(s) < len(prefix) {
		return false
	}
	return s[:len(prefix)] == prefix
}
