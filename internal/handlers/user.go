package handlers

import (
	"cippus-backend/internal/services"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) PatchMe(ctx *gin.Context) {
	userID := uint(ctx.GetFloat64("userID"))
	if userID == 0 {
		ctx.JSON(400, gin.H{"error": "Invalid Request"})
		return
	}

	req := &PatchMeRequest{}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid Request"})
		return
	}

	userMap := map[string]interface{}{}
	if req.Username != nil {
		userMap["username"] = *req.Username
	}
	if req.Bio != nil {
		userMap["bio"] = *req.Bio
	}

	if req.NotifReply != nil {
		userMap["notif_reply"] = *req.NotifReply
	}
	if req.NotifFollows != nil {
		userMap["notif_follow"] = *req.NotifFollows
	}
	if req.NotifLike != nil {
		userMap["notif_like"] = *req.NotifLike
	}
	if req.NotifMention != nil {
		userMap["notif_mention"] = *req.NotifMention
	}
	if req.NotifDM != nil {
		userMap["notif_dm"] = *req.NotifDM
	}
	if req.NotifPostApproved != nil {
		userMap["notif_post_approved"] = *req.NotifPostApproved
	}

	if len(userMap) == 0 {
		ctx.JSON(400, gin.H{"error": "Invalid Request"})
		return
	}

	err = h.service.PatchMe(userID, userMap)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Cannot update user"})
		return
	}

	ctx.JSON(200, gin.H{})
}
