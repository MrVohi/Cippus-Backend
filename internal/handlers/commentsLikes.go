package handlers

import (
	"cippus-backend/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CommentLikeHandler struct {
	service *services.CommentLikeService
}

func NewCommentLikeHandler(service *services.CommentLikeService) *CommentLikeHandler {
	return &CommentLikeHandler{service: service}
}

func (h *CommentLikeHandler) ToggleLikeHandler(ctx *gin.Context) {
	strID := ctx.Param("cid")

	if strID == "" {
		ctx.JSON(400, gin.H{"error": "Invalid request."})
		return
	}

	id, err := strconv.ParseUint(strID, 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Cannot find Comment ID."})
		return
	}

	userID := uint(ctx.GetFloat64("userID"))
	if userID == 0 {
		ctx.JSON(400, gin.H{"error": "Invalid Request"})
		return
	}

	err = h.service.ToggleLike(uint(id), userID)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Cannot add a like to Comment"})
		return
	}

	count, userLiked, err := h.service.GetLikesForComment(uint(id), userID)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Cannot fetch Comment's likes"})
		return
	}

	ctx.JSON(200, gin.H{"likes": count, "user_liked": userLiked})
}

func (h *CommentLikeHandler) GetLikesHandler(ctx *gin.Context) {
	strID := ctx.Param("cid")

	if strID == "" {
		ctx.JSON(400, gin.H{"error": "Invalid request."})
		return
	}

	id, err := strconv.ParseUint(strID, 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Cannot find comment ID."})
		return
	}

	userID := uint(ctx.GetFloat64("userID"))

	count, userLiked, err := h.service.GetLikesForComment(uint(id), userID)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Cannot fetch comment's likes"})
		return
	}

	ctx.JSON(200, gin.H{"likes": count, "user_liked": userLiked})
}