package handlers

import (
	"cippus-backend/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *MinioHandler) UploadImageHandler(ctx *gin.Context) {
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
	post, err := h.postServices.GetPostsById(uint(id))
	if err != nil {
		if err.Error() == "post not found" {
			ctx.JSON(404, gin.H{"error": "Could not find post"})
			return
		}
		ctx.JSON(500, gin.H{"error": "Could not find post"})
		return
	}

	if post.UserId != userID {
		ctx.JSON(403, gin.H{"error": "Could not find post"})
		return
	}

	file, err := ctx.FormFile("image")
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Could not find file"})
		return
	}

	opened, err := file.Open()
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Could not open file"})
		return
	}

	defer opened.Close()

	url, err := h.service.UploadFile(opened, file.Size, file.Filename)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Could not upload file"})
		return
	}

	err = h.postServices.UpdatePostImage(post.ID, url)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Could not update post"})
		return
	}

	ctx.JSON(200, gin.H{"imageUrl": url})
}

func (h *MinioHandler) DeleteImageHandler(ctx *gin.Context) {
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

	post, err := h.postServices.GetPostsById(uint(id))
	if err != nil {
		if err.Error() == "post not found" {
			ctx.JSON(404, gin.H{"error": "Could not find post"})
			return
		}
		ctx.JSON(500, gin.H{"error": "Could not find post"})
		return
	}

	if post.UserId != userID {
		ctx.JSON(403, gin.H{"error": "Could not find post"})
		return
	}

	if post.ImageURL == "" {
		ctx.JSON(400, gin.H{"error": "No image to delete"})
		return
	}

	err = h.service.DeleteFile(post.ImageURL)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Could not delete file"})
		return
	}

	err = h.postServices.UpdatePostImage(uint(id), "")
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Could not update post"})
		return
	}

	ctx.JSON(200, gin.H{})
}

func NewMinioHandler (service *services.MinioService, postServices *services.PostService) *MinioHandler {
	return &MinioHandler{service: service, postServices: postServices}
}
