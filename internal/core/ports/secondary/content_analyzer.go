package secondary

import (
	"context"

	"github.com/reality-filter/internal/core/domain"
)

// ContentAnalyzer defines the secondary port for content analysis services
type ContentAnalyzer interface {
	// AnalyzeSentiment performs sentiment analysis
	AnalyzeSentiment(ctx context.Context, text string) (float64, error)

	// ExtractEntities extracts named entities
	ExtractEntities(ctx context.Context, text string) ([]domain.Entity, error)

	// DetectBias detects bias in content
	DetectBias(ctx context.Context, text string) ([]domain.Flag, error)
}
