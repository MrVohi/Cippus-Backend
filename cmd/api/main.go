package main

import (
	"log"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"cippus-backend/config"
)

func main() {
	router := gin.Default()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error while loading .env file: ", err)
	}
	cfg := config.Load()
	config.InitLogger(cfg.LogLevel)
	_, err = config.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Error while connecting to db: ", err)
	} else {
		slog.Info("Connection to database succesfull!")
		router.Run(cfg.Port)
	}
}
