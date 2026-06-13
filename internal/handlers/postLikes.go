package handlers

import (
	"cippus-backend/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PostLikeHandler struct {
	service *services.PostLikeService
}

func NewPostLikeHandler(service *services.PostLikeService) *PostLikeHandler {
	return &PostLikeHandler{service: service}
}

func (h *PostLikeHandler) ToggleLikeHandler(ctx *gin.Context) {
	strID := ctx.Param("id")

	if strID == "" {
		ctx.JSON(400, gin.H{"error": "Invalid request."})
		return
	}

	id, err := strconv.ParseUint(strID, 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Cannot find post ID."})
		return
	}

	userID := uint(ctx.GetFloat64("userID"))
	if userID == 0 {
		ctx.JSON(400, gin.H{"error": "Invalid Request"})
		return
	}

	err = h.service.ToggleLike(uint(id), userID)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Cannot add a like to post"})
		return
	}

	count, userLiked, err := h.service.GetLikesForPost(uint(id), userID)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Cannot fetch post's likes"})
		return
	}

	ctx.JSON(200, gin.H{"likes": count, "user_liked": userLiked})
}

func (h *PostLikeHandler) GetLikesHandler(ctx *gin.Context) {
	strID := ctx.Param("id")

	if strID == "" {
		ctx.JSON(400, gin.H{"error": "Invalid request."})
		return
	}

	id, err := strconv.ParseUint(strID, 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Cannot find post ID."})
		return
	}

	userID := uint(ctx.GetFloat64("userID"))

	count, userLiked, err := h.service.GetLikesForPost(uint(id), userID)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Cannot fetch post's likes"})
		return
	}

	ctx.JSON(200, gin.H{"likes": count, "user_liked": userLiked})
}
