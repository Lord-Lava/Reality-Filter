package application

import (
	"context"
	"fmt"

	"github.com/reality-filter/internal/core/domain"
	"github.com/reality-filter/internal/core/ports/primary"
	"github.com/reality-filter/internal/core/ports/secondary"
)

// ArticleAnalyzerService implements the ArticleAnalyzer port
type ArticleAnalyzerService struct {
	repository      secondary.ArticleRepository
	cache           secondary.ArticleCache
	factChecker     secondary.FactChecker
	contentAnalyzer secondary.ContentAnalyzer
	eventPublisher  secondary.EventPublisher
}

// Ensure ArticleAnalyzerService implements primary.ArticleAnalyzer
var _ primary.ArticleAnalyzer = (*ArticleAnalyzerService)(nil)

// NewArticleAnalyzerService creates a new instance of ArticleAnalyzerService
func NewArticleAnalyzerService(
	repository secondary.ArticleRepository,
	cache secondary.ArticleCache,
	factChecker secondary.FactChecker,
	contentAnalyzer secondary.ContentAnalyzer,
	eventPublisher secondary.EventPublisher,
) *ArticleAnalyzerService {
	return &ArticleAnalyzerService{
		repository:      repository,
		cache:           cache,
		factChecker:     factChecker,
		contentAnalyzer: contentAnalyzer,
		eventPublisher:  eventPublisher,
	}
}

// AnalyzeArticle performs comprehensive analysis on an article
func (s *ArticleAnalyzerService) AnalyzeArticle(ctx context.Context, article *domain.Article) error {
	// Step 1: Analyze sentiment
	sentiment, err := s.contentAnalyzer.AnalyzeSentiment(ctx, article.Content)
	if err != nil {
		return fmt.Errorf("failed to analyze sentiment: %w", err)
	}

	// Step 2: Extract entities
	entities, err := s.contentAnalyzer.ExtractEntities(ctx, article.Content)
	if err != nil {
		return fmt.Errorf("failed to extract entities: %w", err)
	}

	// Step 3: Detect bias
	biasFlags, err := s.contentAnalyzer.DetectBias(ctx, article.Content)
	if err != nil {
		return fmt.Errorf("failed to detect bias: %w", err)
	}

	// Step 4: Check facts
	factFlags, err := s.factChecker.CheckFacts(ctx, article)
	if err != nil {
		return fmt.Errorf("failed to check facts: %w", err)
	}

	// Step 5: Get source reputation
	sourceScore, err := s.factChecker.GetSourceReputation(ctx, article.Source)
	if err != nil {
		return fmt.Errorf("failed to get source reputation: %w", err)
	}

	// Update article metadata
	article.UpdateMetadata(domain.ArticleMetadata{
		Entities:    entities,
		Sentiment:   sentiment,
		Language:    "en",                       // TODO: Implement language detection
		WordCount:   len(article.Content),       // TODO: Implement proper word counting
		ReadingTime: len(article.Content) / 200, // Rough estimate: 200 words per minute
	})

	// Add all detected flags
	for _, flag := range biasFlags {
		article.AddFlag(flag.Type, flag.Confidence, flag.Details, "bias_detector")
	}
	for _, flag := range factFlags {
		article.AddFlag(flag.Type, flag.Confidence, flag.Details, "fact_checker")
	}

	// Calculate final credibility score (simple weighted average)
	credibilityScore := calculateCredibilityScore(sourceScore, sentiment, len(article.Flags))
	article.UpdateScore(credibilityScore)

	// Update article status
	if len(article.Flags) > 0 {
		article.UpdateStatus(domain.ArticleStatusFlagged)
	} else {
		article.UpdateStatus(domain.ArticleStatusAnalyzed)
	}

	// Persist the results
	if err := s.repository.Update(ctx, article); err != nil {
		return fmt.Errorf("failed to update article: %w", err)
	}

	// Update cache
	if err := s.cache.Set(ctx, article); err != nil {
		// Log error but don't fail the operation
		fmt.Printf("failed to update cache: %v\n", err)
	}

	// Publish events
	if err := s.eventPublisher.PublishArticleAnalyzed(ctx, article); err != nil {
		// Log error but don't fail the operation
		fmt.Printf("failed to publish article analyzed event: %v\n", err)
	}

	if article.Status == domain.ArticleStatusFlagged {
		if err := s.eventPublisher.PublishArticleFlagged(ctx, article); err != nil {
			fmt.Printf("failed to publish article flagged event: %v\n", err)
		}
	}

	return nil
}

