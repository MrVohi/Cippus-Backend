package config

import (
	"log"

	"github.com/caarlos0/env/v11"
)

type Config struct{
	// Server
	Port string `env:"PORT"`
	LogLevel string `env:"LOG_LEVEL"`
	FrontendURL string `env:"FRONTEND_URL,required"`

	// Auth
	JWTSecret string `env:"JWT_SECRET,required"`
	GithubClientID string `env:"GITHUB_CLIENT_ID,required"`
	GithubClientSecret string `env:"GITHUB_CLIENT_SECRET,required"`
	GithubRedirectURL string `env:"GITHUB_REDIRECT_URL,required"`
	GoogleClientID string `env:"GOOGLE_CLIENT_ID,required"`
	GoogleClientSecret string `env:"GOOGLE_CLIENT_SECRET,required"`
	GoogleRedirectURL string `env:"GOOGLE_REDIRECT_URL,required"`

	// Database
	DatabaseURL string `env:"DATABASE_URL,required"`

	// RabbitMQ
	RabbitmqURL string `env:"RABBITMQ_URL,required"`

	// Ollama
	OllamaBaseURL string `env:"OLLAMA_BASE_URL,required"`
	OllamaModel string `env:"OLLAMA_MODEL"`
	OllamaEmbedModel string `env:"OLLAMA_EMBED_MODEL"`
	EmbeddingDim string `env:"EMBEDDING_DIM"`

	// Object storage
	StorageEndpoint string `env:"STORAGE_ENDPOINT"`
	StorageAccessKey string `env:"STORAGE_ACCESS_KEY,required"`
	StorageSecretKey string `env:"STORAGE_SECRET_KEY,required"`
	StorageBucket string `env:"STORAGE_BUCKET"`

	// Web push
	VapidPublicKey string `env:"VAPID_PUBLIC_KEY,required"`
	VapidPrivateKey string `env:"VAPID_PRIVATE_KEY,required"`
	VapidSubject string `env:"VAPID_SUBJECT"`
	
	// Email
	ResendApiKey string `env:"RESEND_API_KEY,required"`
	EmailFrom string `env:"EMAIL_FROM"`

	// Captcha
	RecaptchaSecret string `env:"RECAPTCHA_SECRET,required"`
}

func Load() Config{
	cfg := Config{}
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal("Failed to read env!")
	}
	return cfg
}