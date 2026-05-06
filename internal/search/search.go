package search

import (
	"context"
	"fmt"
)

// Result represents a single fuzzy-search hit with rendering metadata.
type Result struct {
	Path        string
	Highlighted string
	Score       int
}

// PathLister can enumerate secret paths from a backend.
type PathLister interface {
	ListPaths(ctx context.Context, mount string) ([]string, error)
}

// Service combines path listing with fuzzy search and highlight.
type Service struct {
	lister PathLister
}

// NewService creates a Service backed by the given PathLister.
func NewService(lister PathLister) *Service {
	return &Service{lister: lister}
}

// Search lists all paths under mount and returns fuzzy-matched, highlighted
// results sorted by descending score. An empty query returns all paths with
// score 0.
func (s *Service) Search(ctx context.Context, mount, query string) ([]Result, error) {
	paths, err := s.lister.ListPaths(ctx, mount)
	if err != nil {
		return nil, fmt.Errorf("search: list paths: %w", err)
	}

	if query == "" {
		results := make([]Result, len(paths))
		for i, p := range paths {
			results[i] = Result{Path: p, Highlighted: p, Score: 0}
		}
		return results, nil
	}

	fuzzyHits := Fuzzy(paths, query)
	results := make([]Result, len(fuzzyHits))
	for i, hit := range fuzzyHits {
		highlighted := Highlight([]string{hit.Path}, query)
		hl := hit.Path
		if len(highlighted) > 0 {
			hl = highlighted[0]
		}
		results[i] = Result{
			Path:        hit.Path,
			Highlighted: hl,
			Score:       hit.Score,
		}
	}
	return results, nil
}
