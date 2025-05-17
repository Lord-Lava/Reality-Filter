package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/reality-filter/internal/core/domain"
	"github.com/reality-filter/internal/core/ports/primary"
)

// Handler handles HTTP requests for the article analysis API
type Handler struct {
	analyzer primary.ArticleAnalyzer
	manager  primary.ArticleManager
}

// NewHandler creates a new HTTP handler
func NewHandler(analyzer primary.ArticleAnalyzer, manager primary.ArticleManager) *Handler {
	return &Handler{
		analyzer: analyzer,
		manager:  manager,
	}
}

// RegisterRoutes registers the HTTP routes with the Gin engine
func (h *Handler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	{
		api.POST("/articles", h.CreateArticle)
		api.GET("/articles/:id", h.GetArticle)
		api.POST("/articles/:id/analyze", h.AnalyzeArticle)
		api.GET("/articles/:id/analysis", h.GetAnalysisResult)
		api.POST("/articles/:id/reprocess", h.ReprocessArticle)
		api.GET("/articles/flagged", h.ListFlaggedArticles)
	}
}

// CreateArticle godoc
// @Summary Create a new article
// @Description Submit a new article for analysis
// @Tags Articles
// @Accept json
// @Produce json
// @Param article body domain.Article true "Article to create"
// @Success 201 {object} map[string]interface{} "Returns article ID and status"
// @Failure 400 {object} map[string]string "Invalid request payload"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /articles [post]
func (h *Handler) CreateArticle(c *gin.Context) {
	var request struct {
		Title   string   `json:"title" binding:"required"`
		Content string   `json:"content" binding:"required"`
		Source  string   `json:"source" binding:"required"`
		Author  string   `json:"author" binding:"required"`
		Tags    []string `json:"tags"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	article := domain.NewArticle(
		request.Title,
		request.Content,
		request.Source,
		request.Author,
		request.Tags,
	)

	if err := h.manager.CreateArticle(c.Request.Context(), article); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"articleId": article.ID,
		"status":    article.Status,
	})
}

// GetArticle godoc
// @Summary Get article details
// @Description Retrieve details of a specific article
// @Tags Articles
// @Accept json
// @Produce json
// @Param id path string true "Article ID"
// @Success 200 {object} domain.Article
// @Failure 404 {object} map[string]string "Article not found"
// @Router /articles/{id} [get]
func (h *Handler) GetArticle(c *gin.Context) {
	articleID := c.Param("id")

	article, err := h.manager.GetArticle(c.Request.Context(), articleID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}

	c.JSON(http.StatusOK, article)
}

// AnalyzeArticle godoc
// @Summary Analyze an article
// @Description Trigger analysis of an existing article
// @Tags Analysis
// @Accept json
// @Produce json
// @Param id path string true "Article ID"
// @Success 200 {object} map[string]interface{} "Analysis results"
// @Failure 404 {object} map[string]string "Article not found"
// @Failure 500 {object} map[string]string "Analysis failed"
// @Router /articles/{id}/analyze [post]
func (h *Handler) AnalyzeArticle(c *gin.Context) {
	articleID := c.Param("id")

	article, err := h.manager.GetArticle(c.Request.Context(), articleID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}

	if err := h.analyzer.AnalyzeArticle(c.Request.Context(), article); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"articleId": article.ID,
		"score":     article.Score,
		"flags":     article.Flags,
		"status":    article.Status,
	})
}

// GetAnalysisResult godoc
// @Summary Get analysis results
// @Description Retrieve analysis results for an article
// @Tags Analysis
// @Accept json
// @Produce json
// @Param id path string true "Article ID"
// @Success 200 {object} map[string]interface{} "Analysis results with metadata"
// @Failure 404 {object} map[string]string "Analysis result not found"
// @Router /articles/{id}/analysis [get]
func (h *Handler) GetAnalysisResult(c *gin.Context) {
	articleID := c.Param("id")

	article, err := h.analyzer.GetAnalysisResult(c.Request.Context(), articleID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Analysis result not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"articleId": article.ID,
		"score":     article.Score,
		"flags":     article.Flags,
		"status":    article.Status,
		"metadata":  article.MetaData,
	})
}

// ReprocessArticle godoc
// @Summary Reprocess an article
// @Description Trigger reanalysis of an existing article
// @Tags Analysis
// @Accept json
// @Produce json
// @Param id path string true "Article ID"
// @Success 202 "Reprocessing request accepted"
// @Failure 500 {object} map[string]string "Reprocessing failed"
// @Router /articles/{id}/reprocess [post]
func (h *Handler) ReprocessArticle(c *gin.Context) {
	articleID := c.Param("id")

	if err := h.analyzer.ReprocessArticle(c.Request.Context(), articleID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusAccepted)
}

// ListFlaggedArticles godoc
// @Summary List flagged articles
// @Description Retrieve a list of articles that have been flagged during analysis
// @Tags Articles
// @Accept json
// @Produce json
// @Param limit query int false "Maximum number of articles to return (default: 10)"
// @Param offset query int false "Number of articles to skip (default: 0)"
// @Success 200 {object} map[string]interface{} "List of flagged articles"
// @Failure 400 {object} map[string]string "Invalid query parameters"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /articles/flagged [get]
func (h *Handler) ListFlaggedArticles(c *gin.Context) {
	limit := 10 // Default limit
	offset := 0 // Default offset

	if limitParam := c.Query("limit"); limitParam != "" {
		if _, err := fmt.Sscanf(limitParam, "%d", &limit); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
			return
		}
	}

	if offsetParam := c.Query("offset"); offsetParam != "" {
		if _, err := fmt.Sscanf(offsetParam, "%d", &offset); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset parameter"})
			return
		}
	}

	articles, err := h.manager.ListFlaggedArticles(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"articles": articles,
		"limit":    limit,
		"offset":   offset,
	})
}
