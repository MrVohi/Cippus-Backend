package main

import (
	"cippus-backend/config"
	"cippus-backend/internal/handlers"
	"cippus-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func authRoutes(api *gin.RouterGroup, authHandler *handlers.AuthHandler) {
	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.RegisterHandler)
		auth.POST("/login", authHandler.LoginHandler)
		auth.POST("/refresh", authHandler.RefreshHandler)
		auth.POST("/password-reset/request", authHandler.PasswordResetRequestHandler)
		auth.POST("/password-reset/confirm", authHandler.PasswordResetConfirmHandler)
	}
}

func privateRoutes(api *gin.RouterGroup, authHandler *handlers.AuthHandler, cfg *config.Config) {
	private := api.Group("/").Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		private.POST("/auth/logout", authHandler.LogoutHandler)
	}
}

func setupRoutes(router *gin.Engine, authHandler *handlers.AuthHandler, cfg *config.Config) {
	api := router.Group("/api/v1")

	authRoutes(api, authHandler)
	privateRoutes(api, authHandler, cfg)
	
}