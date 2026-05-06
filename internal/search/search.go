package search

import (
	"context"
	"fmt"

	"github.com/your-org/vaultlens/internal/audit"
	"github.com/your-org/vaultlens/internal/vault"
)

// Service combines Vault listing with fuzzy search and audit logging.
type Service struct {
	client *vault.Client
	logger *audit.Logger
}

// NewService creates a new search Service.
func NewService(client *vault.Client, logger *audit.Logger) *Service {
	return &Service{client: client, logger: logger}
}

// Search lists all secrets under mountPath, applies fuzzy search with query,
// logs the operation, and returns matching results.
func (s *Service) Search(ctx context.Context, mountPath, query string) ([]Result, error) {
	paths, err := s.client.ListSecrets(ctx, mountPath)
	if err != nil {
		s.logger.Log(audit.Entry{
			Operation: "search",
			Path:      mountPath,
			Success:   false,
			Error:     err.Error(),
		})
		return nil, fmt.Errorf("search: listing secrets at %q: %w", mountPath, err)
	}

	results := Fuzzy(paths, query)

	s.logger.Log(audit.Entry{
		Operation: "search",
		Path:      mountPath,
		Success:   true,
		Details:   fmt.Sprintf("query=%q matched=%d total=%d", query, len(results), len(paths)),
	})

	return results, nil
}
