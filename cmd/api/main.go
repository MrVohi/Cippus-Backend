package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"cippus-backend/config"
	"cippus-backend/internal/handlers"
	"cippus-backend/internal/middleware"
	"cippus-backend/internal/services"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error while loading .env file: ", err)
	}
	cfg := config.Load()
	config.InitLogger(cfg.LogLevel)

	if cfg.LogLevel != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.SlogLogger())

	db, err := config.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Error while connecting to db: ", err)
	}

	authService := services.NewAuthService(db, cfg.JWTSecret, cfg.ResendApiKey, cfg.RecaptchaSecret)
	authHandler := handlers.NewAuthHandler(authService)

	userService := services.NewUserService(db)
	userHandler := handlers.NewUserHandler(userService)

	postService := services.NewPostService(db)
	postHandler := handlers.NewPostHandler(postService)

	categoryService := services.NewCategoryService(db)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	projectService := services.NewProjectService(db)
	projectHandler := handlers.NewProjectHandler(projectService)

	minioService, err := services.NewMinioService(cfg.StorageEndpoint, cfg.StorageAccessKey, cfg.StorageSecretKey, cfg.StorageBucket)
	if err != nil {
		log.Fatal(err)
	}
	minioHandler := handlers.NewMinioHandler(minioService, postService)

	setupRoutes(router, authHandler, userHandler, postHandler, minioHandler, categoryHandler, projectHandler, &cfg)

	if err := router.Run(cfg.Port); err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}
