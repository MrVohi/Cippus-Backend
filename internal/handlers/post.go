package handlers

import (
	"cippus-backend/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

func parseOptionalUint(raw string) (*uint, error) {
	if raw == "" {
		return nil, nil
	}

	parsed, err := strconv.ParseUint(raw, 10, 64)
	result := uint(parsed)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func parseOptionalBool(raw string) (*bool, error) {
	if raw == "" {
		return nil, nil
	}

	parsed, err := strconv.ParseBool(raw)
	if err != nil {
		return nil, err
	}

	return &parsed, nil
}

func (h *PostHandler) GetPostsHandler(ctx *gin.Context) {
	filter := services.GetPostsFilter{}

	category, err1 := parseOptionalUint(ctx.Query("category_id"))
	author, err2 := parseOptionalUint(ctx.Query("author_id"))
	stuck, err3 := parseOptionalBool(ctx.Query("stuck"))

	if err1 != nil || err2 != nil || err3 != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request."})
		return
	}

	filter.CategoryID = category
	filter.AuthorID = author
	filter.Stuck = stuck
	filter.Sort = ctx.Query("sort")

	q := ctx.Query("q")
	if q != "" {
		filter.Q = &q
	}

	posts, err := h.service.GetPosts(filter)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Cannot find post for the moment."})
		return
	}
	resp := make([]PostResponse, len(posts))
	for i, p := range posts {
		resp[i] = toPostResponse(p)
	}
	ctx.JSON(200, gin.H{"post": resp})
}

func (h *PostHandler) GetPostsByIdHandler(ctx *gin.Context) {
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

	post, err := h.service.GetPostsById(uint(id))
	if err != nil {
		if err.Error() == "post not found" {
			ctx.JSON(404, gin.H{"error": "Could not find post"})
			return
		} else {
			ctx.JSON(500, gin.H{"error": "Could not find post"})
			return
		}
	}
	ctx.JSON(200, gin.H{"post": toPostResponse(post)})
}

func (h *PostHandler) CreatePostHandler(ctx *gin.Context) {
	title := ctx.PostForm("title")
	content := ctx.PostForm("content")
	if title == "" || content == "" {
		ctx.JSON(400, gin.H{"error": "Invalid Request"})
		return
	}

	var categoryIDs []uint
	for _, s := range ctx.PostFormArray("category_ids") {
		id, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			ctx.JSON(400, gin.H{"error": "Invalid category ID"})
			return
		}
		categoryIDs = append(categoryIDs, uint(id))
	}

	req := services.PostInput{
		Title:       title,
		Content:     content,
		CategoryIDs: categoryIDs,
	}

	userID := uint(ctx.GetFloat64("userID"))
	post, err := h.service.CreatePost(userID, req)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Could not create post"})
		return
	}

	ctx.JSON(201, gin.H{"post": toPostResponse(post)})
}

func (h *PostHandler) UpdatePostHandler(ctx *gin.Context) {
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

	title := ctx.PostForm("title")
	content := ctx.PostForm("content")

	var categoryIDs []uint
	for _, s := range ctx.PostFormArray("category_ids") {
		catID, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			ctx.JSON(400, gin.H{"error": "Invalid category ID"})
			return
		}
		categoryIDs = append(categoryIDs, uint(catID))
	}

	req := services.PostInput{
		Title:       title,
		Content:     content,
		CategoryIDs: categoryIDs,
	}

	userID := uint(ctx.GetFloat64("userID"))
	post, err := h.service.UpdatePost(uint(id), userID, ctx.GetString("userRole"), req)

	if err != nil {
		if err.Error() == "post not found" {
			ctx.JSON(404, gin.H{"error": "Could not find post"})
			return
		}
		if err.Error() == "forbidden" {
			ctx.JSON(403, gin.H{"error": "Could not update post"})
			return
		} else {
			ctx.JSON(500, gin.H{"error": "Could not update post"})
			return
		}
	}

	ctx.JSON(200, gin.H{"post": toPostResponse(post)})
}

func (h *PostHandler) DeletePostHandler(ctx *gin.Context) {
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
	err = h.service.DeletePost(uint(id), userID, ctx.GetString("userRole"))

	if err != nil {
		if err.Error() == "post not found" {
			ctx.JSON(404, gin.H{"error": "Could not find post"})
			return
		}
		if err.Error() == "forbidden" {
			ctx.JSON(403, gin.H{"error": "Could not delete post"})
			return
		} else {
			ctx.JSON(500, gin.H{"error": "Could not delete post"})
			return
		}
	}
	ctx.JSON(200, gin.H{})
}

func NewPostHandler(service *services.PostService) *PostHandler {
	return &PostHandler{service: service}
}
