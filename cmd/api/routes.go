package main

import (
	"cippus-backend/config"
	"cippus-backend/internal/handlers"
	"cippus-backend/internal/middleware"
	"cippus-backend/internal/models"
	"cippus-backend/internal/ws"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func authRoutes(api *gin.RouterGroup, authHandler *handlers.AuthHandler, cfg *config.Config) {
	public := api.Group("/auth")
	private := api.Group("/").Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		public.POST("/register", middleware.RateLimiter(3), authHandler.RegisterHandler)
		public.POST("/login", middleware.RateLimiter(5), authHandler.LoginHandler)
		public.POST("/refresh", authHandler.RefreshHandler)
		public.POST("/password-reset/request", middleware.RateLimiter(3), authHandler.PasswordResetRequestHandler)
		public.POST("/password-reset/confirm", middleware.RateLimiter(5), authHandler.PasswordResetConfirmHandler)
		private.POST("/auth/logout", authHandler.LogoutHandler)
	}
}

func postRoutes(api *gin.RouterGroup, postHandler *handlers.PostHandler, cfg *config.Config) {
	publicPosts := api.Group("/logs")
	privatePosts := api.Group("/logs").Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		publicPosts.GET("", postHandler.GetPostsHandler)
		publicPosts.GET("/:id", postHandler.GetPostsByIdHandler)
		privatePosts.POST("", postHandler.CreatePostHandler)
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

func categoryRoutes(api *gin.RouterGroup, categoryHandler *handlers.CategoryHandler, cfg *config.Config) {
	publicCategories := api.Group("/categories")
	privateCategories := api.Group("/categories").Use(middleware.AuthMiddleware(cfg.JWTSecret)).Use(middleware.RequireRole(models.RoleModerator, models.RoleAdmin))
	{
		publicCategories.GET("", categoryHandler.GetCategoriesHandler)
		privateCategories.POST("", categoryHandler.CreateCategoryHandler)
		privateCategories.PATCH("/:id", categoryHandler.UpdateCategoryHandler)
		privateCategories.DELETE("/:id", categoryHandler.DeleteCategoryHandler)
	}
}

func projectRoutes(api *gin.RouterGroup, projectHandler *handlers.ProjectHandler, cfg *config.Config) {
	publicProjects := api.Group("/projects")
	privateProjects := api.Group("/projects").Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		api.GET("/users/:userId/projects", projectHandler.GetProjectsByUserHandler)
		publicProjects.GET("/:id", projectHandler.GetProjectsByIdHandler)
		privateProjects.POST("", projectHandler.CreateProjectHandler)
		privateProjects.PATCH("/:id", projectHandler.UpdateProjectHandler)
		privateProjects.DELETE("/:id", projectHandler.DeleteProjectHandler)
		privateProjects.POST("/:id/posts/:postId", projectHandler.AddPostToProjectHandler)
		privateProjects.DELETE("/:id/posts/:postId", projectHandler.DeletePostFromProjectHandler)
	}
}

func messagesRoutes(api *gin.RouterGroup, messageHandler *handlers.MessageHandler, cfg *config.Config) {
	messages := api.Group("/messages").Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		api.GET("/ws", ws.ServeWS(messageHandler.Hub, cfg.JWTSecret))
		messages.POST("", messageHandler.Send)
		messages.GET("", messageHandler.GetConversations)
		messages.GET("/:userID", messageHandler.GetThread)
		messages.PUT("/:messageID/read", messageHandler.MarkAsRead)
	}
}

func pushRoutes(api *gin.RouterGroup, pushHandler *handlers.PushHandler, cfg *config.Config) {
	pushs := api.Group("/push").Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		api.GET("/push/vapid-public-key", pushHandler.GetVapidKey)
		pushs.POST("/subscribe", pushHandler.Subscribe)
		pushs.DELETE("/subscribe", pushHandler.Unsubscribe)
	}
}

func userRoutes(api *gin.RouterGroup, authHandler *handlers.AuthHandler, userHandler *handlers.UserHandler, cfg *config.Config) {
	public := api.Group("/users")
	private := api.Group("/users").Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		private.PATCH("/me", userHandler.PatchMe)
		private.PATCH("/me/avatar", userHandler.UploadAvatar)
		private.DELETE("/me/avatar", userHandler.DeleteAvatar)
		private.GET("/me", userHandler.GetMe)
		public.GET("/:userId", userHandler.GetUser)
	}
}

func notificationRoutes(api *gin.RouterGroup, notificationHandler *handlers.NotificationHandler, cfg *config.Config) {
	notifications := api.Group("/notifications").Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		notifications.GET("", notificationHandler.GetNotifications)
	}
}

