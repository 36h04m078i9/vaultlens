package export

import (
	"time"

	"github.com/yourorg/vaultlens/internal/bookmark"
)

// Builder converts bookmark entries into export Records.
type Builder struct {
	now func() time.Time
}

// NewBuilder returns a Builder that stamps records with the current time.
func NewBuilder() *Builder {
	return &Builder{now: time.Now}
}

// FromBookmarks converts a slice of bookmarks to export Records.
func (b *Builder) FromBookmarks(bms []bookmark.Bookmark) []Record {
	records := make([]Record, 0, len(bms))
	for _, bm := range bms {
		records = append(records, Record{
			Path:       bm.Path,
			Note:       bm.Note,
			ExportedAt: b.now(),
		})
	}
	return records
}

// FromPaths converts a plain list of secret paths to export Records with no notes.
func (b *Builder) FromPaths(paths []string) []Record {
	records := make([]Record, 0, len(paths))
	for _, p := range paths {
		records = append(records, Record{
			Path:       p,
			ExportedAt: b.now(),
		})
	}
	return records
}
