package services

import (
	"cippus-backend/internal/models"

	"gorm.io/gorm"
)

type AuthService struct {
	db              *gorm.DB
	secret          string
	resendApiKey    string
	recaptchaSecret string
}

type AuthResult struct {
	AccessToken  string
	RefreshToken string
	User         models.User
}
