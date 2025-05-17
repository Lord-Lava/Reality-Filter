package secondary

import (
	"context"

	"github.com/reality-filter/internal/core/domain"
)

// FactChecker defines the secondary port for external fact-checking services
type FactChecker interface {
	// CheckFacts verifies facts in an article
	CheckFacts(ctx context.Context, article *domain.Article) ([]domain.Flag, error)

	// GetSourceReputation gets the reputation score of a news source
	GetSourceReputation(ctx context.Context, source string) (float64, error)
}
