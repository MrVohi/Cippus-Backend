package handlers

import (
	"log/slog"
	"strings"

	"cippus-backend/internal/services"

	"github.com/gin-gonic/gin"
)

func (h *AuthHandler) RegisterHandler(ctx *gin.Context) {
	req := RegisterRequest{}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid Request"})
		return
	}

	ok := h.service.VerifyRecaptcha(req.CaptchaToken)
	if !ok {
		ctx.JSON(400, gin.H{"error": "Captcha failed"})
		return
	}

	result, err := h.service.Register(req.Email, req.Username, req.Password)
	if err != nil {
		switch {
		case err.Error() == "email already exists":
			ctx.JSON(409, gin.H{"error": "Email already exists!"})
		case strings.HasPrefix(err.Error(), "username already exists"):
			suggestion := strings.TrimPrefix(err.Error(), "username already exists: ")
			ctx.JSON(409, gin.H{
				"error":      "Username already taken!",
				"suggestion": suggestion,
			})
		default:
			ctx.JSON(500, gin.H{"error": "Cannot register user for the moment."})
		}
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

	refresh, access, user, err := h.service.Refresh(token)
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
	ctx.JSON(200, gin.H{"user": user, "accessToken": access})
}

func (h *AuthHandler) LogoutHandler(ctx *gin.Context) {
	userID := uint(ctx.GetFloat64("userID"))
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

func (h *AuthHandler) PasswordResetRequestHandler(ctx *gin.Context) {
	req := PasswordResetRequestRequest{}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid Request"})
		return
	}

	if err := h.service.PasswordResetRequest(req.Email); err != nil {
		slog.Warn("password reset request failed", "err", err)
	}
	ctx.JSON(200, gin.H{})
}

func (h *AuthHandler) PasswordResetConfirmHandler(ctx *gin.Context) {
	req := PasswordResetConfirmRequest{}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid Request"})
		return
	}

	err = h.service.PasswordResetConfirm(req.Token, req.NewPassword)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Cannot confirm new password for the moment."})
		return
	}
	ctx.JSON(200, gin.H{})
}

func NewAuthHandler(service *services.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}
