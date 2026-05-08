package diff

import (
	"fmt"
	"strings"
)

// Formatter controls how diff output is rendered as text.
type Formatter struct {
	MaskValues bool
}

const masked = "[redacted]"

// Format renders a slice of Changes into a human-readable string.
func (f Formatter) Format(changes []Change) string {
	if len(changes) == 0 {
		return "no changes detected"
	}

	var sb strings.Builder
	for _, c := range changes {
		switch c.Kind {
		case Added:
			nv := f.maybeRedact(c.NewValue)
			fmt.Fprintf(&sb, "+ %s = %s\n", c.Path, nv)
		case Removed:
			ov := f.maybeRedact(c.OldValue)
			fmt.Fprintf(&sb, "- %s = %s\n", c.Path, ov)
		case Changed:
			ov := f.maybeRedact(c.OldValue)
			nv := f.maybeRedact(c.NewValue)
			fmt.Fprintf(&sb, "~ %s: %s -> %s\n", c.Path, ov, nv)
		}
	}

	s := Summary(changes)
	fmt.Fprintf(&sb, "\nsummary: +%d added, -%d removed, ~%d changed",
		s[Added], s[Removed], s[Changed])

	return sb.String()
}

func (f Formatter) maybeRedact(v string) string {
	if f.MaskValues {
		return masked
	}
	return v
}
