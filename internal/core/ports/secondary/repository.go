package secondary

import (
	"context"

	"github.com/reality-filter/internal/core/domain"
)

// ArticleRepository defines the secondary port for article persistence
type ArticleRepository interface {
	// Save persists an article
	Save(ctx context.Context, article *domain.Article) error

	// FindByID retrieves an article by ID
	FindByID(ctx context.Context, id string) (*domain.Article, error)

	// FindFlagged retrieves flagged articles with pagination
	FindFlagged(ctx context.Context, limit, offset int) ([]*domain.Article, error)

	// Update updates an existing article
	Update(ctx context.Context, article *domain.Article) error
}
