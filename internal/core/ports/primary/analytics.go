package primary

import (
	"context"

	"github.com/reality-filter/internal/core/domain"
)

// AnalyticsProvider defines the primary port for analytics operations
type AnalyticsProvider interface {
	// GetSourceStats retrieves statistics about article sources
	GetSourceStats(ctx context.Context, timeRange string) (map[string]int, error)

	// GetFlagStats retrieves statistics about article flags
	GetFlagStats(ctx context.Context, timeRange string) (map[domain.FlagType]int, error)

	// GetTrendingTopics retrieves trending topics from articles
	GetTrendingTopics(ctx context.Context, limit int) ([]string, error)
}
