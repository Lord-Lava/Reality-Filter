package mongodb

import (
	"context"
	"time"

	"github.com/reality-filter/internal/core/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ArticleRepository implements the secondary.ArticleRepository interface using MongoDB
type ArticleRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
}

// NewArticleRepository creates a new MongoDB article repository
func NewArticleRepository(client *mongo.Client, database string) *ArticleRepository {
	collection := client.Database(database).Collection("articles")
	return &ArticleRepository{
		client:     client,
		collection: collection,
	}
}

// Save persists an article
func (r *ArticleRepository) Save(ctx context.Context, article *domain.Article) error {
	article.UpdatedAt = time.Now()
	if article.CreatedAt.IsZero() {
		article.CreatedAt = article.UpdatedAt
	}

	opts := options.Update().SetUpsert(true)
	filter := bson.M{"_id": article.ID}
	update := bson.M{"$set": article}

	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	return err
}

// FindByID retrieves an article by ID
func (r *ArticleRepository) FindByID(ctx context.Context, id string) (*domain.Article, error) {
	var article domain.Article
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&article)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &article, nil
}

// FindFlagged retrieves flagged articles with pagination
func (r *ArticleRepository) FindFlagged(ctx context.Context, limit, offset int) ([]*domain.Article, error) {
	opts := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.D{{Key: "updated_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{"status": domain.ArticleStatusFlagged}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var articles []*domain.Article
	if err = cursor.All(ctx, &articles); err != nil {
		return nil, err
	}
	return articles, nil
}

// Update updates an existing article
func (r *ArticleRepository) Update(ctx context.Context, article *domain.Article) error {
	article.UpdatedAt = time.Now()

	filter := bson.M{"_id": article.ID}
	update := bson.M{"$set": article}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}
