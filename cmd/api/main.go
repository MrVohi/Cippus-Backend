package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"cippus-backend/config"
)

func main() {
	router := gin.Default()
	err := godotenv.Load()
  	if err != nil {
    	log.Fatal("Error loading .env file")
  	}
	cfg := config.Load()
	config.InitLogger(cfg.LogLevel)
	router.Run(cfg.Port)
}