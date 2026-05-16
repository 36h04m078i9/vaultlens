package snapshot

import (
	"fmt"
	"strings"
)

// FormatRollbackRecords returns a human-readable summary of rollback records.
func FormatRollbackRecords(records []RollbackRecord, showErrors bool) string {
	if len(records) == 0 {
		return "(no paths matched)"
	}

	var sb strings.Builder
	succeeded := 0
	failed := 0
	skipped := 0

	for _, r := range records {
		prefix := "[restored]"
		switch {
		case r.DryRun:
			prefix = "[dry-run] "
			skipped++
		case r.Err != "":
			prefix = "[failed] "
			failed++
		default:
			succeeded++
		}
		sb.WriteString(fmt.Sprintf("%s %s\n", prefix, r.Path))
		if showErrors && r.Err != "" {
			sb.WriteString(fmt.Sprintf("         error: %s\n", r.Err))
		}
	}

	sb.WriteString(fmt.Sprintf("\nsummary: %d restored, %d failed, %d dry-run\n",
		succeeded, failed, skipped))
	return sb.String()
}
