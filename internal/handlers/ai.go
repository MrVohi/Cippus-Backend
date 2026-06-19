package handlers

import (
	"cippus-backend/internal/services"

	"github.com/gin-gonic/gin"
)

type AIHandler struct{
	es *services.EmbeddingService
}

type ImproveRequest struct {
	Improved string `json:"text" binding:"required,min=1"`
}

type ImproveResponse struct { 
	Improved string `json:"improved"`
}

func NewAIHandler(es *services.EmbeddingService) *AIHandler {
	return &AIHandler{es: es}
}

func (h *AIHandler) ImproveHandler(ctx *gin.Context) {
	var req ImproveRequest 
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(400, gin.H{"error":"AI Service unavailable"})
		return
	}

	improved, err := h.es.ImproveText(req.Improved)
	if err != nil {
		ctx.JSON(500, gin.H{"error":"AI service unavailable"})
		return
	}

	ctx.JSON(200, ImproveResponse{ Improved: improved})
}

