package ui

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/vaultlens/internal/search"
)

func newTestRenderer() (*Renderer, *bytes.Buffer) {
	buf := &bytes.Buffer{}
	r := NewRenderer(buf)
	r.colors = false // deterministic output in tests
	return r, buf
}

func TestRenderResultsEmpty(t *testing.T) {
	r, buf := newTestRenderer()
	r.RenderResults(nil, "foo")
	if !strings.Contains(buf.String(), "No results found") {
		t.Errorf("expected 'No results found', got: %q", buf.String())
	}
}

func TestRenderResultsCount(t *testing.T) {
	r, buf := newTestRenderer()
	results := []search.Result{
		{Path: "secret/a", Highlighted: "secret/[a]", Score: 10},
		{Path: "secret/ab", Highlighted: "secret/[a]b", Score: 8},
	}
	r.RenderResults(results, "a")
	out := buf.String()
	if !strings.Contains(out, "2 result(s)") {
		t.Errorf("expected result count header, got: %q", out)
	}
}

func TestRenderResultsContainsPaths(t *testing.T) {
	r, buf := newTestRenderer()
	results := []search.Result{
		{Path: "kv/prod/db", Highlighted: "kv/prod/[d]b", Score: 7},
	}
	r.RenderResults(results, "d")
	out := buf.String()
	if !strings.Contains(out, "kv/prod/db") {
		t.Errorf("expected path in output, got: %q", out)
	}
}

func TestRenderSecret(t *testing.T) {
	r, buf := newTestRenderer()
	r.RenderSecret("secret/myapp/config", map[string]interface{}{
		"host": "localhost",
		"port": 5432,
	})
	out := buf.String()
	if !strings.Contains(out, "secret/myapp/config") {
		t.Errorf("expected path in secret output, got: %q", out)
	}
	if !strings.Contains(out, "localhost") {
		t.Errorf("expected secret value in output, got: %q", out)
	}
}

func TestRenderResultsScoreShown(t *testing.T) {
	r, buf := newTestRenderer()
	results := []search.Result{
		{Path: "kv/x", Highlighted: "kv/[x]", Score: 42},
	}
	r.RenderResults(results, "x")
	if !strings.Contains(buf.String(), "42") {
		t.Errorf("expected score in output")
	}
}
