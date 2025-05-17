package secondary

import (
	"context"

	"github.com/reality-filter/internal/core/domain"
)

// EventPublisher defines the secondary port for event publishing
type EventPublisher interface {
	// PublishArticleAnalyzed publishes an article analyzed event
	PublishArticleAnalyzed(ctx context.Context, article *domain.Article) error

	// PublishArticleFlagged publishes an article flagged event
	PublishArticleFlagged(ctx context.Context, article *domain.Article) error
}
