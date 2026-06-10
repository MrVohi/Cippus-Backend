package handlers

import (
	"cippus-backend/internal/models"
	"cippus-backend/internal/services"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	service *services.ReportService
}

func NewReportHandler(service *services.ReportService) *ReportHandler {
	return &ReportHandler{service: service}
}

func (h *ReportHandler) CreateReportHandler(ctx *gin.Context) {
	userID := uint(ctx.GetFloat64("userID"))
	if userID == 0 {
		ctx.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	var req struct {
		ContentType string `json:"content_type"`
		ContentID   uint   `json:"content_id"`
		Reason      string `json:"reason"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	contentType := models.ContentType(req.ContentType)

	report, err := h.service.CreateReport(userID, contentType, req.ContentID, req.Reason)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Cannot create report"})
		return
	}

	ctx.JSON(201, gin.H{"report": report})
}
