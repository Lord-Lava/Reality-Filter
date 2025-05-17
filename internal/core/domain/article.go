package domain

import (
	"time"

	"github.com/google/uuid"
)

// Article represents the core domain entity for a news article
type Article struct {
	ID        uuid.UUID
	Title     string
	Content   string
	Source    string
	Author    string
	Tags      []string
	CreatedAt time.Time
	UpdatedAt time.Time
	Score     float64
	Flags     []Flag
	Status    ArticleStatus
	MetaData  ArticleMetadata
}

// ArticleMetadata contains extracted information about the article
type ArticleMetadata struct {
	Entities    []Entity
	Sentiment   float64
	Language    string
	WordCount   int
	ReadingTime int // in minutes
}

// Entity represents a named entity in the article content
type Entity struct {
	Type  EntityType
	Value string
}

// EntityType represents different types of named entities
type EntityType string

const (
	EntityTypePerson  EntityType = "PERSON"
	EntityTypePlace   EntityType = "PLACE"
	EntityTypeDate    EntityType = "DATE"
	EntityTypeOrg     EntityType = "ORGANIZATION"
	EntityTypeProduct EntityType = "PRODUCT"
)

// Flag represents issues detected in the article
type Flag struct {
	Type       FlagType
	Confidence float64
	Details    string
	DetectedAt time.Time
	DetectedBy string
}

// FlagType represents different types of issues that can be detected
type FlagType string

const (
	FlagTypeClickbait    FlagType = "CLICKBAIT"
	FlagTypeMisleading   FlagType = "MISLEADING"
	FlagTypeBiased       FlagType = "BIASED"
	FlagTypeUnverified   FlagType = "UNVERIFIED"
	FlagTypeFactualError FlagType = "FACTUAL_ERROR"
	FlagTypeHateSpeech   FlagType = "HATE_SPEECH"
	FlagTypeSpam         FlagType = "SPAM"
)

// ArticleStatus represents the current state of an article
type ArticleStatus string

const (
	ArticleStatusPending  ArticleStatus = "PENDING"
	ArticleStatusAnalyzed ArticleStatus = "ANALYZED"
	ArticleStatusFlagged  ArticleStatus = "FLAGGED"
	ArticleStatusVerified ArticleStatus = "VERIFIED"
	ArticleStatusRejected ArticleStatus = "REJECTED"
)

// NewArticle creates a new Article instance with default values
func NewArticle(title, content, source, author string, tags []string) *Article {
	now := time.Now()
	return &Article{
		ID:        uuid.New(),
		Title:     title,
		Content:   content,
		Source:    source,
		Author:    author,
		Tags:      tags,
		CreatedAt: now,
		UpdatedAt: now,
		Status:    ArticleStatusPending,
		Flags:     make([]Flag, 0),
		MetaData: ArticleMetadata{
			Entities: make([]Entity, 0),
		},
	}
}

// AddFlag adds a new flag to the article
func (a *Article) AddFlag(flagType FlagType, confidence float64, details, detectedBy string) {
	flag := Flag{
		Type:       flagType,
		Confidence: confidence,
		Details:    details,
		DetectedAt: time.Now(),
		DetectedBy: detectedBy,
	}
	a.Flags = append(a.Flags, flag)

	// Update article status if it gets flagged
	if len(a.Flags) > 0 {
		a.Status = ArticleStatusFlagged
	}

	a.UpdatedAt = time.Now()
}

// UpdateScore updates the credibility score of the article
func (a *Article) UpdateScore(score float64) {
	a.Score = score
	a.UpdatedAt = time.Now()
}

// UpdateStatus changes the current status of the article
func (a *Article) UpdateStatus(status ArticleStatus) {
	a.Status = status
	a.UpdatedAt = time.Now()
}

// UpdateMetadata updates the article's metadata
func (a *Article) UpdateMetadata(metadata ArticleMetadata) {
	a.MetaData = metadata
	a.UpdatedAt = time.Now()
}
