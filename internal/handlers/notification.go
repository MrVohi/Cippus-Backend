package handlers

import (
	"cippus-backend/internal/services"

	"github.com/gin-gonic/gin"
)

func NewNotificationHandler(service *services.NotificationService) *NotificationHandler {
	return &NotificationHandler{service: service}
}

func (h *NotificationHandler) GetNotifications(ctx *gin.Context) {
	userID := ctx.GetUint("userID")
	notifications, err := h.service.GetNotifications(userID)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Cannot load notifications"})
		return
	}
	ctx.JSON(200, notifications)
}
