package secondary

import (
	"context"

	"github.com/reality-filter/internal/core/domain"
)

// AnalyticsStore defines the secondary port for analytics data storage
type AnalyticsStore interface {
	// StoreArticleEvent stores an article-related event
	StoreArticleEvent(ctx context.Context, articleID string, eventType string, metadata map[string]interface{}) error

	// GetSourceStats retrieves source statistics
	GetSourceStats(ctx context.Context, timeRange string) (map[string]int, error)

	// GetFlagStats retrieves flag statistics
	GetFlagStats(ctx context.Context, timeRange string) (map[domain.FlagType]int, error)
}
