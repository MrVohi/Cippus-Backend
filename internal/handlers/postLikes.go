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

func (h *PostLikeHandler) LikePostHandler(ctx *gin.Context) {
	postID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid post ID"})
		return
	}

	userID := uint(ctx.GetFloat64("userID"))
	if userID == 0 {
		ctx.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	if err := h.service.LikePost(uint(postID), userID); err != nil {
		ctx.JSON(500, gin.H{"error": "Cannot like post"})
		return
	}

	count, err := h.service.GetLikeCount(uint(postID))
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Cannot get like count"})
		return
	}

	ctx.JSON(200, gin.H{"likes": count})
}

func (h *PostLikeHandler) GetLikesHandler(ctx *gin.Context) {
	postID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid post ID"})
		return
	}

	count, err := h.service.GetLikeCount(uint(postID))
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Cannot get like count"})
		return
	}

	ctx.JSON(200, gin.H{"likes": count})
}
