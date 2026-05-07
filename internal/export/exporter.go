// Package export provides functionality to export secret paths and metadata
// from VaultLens sessions in various formats for offline review or sharing.
package export

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// Record represents a single exportable secret path entry.
type Record struct {
	Path      string    `json:"path"`
	Note      string    `json:"note,omitempty"`
	ExportedAt time.Time `json:"exported_at"`
}

// Format defines the output format for exports.
type Format string

const (
	FormatJSON Format = "json"
	FormatCSV  Format = "csv"
)

// Exporter writes records to an io.Writer in the configured format.
type Exporter struct {
	format Format
	w      io.Writer
}

// NewExporter creates an Exporter that writes to w using the given format.
func NewExporter(w io.Writer, format Format) *Exporter {
	return &Exporter{format: format, w: w}
}

// Export writes all records to the underlying writer.
func (e *Exporter) Export(records []Record) error {
	switch e.format {
	case FormatJSON:
		return e.writeJSON(records)
	case FormatCSV:
		return e.writeCSV(records)
	default:
		return fmt.Errorf("export: unsupported format %q", e.format)
	}
}

func (e *Exporter) writeJSON(records []Record) error {
	enc := json.NewEncoder(e.w)
	enc.SetIndent("", "  ")
	return enc.Encode(records)
}

func (e *Exporter) writeCSV(records []Record) error {
	w := csv.NewWriter(e.w)
	if err := w.Write([]string{"path", "note", "exported_at"}); err != nil {
		return fmt.Errorf("export: writing csv header: %w", err)
	}
	for _, r := range records {
		row := []string{r.Path, r.Note, r.ExportedAt.UTC().Format(time.RFC3339)}
		if err := w.Write(row); err != nil {
			return fmt.Errorf("export: writing csv row: %w", err)
		}
	}
	w.Flush()
	return w.Error()
}
