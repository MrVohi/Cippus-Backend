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

	comments, err := h.service.GetCommentsByPostID(uint(id))
	if err != nil {
		ctx.JSON(500, gin.H{"error" : "Cannot find comments for this post"})
		return
	}

	resp := make([]CommentResponse, len(comments))
	for i, c := range comments {
		resp[i] = toCommentResponse(c)
	}
	ctx.JSON(200, gin.H{"comments": resp})
	return
}

func (h *CommentHandler) CreateCommentHandler(ctx *gin.Context) {
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

	type CommentInput struct{
		Content string `json:"content"`
	}
	req := CommentInput{}
	err = ctx.ShouldBindJSON(&req)

	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid Request"})
		return
	}

	comment, err := h.service.CreateComment(uint(id), userID, req.Content)
	if err != nil {
		ctx.JSON(500, gin.H{"error" : "Cannot create comments for this post"})
		return
	}

	ctx.JSON(201, gin.H{"comment": toCommentResponse(comment)})
	return
}

func (h *CommentHandler) UpdateCommentHandler(ctx *gin.Context) {
	strID := ctx.Param("cid")

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

	type CommentInput struct{
		Content string `json:"content"`
	}
	req := CommentInput{}
	err = ctx.ShouldBindJSON(&req)

	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid Request"})
		return
	}

	comment, err := h.service.UpdateComment(uint(id), userID, req.Content)
	if err != nil {
		ctx.JSON(403, gin.H{"error" : "Cannot update comments for this post"})
		return
	}

	ctx.JSON(200, gin.H{"comment": toCommentResponse(comment)})
	return
}

func (h *CommentHandler) DeleteCommentHandler(ctx *gin.Context) {
	strID := ctx.Param("cid")

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

	userRole := ctx.GetString("userRole")
	err = h.service.DeleteComment(uint(id), userID, userRole)
	if err != nil {
		ctx.JSON(500, gin.H{"error" : "Cannot update comments for this post"})
		return
	}

	ctx.JSON(200, gin.H{})
	return
} 