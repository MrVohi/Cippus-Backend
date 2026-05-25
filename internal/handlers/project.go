package handlers

import (
	"cippus-backend/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *ProjectHandler) CreateProjectHandler(ctx *gin.Context) {
	req := services.ProjectInput{}
	err := ctx.ShouldBindJSON(&req)

	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid Request"})
		return
	}

	userID := uint(ctx.GetFloat64("userID"))
	project, err := h.service.CreateProject(userID, req)

	if err != nil {
		ctx.JSON(500, gin.H{"error": "Could not create project"})
		return
	}

	ctx.JSON(201, gin.H{"project": project})
}

func (h *ProjectHandler) GetProjectsByUserHandler(ctx *gin.Context) {
	strUserID := ctx.Param("userId")
	userID, err := strconv.ParseUint(strUserID, 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Cannot find project id"})
		return
	}
	project, err := h.service.GetProjectsByUser(uint(userID))
	if err != nil {
		ctx.JSON(400, gin.H{"error": err})
		return
	}

	ctx.JSON(200, gin.H{"project": project})
}

func (h *ProjectHandler) GetProjectsByIdHandler(ctx *gin.Context) {
	strID := ctx.Param("id")
	id, err := strconv.ParseUint(strID, 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Cannot find project id"})
		return
	}
	project, err := h.service.GetProjectById(uint(id))
	if err != nil {
		ctx.JSON(400, gin.H{"error": err})
		return
	}

	ctx.JSON(200, gin.H{"project": project})
}

func (h *ProjectHandler) UpdateProjectHandler(ctx *gin.Context) {
	req := services.ProjectInput{}
	err := ctx.ShouldBindJSON(&req)

	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid Request"})
		return
	}

	strID := ctx.Param("id")
	id, err := strconv.ParseUint(strID, 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Cannot find project id"})
		return
	}

	userID := uint(ctx.GetFloat64("userID"))
	project, err := h.service.UpdateProject(uint(id), userID, ctx.GetString("userRole"), req)

	if err != nil {
		if err.Error() == "could not find project" {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "forbidden" {
			ctx.JSON(403, gin.H{"error": "Could not find project"})
			return
		}
		ctx.JSON(500, gin.H{"error": "Could not find project"})
		return
	}

	ctx.JSON(200, gin.H{"project": project})
}

func (h *ProjectHandler) DeleteProjectHandler(ctx *gin.Context) {
	strID := ctx.Param("id")
	id, err := strconv.ParseUint(strID, 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Cannot find project id"})
		return
	}

	userID := uint(ctx.GetFloat64("userID"))
	err = h.service.DeleteProject(uint(id), userID, ctx.GetString("userRole"))

	if err != nil {
		if err.Error() == "could not find project" {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "forbidden" {
			ctx.JSON(403, gin.H{"error": "Could not delete project"})
			return
		}
		ctx.JSON(500, gin.H{"error": "Could not delete project"})
		return
	}

	ctx.JSON(200, gin.H{})
}

func (h *ProjectHandler) AddPostToProjectHandler(ctx *gin.Context) {
	strProjectID := ctx.Param("id")
	projectID, err := strconv.ParseUint(strProjectID, 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Cannot find project id"})
		return
	}

	strPostID := ctx.Param("postId")
	postID, err := strconv.ParseUint(strPostID, 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Cannot find post id"})
		return
	}

	userID := uint(ctx.GetFloat64("userID"))

	err = h.service.AddPostToProject(uint(projectID), uint(postID), userID)

	if err != nil {
		if err.Error() == "could not find project" || err.Error() == "could not find post" {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "forbidden" {
			ctx.JSON(403, gin.H{"error": "Could not add post to project"})
			return
		}
		ctx.JSON(500, gin.H{"error": "Could not add post to project"})
		return
	}

	ctx.JSON(200, gin.H{})
}

func (h *ProjectHandler) DeletePostFromProjectHandler(ctx *gin.Context) {
	strProjectID := ctx.Param("id")
	projectID, err := strconv.ParseUint(strProjectID, 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Cannot find project id"})
		return
	}

	strPostID := ctx.Param("postId")
	postID, err := strconv.ParseUint(strPostID, 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Cannot find post id"})
		return
	}

	userID := uint(ctx.GetFloat64("userID"))

	err = h.service.RemovePostFromProject(uint(projectID), uint(postID), userID)

	if err != nil {
		if err.Error() == "could not find project" || err.Error() == "could not find post" {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "forbidden" {
			ctx.JSON(403, gin.H{"error": "Could not delete post to project"})
			return
		}
		ctx.JSON(500, gin.H{"error": "Could not delete post to project"})
		return
	}

	ctx.JSON(200, gin.H{})
}

func NewProjectHandler(service *services.ProjectService) *ProjectHandler {
	return &ProjectHandler{service: service}
}
