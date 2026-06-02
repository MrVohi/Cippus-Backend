package handlers

import (
	"cippus-backend/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *CategoryHandler) GetCategoriesHandler(ctx *gin.Context) {
	categories, err := h.service.GetCategories()
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Cannot find category for the moment."})
		return
	}
	resp := make([]CategoryResponse, len(categories))
	for i, c := range categories {
		resp[i] = CategoryResponse{CategoryID: c.ID, Name: c.Name, Description: c.Description}
	}
	ctx.JSON(200, gin.H{"category": resp})
}

func (h *CategoryHandler) CreateCategoryHandler(ctx *gin.Context) {
	req := services.CategoryInput{}
	err := ctx.ShouldBindJSON(&req)

	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid Request"})
		return
	}

	category, err := h.service.CreateCategory(req)

	if err != nil {
		ctx.JSON(500, gin.H{"error": "Could not create category"})
		return
	}

	resp := CategoryResponse{CategoryID: category.ID, Name: category.Name, Description: category.Description}
	ctx.JSON(201, gin.H{"category": resp})
}

func (h *CategoryHandler) UpdateCategoryHandler(ctx *gin.Context) {
	strID := ctx.Param("id")

	if strID == "" {
		ctx.JSON(400, gin.H{"error": "Invalid request."})
		return
	}

	id, err := strconv.ParseUint(strID, 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Cannot find category ID."})
		return
	}

	req := services.CategoryInput{}
	err = ctx.ShouldBindJSON(&req)

	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid Request"})
		return
	}

	category, err := h.service.UpdateCategory(uint(id), req)

	if err != nil {
		if err.Error() == "category not found" {
			ctx.JSON(404, gin.H{"error": "Could not find category"})
			return
		} else {
			ctx.JSON(500, gin.H{"error": "Could not update category"})
			return
		}
	}

	resp := CategoryResponse{CategoryID: category.ID, Name: category.Name, Description: category.Description}
	ctx.JSON(200, gin.H{"category": resp})
}

func (h *CategoryHandler) DeleteCategoryHandler(ctx *gin.Context) {
	strID := ctx.Param("id")

	if strID == "" {
		ctx.JSON(400, gin.H{"error": "Invalid request."})
		return
	}

	id, err := strconv.ParseUint(strID, 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Cannot find category ID."})
		return
	}

	err = h.service.DeleteCategory(uint(id))

	if err != nil {
		if err.Error() == "category not found" {
			ctx.JSON(404, gin.H{"error": "Could not find category"})
			return
		} else {
			ctx.JSON(500, gin.H{"error": "Could not delete category"})
			return
		}
	}
	ctx.JSON(200, gin.H{})
}

func NewCategoryHandler(service *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}