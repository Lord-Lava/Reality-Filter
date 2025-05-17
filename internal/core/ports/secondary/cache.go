package secondary

import (
	"context"

	"github.com/reality-filter/internal/core/domain"
)

// ArticleCache defines the secondary port for article caching
type ArticleCache interface {
	// Set stores an article in cache
	Set(ctx context.Context, article *domain.Article) error

	// Get retrieves an article from cache
	Get(ctx context.Context, id string) (*domain.Article, error)

	// Delete removes an article from cache
	Delete(ctx context.Context, id string) error
}
