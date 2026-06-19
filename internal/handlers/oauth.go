package handlers

import (
	"cippus-backend/internal/models"
	"cippus-backend/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

func NewOAuthHandler(service *services.OAuthService) *OAuthHandler {
	return &OAuthHandler{service: service}
}

func (h *OAuthHandler) HandleGoogleInit(ctx *gin.Context) {
	state, err := h.service.GenerateState()
	if err != nil {
		ctx.JSON(500, gin.H{"error": "couldn't init google oauth"})
		return
	}

	ctx.SetCookie(
		"oauth_state",
		state,
		600,
		"/",
		"",
		true,
		true,
	)

	url := h.service.GoogleConfig.AuthCodeURL(state)
	ctx.Redirect(302, url)
}

func (h *OAuthHandler) HandleGithubInit(ctx *gin.Context) {
	state, err := h.service.GenerateState()
	if err != nil {
		ctx.JSON(500, gin.H{"error": "couldn't init github oauth"})
		return
	}

	ctx.SetCookie(
		"oauth_state",
		state,
		600,
		"/",
		"",
		true,
		true,
	)

	url := h.service.GithubConfig.AuthCodeURL(state)
	ctx.Redirect(302, url)
}

func (h *OAuthHandler) HandleGoogleCallback(ctx *gin.Context) {
	cookieState, err := ctx.Cookie("oauth_state")
	if err != nil {
		ctx.JSON(400, gin.H{"error": "missing state cookie"})
		return
	}
	queryState := ctx.Query("state")
	if cookieState != queryState {
		ctx.JSON(400, gin.H{"error": "invalid state"})
		return
	}

	code := ctx.Query("code")
	profile, err := h.service.FetchGoogleProfile(ctx.Request.Context(), code)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "cannot fetch google profile"})
		return
	}

	user, err := h.service.FindOrCreateUserFromOAuth(
		models.ProviderGoogle,
		profile.Sub,
		profile.Email,
		profile.Name,
		profile.Picture,
	)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "cannot fetch google profile"})
		return
	}

	accessToken, err := h.service.IssueTokens(user.ID, user.Role)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "cannot issue tokens"})
		return
	}
	ctx.Redirect(302, h.service.FrontendURL+"/auth/oauth/callback?token="+accessToken)
}

func (h *OAuthHandler) HandleGithubCallback(ctx *gin.Context) {
	cookieState, err := ctx.Cookie("oauth_state")
	if err != nil {
		ctx.JSON(400, gin.H{"error": "missing state cookie"})
		return
	}
	queryState := ctx.Query("state")
	if cookieState != queryState {
		ctx.JSON(400, gin.H{"error": "invalid state"})
		return
	}

	code := ctx.Query("code")
	profile, err := h.service.FetchGithubProfile(ctx.Request.Context(), code)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "cannot fetch github profile"})
		return
	}

	user, err := h.service.FindOrCreateUserFromOAuth(
		models.ProviderGithub,
		strconv.Itoa(profile.ID),
		profile.Email,
		profile.Login,
		profile.AvatarURL,
	)

	if err != nil {
		ctx.JSON(500, gin.H{"error": "cannot create user from github"})
		return
	}

	accessToken, err := h.service.IssueTokens(user.ID, user.Role)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "cannot issue tokens"})
		return
	}
	ctx.Redirect(302, h.service.FrontendURL+"/auth/oauth/callback?token="+accessToken)

}
