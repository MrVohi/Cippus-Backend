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
	userID := ctx.GetUint("userID")
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