func searchRoutes(api *gin.RouterGroup, searchHandler *handlers.SearchHandler, cfg *config.Config) {
	search := api.Group("")
	search.GET("/search", searchHandler.SearchHandler)
	search.GET("/logs/:id/similar", searchHandler.GetSimilarHandler)
}

func commentsRoutes(api *gin.RouterGroup, commentHandler *handlers.CommentHandler, cfg *config.Config) {
	public := api.Group("/logs")
	private := api.Group("/logs").Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		public.GET("/:id/comments", commentHandler.GetCommentsByPostIDHandler)
		private.POST("/:id/comments", commentHandler.CreateCommentHandler)
		private.PATCH("/:id/comments/:cid", commentHandler.UpdateCommentHandler)
		private.DELETE("/:id/comments/:cid", commentHandler.DeleteCommentHandler)
	}
}

func postLikeRoutes(api *gin.RouterGroup, postLikeHandler *handlers.PostLikeHandler, cfg *config.Config) {
	optAuth := api.Group("/logs").Use(middleware.OptionalAuthMiddleware(cfg.JWTSecret))
	private := api.Group("/logs").Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		optAuth.GET("/:id/likes", postLikeHandler.GetLikesHandler)
		private.POST("/:id/likes", postLikeHandler.ToggleLikeHandler)
	}
}

func commentLikeRoutes(api *gin.RouterGroup, commentLikeHandler *handlers.CommentLikeHandler, cfg *config.Config) {
	optAuth := api.Group("/logs").Use(middleware.OptionalAuthMiddleware(cfg.JWTSecret))
	private := api.Group("/logs").Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		optAuth.GET("/:id/comments/:cid/likes", commentLikeHandler.GetLikesHandler)
		private.POST("/:id/comments/:cid/likes", commentLikeHandler.ToggleLikeHandler)
	}
}

func oauthRoutes(api *gin.RouterGroup, oauthHandler *handlers.OAuthHandler, cfg *config.Config) {
	public := api.Group("/auth")

	{
		public.GET("/google", oauthHandler.HandleGoogleInit)
		public.GET("/github", oauthHandler.HandleGithubInit)
		public.GET("/google/callback", oauthHandler.HandleGoogleCallback)
		public.GET("/github/callback", oauthHandler.HandleGithubCallback)
	}
}

func aiRoutes(api *gin.RouterGroup, aiHandler *handlers.AIHandler, cfg *config.Config) {
	private := api.Group("/ai").Use(middleware.AuthMiddleware(cfg.JWTSecret))

	{
		private.POST("/improve", middleware.RateLimiter(5), aiHandler.ImproveHandler)
	}
}

func setupRoutes(router *gin.Engine, authHandler *handlers.AuthHandler, userHandler *handlers.UserHandler, postHandler *handlers.PostHandler, minioHandler *handlers.MinioHandler, categoryHandler *handlers.CategoryHandler, projectHandler *handlers.ProjectHandler, messageHandler *handlers.MessageHandler, pushHandler *handlers.PushHandler, notificationHandler *handlers.NotificationHandler, searchHandler *handlers.SearchHandler, commentHandler *handlers.CommentHandler, postLikeHandler *handlers.PostLikeHandler, commentLikeHandler *handlers.CommentLikeHandler, oauthHandler *handlers.OAuthHandler, aiHandler *handlers.AIHandler, cfg *config.Config) {
	corsCfg := cors.DefaultConfig()
	corsCfg.AllowOrigins = []string{cfg.FrontendURL}
	corsCfg.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	corsCfg.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	corsCfg.AllowCredentials = true

	corsMw := cors.New(corsCfg)
	router.Use(corsMw)

	api := router.Group("/api/v1")

	authRoutes(api, authHandler, cfg)
	userRoutes(api, authHandler, userHandler, cfg)
	postRoutes(api, postHandler, cfg)
	imageRoutes(api, minioHandler, cfg)
	categoryRoutes(api, categoryHandler, cfg)
	projectRoutes(api, projectHandler, cfg)
	messagesRoutes(api, messageHandler, cfg)
	pushRoutes(api, pushHandler, cfg)
	notificationRoutes(api, notificationHandler, cfg)
	searchRoutes(api, searchHandler, cfg)
	commentsRoutes(api, commentHandler, cfg)
	postLikeRoutes(api, postLikeHandler, cfg)
	commentLikeRoutes(api, commentLikeHandler, cfg)
	oauthRoutes(api, oauthHandler, cfg)
	aiRoutes(api, aiHandler, cfg)
}
