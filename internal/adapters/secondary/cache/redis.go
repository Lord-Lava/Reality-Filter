package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/reality-filter/internal/core/domain"
)

// RedisCache implements the ArticleCache port using Redis
type RedisCache struct {
	client *redis.Client
	ttl    time.Duration
}

// NewRedisCache creates a new Redis cache instance
func NewRedisCache(addr string, password string, db int, ttl time.Duration) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &RedisCache{
		client: client,
		ttl:    ttl,
	}
}

// Set stores an article in Redis
func (c *RedisCache) Set(ctx context.Context, article *domain.Article) error {
	data, err := json.Marshal(article)
	if err != nil {
		return fmt.Errorf("failed to marshal article: %w", err)
	}

	key := fmt.Sprintf("article:%s", article.ID)
	if err := c.client.Set(ctx, key, data, c.ttl).Err(); err != nil {
		return fmt.Errorf("failed to set article in cache: %w", err)
	}

	return nil
}

// Get retrieves an article from Redis
func (c *RedisCache) Get(ctx context.Context, id string) (*domain.Article, error) {
	key := fmt.Sprintf("article:%s", id)
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("article not found in cache")
		}
		return nil, fmt.Errorf("failed to get article from cache: %w", err)
	}

	var article domain.Article
	if err := json.Unmarshal(data, &article); err != nil {
		return nil, fmt.Errorf("failed to unmarshal article: %w", err)
	}

	return &article, nil
}

// Delete removes an article from Redis
func (c *RedisCache) Delete(ctx context.Context, id string) error {
	key := fmt.Sprintf("article:%s", id)
	if err := c.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete article from cache: %w", err)
	}

	return nil
}

// Close closes the Redis connection
func (c *RedisCache) Close() error {
	return c.client.Close()
}
