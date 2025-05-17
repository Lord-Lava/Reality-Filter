package primary

import (
	"context"

	"github.com/reality-filter/internal/core/domain"
)

// ArticleManager defines the primary port for article management operations
type ArticleManager interface {
	// CreateArticle creates a new article in the system
	CreateArticle(ctx context.Context, article *domain.Article) error

	// GetArticle retrieves an article by ID
	GetArticle(ctx context.Context, articleID string) (*domain.Article, error)

	// UpdateArticleStatus updates the status of an article
	UpdateArticleStatus(ctx context.Context, articleID string, status domain.ArticleStatus) error

	// ListFlaggedArticles retrieves a list of flagged articles
	ListFlaggedArticles(ctx context.Context, limit, offset int) ([]*domain.Article, error)
}
