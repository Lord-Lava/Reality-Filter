package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/reality-filter/internal/core/domain"
)

const (
	// ArticleCacheTTL is the time-to-live for cached articles
	ArticleCacheTTL = 1 * time.Hour
)

// ArticleCache implements the secondary.ArticleCache interface using Redis
type ArticleCache struct {
	client *redis.Client
}

// NewArticleCache creates a new Redis article cache
func NewArticleCache(client *redis.Client) *ArticleCache {
	return &ArticleCache{
		client: client,
	}
}

// Set stores an article in cache
func (c *ArticleCache) Set(ctx context.Context, article *domain.Article) error {
	data, err := json.Marshal(article)
	if err != nil {
		return err
	}

	key := c.getKey(article.ID.String())
	return c.client.Set(ctx, key, data, ArticleCacheTTL).Err()
}

// Get retrieves an article from cache
func (c *ArticleCache) Get(ctx context.Context, id string) (*domain.Article, error) {
	key := c.getKey(id)
	data, err := c.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var article domain.Article
	if err := json.Unmarshal(data, &article); err != nil {
		return nil, err
	}
	return &article, nil
}

// Delete removes an article from cache
func (c *ArticleCache) Delete(ctx context.Context, id string) error {
	key := c.getKey(id)
	return c.client.Del(ctx, key).Err()
}

// getKey returns the cache key for an article ID
func (c *ArticleCache) getKey(id string) string {
	return "article:" + id
}
