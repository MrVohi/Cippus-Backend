package handlers

import (
	"cippus-backend/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	service *services.CommentService
}

func NewCommentHandler(service *services.CommentService) *CommentHandler {
	return &CommentHandler{service: service}
}

func (h *CommentHandler) GetCommentsByPostIDHandler(ctx *gin.Context) {
	postID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid post ID"})
		return
	}

	comments, err := h.service.GetCommentsByPostID(uint(postID))
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Cannot get comments"})
		return
	}

	ctx.JSON(200, gin.H{"comments": comments})
}

func (h *CommentHandler) CreateCommentHandler(ctx *gin.Context) {
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

	var req struct {
		Content string `json:"content"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	comment, err := h.service.CreateComment(uint(postID), userID, req.Content)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Cannot create comment"})
		return
	}

	ctx.JSON(201, gin.H{"comment": comment})
}

func (h *CommentHandler) DeleteCommentHandler(ctx *gin.Context) {
	commentID, err := strconv.Atoi(ctx.Param("commentId"))
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid comment ID"})
		return
	}

	userID := uint(ctx.GetFloat64("userID"))
	if userID == 0 {
		ctx.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	if err := h.service.DeleteComment(uint(commentID), userID); err != nil {
		ctx.JSON(500, gin.H{"error": "Cannot delete comment"})
		return
	}

	ctx.JSON(200, gin.H{})
}
