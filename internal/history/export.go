package history

import (
	"time"

	"github.com/youorg/vaultlens/internal/export"
)

// ExportOptions controls how history entries are exported.
type ExportOptions struct {
	// Limit caps the number of entries included; 0 means all.
	Limit int
	// Since filters entries accessed at or after this time. Zero means no filter.
	Since time.Time
}

// ToExportRecords converts history entries to export.Record slices,
// applying any options provided. The most-recent entries are first.
func ToExportRecords(h *History, opts ExportOptions) []export.Record {
	entries := h.List()

	var filtered []Entry
	for _, e := range entries {
		if !opts.Since.IsZero() && e.AccessedAt.Before(opts.Since) {
			continue
		}
		filtered = append(filtered, e)
	}

	if opts.Limit > 0 && len(filtered) > opts.Limit {
		filtered = filtered[:opts.Limit]
	}

	records := make([]export.Record, 0, len(filtered))
	for _, e := range filtered {
		records = append(records, export.Record{
			Path:      e.Path,
			Note:      "",
			Timestamp: e.AccessedAt,
		})
	}
	return records
}
