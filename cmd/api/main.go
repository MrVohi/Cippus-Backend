package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"cippus-backend/config"
	"cippus-backend/internal/handlers"
	"cippus-backend/internal/middleware"
	"cippus-backend/internal/services"
	"cippus-backend/internal/ws"
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
	router.Use(func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 4<<20)
		c.Next()
	})

	db, err := config.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Error while connecting to db: ", err)
	}

	authService := services.NewAuthService(db, cfg.JWTSecret, cfg.FrontendURL, cfg.ResendApiKey, cfg.RecaptchaSecret)
	authHandler := handlers.NewAuthHandler(authService)

	userService := services.NewUserService(db)
	userHandler := handlers.NewUserHandler(userService)

	embeddingService := services.NewEmbeddingService(cfg)
	postService := services.NewPostService(db, embeddingService)
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

	rabbit, err := services.NewRabbitPublisher(cfg.RabbitmqURL)
	if err != nil {
		log.Fatal("Error while connecting to RabbitMQ: ", err)
	}
	defer rabbit.Close()

	notificationService := services.NewNotificationService(db, rabbit)
	notificationHandler := handlers.NewNotificationHandler(notificationService)

	messageService := services.NewMessageService(db)
	hub := ws.NewHub()
	go hub.Run()
	messageHandler := handlers.NewMessageHandler(messageService, hub, notificationService)

	pushService := services.NewPushService(db, cfg.VapidPublicKey, cfg.VapidPrivateKey, cfg.VapidSubject)
	pushHandler := handlers.NewPushHandler(pushService, cfg.VapidPublicKey)

	searchService := services.NewSearchService(db, embeddingService)
	searchHandler := handlers.NewSearchHandler(searchService)

	setupRoutes(router, authHandler, userHandler, postHandler, minioHandler, categoryHandler, projectHandler, messageHandler, pushHandler, notificationHandler, searchHandler, &cfg)

	srv := &http.Server{
		Addr:              cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server failed", "err", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("graceful shutdown failed", "err", err)
	}
}
