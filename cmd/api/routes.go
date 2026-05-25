package main

import (
	"cippus-backend/config"
	"cippus-backend/internal/handlers"
	"cippus-backend/internal/middleware"

	"github.com/gin-contrib/cors"
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

func postRoutes(api *gin.RouterGroup, postHandler *handlers.PostHandler, cfg *config.Config) {
	publicPosts := api.Group("/logs")
	privatePosts := api.Group("/logs").Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		publicPosts.GET("/", postHandler.GetPostsHandler)
		publicPosts.GET("/:id", postHandler.GetPostsByIdHandler)
		privatePosts.POST("/", postHandler.CreatePostHandler)
		privatePosts.PATCH("/:id", postHandler.UpdatePostHandler)
		privatePosts.DELETE("/:id", postHandler.DeletePostHandler)		
	}
}

func imageRoutes(api *gin.RouterGroup, minioHandler *handlers.MinioHandler, cfg *config.Config) {
	images := api.Group("/logs").Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		images.POST("/:id/image", minioHandler.UploadImageHandler)
		images.DELETE("/:id/image", minioHandler.DeleteImageHandler)
	}
}


func privateRoutes(api *gin.RouterGroup, authHandler *handlers.AuthHandler, userHandler *handlers.UserHandler, cfg *config.Config) {
	private := api.Group("/").Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		private.POST("/auth/logout", authHandler.LogoutHandler)
		private.PATCH("/users/me", userHandler.PatchMe)
	}
}

func setupRoutes(router *gin.Engine, authHandler *handlers.AuthHandler, userHandler *handlers.UserHandler, postHandler *handlers.PostHandler, minioHandler *handlers.MinioHandler, cfg *config.Config) {
	corsCfg := cors.DefaultConfig()
	corsCfg.AllowOrigins = []string{"http://localhost:3000"}
	corsCfg.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	corsCfg.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	corsCfg.AllowCredentials = true

	corsMw := cors.New(corsCfg)
	router.Use(corsMw)

	api := router.Group("/api/v1")

	authRoutes(api, authHandler)
	privateRoutes(api, authHandler, userHandler, cfg)
	postRoutes(api, postHandler, cfg)
	imageRoutes(api, minioHandler, cfg)

}
