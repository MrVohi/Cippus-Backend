package handlers

import (
	"cippus-backend/internal/services"

	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthHandler struct {
	service *services.AuthService
}

func (h *AuthHandler) RegisterHandler(ctx *gin.Context) {
	req := RegisterRequest{}
	err := ctx.ShouldBindJSON(&req)

	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid Request"})
		return
	}

	result, authErr := h.service.Register(req.Email, req.Username, req.Password)
	if authErr != nil {
		if authErr.Error() == "email already exists" {
			ctx.JSON(409, gin.H{"error": "Email already exists!"})
			return
		}
		ctx.JSON(500, gin.H{"error": "Cannot register user for the moment."})
		return
	}

	ctx.SetCookie(
		"refresh_token",
		string(result.RefreshToken),
		7*24*3600,
		"/",
		"",
		true,
		true,
	)
	ctx.JSON(201, gin.H{"accessToken": result.AccessToken, "user": result.User})
}
