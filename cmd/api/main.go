package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"cippus-backend/config"
	"cippus-backend/internal/handlers"
	"cippus-backend/internal/services"
)

func main() {
	router := gin.Default()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error while loading .env file: ", err)
	}
	cfg := config.Load()
	config.InitLogger(cfg.LogLevel)
	db, err := config.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Error while connecting to db: ", err)
	}

	authService := services.NewAuthService(db, cfg.JWTSecret, cfg.ResendApiKey, cfg.RecaptchaSecret)
	authHandler := handlers.NewAuthHandler(authService)

	setupRoutes(router, authHandler, &cfg)

	if err := router.Run(cfg.Port); err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}