// GetAnalysisResult retrieves the analysis result for an article
func (s *ArticleAnalyzerService) GetAnalysisResult(ctx context.Context, articleID string) (*domain.Article, error) {
	// Try cache first
	article, err := s.cache.Get(ctx, articleID)
	if err == nil {
		return article, nil
	}

	// If not in cache, get from repository
	article, err = s.repository.FindByID(ctx, articleID)
	if err != nil {
		return nil, fmt.Errorf("failed to find article: %w", err)
	}

	// Update cache for next time
	if err := s.cache.Set(ctx, article); err != nil {
		fmt.Printf("failed to update cache: %v\n", err)
	}

	return article, nil
}

// ReprocessArticle triggers reanalysis of an existing article
func (s *ArticleAnalyzerService) ReprocessArticle(ctx context.Context, articleID string) error {
	article, err := s.repository.FindByID(ctx, articleID)
	if err != nil {
		return fmt.Errorf("failed to find article: %w", err)
	}

	// Clear existing analysis results
	article.Flags = make([]domain.Flag, 0)
	article.Score = 0
	article.Status = domain.ArticleStatusPending
	article.MetaData = domain.ArticleMetadata{
		Entities: make([]domain.Entity, 0),
	}

	// Perform fresh analysis
	return s.AnalyzeArticle(ctx, article)
}

// calculateCredibilityScore calculates the final credibility score
func calculateCredibilityScore(sourceScore, sentiment float64, numFlags int) float64 {
	// Simple weighted average:
	// - Source reputation: 40%
	// - Sentiment extremity penalty: 20% (neutral sentiment is better)
	// - Flag penalty: 40% (more flags = lower score)

	// Normalize sentiment to a 0-1 scale where 0.5 is neutral
	sentimentScore := 1.0 - abs(sentiment-0.5)*2

	// Calculate flag penalty (0 flags = 1.0, 5+ flags = 0.0)
	flagPenalty := max(0.0, 1.0-float64(numFlags)/5.0)

	// Weighted average
	score := sourceScore*0.4 + sentimentScore*0.2 + flagPenalty*0.4

	return max(0.0, min(1.0, score))
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func max(x, y float64) float64 {
	if x > y {
		return x
	}
	return y
}

func min(x, y float64) float64 {
	if x < y {
		return x
	}
	return y
}

// CreateArticle implements the ArticleManager interface
func (s *ArticleAnalyzerService) CreateArticle(ctx context.Context, article *domain.Article) error {
	return s.repository.Save(ctx, article)
}

// GetArticle implements the ArticleManager interface
func (s *ArticleAnalyzerService) GetArticle(ctx context.Context, id string) (*domain.Article, error) {
	return s.repository.FindByID(ctx, id)
}

// ListFlaggedArticles implements the ArticleManager interface
func (s *ArticleAnalyzerService) ListFlaggedArticles(ctx context.Context, limit, offset int) ([]*domain.Article, error) {
	return s.repository.FindFlagged(ctx, limit, offset)
}

// UpdateArticleStatus implements the ArticleManager interface
func (s *ArticleAnalyzerService) UpdateArticleStatus(ctx context.Context, id string, status domain.ArticleStatus) error {
	article, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return err
	}
	article.UpdateStatus(status)
	return s.repository.Update(ctx, article)
}
