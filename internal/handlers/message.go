package handlers

import (
	"cippus-backend/internal/models"
	"cippus-backend/internal/services"
	"cippus-backend/internal/ws"
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"
)

func NewMessageHandler(service *services.MessageService, hub *ws.Hub, notif *services.NotificationService) *MessageHandler {
	return &MessageHandler{service: service, Hub: hub, notif: notif}
}

func (h *MessageHandler) Send(ctx *gin.Context) {
	senderID := ctx.GetUint("userID")

	var body struct {
		ReceiverID uint
		Content    string
	}
	if ctx.ShouldBindJSON(&body) != nil {
		ctx.JSON(400, gin.H{"error": "Invalid Request"})
		return
	}

	msg, err := h.service.SendMessage(senderID, body.ReceiverID, body.Content)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Cannot send message"})
		return
	}

	payload, _ := json.Marshal(msg)
	h.Hub.Send <- ws.OutboundMsg{RecipientID: body.ReceiverID, Payload: payload}

	h.notif.Notify(body.ReceiverID, senderID, models.NotificationDM, "", 0)

	ctx.JSON(201, msg)
}

func (h *MessageHandler) GetConversations(ctx *gin.Context) {
	userID := ctx.GetUint("userID")
	conversations, err := h.service.GetConversations(userID)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Cannot load conversations"})
		return
	}
	ctx.JSON(200, conversations)
}

func (h *MessageHandler) GetThread(ctx *gin.Context) {
	userID := ctx.GetUint("userID")
	otherID, err := strconv.ParseUint(ctx.Param("userID"), 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid Request"})
		return
	}
	messages, err := h.service.GetThread(userID, uint(otherID))
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Cannot load thread"})
		return
	}
	ctx.JSON(200, messages)
}

func (h *MessageHandler) MarkAsRead(ctx *gin.Context) {
	userID := ctx.GetUint("userID")
	messageID, err := strconv.ParseUint(ctx.Param("messageID"), 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request"})
		return
	}
	err = h.service.MarkAsRead(uint(messageID), userID)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Cannot mark message as read"})
		return
	}
	ctx.JSON(200, gin.H{"ok": true})
}
