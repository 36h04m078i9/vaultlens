package ui

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/yourorg/vaultlens/internal/search"
)

// Renderer handles terminal output for secret search results.
type Renderer struct {
	out    io.Writer
	colors bool
}

// ANSI color codes.
const (
	colorReset  = "\033[0m"
	colorBold   = "\033[1m"
	colorCyan   = "\033[36m"
	colorYellow = "\033[33m"
	colorGray   = "\033[90m"
)

// NewRenderer creates a Renderer writing to w.
// Colors are enabled when w is os.Stdout and the NO_COLOR env var is unset.
func NewRenderer(w io.Writer) *Renderer {
	colors := w == os.Stdout && os.Getenv("NO_COLOR") == ""
	return &Renderer{out: w, colors: colors}
}

// RenderResults prints a formatted list of search results.
func (r *Renderer) RenderResults(results []search.Result, query string) {
	if len(results) == 0 {
		r.printf("%sNo results found for %q.%s\n", colorGray, query, colorReset)
		return
	}
	r.printf("%s%d result(s) for %q:%s\n\n", colorBold, len(results), query, colorReset)
	for i, res := range results {
		r.renderResult(i+1, res)
	}
}

// RenderSecret prints key/value pairs for a single secret path.
func (r *Renderer) RenderSecret(path string, data map[string]interface{}) {
	r.printf("%s%s%s\n", colorCyan, path, colorReset)
	for k, v := range data {
		r.printf("  %s%-24s%s %v\n", colorYellow, k, colorReset, v)
	}
	r.printf("\n")
}

func (r *Renderer) renderResult(idx int, res search.Result) {
	highlighted := res.Highlighted
	if !r.colors {
		// Strip bracket markers used by highlight package.
		highlighted = strings.NewReplacer("[", "", "]", "").Replace(highlighted)
	} else {
		highlighted = strings.ReplaceAll(highlighted, "[", colorBold+colorCyan)
		highlighted = strings.ReplaceAll(highlighted, "]", colorReset)
	}
	r.printf("  %s%d.%s %s %s(score: %d)%s\n", colorBold, idx, colorReset, highlighted, colorGray, res.Score, colorReset)
}

func (r *Renderer) printf(format string, args ...interface{}) {
	if !r.colors {
		format = stripANSI(format)
	}
	fmt.Fprintf(r.out, format, args...)
}

// stripANSI removes ANSI escape sequences from a format string.
func stripANSI(s string) string {
	replacer := strings.NewReplacer(
		colorReset, "", colorBold, "", colorCyan, "",
		colorYellow, "", colorGray, "",
	)
	return replacer.Replace(s)
}
