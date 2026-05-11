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

type LoginRequest struct {
	Email    string `json:"email"`
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
		result.RefreshToken,
		7*24*3600,
		"/",
		"",
		true,
		true,
	)
	ctx.JSON(201, gin.H{"accessToken": result.AccessToken, "user": result.User})
}

func (h *AuthHandler) LoginHandler(ctx *gin.Context) {
	req := LoginRequest{}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid Request"})
		return
	}

	result, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		if err.Error() == "invalid credentials" {
			ctx.JSON(401, gin.H{"error": "Invalid credentials!"})
			return
		}
		ctx.JSON(500, gin.H{"error": "Cannot login user for the moment."})
		return
	}

	ctx.SetCookie(
		"refresh_token",
		result.RefreshToken,
		7*24*3600,
		"/",
		"",
		true,
		true,
	)
	ctx.JSON(200, gin.H{"accessToken": result.AccessToken, "user": result.User})

}

func (h *AuthHandler) RefreshHandler(ctx *gin.Context) {
	token, err := ctx.Cookie("refresh_token")
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid Request"})
		return
	}

	refresh, access, err := h.service.Refresh(token)
	if err != nil {
		ctx.JSON(401, gin.H{"error": "Cannot refresh status for the moment."})
		return
	}

	ctx.SetCookie(
		"refresh_token",
		refresh,
		7*24*3600,
		"/",
		"",
		true,
		true,
	)
	ctx.JSON(200, gin.H{"accessToken": access})
}

func (h *AuthHandler) LogoutHandler(ctx *gin.Context) {
	userID := ctx.GetUint("userID")
	if userID == 0 {
		ctx.JSON(400, gin.H{"error": "Invalid Request"})
		return
	}

	err := h.service.Logout(userID)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Cannot refresh status for the moment."})
		return
	}

	ctx.SetCookie("refresh_token", "", -1, "/", "", true, true)
	ctx.JSON(200, gin.H{})
}