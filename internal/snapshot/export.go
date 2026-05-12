package snapshot

import (
	"time"

	"github.com/youorg/vaultlens/internal/export"
)

// ExportOptions controls which snapshots are included in the export.
type ExportOptions struct {
	// IDs restricts the export to specific snapshot IDs. If empty, all snapshots are included.
	IDs []string
	// Since filters out snapshots taken before this time.
	Since time.Time
	// IncludeSecrets controls whether secret key/value pairs are written to the export.
	IncludeSecrets bool
}

// ToExportRecords converts snapshots from the store into export.Record values
// that can be passed to an export.Exporter.
func ToExportRecords(store *Store, opts ExportOptions) ([]export.Record, error) {
	if store == nil {
		return nil, ErrNilStore
	}

	snaps, err := store.List()
	if err != nil {
		return nil, err
	}

	allowed := make(map[string]struct{}, len(opts.IDs))
	for _, id := range opts.IDs {
		allowed[id] = struct{}{}
	}

	var records []export.Record
	for _, snap := range snaps {
		if !opts.Since.IsZero() && snap.TakenAt.Before(opts.Since) {
			continue
		}
		if len(allowed) > 0 {
			if _, ok := allowed[snap.ID]; !ok {
				continue
			}
		}

		for path, secret := range snap.Secrets {
			r := export.Record{
				Path:      path,
				Timestamp: snap.TakenAt,
				Meta: map[string]string{
					"snapshot_id": snap.ID,
					"snapshot_tag": snap.Tag,
				},
			}
			if opts.IncludeSecrets {
				r.Data = secret
			}
			records = append(records, r)
		}
	}
	return records, nil
}
