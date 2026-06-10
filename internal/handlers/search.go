package handlers

import (
	"cippus-backend/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *SearchHandler) SearchHandler(ctx *gin.Context) {
	q := ctx.Query("q")
	limitStr := ctx.DefaultQuery("limit", "10")
	parsed, err := strconv.Atoi(limitStr)
	if err != nil {
		parsed = 10
	}
	limit := max(1, min(parsed, 50))

	if q == "" {
		ctx.JSON(400, gin.H{"error": "Bad request"})
		return
	}

	results, err := h.service.SemanticSearch(q, limit)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Search failed"})
		return
	}

	resp := make([]SearchPostResponse, len(results))
	for i, p := range results {
		resp[i] = toSearchPostResponse(p)
	}
	ctx.JSON(200, gin.H{"results": resp})
}

func (h *SearchHandler) GetSimilarHandler(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, gin.H{"error": "invalid id"})
		return
	}
	limitStr := ctx.DefaultQuery("limit", "5")
	parsed, err := strconv.Atoi(limitStr)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Could not find original post"})
		return
	}
	limit := max(1, min(parsed, 20))

	results, err := h.service.GetSimilar(uint(id), limit)
	if err != nil {
		if err.Error() == "no embedding" {
			ctx.JSON(404, gin.H{"error": "Could not find similar posts"})
			return
		} else {
			ctx.JSON(500, gin.H{"error": "Cannot find similar posts"})
			return
		}
	}

	resp := make([]SearchPostResponse, len(results))
	for i, p := range results {
		resp[i] = toSearchPostResponse(p)
	}
	ctx.JSON(200, gin.H{"results": resp})
}

func NewSearchHandler(searchService *services.SearchService) *SearchHandler {
	return &SearchHandler{service: searchService}
}
