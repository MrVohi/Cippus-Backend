package handlers

import (
	"cippus-backend/internal/services"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/gin-gonic/gin"
)

var googleOauthConfig = &oauth2.Config{
	RedirectURL:  "http://localhost:8080/api/v1/auth/google/callback",
	ClientID:     "867066574781-eijfcja8qq7o6suh7v6hqqpinrbrh0dv.apps.googleusercontent.com",
	ClientSecret: "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",
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
	Endpoint:     google.Endpoint,
}

const oauthStateString = "random_state_string"

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

	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie(
		"refresh_token",
		result.RefreshToken,
		7*24*3600,
		"/",
		"",
		true,
		true,
	)
	ctx.JSON(201, gin.H{"accessToken": result.AccessToken, "user": toUserResponse(result.User)})
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

	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie(
		"refresh_token",
		result.RefreshToken,
		7*24*3600,
		"/",
		"",
		true,
		true,
	)
	ctx.JSON(200, gin.H{"accessToken": result.AccessToken, "user": toUserResponse(result.User)})

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

	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie(
		"refresh_token",
		refresh,
		7*24*3600,
		"/",
		"",
		true,
		true,
	)
	ctx.JSON(200, gin.H{"user": toUserResponse(user), "accessToken": access})
}

func (h *AuthHandler) GoogleLoginHandler(c *gin.Context) {
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *AuthHandler) GoogleCallbackHandler(c *gin.Context) {
	state := c.Query("state")
	if state != oauthStateString {
		c.JSON(http.StatusBadRequest, gin.H{"error": "État de sécurité invalide"})
		return
	}
	code := c.Query("code")
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Échec de l'échange de token"})
		return
	}
	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de récupérer les infos utilisateur"})
		return
	}
	defer response.Body.Close()

	var googleUser struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}

	if err := json.NewDecoder(response.Body).Decode(&googleUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Échec du décodage des infos"})
		return
	}

	result, err := h.service.LoginOrCreateWithGoogle(googleUser.Email, googleUser.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de vous connecter via Google pour le moment."})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		"refresh_token",
		result.RefreshToken,
		7*24*3600,
		"/",
		"",
		true,
		true,
	)
	frontendRedirectURL := "http://localhost:3000/auth/success?token=" + result.AccessToken
	c.Redirect(http.StatusTemporaryRedirect, frontendRedirectURL)
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

	ctx.SetSameSite(http.SameSiteLaxMode)
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
