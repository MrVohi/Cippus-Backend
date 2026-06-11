package handlers

import (
	"cippus-backend/internal/services"

	"github.com/gin-gonic/gin"
)

func NewPushHandler(service *services.PushService, vapidPublic string) *PushHandler {
	return &PushHandler{service: service, vapidPublic: vapidPublic}
}

func (h *PushHandler) GetVapidKey(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"publicKey": h.vapidPublic})
}

func (h *PushHandler) Subscribe(ctx *gin.Context) {
	userID := ctx.GetUint("userID")

	var body struct {
		Endpoint string
		P256dh   string
		Auth     string
	}
	if ctx.ShouldBindJSON(&body) != nil {
		ctx.JSON(400, gin.H{"error": "Invalid Request"})
		return
	}

	err := h.service.Subscribe(userID, body.Endpoint, body.P256dh, body.Auth)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Invalid Request"})
		return
	}

	ctx.JSON(201, gin.H{"ok": true})
}

func (h *PushHandler) Unsubscribe(ctx *gin.Context) {
	userID := ctx.GetUint("userID")

	var body struct{ Endpoint string }
	if ctx.ShouldBindJSON(&body) != nil {
		ctx.JSON(400, gin.H{"error": "Invalid Request"})
		return
	}

	err := h.service.Unsubscribe(userID, body.Endpoint)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Invalid Request"})
		return
	}

	ctx.JSON(200, gin.H{"ok": true})
}
