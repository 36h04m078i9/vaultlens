package snapshot

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
)

// RestoreFormatOptions controls how restore results are rendered.
type RestoreFormatOptions struct {
	// DryRun indicates the output should be prefixed with [DRY-RUN].
	DryRun bool
	// ShowData includes secret keys (not values) in the output.
	ShowData bool
}

// FormatRestoreRecords writes a human-readable table of restore records to w.
func FormatRestoreRecords(w io.Writer, records []RestoreRecord, opts RestoreFormatOptions) error {
	if len(records) == 0 {
		_, err := fmt.Fprintln(w, "No secrets restored.")
		return err
	}

	prefix := ""
	if opts.DryRun {
		prefix = "[DRY-RUN] "
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintf(tw, "%sPATH\tKEYS\tRESTORED AT\n", prefix)
	_, _ = fmt.Fprintf(tw, "%s----\t----\t-----------\n", prefix)

	for _, rec := range records {
		keys := "-"
		if opts.ShowData && len(rec.Data) > 0 {
			kl := make([]string, 0, len(rec.Data))
			for k := range rec.Data {
				kl = append(kl, k)
			}
			keys = strings.Join(kl, ",")
		}
		_, _ = fmt.Fprintf(tw, "%s%s\t%s\t%s\n",
			prefix,
			rec.Path,
			keys,
			rec.RestoredAt.Format("2006-01-02T15:04:05Z"),
		)
	}
	return tw.Flush()
}
