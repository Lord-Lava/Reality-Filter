package primary

import (
	"context"

	"github.com/reality-filter/internal/core/domain"
)

// ArticleAnalyzer defines the primary port for article analysis operations
type ArticleAnalyzer interface {
	// AnalyzeArticle performs full analysis on an article
	AnalyzeArticle(ctx context.Context, article *domain.Article) error

	// GetAnalysisResult retrieves the analysis result for an article
	GetAnalysisResult(ctx context.Context, articleID string) (*domain.Article, error)

	// ReprocessArticle triggers reanalysis of an existing article
	ReprocessArticle(ctx context.Context, articleID string) error
}
