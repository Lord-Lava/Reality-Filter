// @title Reality Filter API
// @version 1.0
// @description API for analyzing and fact-checking news articles
// @host localhost:8080
// @BasePath /api/v1
// @contact.name Reality Filter Team

package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/reality-filter/docs"
	"github.com/reality-filter/internal/adapters/primary/http/handler"
	"github.com/reality-filter/internal/adapters/secondary/mongodb"
	redisadapter "github.com/reality-filter/internal/adapters/secondary/redis"
	"github.com/reality-filter/internal/application"
	"github.com/reality-filter/internal/core/domain"
	"github.com/reality-filter/pkg/config"
	"github.com/reality-filter/pkg/logger"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

// @securityDefinitions.basic BasicAuth

func init() {
	docs.SwaggerInfo.Title = "Reality Filter API"
	docs.SwaggerInfo.Description = "API for analyzing and fact-checking news articles"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	logConfig := logger.Config{
		LogLevel:         "debug",   // Set to debug during development
		Development:      true,      // Enable development mode for more verbose logging
		Encoding:         "console", // Use console encoding for readable logs during development
		OutputPaths:      []string{"stdout", "logs/reality-filter.log"},
		ErrorOutputPaths: []string{"stderr"},
	}

	if err := logger.Init(logConfig); err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
}

func main() {
	defer logger.Sync()

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Failed to load configuration", zap.Error(err))
	}

	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.GetMongoDBConfig().GetURI()))
	if err != nil {
		logger.Fatal("Failed to connect to MongoDB", zap.Error(err))
	}
	defer mongoClient.Disconnect(context.Background())

	redisConfig := cfg.GetRedisConfig()
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisConfig.GetAddr(),
		Password: redisConfig.GetPassword(),
		DB:       redisConfig.GetDB(),
	})
	defer redisClient.Close()

	repository := mongodb.NewArticleRepository(mongoClient, cfg.MongoDB.Database)
	cache := redisadapter.NewArticleCache(redisClient)

	// TODO: Implement these interfaces
	var (
		factChecker     = &mockFactChecker{}     // Replace with actual implementation
		contentAnalyzer = &mockContentAnalyzer{} // Replace with actual implementation
		eventPublisher  = &mockEventPublisher{}  // Replace with actual implementation
	)

	analyzer := application.NewArticleAnalyzerService(
		repository,
		cache,
		factChecker,
		contentAnalyzer,
		eventPublisher,
	)

	handler := handler.NewHandler(analyzer, analyzer) // Using analyzer as both ArticleAnalyzer and ArticleManager

	gin.SetMode(gin.ReleaseMode)
	router := gin.New() // Use New() instead of Default() to avoid using the default logger

	// Use our custom logger for Gin
	router.Use(gin.Recovery())
	router.Use(func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		logger.Info("HTTP Request",
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", time.Since(start)),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
		)
	})

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.URL("http://localhost:8080/swagger/doc.json"),
		ginSwagger.DefaultModelsExpandDepth(-1),
	))

	handler.RegisterRoutes(router)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		logger.Info("Starting server", zap.String("address", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited successfully")
}

// Mock implementations for remaining interfaces
type mockFactChecker struct{}

func (m *mockFactChecker) CheckFacts(ctx context.Context, article *domain.Article) ([]domain.Flag, error) {
	return nil, nil
}

func (m *mockFactChecker) GetSourceReputation(ctx context.Context, source string) (float64, error) {
	return 0.8, nil
}

type mockContentAnalyzer struct{}

func (m *mockContentAnalyzer) AnalyzeSentiment(ctx context.Context, text string) (float64, error) {
	return 0.5, nil
}

func (m *mockContentAnalyzer) ExtractEntities(ctx context.Context, text string) ([]domain.Entity, error) {
	return nil, nil
}

func (m *mockContentAnalyzer) DetectBias(ctx context.Context, text string) ([]domain.Flag, error) {
	return nil, nil
}

type mockEventPublisher struct{}

func (m *mockEventPublisher) PublishArticleAnalyzed(ctx context.Context, article *domain.Article) error {
	return nil
}

func (m *mockEventPublisher) PublishArticleFlagged(ctx context.Context, article *domain.Article) error {
	return nil
}
